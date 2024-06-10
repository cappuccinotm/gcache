// Package gcache provides caching interceptors for gRPC clients.
package gcache

import (
	"context"
	"crypto/sha1" //nolint: gosec // we use sha1 for hashing
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"sync"

	lru "github.com/hashicorp/golang-lru/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Interceptor is a cache interceptor.
// It looks over the ETag header and caches the response ONLY if the ETag is present.
type Interceptor struct {
	store  Store
	logger *slog.Logger
	codec  encoding.Codec
	filter *regexp.Regexp
}

// NewInterceptor makes a new Interceptor.
func NewInterceptor(opts ...Option) *Interceptor {
	c := &Interceptor{
		codec:  RawBytesCodec{},
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		filter: regexp.MustCompile(`.*`),
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.store == nil { // lazy init for LRU
		l, _ := lru.New[string, Entry](1024)
		c.store = NewLRU(l)
	}

	return c
}

// UnaryServerInterceptor returns a new unary server interceptor that caches the response.
// It doesn't use ETag header, but Cache-Control header.
func (c *Interceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		if !c.filter.MatchString(info.FullMethod) {
			return handler(ctx, req)
		}

		if inMD, ok := metadata.FromIncomingContext(ctx); ok {
			if h := inMD.Get("Cache-Control"); len(h) > 0 && h[0] == "no-cache" {
				return handler(ctx, req)
			}
		}

		key, err := c.key(info.FullMethod, req)
		if err != nil {
			c.logger.WarnContext(ctx, "gcache: failed to produce a key, skipping cache",
				slog.Any(ErrKey, err))
			return handler(ctx, req)
		}

		if e, ok := c.store.Get(ctx, key); ok {
			if resp, err = c.buildResponse(info, e); err == nil {
				return resp, nil
			}

			c.logger.WarnContext(ctx, "gcache: failed to unmarshal response, retrieving from handler",
				slog.Any(ErrKey, err))
		}

		if resp, err = handler(ctx, req); err != nil {
			return nil, err
		}

		bts, err := c.codec.Marshal(resp)
		if err != nil {
			c.logger.WarnContext(ctx, "gcache: failed to marshal response, value won't be cached",
				slog.Any(ErrKey, err))
			return resp, nil
		}

		c.store.Set(ctx, key, Entry{Value: bts})
		return resp, nil
	}
}

// UnaryClientInterceptor returns a new unary client interceptor that caches the response.
func (c *Interceptor) UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if !c.filter.MatchString(method) {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		key, err := c.key(method, req)
		if err != nil {
			c.logger.WarnContext(ctx, "gcache: failed to produce a key, skipping cache",
				slog.Any(ErrKey, err))
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		cachedValue, ok := c.store.Get(ctx, key)
		if ok {
			outMD, ok := metadata.FromOutgoingContext(ctx)
			if !ok {
				outMD = metadata.MD{}
			}
			outMD.Set("If-None-Match", cachedValue.ETag)
			ctx = metadata.NewOutgoingContext(ctx, outMD)
		}

		inMD := &metadata.MD{}
		opts = append(opts, grpc.Header(inMD), grpc.ForceCodec(c.codec))

		var raw []byte
		switch err = invoker(ctx, method, req, &raw, cc, opts...); {
		case notChanged(ctx, err, inMD):
			if err = c.codec.Unmarshal(cachedValue.Value, reply); err != nil {
				return fmt.Errorf("unmarshal cached response: %w", err)
			}
		case err != nil:
			return fmt.Errorf("call invoker: %w", err)
		}

		if err = c.codec.Unmarshal(raw, reply); err != nil {
			return fmt.Errorf("unmarshal response: %w", err)
		}

		if etag := inMD.Get("ETag"); len(etag) != 0 {
			c.store.Set(ctx, key, Entry{Value: raw, ETag: etag[0]})
		} else {
			c.store.Remove(ctx, key)
		}

		return nil
	}
}

var responseCache sync.Map

func (c *Interceptor) responseType(fullMethodName string, srv any) (any, error) {
	if v, ok := responseCache.Load(fullMethodName); ok {
		return reflect.New(v.(reflect.Type)).Interface(), nil
	}

	parts := strings.Split(fullMethodName, "/")
	method := parts[len(parts)-1]

	typ := reflect.TypeOf(srv)

	m, ok := typ.MethodByName(method)
	if !ok {
		return nil, fmt.Errorf("method %s not found in type %s", method, typ.Name())
	}

	respTyp := m.Type.Out(0).Elem()
	responseCache.Store(fullMethodName, respTyp)
	return reflect.New(respTyp).Interface(), nil
}

func (c *Interceptor) buildResponse(info *grpc.UnaryServerInfo, e Entry) (any, error) {
	out, err := c.responseType(info.FullMethod, info.Server)
	if err != nil {
		return nil, fmt.Errorf("get response type: %w", err)
	}

	if err = c.codec.Unmarshal(e.Value, out); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return out, nil
}

func (c *Interceptor) key(method string, req interface{}) (string, error) {
	bts, err := c.codec.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	return fmt.Sprintf("%s{%x}", method, hash(bts)), nil
}

func notChanged(ctx context.Context, err error, inMD *metadata.MD) bool {
	if err == nil {
		return false
	}

	outMD, _ := metadata.FromOutgoingContext(ctx)

	if status.Code(err) != codes.Aborted {
		return false
	}

	if !slices.Equal(outMD.Get("if-none-match"), inMD.Get("etag")) {
		return false
	}

	return true
}

func hash(bts []byte) []byte { h := sha1.Sum(bts); return h[:] } //nolint: gosec // we use sha1 for hashing
