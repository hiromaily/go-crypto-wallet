package erc20

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/contract"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// ERC20er ABI Token Interface
type ERC20er interface {
	GetBalance(hexAddr string) (*big.Int, error)
}

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

//func (e *ERC20) getOption(
//	ctx context.Context,
//	isPending bool,
//	fromAddr common.Address,
//	blockNumber *big.Int) *bind.CallOpts {
//
//	opts := bind.CallOpts{}
//	if ctx != nil {
//		opts.Context = ctx
//	}
//	opts.Pending = isPending
//	opts.From = fromAddr
//	if blockNumber != nil {
//		opts.BlockNumber = blockNumber
//	}
//	return &opts
//}

func (e *ERC20) GetBalance(hexAddr string) (*big.Int, error) {
	balance, err := e.tokenClient.BalanceOf(nil, common.HexToAddress(hexAddr))
	if err != nil {
		return nil, errors.Wrapf(err, "fail to call e.contract.BalanceOf(%s)", hexAddr)
	}
	return balance, nil
}
