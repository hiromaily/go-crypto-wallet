package xrp

import (
	"google.golang.org/grpc"
)

// RippleAPI it RippleAPI client
type RippleAPI struct {
	txClient      RippleTransactionAPIClient
	accountClient RippleAccountAPIClient
	addressClient RippleAddressAPIClient
	conn          *grpc.ClientConn
}

// NewRippleAPI creates Ripple API object
func NewRippleAPI(
	conn *grpc.ClientConn,
) *RippleAPI {
	return &RippleAPI{
		txClient:      NewRippleTransactionAPIClient(conn),
		accountClient: NewRippleAccountAPIClient(conn),
		addressClient: NewRippleAddressAPIClient(conn),
		conn:          conn,
	}
}

// Close disconnect to server
func (r *RippleAPI) Close() {
	if r.conn != nil {
		r.conn.Close()
	}
}

// func (r *RippleAPI) APIClient() pb.RippleAPIClient {
//	return r.client
//}
