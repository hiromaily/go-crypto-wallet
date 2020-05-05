package address

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	models "github.com/hiromaily/go-bitcoin/pkg/models/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/coin"
)

// AddressFormat is address csv format
type AddressFormat struct {
	CoinTypeCode      coin.CoinTypeCode
	AccountType       account.AccountType
	P2PKHAddress      string
	P2SHSegwitAddress string
	FullPublicKey     string
	MultisigAddress   string
	Idx               string
}

// CreateLine creates line for csv
func CreateLine(accountKeyItem *models.AccountKey) []string {
	return []string{
		accountKeyItem.Coin,
		accountKeyItem.Account,
		accountKeyItem.P2PKHAddress,
		accountKeyItem.P2SHSegwitAddress,
		accountKeyItem.FullPublicKey,
		accountKeyItem.MultisigAddress,
		strconv.Itoa(int(accountKeyItem.Idx)),
	}
}

// ConvertLine converts line to AddressFormat
func ConvertLine(coinTypeCode coin.CoinTypeCode, line []string) (*AddressFormat, error) {
	if len(line) != 7 {
		return nil, errors.New("csv format is invalid")
	}

	// validate
	if !coin.ValidateCoinTypeCode(line[0]) || coin.CoinTypeCode(line[0]) != coinTypeCode {
		return nil, errors.Errorf("coinTypeCode is invalid. got %s, want %s", line[0], coinTypeCode.String())
	}
	if !account.ValidateAccountType(line[1]) {
		return nil, errors.Errorf("account is invalid: %s", line[1])
	}

	return &AddressFormat{
		CoinTypeCode:      coin.CoinTypeCode(line[0]),
		AccountType:       account.AccountType(line[1]),
		P2PKHAddress:      line[2],
		P2SHSegwitAddress: line[3],
		FullPublicKey:     line[4],
		MultisigAddress:   line[5],
		Idx:               line[6],
	}, nil
}
