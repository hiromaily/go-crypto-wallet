package btc

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainBTC "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/btc"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
)

// BitcoinCompatible defines methods common to both BTC and BCH.
// BCH implementations should use this interface instead of the full Bitcoiner interface.
//
// This interface excludes BTC-only features:
//   - PSBT (Partially Signed Bitcoin Transactions) - BIP174
//   - Descriptors - BIP380, BIP381, BIP382, BIP386
//   - MuSig2 - BIP327 (Schnorr-based multisig)
//
// BCH uses:
//   - Raw Transaction Hex format (not PSBT)
//   - SIGHASH_FORKID (0x41) for transaction signing
//   - ImportAddress instead of ImportDescriptors
//   - P2SH multisig (no P2WSH, P2TR)
//
//nolint:iface // This interface is a subset of Bitcoiner for BCH compatibility
type BitcoinCompatible interface {
	// public_account.go -> wrapper of GetAddressInfo to return account
	GetAccount(addr string) (string, error)

	// address.go
	GetAddressInfo(addr string) (*dtobtc.AddressInfo, error)
	GetAddressesByLabel(labelName string) ([]btcutil.Address, error)
	ValidateAddress(addr string) (*dtobtc.ValidateAddressResult, error)
	DecodeAddress(addr string) (btcutil.Address, error)

	// balance.go
	GetBalance() (btcutil.Amount, error)
	GetBalanceByListUnspent(confirmationNum uint64) (btcutil.Amount, error)
	GetBalanceByAccount(accountType domainAccount.AccountType, confirmationNum uint64) (btcutil.Amount, error)

	// bitcoin.go
	Close()
	GetChainConf() *chaincfg.Params
	SetChainConf(conf *chaincfg.Params)
	SetChainConfNet(btcNet wire.BitcoinNet)
	ConfirmationBlock() uint64
	FeeRangeMax() float64
	FeeRangeMin() float64
	Version() domainBTC.Version
	CoinTypeCode() domainCoin.CoinTypeCode

	// fee.go
	GetTransactionFee(tx *wire.MsgTx) (btcutil.Amount, error)
	GetFee(tx *wire.MsgTx, adjustmentFee float64) (btcutil.Amount, error)

	// import.go (traditional methods - works for both BTC and BCH)
	ImportPrivKey(privKeyWIF *btcutil.WIF) error
	ImportPrivKeyLabel(privKeyWIF *btcutil.WIF, label string) error
	ImportPrivKeyWithoutReScan(privKeyWIF *btcutil.WIF, label string) error
	ImportAddress(pubkey string) error
	ImportAddressWithoutReScan(pubkey string) error
	ImportAddressWithLabel(addr, label string, rescan bool) error

	// label.go
	SetLabel(addr, label string) error

	// multisig.go (P2SH works for both, but BCH should not use P2WSH or P2TR address types)
	AddMultisigAddress(
		requiredSigs int, addresses []string, accountName string, addressType domainBTC.AddressType,
	) (*dtobtc.MultisigAddress, error)

	// transaction.go (Raw TX operations - works for both BTC and BCH)
	ToHex(tx *wire.MsgTx) (string, error)
	ToMsgTx(txHex string) (*wire.MsgTx, error)
	GetTransactionByTxID(txID string) (*dtobtc.TransactionResult, error)
	GetTxOutByTxID(txID string, index uint32) (string, error)
	DecodeRawTransaction(hexTx string) (string, error)
	GetRawTransactionByHex(strHashTx string) (*btcutil.Tx, error)
	CreateRawTransaction(
		inputs []dtobtc.TransactionInput, outputs map[btcutil.Address]btcutil.Amount,
	) (*wire.MsgTx, error)
	FundRawTransaction(hex string) (string, error)
	SignRawTransaction(tx *wire.MsgTx, prevtxs []dtobtc.PreviousTx) (*wire.MsgTx, bool, error)
	SignRawTransactionWithKey(
		tx *wire.MsgTx, privKeysWIF []string, prevtxs []dtobtc.PreviousTx,
	) (*wire.MsgTx, bool, error)
	SendTransactionByHex(hex string) (*chainhash.Hash, error)
	SendTransactionByByte(rawTx []byte) (*chainhash.Hash, error)
	Sign(tx *wire.MsgTx, strPrivateKey string) (string, error)

	// unspent.go
	ListUnspent(confirmationNum uint64) ([]dtobtc.UnspentOutput, error)
	ListUnspentByAccount(
		accountType domainAccount.AccountType, confirmationNum uint64,
	) ([]dtobtc.UnspentOutput, error)
	GetUnspentListAddrs(
		unspentList []dtobtc.UnspentOutput, accountType domainAccount.AccountType,
	) []string
	LockUnspent(tx *dtobtc.UnspentOutput) error
	UnlockUnspent() error
}
