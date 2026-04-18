// Package middleware provides gRPC interceptors for capturing traffic.
package middleware

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/grpcmon/internal/capture"
)

// UnaryInterceptor returns a gRPC unary interceptor that records each call
// into the provided Store.
func UnaryInterceptor(store *capture.Store) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		latency := time.Since(start)

		st, _ := status.FromError(err)
		entry := capture.Entry{
			ID:        capture.NewID(),
			Method:    method,
			Request:   req,
			Response:  reply,
			Status:    st.Code(),
			Latency:   latency,
			Timestamp: start,
		}
		store.Add(entry)
		return err
	}
}
