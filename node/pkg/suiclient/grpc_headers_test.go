package suiclient

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestParseGrpcHeaderSpecs(t *testing.T) {
	tests := []struct {
		name    string
		specs   []string
		want    metadata.MD
		wantErr bool
	}{
		{
			name:  "empty",
			specs: nil,
			want:  metadata.MD{},
		},
		{
			name:  "valid headers",
			specs: []string{"X-API-Key=secret", "chain-route=sui-mainnet"},
			want: metadata.MD{
				"x-api-key":   []string{"secret"},
				"chain-route": []string{"sui-mainnet"},
			},
		},
		{
			name:    "duplicate key",
			specs:   []string{"X-API-Key=secret", "x-api-key=second"},
			wantErr: true,
		},
		{
			name:  "trims whitespace",
			specs: []string{" x-api-key = secret "},
			want: metadata.MD{
				"x-api-key": []string{"secret"},
			},
		},
		{
			name:    "missing separator",
			specs:   []string{"x-api-key"},
			wantErr: true,
		},
		{
			name:    "empty key",
			specs:   []string{"=secret"},
			wantErr: true,
		},
		{
			name:    "empty value",
			specs:   []string{"x-api-key="},
			wantErr: true,
		},
		{
			name:    "invalid key",
			specs:   []string{"x api key=secret"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseGrpcHeaderSpecs(tt.specs)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGrpcHeaderInterceptorsAppendMetadata(t *testing.T) {
	headers := metadata.MD{
		"x-api-key":   []string{"secret"},
		"chain-route": []string{"sui-mainnet"},
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "existing-header", "existing-value")

	var unaryMetadata metadata.MD
	unaryInvoker := func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		var ok bool
		unaryMetadata, ok = metadata.FromOutgoingContext(ctx)
		require.True(t, ok)
		return nil
	}

	err := grpcHeaderUnaryInterceptor(headers)(ctx, "/sui.rpc.v2.LedgerService/GetCheckpoint", nil, nil, nil, unaryInvoker)
	require.NoError(t, err)
	require.Equal(t, []string{"existing-value"}, unaryMetadata.Get("existing-header"))
	require.Equal(t, []string{"secret"}, unaryMetadata.Get("x-api-key"))
	require.Equal(t, []string{"sui-mainnet"}, unaryMetadata.Get("chain-route"))

	var streamMetadata metadata.MD
	streamer := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		var ok bool
		streamMetadata, ok = metadata.FromOutgoingContext(ctx)
		require.True(t, ok)
		return nil, nil
	}

	_, err = grpcHeaderStreamInterceptor(headers)(ctx, &grpc.StreamDesc{}, nil, "/sui.rpc.v2.SubscriptionService/SubscribeCheckpoints", streamer)
	require.NoError(t, err)
	require.Equal(t, []string{"existing-value"}, streamMetadata.Get("existing-header"))
	require.Equal(t, []string{"secret"}, streamMetadata.Get("x-api-key"))
	require.Equal(t, []string{"sui-mainnet"}, streamMetadata.Get("chain-route"))
}

func TestGrpcHeaderDialOptions(t *testing.T) {
	opts, err := GrpcHeaderDialOptions([]string{"x-api-key=secret"})
	require.NoError(t, err)
	require.Len(t, opts, 2)

	opts, err = GrpcHeaderDialOptions(nil)
	require.NoError(t, err)
	require.Empty(t, opts)
}
