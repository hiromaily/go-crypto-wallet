package btc

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/quagmt/udecimal"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainBitcoin "github.com/hiromaily/go-crypto-wallet/internal/domain/bitcoin"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
)

// Bitcoiner is the full Bitcoin/BitcoinCash interface.
//
// Deprecated: This interface is too large and violates the Interface Segregation Principle.
// New code should use the small, focused interfaces from interfaces_small.go instead:
//
//   - ChainConfigProvider: Chain configuration and coin type
//   - AmountConverter: Amount format conversion
//   - UTXOProvider: UTXO queries
//   - RawTransactionCreator: Raw transaction creation
//   - AddressOperator: Address operations
//   - BalanceChecker: Balance queries
//   - TransactionSender: Transaction broadcast
//   - TransactionMonitor: Transaction monitoring
//   - AddressImporter: Address import
//   - DescriptorManager: Descriptor operations (BTC only)
//   - PSBTHandler: PSBT operations (BTC only)
//   - RawTransactionSigner: Raw transaction signing
//   - MultisigManager: Multisig address management
//
// For chain-specific interfaces:
//   - BitcoinCompatible: Methods common to BTC and BCH
//   - BTCOnly: Methods specific to BTC (PSBT, Descriptors, MuSig2)
//   - BCHer: Interface for BCH operations
//
// See:
//   - interfaces_small.go: Small, focused interfaces (recommended)
//   - bitcoin_compatible.go: BitcoinCompatible interface
//   - btc_only.go: BTCOnly interface
//   - bch_only.go: BCHer interface
type Bitcoiner interface {
	// public_account.go -> wrapper of GetAddressInfo to return account
	GetAccount(addr string) (string, error)

	// address.go
	GetAddressInfo(addr string) (*dtobtc.AddressInfo, error)
	GetAddressesByLabel(labelName string) ([]btcutil.Address, error)
	ValidateAddress(addr string) (*dtobtc.ValidateAddressResult, error)
	DecodeAddress(addr string) (btcutil.Address, error)

	// amount.go
	AmountString(amt btcutil.Amount) string
	AmountToDecimal(amt btcutil.Amount) (udecimal.Decimal, error)
	FloatToDecimal(f float64) (udecimal.Decimal, error)
	FloatToAmount(f float64) (btcutil.Amount, error)
	StrToAmount(s string) (btcutil.Amount, error)
	StrSatoshiToAmount(s string) (btcutil.Amount, error)

	// balance.go
	GetBalance() (btcutil.Amount, error)
	GetBalanceByListUnspent(confirmationNum uint64) (btcutil.Amount, error)
	GetBalanceByAccount(accountType domainAccount.AccountType, confirmationNum uint64) (btcutil.Amount, error)

	// block.go
	GetBlockCount() (int64, error)

	// bitcoin.go
	Close()
	GetChainConf() *chaincfg.Params
	SetChainConf(conf *chaincfg.Params)
	SetChainConfNet(btcNet wire.BitcoinNet)
	ConfirmationBlock() uint64
	FeeRangeMax() float64
	FeeRangeMin() float64
	Version() domainBitcoin.Version
	CoinTypeCode() domainCoin.CoinTypeCode

	// fee.go
	EstimateSmartFee() (float64, error)
	GetTransactionFee(tx *wire.MsgTx) (btcutil.Amount, error)
	GetFee(tx *wire.MsgTx, adjustmentFee float64) (btcutil.Amount, error)

	// import.go
	ImportPrivKey(privKeyWIF *btcutil.WIF) error
	ImportPrivKeyLabel(privKeyWIF *btcutil.WIF, label string) error
	ImportPrivKeyWithoutReScan(privKeyWIF *btcutil.WIF, label string) error
	ImportAddress(pubkey string) error
	ImportAddressWithoutReScan(pubkey string) error
	ImportAddressWithLabel(addr, label string, rescan bool) error
	ImportDescriptors(requests []dtobtc.ImportDescriptorsRequest) ([]dtobtc.ImportDescriptorsResponse, error)
	ImportMulti(
		requests []dtobtc.ImportMultiRequest,
		options *dtobtc.ImportMultiOptions,
	) ([]dtobtc.ImportMultiResponse, error)

	// descriptor_info.go
	GetDescriptorInfo(descriptor string) (*dtobtc.DescriptorInfo, error)
	ListDescriptors(privateDescriptors bool) (*dtobtc.ListDescriptorsResult, error)

	// label.go
	SetLabel(addr, label string) error
	// GetReceivedByLabelAndMinConf(accountName string, minConf int) (btcutil.Amount, error)

	// logging.go
	Logging() (*dtobtc.LoggingResult, error)

	// multisig.go
	AddMultisigAddress(
		requiredSigs int, addresses []string, accountName string, addressType domainBitcoin.AddressType,
	) (*dtobtc.MultisigAddress, error)

	// network.go
	GetNetworkInfo() (*dtobtc.NetworkInfo, error)
	GetBlockchainInfo() (*dtobtc.BlockchainInfo, error)

	// transaction.go
	ToHex(tx *wire.MsgTx) (string, error)
	ToMsgTx(txHex string) (*wire.MsgTx, error)
	GetTransactionByTxID(txID string) (*dtobtc.TransactionResult, error)
	GetTxOutByTxID(txID string, index uint32) (*btcjson.GetTxOutResult, error)
	DecodeRawTransaction(hexTx string) (*dtobtc.RawTransaction, error)
	GetRawTransactionByHex(strHashTx string) (*btcutil.Tx, error)
	CreateRawTransaction(
		inputs []btcjson.TransactionInput, outputs map[btcutil.Address]btcutil.Amount,
	) (*wire.MsgTx, error)
	FundRawTransaction(hex string) (*dtobtc.FundRawTransactionResult, error)
	SignRawTransaction(tx *wire.MsgTx, prevtxs []dtobtc.PreviousTx) (*wire.MsgTx, bool, error)
	SignRawTransactionWithKey(
		tx *wire.MsgTx, privKeysWIF []string, prevtxs []dtobtc.PreviousTx,
	) (*wire.MsgTx, bool, error)
	SendTransactionByHex(hex string) (*chainhash.Hash, error)
	SendTransactionByByte(rawTx []byte) (*chainhash.Hash, error)
	Sign(tx *wire.MsgTx, strPrivateKey string) (string, error)

	// psbt.go (BIP174 Partially Signed Bitcoin Transaction support)
	CreatePSBT(msgTx *wire.MsgTx, prevTxs []dtobtc.PreviousTx, senderAccount domainAccount.AccountType) (string, error)
	ParsePSBT(psbtBase64 string) (*dtobtc.ParsedPSBT, error)
	ValidatePSBT(psbtBase64 string) error
	SignPSBTWithKey(psbtBase64 string, wifs []string) (string, bool, error)
	WalletProcessPsbt(psbtBase64 string, sign bool) (string, bool, error) // RPC-based signing for descriptor wallets
	FinalizePSBT(psbtBase64 string) (string, error)
	ExtractTransaction(psbtBase64 string) (*wire.MsgTx, error)
	IsPSBTComplete(psbtBase64 string) (bool, error)
	GetPSBTFee(psbtBase64 string) (int64, error)

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

	// wallet.go
	BackupWallet(fileName string) error
	DumpWallet(fileName string) error
	ImportWallet(fileName string) error
	EncryptWallet(passphrase string) error
	WalletLock() error
	WalletPassphrase(passphrase string, timeoutSecs int64) error
	WalletPassphraseChange(old, newPass string) error
	LoadWallet(fileName string) error
	UnLoadWallet(fileName string) error
	CreateWallet(fileName string, disablePrivKey bool) error
}
