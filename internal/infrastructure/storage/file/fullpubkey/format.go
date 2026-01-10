package fullpubkey

import (
	"fmt"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
)

// FullPubKeyFormat is fullpubkey csv format
type FullPubKeyFormat struct {
	CoinTypeCode   domainCoin.CoinTypeCode
	AuthType       domainAccount.AuthType
	FullPubKey     string // Compressed public key (legacy format, 33 bytes hex)
	ExtendedPubKey string // Extended public key (xpub/tpub format) - optional, new format
	Fingerprint    string // Master key fingerprint (8 hex chars) - optional, new format
	DerivationPath string // BIP32 derivation path (e.g., m/48'/0'/0'/2') - optional, new format
}

// CreateLine creates line for csv (legacy format with compressed pubkey only)
func CreateLine(coinTypeCode domainCoin.CoinTypeCode, authType domainAccount.AuthType, fullPubKey string) string {
	// 0: coinTypeCode
	// 1: authType
	// 2: fullPubKey (compressed)
	return fmt.Sprintf("%s,%s,%s\n", coinTypeCode.String(), authType.String(), fullPubKey)
}

// CreateExtendedLine creates line for csv with extended public key format
func CreateExtendedLine(
	coinTypeCode domainCoin.CoinTypeCode,
	authType domainAccount.AuthType,
	extendedPubKey string,
	fingerprint string,
	derivationPath string,
) string {
	// 0: coinTypeCode
	// 1: authType
	// 2: extendedPubKey (xpub/tpub format)
	// 3: fingerprint (8 hex chars)
	// 4: derivationPath (e.g., m/48'/0'/0'/2')
	return fmt.Sprintf("%s,%s,%s,%s,%s\n",
		coinTypeCode.String(),
		authType.String(),
		extendedPubKey,
		fingerprint,
		derivationPath,
	)
}

// ConvertLine converts line to FullPubKeyFormat
// Supports both legacy format (3 fields) and extended format (5 fields)
func ConvertLine(coinTypeCode domainCoin.CoinTypeCode, line []string) (*FullPubKeyFormat, error) {
	// Support legacy format (3 fields) and extended format (5 fields)
	if len(line) != 3 && len(line) != 5 {
		return nil, fmt.Errorf("csv format is invalid: expected 3 or 5 fields, got %d", len(line))
	}

	// validate
	if !domainCoin.IsCoinTypeCode(line[0]) || domainCoin.CoinTypeCode(line[0]) != coinTypeCode {
		return nil, fmt.Errorf("coinTypeCode is invalid. got %s, want %s", line[0], coinTypeCode.String())
	}
	if !domainAccount.ValidateAuthType(line[1]) {
		return nil, fmt.Errorf("auth account is invalid: %s", line[1])
	}

	format := &FullPubKeyFormat{
		CoinTypeCode: domainCoin.CoinTypeCode(line[0]),
		AuthType:     domainAccount.AuthType(line[1]),
	}

	if len(line) == 3 {
		// Legacy format: compressed public key only
		format.FullPubKey = line[2]
	} else {
		// Extended format: xpub, fingerprint, derivation path
		format.ExtendedPubKey = line[2]
		format.Fingerprint = line[3]
		format.DerivationPath = line[4]
	}

	return format, nil
}
