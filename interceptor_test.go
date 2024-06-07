package gcache

import (
	"context"
	"log/slog"
	"regexp"
	"testing"

	"github.com/cappuccinotm/gcache/internal/tspb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

// hash that is generated for the empty request
const emptyReqKey = "/com.github.cappuccinotm.gcache.example.TestService/Test{da39a3ee5e6b4b0d3255bfef95601890afd80709}"

func TestInterceptor_UnaryClientInterceptor(t *testing.T) {
	t.Run("filtered out, must not be cached", func(t *testing.T) {
		addr := tspb.Run(t, tspb.MockTestService{
			TestFunc: func(ctx context.Context, in *tspb.TestRequest) (*tspb.TestResponse, error) {
				err := grpc.SendHeader(ctx, metadata.Pairs("ETag", "must-not-be-cached"))
				require.NoError(t, err)
				return &tspb.TestResponse{Value: "must-not-be-cached"}, nil
			},
		})

		icptr := NewInterceptor(WithFilter(regexp.MustCompile(`blah-will-never-match`)))

		cc, err := grpc.NewClient(addr,
			grpc.WithUnaryInterceptor(icptr.UnaryClientInterceptor()),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		require.NoError(t, err)

		cl := tspb.NewTestServiceClient(cc)

		resp, err := cl.Test(context.Background(), &tspb.TestRequest{})
		require.NoError(t, err)
		assert.Equal(t, "must-not-be-cached", resp.Value)

		_, ok := icptr.store.Get(context.Background(), emptyReqKey)
		require.False(t, ok)
	})

	t.Run("no value stored, cache new", func(t *testing.T) {
		addr := tspb.Run(t, tspb.MockTestService{
			TestFunc: func(ctx context.Context, in *tspb.TestRequest) (*tspb.TestResponse, error) {
				err := grpc.SendHeader(ctx, metadata.Pairs("ETag", "must-be-cached"))
				require.NoError(t, err)
				return &tspb.TestResponse{Value: "must-be-cached"}, nil
			},
		})

		icptr := NewInterceptor()

		cc, err := grpc.NewClient(addr,
			grpc.WithUnaryInterceptor(icptr.UnaryClientInterceptor()),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		require.NoError(t, err)

		cl := tspb.NewTestServiceClient(cc)

		resp, err := cl.Test(context.Background(), &tspb.TestRequest{})
		require.NoError(t, err)
		assert.Equal(t, "must-be-cached", resp.Value)

		e, ok := icptr.store.Get(context.Background(), emptyReqKey)
		require.True(t, ok)
		assert.Equal(t, "must-be-cached", e.ETag)
	})

	t.Run("no value stored, server returned no etag", func(t *testing.T) {
		addr := tspb.Run(t, tspb.MockTestService{
			TestFunc: func(ctx context.Context, in *tspb.TestRequest) (*tspb.TestResponse, error) {
				return &tspb.TestResponse{Value: "must-be-cached"}, nil
			},
		})

		icptr := NewInterceptor()

		cc, err := grpc.NewClient(addr,
			grpc.WithUnaryInterceptor(icptr.UnaryClientInterceptor()),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		require.NoError(t, err)

		cl := tspb.NewTestServiceClient(cc)

		resp, err := cl.Test(context.Background(), &tspb.TestRequest{})
		require.NoError(t, err)
		assert.Equal(t, "must-be-cached", resp.Value)

		_, ok := icptr.store.Get(context.Background(), emptyReqKey)
		require.False(t, ok)
	})

	t.Run("value stored, use cached one", func(t *testing.T) {
		addr := tspb.Run(t, tspb.MockTestService{
			TestFunc: func(ctx context.Context, in *tspb.TestRequest) (*tspb.TestResponse, error) {
				assert.Equal(t, "use-cached", ETag(ctx))
				return nil, NotChanged(ctx, "use-cached")
			},
		})

		icptr := NewInterceptor()

		bts, err := RawBytesCodec{}.Marshal(&tspb.TestResponse{Value: "use-cached"})
		require.NoError(t, err)

		icptr.store.Set(context.Background(), emptyReqKey, Entry{Value: bts, ETag: "use-cached"})

		cc, err := grpc.NewClient(addr,
			grpc.WithUnaryInterceptor(icptr.UnaryClientInterceptor()),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		require.NoError(t, err)

		cl := tspb.NewTestServiceClient(cc)

		resp, err := cl.Test(context.Background(), &tspb.TestRequest{})
		require.NoError(t, err)
		assert.Equal(t, "use-cached", resp.Value)
	})

	t.Run("value stored, outdated", func(t *testing.T) {
		addr := tspb.Run(t, tspb.MockTestService{
			TestFunc: func(ctx context.Context, in *tspb.TestRequest) (*tspb.TestResponse, error) {
				err := grpc.SendHeader(ctx, metadata.Pairs("ETag", "update"))
				require.NoError(t, err)
				return &tspb.TestResponse{Value: "update"}, nil
			},
		})

		icptr := NewInterceptor()

		bts, err := RawBytesCodec{}.Marshal(&tspb.TestResponse{Value: "use-cached"})
		require.NoError(t, err)

		icptr.store.Set(context.Background(), emptyReqKey, Entry{Value: bts, ETag: "use-cached"})

		cc, err := grpc.NewClient(addr,
			grpc.WithUnaryInterceptor(icptr.UnaryClientInterceptor()),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		require.NoError(t, err)

		cl := tspb.NewTestServiceClient(cc)

		resp, err := cl.Test(context.Background(), &tspb.TestRequest{})
		require.NoError(t, err)
		assert.Equal(t, "update", resp.Value)

		e, ok := icptr.store.Get(context.Background(), emptyReqKey)
		require.True(t, ok)

		assert.Equal(t, "update", e.ETag)
	})

	t.Run("value stored, server returned no etag", func(t *testing.T) {
		addr := tspb.Run(t, tspb.MockTestService{
			TestFunc: func(ctx context.Context, in *tspb.TestRequest) (*tspb.TestResponse, error) {
				return &tspb.TestResponse{Value: "update"}, nil
			},
		})

		icptr := NewInterceptor()

		bts, err := RawBytesCodec{}.Marshal(&tspb.TestResponse{Value: "use-cached"})
		require.NoError(t, err)

		icptr.store.Set(context.Background(), emptyReqKey, Entry{Value: bts, ETag: "use-cached"})

		cc, err := grpc.NewClient(addr,
			grpc.WithUnaryInterceptor(icptr.UnaryClientInterceptor()),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		require.NoError(t, err)

		cl := tspb.NewTestServiceClient(cc)

		resp, err := cl.Test(context.Background(), &tspb.TestRequest{})
		require.NoError(t, err)
		assert.Equal(t, "update", resp.Value)

		_, ok := icptr.store.Get(context.Background(), emptyReqKey)
		require.False(t, ok)
	})
}

func TestInterceptor_UnaryServerInterceptor(t *testing.T) {
	t.Run("filtered out, must not be cached", func(t *testing.T) {
		icptr := NewInterceptor(WithFilter(regexp.MustCompile(`blah-will-never-match`)))

		addr := tspb.Run(t, tspb.MockTestService{
			TestFunc: func(ctx context.Context, in *tspb.TestRequest) (*tspb.TestResponse, error) {
				return &tspb.TestResponse{Value: "must-not-be-cached"}, nil
			}},
			grpc.UnaryInterceptor(icptr.UnaryServerInterceptor()),
		)

		cc, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)

		cl := tspb.NewTestServiceClient(cc)

		resp, err := cl.Test(context.Background(), &tspb.TestRequest{})
		require.NoError(t, err)
		assert.Equal(t, "must-not-be-cached", resp.Value)

		_, ok := icptr.store.Get(context.Background(), emptyReqKey)
		require.False(t, ok)
	})

	t.Run("no value stored, cache new", func(t *testing.T) {
		icptr := NewInterceptor()

		addr := tspb.Run(t, tspb.MockTestService{
			TestFunc: func(ctx context.Context, in *tspb.TestRequest) (*tspb.TestResponse, error) {
				return &tspb.TestResponse{Value: "must-not-be-cached"}, nil
			},
		}, grpc.UnaryInterceptor(icptr.UnaryServerInterceptor()))

		cc, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)

		cl := tspb.NewTestServiceClient(cc)

		resp, err := cl.Test(context.Background(), &tspb.TestRequest{})
		require.NoError(t, err)
		assert.Equal(t, "must-not-be-cached", resp.Value)

		e, ok := icptr.store.Get(context.Background(), emptyReqKey)
		require.True(t, ok)

		expected, err := proto.Marshal(&tspb.TestResponse{Value: "must-not-be-cached"})
		require.NoError(t, err)
		assert.Equal(t, expected, e.Value)
	})

	t.Run("call thrice", func(t *testing.T) {
		// first time - we just cache value, we don't get the response of the handler
		// second time - we seek for the response type and cache it internally
		// third time - we use the cached value AND the cached response type
		icptr := NewInterceptor()

		addr := tspb.Run(t, tspb.MockTestService{
			TestFunc: func(ctx context.Context, in *tspb.TestRequest) (*tspb.TestResponse, error) {
				return &tspb.TestResponse{Value: "must-not-be-cached"}, nil
			},
		}, grpc.UnaryInterceptor(icptr.UnaryServerInterceptor()))

		cc, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)

		cl := tspb.NewTestServiceClient(cc)

		resp, err := cl.Test(context.Background(), &tspb.TestRequest{})
		require.NoError(t, err)
		assert.Equal(t, "must-not-be-cached", resp.Value)

		resp, err = cl.Test(context.Background(), &tspb.TestRequest{})
		require.NoError(t, err)
		assert.Equal(t, "must-not-be-cached", resp.Value)

		resp, err = cl.Test(context.Background(), &tspb.TestRequest{})
		require.NoError(t, err)
		assert.Equal(t, "must-not-be-cached", resp.Value)
	})

	t.Run("value stored, use cached one", func(t *testing.T) {
		icptr := NewInterceptor(WithLogger(slog.Default()))

		bts, err := RawBytesCodec{}.Marshal(&tspb.TestResponse{Value: "use-cached"})
		require.NoError(t, err)

		icptr.store.Set(context.Background(), emptyReqKey, Entry{Value: bts})

		addr := tspb.Run(t, tspb.MockTestService{
			TestFunc: func(ctx context.Context, in *tspb.TestRequest) (*tspb.TestResponse, error) {
				require.Fail(t, "must not be called")
				return nil, nil
			},
		}, grpc.UnaryInterceptor(icptr.UnaryServerInterceptor()))

		cc, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)

		cl := tspb.NewTestServiceClient(cc)

		resp, err := cl.Test(context.Background(), &tspb.TestRequest{})
		require.NoError(t, err)

		assert.Equal(t, "use-cached", resp.Value)
	})

	t.Run("value stored, client sent no-cache", func(t *testing.T) {
		icptr := NewInterceptor(WithLogger(slog.Default()))

		bts, err := RawBytesCodec{}.Marshal(&tspb.TestResponse{Value: "must-not-be-responded"})
		require.NoError(t, err)

		icptr.store.Set(context.Background(), emptyReqKey, Entry{Value: bts})

		addr := tspb.Run(t, tspb.MockTestService{
			TestFunc: func(ctx context.Context, in *tspb.TestRequest) (*tspb.TestResponse, error) {
				return &tspb.TestResponse{Value: "success"}, nil
			},
		}, grpc.UnaryInterceptor(icptr.UnaryServerInterceptor()))

		cc, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)

		cl := tspb.NewTestServiceClient(cc)

		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("Cache-Control", "no-cache"))
		resp, err := cl.Test(ctx, &tspb.TestRequest{})
		require.NoError(t, err)
		assert.Equal(t, "success", resp.Value)
	})
}
