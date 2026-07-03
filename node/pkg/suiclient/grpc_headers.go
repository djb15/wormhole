package suiclient

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// GrpcHeaderDialOptions converts key=value header specs from command line args
// into gRPC client interceptors that attach the headers as outgoing metadata
func GrpcHeaderDialOptions(headerSpecs []string) ([]grpc.DialOption, error) {
	headers, err := parseGrpcHeaderSpecs(headerSpecs)
	if err != nil {
		return nil, err
	}
	if len(headers) == 0 {
		return nil, nil
	}

	headers = headers.Copy()
	return []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(grpcHeaderUnaryInterceptor(headers)),
		grpc.WithChainStreamInterceptor(grpcHeaderStreamInterceptor(headers)),
	}, nil
}

func parseGrpcHeaderSpecs(headerSpecs []string) (metadata.MD, error) {
	headers := metadata.MD{}

	for _, spec := range headerSpecs {
		spec = strings.TrimSpace(spec)
		if spec == "" {
			continue
		}

		key, value, ok := strings.Cut(spec, "=")
		if !ok {
			return nil, fmt.Errorf("invalid gRPC header: expected key=value")
		}

		key = strings.ToLower(strings.TrimSpace(key))
		value = strings.TrimSpace(value)

		if err := validateGrpcMetadataKey(key); err != nil {
			return nil, err
		}
		if value == "" {
			return nil, fmt.Errorf("invalid gRPC header %q: value must not be empty", key)
		}
		if _, exists := headers[key]; exists {
			return nil, fmt.Errorf("duplicate gRPC header key %q", key)
		}

		headers.Append(key, value)
	}

	return headers, nil
}

func validateGrpcMetadataKey(key string) error {
	if key == "" {
		return fmt.Errorf("invalid gRPC header: key must not be empty")
	}

	for _, r := range key {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= '0' && r <= '9':
		case r == '-' || r == '_' || r == '.':
		default:
			return fmt.Errorf("invalid gRPC header key %q: keys may only contain letters, digits, '-', '_' or '.'", key)
		}
	}

	return nil
}

func grpcHeaderUnaryInterceptor(headers metadata.MD) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(appendGrpcHeadersToContext(ctx, headers), method, req, reply, cc, opts...)
	}
}

func grpcHeaderStreamInterceptor(headers metadata.MD) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return streamer(appendGrpcHeadersToContext(ctx, headers), desc, cc, method, opts...)
	}
}

func appendGrpcHeadersToContext(ctx context.Context, headers metadata.MD) context.Context {
	if len(headers) == 0 {
		return ctx
	}

	existing, _ := metadata.FromOutgoingContext(ctx)
	return metadata.NewOutgoingContext(ctx, metadata.Join(existing, headers))
}
