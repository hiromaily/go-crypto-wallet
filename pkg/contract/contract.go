package contract

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

// https://geth.ethereum.org/docs/dapp/native-bindings

func NewContractToken(contractAddr string, cliConn bind.ContractBackend) (*Token, error) {
	token, err := NewToken(common.HexToAddress(contractAddr), cliConn)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call NewToken()")
	}
	// token.Name(nil)
	// token.BalanceOf(nil, xxx)
	return token, nil
}
