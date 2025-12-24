package fullpubkey

import (
	"errors"
	"fmt"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
)

// FullPubKeyFormat is fullpubkey csv format
type FullPubKeyFormat struct {
	CoinTypeCode domainCoin.CoinTypeCode
	AuthType     domainAccount.AuthType
	FullPubKey   string
}

// CreateLine creates line for csv
func CreateLine(coinTypeCode domainCoin.CoinTypeCode, authType domainAccount.AuthType, fullPubKey string) string {
	// 0: coinTypeCode
	// 1: authType
	// 2: fullPubKey
	return fmt.Sprintf("%s,%s,%s\n", coinTypeCode.String(), authType.String(), fullPubKey)
}

// ConvertLine converts line to FullPubKeyFormat
func ConvertLine(coinTypeCode domainCoin.CoinTypeCode, line []string) (*FullPubKeyFormat, error) {
	if len(line) != 3 {
		return nil, errors.New("csv format is invalid")
	}

	// validate
	if !domainCoin.IsCoinTypeCode(line[0]) || domainCoin.CoinTypeCode(line[0]) != coinTypeCode {
		return nil, fmt.Errorf("coinTypeCode is invalid. got %s, want %s", line[0], coinTypeCode.String())
	}
	if !domainAccount.ValidateAuthType(line[1]) {
		return nil, fmt.Errorf("auth account is invalid: %s", line[1])
	}

	return &FullPubKeyFormat{
		CoinTypeCode: domainCoin.CoinTypeCode(line[0]),
		AuthType:     domainAccount.AuthType(line[1]),
		FullPubKey:   line[2],
	}, nil
}
