package xrp

import (
	"google.golang.org/grpc"

	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/xrp/protogen"
)

// RippleAPI it RippleAPI client
type RippleAPI struct {
	txClient      protogen.RippleTransactionAPIClient
	accountClient protogen.RippleAccountAPIClient
	addressClient protogen.RippleAddressAPIClient
	conn          *grpc.ClientConn
}

// NewRippleAPI creates Ripple API object
func NewRippleAPI(
	conn *grpc.ClientConn,
) *RippleAPI {
	return &RippleAPI{
		txClient:      protogen.NewRippleTransactionAPIClient(conn),
		accountClient: protogen.NewRippleAccountAPIClient(conn),
		addressClient: protogen.NewRippleAddressAPIClient(conn),
		conn:          conn,
	}
}

// Close disconnect to server
func (r *RippleAPI) Close() {
	if r.conn != nil {
		_ = r.conn.Close() // Best effort cleanup
	}
}

// func (r *RippleAPI) APIClient() pb.RippleAPIClient {
//	return r.client
//}
