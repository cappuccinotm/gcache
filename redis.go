package gcache

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"time"

	rediscache "github.com/go-redis/cache/v9"
)

// RedisOption is a configuration option.
type RedisOption func(*redisStore)

// WithRedisLogger sets the logger.
func WithRedisLogger(l *slog.Logger) RedisOption {
	return func(r *redisStore) { r.logger = l }
}

// WithRedisTTL sets the TTL.
func WithRedisTTL(ttl time.Duration) RedisOption {
	return func(r *redisStore) { r.ttl = ttl }
}

// WithRedisSkipLocalCache sets the skipLocalCache.
func WithRedisSkipLocalCache(skipLocalCache bool) RedisOption {
	return func(r *redisStore) { r.skipLocalCache = skipLocalCache }
}

type redisStore struct {
	backend        *rediscache.Cache
	logger         *slog.Logger
	ttl            time.Duration
	skipLocalCache bool
}

// NewRedis returns a new redisStore cache store.
func NewRedis(backend *rediscache.Cache, opts ...RedisOption) Store {
	store := &redisStore{
		backend: backend,
		logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
	}

	for _, opt := range opts {
		opt(store)
	}

	return store
}

// Get returns the value for the given key.
func (r *redisStore) Get(ctx context.Context, key string) (e Entry, ok bool) {
	switch err := r.backend.Get(ctx, key, &e); {
	case errors.Is(err, rediscache.ErrCacheMiss):
		return Entry{}, false
	case err != nil:
		r.logger.WarnContext(ctx, "gcache: failed to get from redisStore cache", slog.Any(ErrKey, err))
	}
	return e, true
}

// Set sets the value for the given key.
func (r *redisStore) Set(ctx context.Context, key string, e Entry) {
	item := &rediscache.Item{
		Ctx:            ctx,
		Key:            key,
		Value:          e,
		TTL:            r.ttl,
		SkipLocalCache: r.skipLocalCache,
	}

	if err := r.backend.Set(item); err != nil {
		r.logger.WarnContext(ctx, "gcache: failed to set to redisStore cache", slog.Any(ErrKey, err))
	}
}

// Remove removes the value for the given key.
func (r *redisStore) Remove(ctx context.Context, key string) {
	if err := r.backend.Delete(ctx, key); err != nil {
		r.logger.WarnContext(ctx, "gcache: failed to remove from redisStore cache", slog.Any(ErrKey, err))
	}
}
