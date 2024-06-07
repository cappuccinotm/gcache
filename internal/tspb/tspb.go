// Package tspb contains a mock of gRPC service for testing purposes.
package tspb

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

// MockTestService is a mock for TestServiceServer
type MockTestService struct {
	TestFunc func(ctx context.Context, in *TestRequest) (*TestResponse, error)

	UnimplementedTestServiceServer
}

// Test implements TestServiceServer
func (m *MockTestService) Test(ctx context.Context, in *TestRequest) (*TestResponse, error) {
	if m.TestFunc == nil {
		return m.UnimplementedTestServiceServer.Test(ctx, in)
	}
	return m.TestFunc(ctx, in)
}

// Run runs test server.
func Run(t *testing.T, ts MockTestService, opts ...grpc.ServerOption) (addr string) {
	l, err := net.Listen("tcp", ":0") //nolint:gosec // OK for tests
	require.NoError(t, err)

	s := grpc.NewServer(opts...)
	RegisterTestServiceServer(s, &ts)

	go func() {
		require.NoError(t, s.Serve(l))
	}()

	t.Cleanup(func() { s.GracefulStop() })

	return l.Addr().String()
}
