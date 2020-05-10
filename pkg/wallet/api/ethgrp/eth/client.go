package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

// BalanceAt returns balance of address
func (e *Ethereum) BalanceAt(hexAddr string) (*big.Int, error) {
	account := common.HexToAddress(hexAddr)
	balance, err := e.ethClient.BalanceAt(e.ctx, account, nil)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call ethClient.BalanceAt()")
	}
	return balance, nil
}

// SendRawTx sends raw transaction
func (e *Ethereum) SendRawTx(tx *types.Transaction) error {
	return e.ethClient.SendTransaction(e.ctx, tx)
}
