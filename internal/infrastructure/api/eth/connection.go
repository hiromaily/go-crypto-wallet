package ethereum

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	apiethimpl "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/eth"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/uuid"
)

// NewEthereum creates ethereum instance according to coinType
func NewEthereum(
	rpcClient *ethrpc.Client, conf *config.Ethereum,
	coinTypeCode domainCoin.CoinTypeCode, uuidHandler uuid.UUIDHandler,
) (apieth.Ethereumer, error) {
	client := ethclient.NewClient(rpcClient)

	ethAPI, err := apiethimpl.NewEthereum(
		context.Background(),
		client,
		rpcClient,
		coinTypeCode,
		conf,
		uuidHandler,
	)
	if err != nil {
		return nil, fmt.Errorf("fail to call eth.NewEthereum(): %w", err)
	}
	return ethAPI, err
}
