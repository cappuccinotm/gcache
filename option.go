package gcache

import (
	"log/slog"
	"regexp"

	"google.golang.org/grpc/encoding"
)

// Option is a configuration option.
type Option func(*Interceptor)

// WithLogger sets the logger.
func WithLogger(l *slog.Logger) Option { return func(c *Interceptor) { c.logger = l } }

// WithCodec sets the codec.
func WithCodec(codec encoding.Codec) Option { return func(c *Interceptor) { c.codec = codec } }

// WithStore sets the store.
func WithStore(store Store) Option { return func(c *Interceptor) { c.store = store } }

// WithFilter sets the filter that is used to match the methods that
// must be cached.
func WithFilter(rx *regexp.Regexp) Option { return func(c *Interceptor) { c.filter = rx } }
