package eth

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// BalanceAt returns balance of address
// if wrong address is given, response of balance would be strange like `416778046407207737`
func (e *Ethereum) BalanceAt(ctx context.Context, hexAddr string) (*big.Int, error) {
	account := common.HexToAddress(hexAddr)
	balance, err := e.ethClient.BalanceAt(ctx, account, nil)
	if err != nil {
		return nil, fmt.Errorf("fail to call ethClient.BalanceAt(): %w", err)
	}
	if balance.Uint64() == 416778046407207737 {
		return nil, errors.New("response of balance is strange. 416778046407207737 is returned")
	}
	return balance, nil
}

// SendRawTx sends raw transaction
func (e *Ethereum) SendRawTx(ctx context.Context, tx *types.Transaction) error {
	return e.ethClient.SendTransaction(ctx, tx)
}
