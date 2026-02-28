package eth

import (
	"context"

	"github.com/ethereum/go-ethereum/p2p"

	ethrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth/rpc"
)

// AddPeer requests adding a new remote node to the list of tracked static nodes
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#admin_addpeer
func (e *Ethereum) AddPeer(ctx context.Context, nodeURL string) error {
	return ethrpc.AddPeer(ctx, e.rpcClient, nodeURL)
}

// AdminDataDir returns the absolute path the running Geth node currently uses to store all its databases
// returns like ${HOME}/Library/Ethereum/goerli
func (e *Ethereum) AdminDataDir(ctx context.Context) (string, error) {
	if e.isParity {
		return "", nil
	}
	return ethrpc.AdminDataDir(ctx, e.rpcClient)
}

// NodeInfo gathers and returns a collection of metadata known about the host.
func (e *Ethereum) NodeInfo(ctx context.Context) (*p2p.NodeInfo, error) {
	return ethrpc.NodeInfo(ctx, e.rpcClient)
}

// AdminPeers returns all the information known about the connected remote nodes at the networking granularity.
func (e *Ethereum) AdminPeers(ctx context.Context) ([]*p2p.PeerInfo, error) {
	if e.isParity {
		return nil, nil
	}
	return ethrpc.AdminPeers(ctx, e.rpcClient)
}
