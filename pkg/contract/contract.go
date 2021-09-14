package contract

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

// https://geth.ethereum.org/docs/dapp/native-bindings

// TODO: only required function should be defined as interface
// TokenContructor ABI Token Interface
//type TokenContractor interface {
//}

func NewContractToken(masterAddr string, cliConn bind.ContractBackend) (*Token, error) {
	token, err := NewToken(common.HexToAddress(masterAddr), cliConn)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call NewToken()")
	}
	// token.Name(nil)
	// token.BalanceOf(nil, xxx)
	return token, nil
}
