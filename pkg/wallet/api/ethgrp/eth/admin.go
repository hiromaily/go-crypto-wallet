package eth

import (
	"github.com/ethereum/go-ethereum/p2p"
)

// NodeInfo gathers and returns a collection of metadata known about the host.
//  - FIXME: the method admin_nodeInfo does not exist/is not available
func (e *Ethereum) NodeInfo() (*p2p.NodeInfo, error) {
	var r *p2p.NodeInfo
	err := e.client.CallContext(e.ctx, &r, "admin_nodeInfo")
	if err != nil {
		return nil, err
	}
	return r, err
}
