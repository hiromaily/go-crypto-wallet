package contract

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// https://geth.ethereum.org/docs/dapp/native-bindings

func NewContractToken(contractAddr string, cliConn bind.ContractBackend) (*Token, error) {
	token, err := NewToken(common.HexToAddress(contractAddr), cliConn)
	if err != nil {
		return nil, fmt.Errorf("fail to call NewToken(): %w", err)
	}
	// token.Name(nil)
	// token.BalanceOf(nil, xxx)
	return token, nil
}
