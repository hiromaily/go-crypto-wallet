package address

import (
	"fmt"
	"strconv"

	appdto "github.com/hiromaily/go-crypto-wallet/internal/application/dto"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
)

// CreateLine creates line for csv from BtcAccountKey
func CreateLine(accountKeyItem *sqlcgen.BtcAccountKey) []string {
	taprootAddr := ""
	if accountKeyItem.TaprootAddress.Valid {
		taprootAddr = accountKeyItem.TaprootAddress.String
	}
	return []string{
		string(accountKeyItem.Coin),
		string(accountKeyItem.Account),
		accountKeyItem.P2pkhAddress,
		accountKeyItem.P2shSegwitAddress,
		accountKeyItem.Bech32Address,
		taprootAddr,
		accountKeyItem.FullPublicKey,
		accountKeyItem.MultisigAddress,
		strconv.Itoa(int(accountKeyItem.Idx)),
	}
}

// CreateEthLine creates line for csv from EthAccountKey
func CreateEthLine(accountKeyItem *sqlcgen.EthAccountKey) []string {
	return []string{
		"eth",
		string(accountKeyItem.Account),
		accountKeyItem.Address,
		accountKeyItem.FullPublicKey,
		strconv.Itoa(int(accountKeyItem.Idx)),
	}
}

// ConvertLine converts line to AddressFormat
func ConvertLine(coinTypeCode domainCoin.CoinTypeCode, line []string) (*appdto.AddressFormat, error) {
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

	return &appdto.AddressFormat{
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
