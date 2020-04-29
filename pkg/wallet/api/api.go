package api

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/volatiletech/sqlboiler/types"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/btc"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/coin"
)

// Bitcoiner Bitcoin/BitcoinCash Interface
type Bitcoiner interface {
	//account.go
	GetAccount(addr string) (string, error)

	//address.go
	GetAddressInfo(addr string) (*btc.GetAddressInfoResult, error)
	GetAddressesByLabel(labelName string) ([]btcutil.Address, error)
	ValidateAddress(addr string) (*btc.ValidateAddressResult, error)
	DecodeAddress(addr string) (btcutil.Address, error)

	//amount.go
	AmountString(amt btcutil.Amount) string
	AmountToDecimal(amt btcutil.Amount) types.Decimal
	FloatToDecimal(f float64) types.Decimal
	FloatToAmount(f float64) (btcutil.Amount, error)
	StrToAmount(s string) (btcutil.Amount, error)
	StrSatoshiToAmount(s string) (btcutil.Amount, error)

	//balance.go
	GetBalance() (btcutil.Amount, error)
	GetBalanceByListUnspent() (btcutil.Amount, error)
	GetBalanceByAccount(accountType account.AccountType) (btcutil.Amount, error)

	//block.go
	GetBlockCount() (int64, error)

	//bitcoin.go
	Close()
	GetChainConf() *chaincfg.Params
	SetChainConf(conf *chaincfg.Params)
	SetChainConfNet(btcNet wire.BitcoinNet)
	ConfirmationBlock() uint64
	FeeRangeMax() float64
	FeeRangeMin() float64
	Version() coin.BTCVersion
	CoinTypeCode() coin.CoinTypeCode

	//fee.go
	EstimateSmartFee() (float64, error)
	GetTransactionFee(tx *wire.MsgTx) (btcutil.Amount, error)
	GetFee(tx *wire.MsgTx, adjustmentFee float64) (btcutil.Amount, error)

	//import.go
	ImportPrivKey(privKeyWIF *btcutil.WIF) error
	ImportPrivKeyLabel(privKeyWIF *btcutil.WIF, label string) error
	ImportPrivKeyWithoutReScan(privKeyWIF *btcutil.WIF, label string) error
	ImportAddress(pubkey string) error
	ImportAddressWithoutReScan(pubkey string) error
	ImportAddressWithLabel(address, label string, rescan bool) error

	//info.go
	GetNetworkInfo() (*btc.GetNetworkInfoResult, error)

	//label.go
	SetLabel(addr, label string) error
	//GetReceivedByLabelAndMinConf(accountName string, minConf int) (btcutil.Amount, error)

	//logging.go
	Logging() (*btc.LoggingResult, error)

	//multisig.go
	AddMultisigAddress(requiredSigs int, addresses []string, accountName string, addressType address.AddrType) (*btc.AddMultisigAddressResult, error)

	//transaction.go
	ToHex(tx *wire.MsgTx) (string, error)
	ToMsgTx(txHex string) (*wire.MsgTx, error)
	GetTransactionByTxID(txID string) (*btc.GetTransactionResult, error)
	GetTxOutByTxID(txID string, index uint32) (*btcjson.GetTxOutResult, error)
	DecodeRawTransaction(hexTx string) (*btc.TxRawResult, error)
	GetRawTransactionByHex(strHashTx string) (*btcutil.Tx, error)
	CreateRawTransaction(inputs []btcjson.TransactionInput, outputs map[btcutil.Address]btcutil.Amount) (*wire.MsgTx, error)
	FundRawTransaction(hex string) (*btc.FundRawTransactionResult, error)
	SignRawTransaction(tx *wire.MsgTx, prevtxs []btc.PrevTx) (*wire.MsgTx, bool, error)
	SignRawTransactionWithKey(tx *wire.MsgTx, privKeysWIF []string, prevtxs []btc.PrevTx) (*wire.MsgTx, bool, error)
	SendTransactionByHex(hex string) (*chainhash.Hash, error)
	SendTransactionByByte(rawTx []byte) (*chainhash.Hash, error)
	Sign(tx *wire.MsgTx, strPrivateKey string) (string, error)

	//unspent.go
	ListUnspent() ([]btc.ListUnspentResult, error)
	ListUnspentByAccount(accountType account.AccountType) ([]btc.ListUnspentResult, error)
	GetUnspentListAddrs(unspentList []btc.ListUnspentResult, accountType account.AccountType) []string
	LockUnspent(tx *btc.ListUnspentResult) error
	UnlockUnspent() error

	//wallet.go
	BackupWallet(fileName string) error
	DumpWallet(fileName string) error
	ImportWallet(fileName string) error
	EncryptWallet(passphrase string) error
	WalletLock() error
	WalletPassphrase(passphrase string, timeoutSecs int64) error
	WalletPassphraseChange(old, new string) error
	LoadWallet(fileName string) error
	UnLoadWallet(fileName string) error
	CreateWallet(fileName string, disablePrivKey bool) error
}
