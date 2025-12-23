package ethgrp

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/uuid"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp/eth"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// NewRPCClient try to connect Ethereum node RPC Server to create client instance
func NewRPCClient(conf *config.Ethereum) (*ethrpc.Client, error) {
	url := "http://" + net.JoinHostPort(conf.Host, strconv.Itoa(conf.Port))
	if conf.IPCPath != "" {
		log.Println("IPC connection")
		url = conf.IPCPath
	}

	rpcClient, err := ethrpc.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("fail to call rpc.Dial(): %w", err)
	}
	return rpcClient, nil
}

// NewEthereum creates ethereum instance according to coinType
func NewEthereum(
	rpcClient *ethrpc.Client, conf *config.Ethereum, logger logger.Logger,
	coinTypeCode coin.CoinTypeCode, uuidHandler uuid.UUIDHandler,
) (Ethereumer, error) {
	client := ethclient.NewClient(rpcClient)

	ethAPI, err := eth.NewEthereum(
		context.Background(),
		client,
		rpcClient,
		coinTypeCode,
		conf,
		logger,
		uuidHandler,
	)
	if err != nil {
		return nil, fmt.Errorf("fail to call eth.NewEthereum(): %w", err)
	}
	return ethAPI, err
}
