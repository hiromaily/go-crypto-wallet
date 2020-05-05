package fullpubkey

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/coin"
)

type FullPubKeyFormat struct {
	CoinTypeCode coin.CoinTypeCode
	AuthType     account.AuthType
	FullPubKey   string
}

func CreateLine(coinTypeCode coin.CoinTypeCode, authType account.AuthType, fullPubKey string) string {
	//0: coinTypeCode
	//1: authType
	//2: fullPubKey
	return fmt.Sprintf("%s,%s,%s\n", coinTypeCode.String(), authType.String(), fullPubKey)
}

func ConvertLine(coinTypeCode coin.CoinTypeCode, line []string) (*FullPubKeyFormat, error) {
	if len(line) != 3 {
		return nil, errors.New("csv format is invalid")
	}

	// validate
	if !coin.ValidateCoinTypeCode(line[0]) || coin.CoinTypeCode(line[0]) != coinTypeCode {
		return nil, errors.Errorf("coinTypeCode is invalid. got %s, want %s", line[0], coinTypeCode.String())
	}
	if !account.ValidateAuthType(line[1]) {
		return nil, errors.Errorf("auth account is invalid: %s", line[1])
	}

	return &FullPubKeyFormat{
		CoinTypeCode: coin.CoinTypeCode(line[0]),
		AuthType:     account.AuthType(line[1]),
		FullPubKey:   line[2],
	}, nil
}
