package middleware_test

import (
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/grpcmon/internal/capture"
	"github.com/grpcmon/internal/middleware"
)

func startListener(t *testing.T) string {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	srv := grpc.NewServer()
	go srv.Serve(lis) //nolint:errcheck
	t.Cleanup(srv.Stop)
	return lis.Addr().String()
}

func TestDialWithCapture_ReturnsConn(t *testing.T) {
	addr := startListener(t)
	store := capture.NewStore(10)

	conn, err := middleware.DialWithCapture(
		addr,
		store,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer conn.Close()

	if conn == nil {
		t.Error("expected non-nil connection")
	}
}
