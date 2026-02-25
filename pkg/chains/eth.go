package cryptocurrency

import (
	"fmt"
	"log"
	"net"
	"strconv"

	ethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
)

// NewEthereumRPCClient try to connect Ethereum node RPC Server to create client instance
func NewEthereumRPCClient(conf *config.Ethereum) (*ethrpc.Client, error) {
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
