package ethgrp

import (
	"fmt"

	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/config"
)

// NewRPCClient try to connect Ethereum node RPC Server to create client instance
func NewRPCClient(conf *config.Ethereum) (*ethrpc.Client, error) {
	url := fmt.Sprintf("http://%s:%d", conf.Host, conf.Port)

	rpcClient, err := ethrpc.Dial(url)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call rpc.Dial()")
	}
	return rpcClient, nil
}
