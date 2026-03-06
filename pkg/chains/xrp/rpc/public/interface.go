package public

import "context"

// XRPPublicer defines the interface for XRP public node operations.
// These operations query public information from the XRP network.
type XRPPublicer interface {
	// public_account
	AccountChannels(ctx context.Context, sender, receiver string) (*ResponseAccountChannels, error)
	AccountInfo(ctx context.Context, address string) (*ResponseAccountInfo, error)
	// public_server_info
	ServerInfo(ctx context.Context) (*ResponseServerInfo, error)
}
