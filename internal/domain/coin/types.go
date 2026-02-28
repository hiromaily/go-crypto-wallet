package coin

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
)

// CoinType creates a separate subtree for every cryptocurrency.
// This follows SLIP-0044 standard for HD wallet derivation paths.
// See: https://github.com/satoshilabs/slips/blob/master/slip-0044.md
type CoinType uint32

// Uint32 returns the uint32 representation of the coin type.
func (c CoinType) Uint32() uint32 {
	return uint32(c)
}

// Coin type constants following SLIP-0044 standard
const (
	// CoinTypeBitcoin represents Bitcoin (BIP44 coin type 0)
	CoinTypeBitcoin CoinType = 0

	// CoinTypeTestnet represents testnet for all coins (BIP44 coin type 1)
	CoinTypeTestnet CoinType = 1

	// CoinTypeLitecoin represents Litecoin (BIP44 coin type 2)
	CoinTypeLitecoin CoinType = 2

	// CoinTypeEther represents Ethereum (BIP44 coin type 60)
	CoinTypeEther CoinType = 60

	// CoinTypeRipple represents Ripple/XRP (BIP44 coin type 144)
	CoinTypeRipple CoinType = 144

	// CoinTypeBitcoinCash represents Bitcoin Cash (BIP44 coin type 145)
	CoinTypeBitcoinCash CoinType = 145

	// ERC20 tokens (temporary values, not part of SLIP-0044)
	// TODO: Review these temporary values

	// CoinTypeERC20 represents generic ERC20 tokens
	CoinTypeERC20 CoinType = 9000

	// CoinTypeERC20HYT represents HYT ERC20 token
	CoinTypeERC20HYT CoinType = 9001
)

// CoinTypeCode represents human-readable coin identifiers.
type CoinTypeCode string

// Coin type code constants
const (
	// BTC represents Bitcoin
	BTC CoinTypeCode = "btc"

	// BCH represents Bitcoin Cash
	BCH CoinTypeCode = "bch"

	// LTC represents Litecoin
	LTC CoinTypeCode = "ltc"

	// ETH represents Ethereum
	ETH CoinTypeCode = "eth"

	// XRP represents Ripple
	XRP CoinTypeCode = "xrp"

	// ERC20 represents generic ERC20 tokens
	ERC20 CoinTypeCode = "erc20"

	// HYT represents HYT token (custom ERC20)
	HYT CoinTypeCode = "hyt"
)

// String returns the string representation of the coin type code.
func (c CoinTypeCode) String() string {
	return string(c)
}

// CoinTypeCodeValue maps coin type codes to their SLIP-0044 coin types.
var CoinTypeCodeValue = map[CoinTypeCode]CoinType{
	BTC:   CoinTypeBitcoin,
	BCH:   CoinTypeBitcoinCash,
	LTC:   CoinTypeLitecoin,
	ETH:   CoinTypeEther,
	XRP:   CoinTypeRipple,
	ERC20: CoinTypeERC20,
	HYT:   CoinTypeERC20HYT,
}

// IsCoinTypeCode validates whether the given string is a valid coin type code.
func IsCoinTypeCode(val string) bool {
	_, ok := CoinTypeCodeValue[CoinTypeCode(val)]
	return ok
}

// IsBTCGroup returns true if the coin is part of the Bitcoin group (BTC, BCH).
//
// Deprecated: Use IsBTCCompatible for clarity. This function is kept for backward compatibility.
func IsBTCGroup(val CoinTypeCode) bool {
	return IsBTCCompatible(val)
}

// IsBTCOnly returns true only for BTC.
// Use this to guard BTC-specific features that BCH does not support:
//   - PSBT (BIP174)
//   - Descriptors (BIP380-386)
//   - SegWit (BIP141)
//   - Taproot (BIP340, BIP341, BIP342)
//   - MuSig2 (BIP327)
func IsBTCOnly(val CoinTypeCode) bool {
	return val == BTC
}

// IsBTCCompatible returns true for coins sharing Bitcoin base functionality (BTC, BCH).
// Use this for features common to both BTC and BCH:
//   - Raw transaction creation/signing
//   - P2PKH addresses
//   - P2SH multisig
//   - Basic UTXO operations
func IsBTCCompatible(val CoinTypeCode) bool {
	return val == BTC || val == BCH
}

// IsETHGroup returns true if the coin is part of the Ethereum group (ETH, ERC20 tokens).
func IsETHGroup(val CoinTypeCode) bool {
	return val == ETH || val == ERC20 || IsERC20Token(val.String())
}

// ERC20Token represents ERC20 token identifiers.
type ERC20Token string

// ERC20 token constants
const (
	// TokenHYT represents HYT token
	TokenHYT ERC20Token = "hyt"

	// TokenBAT represents Basic Attention Token
	// TokenBAT ERC20Token = "bat"
)

// String returns the string representation of the ERC20 token.
func (e ERC20Token) String() string {
	return string(e)
}

// ERC20Map maps known ERC20 tokens for validation.
var ERC20Map = map[ERC20Token]bool{
	TokenHYT: true,
	// TokenBAT: true,
}

// IsERC20Token validates whether the given string is a known ERC20 token.
func IsERC20Token(val string) bool {
	_, ok := ERC20Map[ERC20Token(val)]
	return ok
}

// BIP44Purpose is the purpose index for BIP-44 HD wallet derivation paths.
// All BIP-44 derivation paths use purpose = 44 (hardened).
const BIP44Purpose = 44

// BIP44AccountPath returns the BIP44 account-level derivation path string.
// Format: m/44'/coinType'/accountIndex'
//
// This is the path to the account-level xpub/xpriv used for HD wallet key derivation.
// The accountIndex comes from AccountType.BIP44AccountIndex().
// Example for ETH deposit account: m/44'/60'/0'
func BIP44AccountPath(coinTypeCode CoinTypeCode, accountIndex uint32) string {
	coinType := CoinTypeCodeValue[coinTypeCode].Uint32()
	return fmt.Sprintf("m/44'/%d'/%d'", coinType, accountIndex)
}

// GetCoinType returns CoinType based on network configuration
// This function has infrastructure dependency (chaincfg) and remains in this package
func GetCoinType(c CoinTypeCode, conf *chaincfg.Params) CoinType {
	if conf.Name != "mainnet" {
		return CoinTypeTestnet
	}
	if coinType, ok := CoinTypeCodeValue[c]; ok {
		return coinType
	}
	// coinType could not found
	return CoinTypeTestnet
}
