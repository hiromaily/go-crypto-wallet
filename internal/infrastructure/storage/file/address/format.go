package address

import (
	"fmt"
	"strconv"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	models "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/models/rdb"
)

// AddressFormat is address csv format
type AddressFormat struct {
	CoinTypeCode      domainCoin.CoinTypeCode
	AccountType       domainAccount.AccountType
	P2PKHAddress      string
	P2SHSegwitAddress string
	Bech32Address     string
	TaprootAddress    string
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
		accountKeyItem.TaprootAddress,
		accountKeyItem.FullPublicKey,
		accountKeyItem.MultisigAddress,
		strconv.Itoa(int(accountKeyItem.Idx)),
	}
}

// ConvertLine converts line to AddressFormat
func ConvertLine(coinTypeCode domainCoin.CoinTypeCode, line []string) (*AddressFormat, error) {
	// Support both old format (8 fields) and new format (9 fields with Taproot)
	if len(line) != 8 && len(line) != 9 {
		return nil, fmt.Errorf("csv format is invalid: expected 8 or 9 fields, got %d", len(line))
	}

	// validate
	if !domainCoin.IsCoinTypeCode(line[0]) || domainCoin.CoinTypeCode(line[0]) != coinTypeCode {
		return nil, fmt.Errorf("coinTypeCode is invalid. got %s, want %s", line[0], coinTypeCode.String())
	}
	if !domainAccount.ValidateAccountType(line[1]) {
		return nil, fmt.Errorf("account is invalid: %s", line[1])
	}

	// For backward compatibility with old CSV format (without Taproot)
	taprootAddress := ""
	fullPublicKeyIdx := 5
	multisigAddressIdx := 6
	idxIdx := 7

	if len(line) == 9 {
		// New format with Taproot address
		taprootAddress = line[5]
		fullPublicKeyIdx = 6
		multisigAddressIdx = 7
		idxIdx = 8
	}

	return &AddressFormat{
		CoinTypeCode:      domainCoin.CoinTypeCode(line[0]),
		AccountType:       domainAccount.AccountType(line[1]),
		P2PKHAddress:      line[2],
		P2SHSegwitAddress: line[3],
		Bech32Address:     line[4],
		TaprootAddress:    taprootAddress,
		FullPublicKey:     line[fullPublicKeyIdx],
		MultisigAddress:   line[multisigAddressIdx],
		Idx:               line[idxIdx],
	}, nil
}
