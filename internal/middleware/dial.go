package middleware

import (
	"google.golang.org/grpc"

	"github.com/grpcmon/internal/capture"
)

// DialWithCapture dials a gRPC target and attaches the capture interceptor.
func DialWithCapture(target string, store *capture.Store, extra ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(UnaryInterceptor(store)),
	}
	opts = append(opts, extra...)
	//nolint:staticcheck // grpc.Dial is used for broad compatibility
	return grpc.Dial(target, opts...)
}
