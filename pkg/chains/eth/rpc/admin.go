package rpc

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/p2p"
)

// AddPeer requests adding a new remote node to the list of tracked static nodes.
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#admin_addpeer
func (c *rpcClient) AddPeer(ctx context.Context, nodeURL string) error {
	var bRet bool
	err := c.caller.CallContext(ctx, &bRet, "admin_addPeer", nodeURL)
	if err != nil {
		return err
	}
	return err
}

// AdminDataDir returns the absolute path the running Geth node uses to store its databases.
func (c *rpcClient) AdminDataDir(ctx context.Context) (string, error) {
	var dataDir string
	err := c.caller.CallContext(ctx, &dataDir, "admin_datadir")
	if err != nil {
		return "", fmt.Errorf("fail to call rpc.CallContext(admin_datadir): %w", err)
	}
	return dataDir, nil
}

// NodeInfo gathers and returns metadata about the host.
func (c *rpcClient) NodeInfo(ctx context.Context) (*p2p.NodeInfo, error) {
	var r *p2p.NodeInfo
	err := c.caller.CallContext(ctx, &r, "admin_nodeInfo")
	if err != nil {
		return nil, err
	}
	return r, err
}

// AdminPeers returns information about all connected remote nodes.
func (c *rpcClient) AdminPeers(ctx context.Context) ([]*p2p.PeerInfo, error) {
	var peerInfo []*p2p.PeerInfo
	err := c.caller.CallContext(ctx, &peerInfo, "admin_peers")
	if err != nil {
		return nil, err
	}
	return peerInfo, err
}
