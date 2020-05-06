package ethgrp

import (
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
)

// NewRPCClient try to connect Ethereum node RPC Server to create client instance
func NewRPCClient(host string) (*rpc.Client, error) {
	//url := fmt.Sprintf("http://%s:%d", conf.Ethereum.Host, conf.Ethereum.Port)

	rpcClient, err := rpc.Dial(host)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call rpc.Dial()")
	}
	return rpcClient, nil
}
