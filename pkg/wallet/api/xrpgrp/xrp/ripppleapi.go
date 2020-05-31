package xrp

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/hiromaily/ripple-lib-proto/pb/go/rippleapi"
)

// RippleAPI it RippleAPI client
type RippleAPI struct {
	client pb.RippleAPIClient
	conn   *grpc.ClientConn
	logger *zap.Logger
}

// NewRippleAPI creates Ripple API object
func NewRippleAPI(
	conn *grpc.ClientConn,
	logger *zap.Logger) *RippleAPI {

	return &RippleAPI{
		client: pb.NewRippleAPIClient(conn),
		conn:   conn,
		logger: logger,
	}
}

// Close disconnect to server
func (r *RippleAPI) Close() {
	if r.conn != nil {
		r.conn.Close()
	}
}

//func (r *RippleAPI) APIClient() pb.RippleAPIClient {
//	return r.client
//}
