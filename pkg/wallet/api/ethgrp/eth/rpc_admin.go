package eth

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/p2p"
)

// AddPeer requests adding a new remote node to the list of tracked static nodes
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#admin_addpeer
func (e *Ethereum) AddPeer(ctx context.Context, nodeURL string) error {
	var bRet bool
	// TODO: Result needs to be verified
	// The response data type are bytes, but it cannot parse...
	err := e.rpcClient.CallContext(ctx, &bRet, "admin_addPeer", nodeURL)
	if err != nil {
		return err
	}
	return err
}

// AdminDataDir returns the absolute path the running Geth node currently uses to store all its databases
// returns like ${HOME}/Library/Ethereum/goerli
func (e *Ethereum) AdminDataDir(ctx context.Context) (string, error) {
	if e.isParity {
		return "", nil
	}

	var dataDir string
	err := e.rpcClient.CallContext(ctx, &dataDir, "admin_datadir")
	if err != nil {
		return "", fmt.Errorf("fail to call rpc.CallContext(admin_datadir): %w", err)
	}
	return dataDir, nil
}

// NodeInfo gathers and returns a collection of metadata known about the host.
func (e *Ethereum) NodeInfo(ctx context.Context) (*p2p.NodeInfo, error) {
	var r *p2p.NodeInfo
	err := e.rpcClient.CallContext(ctx, &r, "admin_nodeInfo")
	if err != nil {
		return nil, err
	}
	return r, err
}

// AdminPeers returns all the information known about the connected remote nodes at the networking granularity.
func (e *Ethereum) AdminPeers(ctx context.Context) ([]*p2p.PeerInfo, error) {
	if e.isParity {
		return nil, nil
	}

	var peerInfo []*p2p.PeerInfo
	err := e.rpcClient.CallContext(ctx, &peerInfo, "admin_peers")
	if err != nil {
		return nil, err
	}
	return peerInfo, err
}
