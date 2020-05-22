package address

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// AddressFormat is address csv format
type AddressFormat struct {
	CoinTypeCode      coin.CoinTypeCode
	AccountType       account.AccountType
	P2PKHAddress      string
	P2SHSegwitAddress string
	Bech32Address     string
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
		accountKeyItem.Bech32Address,
		accountKeyItem.FullPublicKey,
		accountKeyItem.MultisigAddress,
		strconv.Itoa(int(accountKeyItem.Idx)),
	}
}

// ConvertLine converts line to AddressFormat
func ConvertLine(coinTypeCode coin.CoinTypeCode, line []string) (*AddressFormat, error) {
	if len(line) != 8 {
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
		Bech32Address:     line[4],
		FullPublicKey:     line[5],
		MultisigAddress:   line[6],
		Idx:               line[7],
	}, nil
}
