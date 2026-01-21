package xrp

import (
	"google.golang.org/grpc"

	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/xrp/protogen"
)

// XRPAPI is the XRP API client
type XRPAPI struct {
	txClient      protogen.XRPTransactionAPIClient
	accountClient protogen.XRPAccountAPIClient
	addressClient protogen.XRPAddressAPIClient
	conn          *grpc.ClientConn
}

// NewXRPAPI creates XRP API object
func NewXRPAPI(
	conn *grpc.ClientConn,
) *XRPAPI {
	return &XRPAPI{
		txClient:      protogen.NewXRPTransactionAPIClient(conn),
		accountClient: protogen.NewXRPAccountAPIClient(conn),
		addressClient: protogen.NewXRPAddressAPIClient(conn),
		conn:          conn,
	}
}

// Close disconnect to server
func (r *XRPAPI) Close() {
	if r.conn != nil {
		_ = r.conn.Close() // Best effort cleanup
	}
}
