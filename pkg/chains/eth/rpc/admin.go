package rpc

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/p2p"
)

// AddPeer requests adding a new remote node to the list of tracked static nodes.
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#admin_addpeer
func AddPeer(ctx context.Context, caller RPCCaller, nodeURL string) error {
	var bRet bool
	err := caller.CallContext(ctx, &bRet, "admin_addPeer", nodeURL)
	if err != nil {
		return err
	}
	return err
}

// AdminDataDir returns the absolute path the running Geth node uses to store its databases.
func AdminDataDir(ctx context.Context, caller RPCCaller) (string, error) {
	var dataDir string
	err := caller.CallContext(ctx, &dataDir, "admin_datadir")
	if err != nil {
		return "", fmt.Errorf("fail to call rpc.CallContext(admin_datadir): %w", err)
	}
	return dataDir, nil
}

// NodeInfo gathers and returns metadata about the host.
func NodeInfo(ctx context.Context, caller RPCCaller) (*p2p.NodeInfo, error) {
	var r *p2p.NodeInfo
	err := caller.CallContext(ctx, &r, "admin_nodeInfo")
	if err != nil {
		return nil, err
	}
	return r, err
}

// AdminPeers returns information about all connected remote nodes.
func AdminPeers(ctx context.Context, caller RPCCaller) ([]*p2p.PeerInfo, error) {
	var peerInfo []*p2p.PeerInfo
	err := caller.CallContext(ctx, &peerInfo, "admin_peers")
	if err != nil {
		return nil, err
	}
	return peerInfo, err
}
