package gcache

import (
	"context"
)

// Store is a cache store.
type Store interface {
	Get(ctx context.Context, key string) (e Entry, ok bool)
	Set(ctx context.Context, key string, e Entry)
	Remove(ctx context.Context, key string)
}

// Entry is a cache entry to store.
type Entry struct {
	Value []byte `json:"value"`
	ETag  string `json:"etag"`
}

// LRUBackend specifies interface to be implemented by hashicorp LRU cache backends.
type LRUBackend interface {
	Add(key string, value Entry) (evicted bool)
	Get(key string) (value Entry, ok bool)
	Remove(key string) (present bool)
}

type lruWrapper struct{ backend LRUBackend }

// NewLRU wraps hashicorp/golang-lru/v2 cache implementations to be used as interceptor's store.
func NewLRU(backend LRUBackend) Store { return lruWrapper{backend: backend} }

func (l lruWrapper) Get(_ context.Context, key string) (e Entry, ok bool) { return l.backend.Get(key) }
func (l lruWrapper) Set(_ context.Context, key string, e Entry)           { l.backend.Add(key, e) }
func (l lruWrapper) Remove(_ context.Context, key string)                 { l.backend.Remove(key) }
