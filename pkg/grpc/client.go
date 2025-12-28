package grpc

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewClient creates a new gRPC client connection
// url: gRPC server URL (e.g., "dns:///localhost:50051", "unix:///tmp/grpc.sock")
// If the URL starts with "https://" or "grpcs://", secure credentials will be used.
// Otherwise, insecure credentials will be used.
func NewClient(url string) (*grpc.ClientConn, error) {
	if url == "" {
		return nil, errors.New("url for grpc server is not defined")
	}

	var opts []grpc.DialOption
	if !isSecureURL(url) {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.NewClient(url, opts...)
	if err != nil {
		return nil, fmt.Errorf("fail to call grpc.NewClient: %s: %w", url, err)
	}
	return conn, nil
}

// isSecureURL determines if the URL should use secure credentials
// Returns true for URLs starting with "https://" or "grpcs://"
func isSecureURL(url string) bool {
	urlLower := strings.ToLower(url)
	return strings.HasPrefix(urlLower, "https://") || strings.HasPrefix(urlLower, "grpcs://")
}

