package ethgrp

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/contract"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp/erc20"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp/eth"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// NewRPCClient try to connect Ethereum node RPC Server to create client instance
func NewRPCClient(conf *config.Ethereum) (*ethrpc.Client, error) {
	url := fmt.Sprintf("http://%s:%d", conf.Host, conf.Port)
	if conf.IPCPath != "" {
		log.Println("IPC connection")
		url = conf.IPCPath
	}

	rpcClient, err := ethrpc.Dial(url)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call rpc.Dial()")
	}
	return rpcClient, nil
}

// NewEthereum creates ethereum instance according to coinType
func NewEthereum(rpcClient *ethrpc.Client, conf *config.Ethereum, logger *zap.Logger, coinTypeCode coin.CoinTypeCode) (Ethereumer, error) {
	client := ethclient.NewClient(rpcClient)

	var erc20Obj *erc20.ERC20
	if coinTypeCode == coin.ERC20 {
		tokenClient, err := contract.NewContractToken(conf.ERC20s[conf.ERC20Token].ContractAddress, client)
		if err != nil {
			return nil, errors.Wrap(err, "fail to call contract.NewContractToken()")
		}
		erc20Obj = erc20.NewERC20(
			tokenClient,
			conf.ERC20Token,
			conf.ERC20s[conf.ERC20Token].Name,
			conf.ERC20s[conf.ERC20Token].ContractAddress,
			conf.ERC20s[conf.ERC20Token].MasterAddress,
			logger,
		)
	}

	eth, err := eth.NewEthereum(
		context.Background(),
		client,
		rpcClient,
		erc20Obj,
		coinTypeCode,
		conf,
		logger,
	)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call eth.NewEthereum()")
	}
	return eth, err
}
