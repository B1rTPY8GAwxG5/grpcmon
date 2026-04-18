package middleware_test

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/grpcmon/internal/capture"

)

func invoke(errpc.UnaryInvoker {
	ret context.Context, _ _ interface{}, *grpc.ClientConn, _ ...grpc.CallOption) error {
		return err
	}
}

func TestUnaryInterceptor_RecordsEntry(t *testing.T) {
	store := capture.NewStore(10)
	intercept := middleware.UnaryInterceptor(store)

	err := intercept(context.Background(), "/pkg.Svc/Method", "req", "rep", nil, invoke(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries := store.List()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Method != "/pkg.Svc/Method" {
		t.Errorf("unexpected method: %s", entries[0].Method)
	}
	if entries[0].Status != codes.OK {
		t.Errorf("expected OK, got %v", entries[0].Status)
	}
}

func TestUnaryInterceptor_RecordsError(t *testing.T) {
	store := capture.NewStore(10)
	intercept := middleware.UnaryInterceptor(store)

	grpcErr := status.Error(codes.NotFound, "not found")
	_ = intercept(context.Background(), "/pkg.Svc/Find", nil, nil, nil, invoke(grpcErr))

	entries := store.List()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Status != codes.NotFound {
		t.Errorf("expected NotFound, got %v", entries[0].Status)
	}
}

func TestUnaryInterceptor_SetsTimestamp(t *testing.T) {
	store := capture.NewStore(10)
	intercept := middleware.UnaryInterceptor(store)
	_ = intercept(context.Background(), "/s/M", nil, nil, nil, invoke(nil))

	entry := store.List()[0]
	if entry.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
	if entry.Latency <= 0 {
		t.Error("latency should be positive")
	}
}
