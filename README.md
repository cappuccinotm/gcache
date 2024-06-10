# gcache [![Go Reference](https://pkg.go.dev/badge/github.com/cappuccinotm/gcache.svg)](https://pkg.go.dev/github.com/cappuccinotm/gcache) [![Go](https://github.com/cappuccinotm/gcache/actions/workflows/go.yaml/badge.svg)](https://github.com/cappuccinotm/gcache/actions/workflows/go.yaml) [![codecov](https://codecov.io/gh/cappuccinotm/gcache/graph/badge.svg?token=We5B4lzTNj)](https://codecov.io/gh/cappuccinotm/gcache)
gcache uses a plain old `ETag`, `If-None-Match` (as defined in [RFC7232](https://datatracker.ietf.org/doc/html/rfc7232)), [`Cache-Control`](https://datatracker.ietf.org/doc/html/rfc7234) headers to control the caching of responses. 

The standard is mostly used in caching the static files, but it can be used in caching the dynamic content as well. 

Briefly, the key idea is that the gRPC service responds to the client with the `ETag` header in metadata, client caches the response with this `ETag`, sends the request to the service with that `ETag` of the cached response, and, if the service has detected that the resource hasn't changed (judging by the provided `ETag`), responds the client with `codes.Aborted` and the provided `ETag` (similar with `304 Not Modified` status code).

The package also provides a server-side cache, it looks only for a client's `Cache-Control` header, if it is set to `no-cache`, the server doesn't send the cached response and evaluates the response as usual.

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