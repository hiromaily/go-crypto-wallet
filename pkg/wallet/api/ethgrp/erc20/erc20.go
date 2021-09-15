package erc20

import (
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/contract"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

type ERC20 struct {
	tokenClient     *contract.Token
	token           coin.ERC20Token
	name            string
	contractAddress string
	masterAddress   string
	logger          *zap.Logger
}

func NewERC20(
	tokenClient *contract.Token,
	token coin.ERC20Token,
	name string,
	contractAddress string,
	masterAddress string,
	logger *zap.Logger,
) *ERC20 {
	return &ERC20{
		tokenClient:     tokenClient,
		token:           token,
		name:            name,
		contractAddress: contractAddress,
		masterAddress:   masterAddress,
		logger:          logger,
	}
}
