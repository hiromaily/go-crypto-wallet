package btc

import (
	"errors"
	"fmt"

	"github.com/bookerzzz/grok"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin"
)

func runListUnspent(btc bitcoin.Bitcoiner, acnt string, argsNum int64) error {
	// validator
	if acnt != "" && !domainAccount.ValidateAccountType(acnt) {
		return errors.New("account option [-account] is invalid")
	}

	var confirmationNum uint64
	if argsNum == -1 {
		confirmationNum = btc.ConfirmationBlock()
	} else {
		confirmationNum = uint64(argsNum)
	}

	if acnt != "" {
		// call listunspent
		unspentList, err := btc.ListUnspentByAccount(domainAccount.AccountType(acnt), confirmationNum)
		if err != nil {
			return fmt.Errorf("fail to call btc.ListUnspentByAccount() %w", err)
		}
		grok.Value(unspentList)

		unspentAddrs := btc.GetUnspentListAddrs(unspentList, domainAccount.AccountType(acnt))
		for _, addr := range unspentAddrs {
			grok.Value(addr)
		}
	} else {
		// call listunspent
		// ListUnspentMin doesn't have proper response, label can't be retrieved

		unspentList, err := btc.ListUnspent(confirmationNum)
		if err != nil {
			return fmt.Errorf("fail to call btc.ListUnspent() %w", err)
		}
		grok.Value(unspentList)
	}

	return nil
}
