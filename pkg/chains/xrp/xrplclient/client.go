package xrplclient

import (
	"google.golang.org/grpc"

	"github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/protogen"
)

// Deprecated: XRPLClient is the apps/xrpl-grpc-server client
// TODO: divide XRPLClient into 3 structs
type XRPLClient struct {
	TxClient      protogen.XRPTransactionAPIClient
	AccountClient protogen.XRPAccountAPIClient
	AddressClient protogen.XRPAddressAPIClient
	conn          *grpc.ClientConn
}

// XRPLClient creates XRPLServer Client object
func NewXRPLClient(
	conn *grpc.ClientConn,
) *XRPLClient {
	return &XRPLClient{
		TxClient:      protogen.NewXRPTransactionAPIClient(conn),
		AccountClient: protogen.NewXRPAccountAPIClient(conn),
		AddressClient: protogen.NewXRPAddressAPIClient(conn),
		conn:          conn,
	}
}

// Close disconnect to server
func (r *XRPLClient) Close() {
	if r.conn != nil {
		_ = r.conn.Close() // Best effort cleanup
	}
}
