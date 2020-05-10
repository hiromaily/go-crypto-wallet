package eth

import (
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/pkg/errors"
)

// AddPeer requests adding a new remote node to the list of tracked static nodes
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#admin_addpeer
func (e *Ethereum) AddPeer(nodeURL string) error {
	var bRet bool
	// TODO: Result needs to be verified
	// The response data type are bytes, but it cannot parse...
	err := e.rpcClient.CallContext(e.ctx, &bRet, "admin_addPeer", nodeURL)
	if err != nil {
		return err
	}
	return err
}

// AdminDataDir returns the absolute path the running Geth node currently uses to store all its databases
func (e *Ethereum) AdminDataDir() (string, error) {
	if e.isParity {
		return "", nil
	}

	var dataDir string
	err := e.rpcClient.CallContext(e.ctx, &dataDir, "admin_datadir")
	if err != nil {
		return "", errors.Wrap(err, "fail to call rpc.CallContext(admin_datadir)")
	}
	return dataDir, nil
}

// NodeInfo gathers and returns a collection of metadata known about the host.
func (e *Ethereum) NodeInfo() (*p2p.NodeInfo, error) {
	var r *p2p.NodeInfo
	err := e.rpcClient.CallContext(e.ctx, &r, "admin_nodeInfo")
	if err != nil {
		return nil, err
	}
	return r, err
}

// AdminPeers returns all the information known about the connected remote nodes at the networking granularity.
func (e *Ethereum) AdminPeers() ([]*p2p.PeerInfo, error) {
	if e.isParity {
		return nil, nil
	}

	var peerInfo []*p2p.PeerInfo
	err := e.rpcClient.CallContext(e.ctx, &peerInfo, "admin_peers")
	if err != nil {
		return nil, err
	}
	return peerInfo, err
}
