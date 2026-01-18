package address

import (
	"fmt"
	"strconv"

	appdto "github.com/hiromaily/go-crypto-wallet/internal/application/dto"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainBitcoin "github.com/hiromaily/go-crypto-wallet/internal/domain/bitcoin"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
)

// CreateLine creates line for csv from BTCAccountKey
func CreateLine(accountKeyItem *domainBitcoin.BTCAccountKey) []string {
	taprootAddr := ""
	if accountKeyItem.TaprootAddress != nil {
		taprootAddr = *accountKeyItem.TaprootAddress
	}
	return []string{
		accountKeyItem.CoinTypeCode.String(),
		accountKeyItem.Account.String(),
		accountKeyItem.P2pkhAddress,
		accountKeyItem.P2shSegwitAddress,
		accountKeyItem.Bech32Address,
		taprootAddr,
		accountKeyItem.FullPublicKey,
		accountKeyItem.MultisigAddress,
		accountKeyItem.RedeemScript, // Add redeem script for P2SH multisig
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
	// Support: old format (8 fields), with Taproot (9 fields), with RedeemScript (10 fields)
	if len(line) != 8 && len(line) != 9 && len(line) != 10 {
		return nil, fmt.Errorf("csv format is invalid: expected 8, 9, or 10 fields, got %d", len(line))
	}

	// validate
	if !domainCoin.IsCoinTypeCode(line[0]) || domainCoin.CoinTypeCode(line[0]) != coinTypeCode {
		return nil, fmt.Errorf("coinTypeCode is invalid. got %s, want %s", line[0], coinTypeCode.String())
	}
	if !domainAccount.ValidateAccountType(line[1]) {
		return nil, fmt.Errorf("account is invalid: %s", line[1])
	}

	// For backward compatibility with old CSV formats
	taprootAddress := ""
	redeemScript := ""
	fullPublicKeyIdx := 5
	multisigAddressIdx := 6
	idxIdx := 7

	if len(line) == 9 {
		// Format with Taproot address (no redeemScript)
		taprootAddress = line[5]
		fullPublicKeyIdx = 6
		multisigAddressIdx = 7
		idxIdx = 8
	} else if len(line) == 10 {
		// New format with Taproot address and RedeemScript
		taprootAddress = line[5]
		fullPublicKeyIdx = 6
		multisigAddressIdx = 7
		redeemScript = line[8]
		idxIdx = 9
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
		RedeemScript:      redeemScript,
		Idx:               line[idxIdx],
	}, nil
}
