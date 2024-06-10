# gcache [![Go Reference](https://pkg.go.dev/badge/github.com/cappuccinotm/gcache.svg)](https://pkg.go.dev/github.com/cappuccinotm/gcache) [![Go](https://github.com/cappuccinotm/gcache/actions/workflows/go.yaml/badge.svg)](https://github.com/cappuccinotm/gcache/actions/workflows/go.yaml) [![codecov](https://codecov.io/gh/cappuccinotm/gcache/graph/badge.svg?token=We5B4lzTNj)](https://codecov.io/gh/cappuccinotm/gcache)
gcache is a gRPC caching library that provides a simple way to cache gRPC requests and responses. It is designed to be used with gRPC services that have a high request rate and where caching can improve the performance of the service.

It uses a plain old ETag and If-None-Match headers to cache the responses. The cache is stored in memory and is not persistent (unless you use a persistent cache store).

It also provides a server-side caching implementation.

## Installation
```shell
go get -u github.com/cappuccinotm/gcache
```

## Usage
- [Client-side caching](#client-side-caching)
- [Server-side caching](#server-side-caching)

### Client-side caching
```go
icptr := gcache.NewInterceptor(gcache.WithLogger(slog.Default()))
conn, err := grpc.NewClient("localhost:8080",
    grpc.WithTransportCredentials(insecure.NewCredentials()),
    grpc.WithUnaryInterceptor(icptr.UnaryClientInterceptor()),
)
if err != nil {
    return fmt.Errorf("dial localhost:8080: %w", err)
}

client := order.NewOrderServiceClient(conn)
```

Client-side interceptor seeks for server's `ETag` header in the response, if the server has provided one, it stores the response in the cache. When the client sends a request, the interceptor adds the `If-None-Match` header to the request. If the server responds with code `Aborted` and the `ETag` header, equal to the one that has been sent by client, the interceptor returns the cached response.

### Server-side caching
```go
icptr := gcache.NewInterceptor(gcache.WithLogger(slog.Default()))
server := grpc.NewServer(
    grpc.WithTransportCredentials(insecure.NewCredentials()),
    grpc.WithUnaryInterceptor(icptr.UnaryServerInterceptor()),
)
```

Server-side interceptor always sends the cached response unless client has specifically set the `Cache-Control: no-cache` header.