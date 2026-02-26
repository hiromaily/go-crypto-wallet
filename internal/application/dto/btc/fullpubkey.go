package btc

import (
	"fmt"
	"strconv"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAuth "github.com/hiromaily/go-crypto-wallet/internal/domain/auth"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
)

// FullPubKeyFormat is fullpubkey csv format
type FullPubKeyFormat struct {
	CoinTypeCode   domainCoin.CoinTypeCode
	AuthType       domainAccount.AuthType
	Purpose        uint8  // BIP purpose (44, 49, 84, 86) - 0 means not set (legacy format)
	FullPubKey     string // Compressed public key (legacy format, 33 bytes hex)
	ExtendedPubKey string // Extended public key (xpub/tpub format) - optional, new format
	Fingerprint    string // Master key fingerprint (8 hex chars) - optional, new format
	DerivationPath string // BIP32 derivation path (e.g., m/48'/0'/0'/2') - optional, new format
}

// CreateFullPubKeyLine creates line for csv (legacy format with compressed pubkey only)
func CreateFullPubKeyLine(
	coinTypeCode domainCoin.CoinTypeCode,
	authType domainAccount.AuthType,
	fullPubKey string,
) string {
	// 0: coinTypeCode
	// 1: authType
	// 2: fullPubKey (compressed)
	return fmt.Sprintf("%s,%s,%s\n", coinTypeCode.String(), authType.String(), fullPubKey)
}

// CreateExtendedFullPubKeyLine creates line for csv with extended public key format (6-field format with purpose)
func CreateExtendedFullPubKeyLine(
	coinTypeCode domainCoin.CoinTypeCode,
	authType domainAccount.AuthType,
	purpose uint8,
	extendedPubKey string,
	fingerprint string,
	derivationPath string,
) string {
	// 0: coinTypeCode
	// 1: authType
	// 2: purpose (BIP purpose: 44, 49, 84, 86)
	// 3: extendedPubKey (xpub/tpub format)
	// 4: fingerprint (8 hex chars)
	// 5: derivationPath (e.g., m/48'/0'/0'/2')
	return fmt.Sprintf("%s,%s,%d,%s,%s,%s\n",
		coinTypeCode.String(),
		authType.String(),
		purpose,
		extendedPubKey,
		fingerprint,
		derivationPath,
	)
}

// ConvertFullPubKeyLine converts line to FullPubKeyFormat
// Supports legacy format (3 fields), old extended format (5 fields), and new format (6 fields with purpose)
func ConvertFullPubKeyLine(coinTypeCode domainCoin.CoinTypeCode, line []string) (*FullPubKeyFormat, error) {
	// Support three formats:
	// - Legacy: 3 fields (coin, auth, pubkey)
	// - Old extended: 5 fields (coin, auth, xpub, fingerprint, path) - purpose defaults to 49
	// - New: 6 fields (coin, auth, purpose, xpub, fingerprint, path)
	if len(line) != 3 && len(line) != 5 && len(line) != 6 {
		return nil, fmt.Errorf("csv format is invalid: expected 3, 5, or 6 fields, got %d", len(line))
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

	switch len(line) {
	case 3:
		// Legacy format: compressed public key only
		// Purpose defaults to 0 (not set)
		format.Purpose = 0
		format.FullPubKey = line[2]
	case 5:
		// Old extended format: xpub, fingerprint, derivation path (NO purpose field)
		// Purpose defaults to 49 for backward compatibility
		format.Purpose = 49
		format.ExtendedPubKey = line[2]
		format.Fingerprint = line[3]
		format.DerivationPath = line[4]
	case 6:
		// New format: purpose, xpub, fingerprint, derivation path
		purposeVal, err := strconv.ParseUint(line[2], 10, 8)
		if err != nil {
			return nil, fmt.Errorf("invalid purpose value: %s (must be 44, 49, 84, or 86)", line[2])
		}
		format.Purpose = uint8(purposeVal)
		// Validate purpose value using domain layer validation
		if err := domainAuth.Purpose(format.Purpose).Validate(); err != nil {
			return nil, fmt.Errorf("invalid purpose: %w", err)
		}
		format.ExtendedPubKey = line[3]
		format.Fingerprint = line[4]
		format.DerivationPath = line[5]
	}

	return format, nil
}
