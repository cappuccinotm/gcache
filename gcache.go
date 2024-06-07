package gcache

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ErrKey specifies the error key for logging.
var ErrKey = "error"

// ETag returns the ETag from the context.
func ETag(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if h := md.Get("If-None-Match"); len(h) > 0 {
			return h[0]
		}
	}
	return ""
}

// NotChanged responds to the client with codes.Aborted and the etag,
// meaning that the client should use the cached value.
func NotChanged(ctx context.Context, etag string) error {
	if err := grpc.SetHeader(ctx, metadata.Pairs("ETag", etag)); err != nil {
		slog.WarnContext(ctx, "gcache: failed to set ETag header", slog.Any(ErrKey, err))
		return fmt.Errorf("gcache: failed to set ETag header: %w", err)
	}

	return status.Error(codes.Aborted, "gcache: not changed")
}
