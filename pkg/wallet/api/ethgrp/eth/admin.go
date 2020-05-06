package eth

import (
	"github.com/ethereum/go-ethereum/p2p"
)

// NodeInfo gathers and returns a collection of metadata known about the host.
func (e *Ethereum) NodeInfo() (*p2p.PeerInfo, error) {
	var r *p2p.PeerInfo
	err := e.client.CallContext(e.ctx, &r, "admin_nodeInfo")
	if err != nil {
		return nil, err
	}
	return r, err
}
