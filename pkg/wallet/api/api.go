package api

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/btc"
)

// Bitcoiner Bitcoin/BitcoinCash Interface
type Bitcoiner interface {
	//account.go
	GetAccount(addr string) (string, error)
	SetAccount(addr, account string) error
	GetReceivedByAccountAndMinConf(accountName string, minConf int) (btcutil.Amount, error)

	//address.go
	CreateNewAddress(accountName string) (btcutil.Address, error)
	GetAccountAddress(accountName string) (btcutil.Address, error)
	GetAddressesByLabel(labelName string) ([]btcutil.Address, error)
	GetAddressesByAccount(accountName string) ([]btcutil.Address, error)
	ValidateAddress(addr string) (*btcjson.ValidateAddressWalletResult, error)
	DecodeAddress(addr string) (btcutil.Address, error)
	GetAddressInfo(addr string) (*btc.GetAddressInfoResult, error)

	//amount.go
	AmountString(amt btcutil.Amount) string
	FloatBitToAmount(f float64) (btcutil.Amount, error)
	CastStrBitToAmount(s string) (btcutil.Amount, error)
	CastStrSatoshiToAmount(s string) (btcutil.Amount, error)
	ListAccounts() (map[string]btcutil.Amount, error)
	GetBalanceByAccount(accountName string) (btcutil.Amount, error)

	//block.go
	GetBlockCount() (int64, error)

	//bitcoin.go
	Close()
	GetChainConf() *chaincfg.Params
	SetChainConf(conf *chaincfg.Params)
	SetChainConfNet(btcNet wire.BitcoinNet)
	Client() *rpcclient.Client
	ConfirmationBlock() int
	FeeRangeMax() float64
	FeeRangeMin() float64
	SetVersion(ver enum.BTCVersion)
	Version() enum.BTCVersion
	SetCoinType(coinType enum.CoinType)
	CoinType() enum.CoinType

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
	GetReceivedByLabelAndMinConf(accountName string, minConf int) (btcutil.Amount, error)

	//logging.go
	Logging() (*btc.LoggingResult, error)

	//multisig.go
	AddMultisigAddress(requiredSigs int, addresses []string, accountName string, addressType enum.AddressType) (*btc.AddMultisigAddressResult, error)

	//transaction.go
	ToHex(tx *wire.MsgTx) (string, error)
	ToMsgTx(txHex string) (*wire.MsgTx, error)
	DecodeRawTransaction(hexTx string) (*btcjson.TxRawResult, error)
	GetRawTransactionByHex(strHashTx string) (*btcutil.Tx, error)
	GetTransactionByTxID(txID string) (*btcjson.GetTransactionResult, error)
	GetTxOutByTxID(txID string, index uint32) (*btcjson.GetTxOutResult, error)
	CreateRawTransaction(sendAddr string, amount btcutil.Amount, inputs []btcjson.TransactionInput) (*wire.MsgTx, error)
	CreateRawTransactionWithOutput(inputs []btcjson.TransactionInput, outputs map[btcutil.Address]btcutil.Amount) (*wire.MsgTx, error)
	FundRawTransaction(hex string) (*btc.FundRawTransactionResult, error)
	SignRawTransaction(tx *wire.MsgTx, prevtxs []btc.PrevTx) (*wire.MsgTx, bool, error)
	SignRawTransactionWithKey(tx *wire.MsgTx, privKeysWIF []string, prevtxs []btc.PrevTx) (*wire.MsgTx, bool, error)
	SendTransactionByHex(hex string) (*chainhash.Hash, error)
	SendTransactionByByte(rawTx []byte) (*chainhash.Hash, error)
	Sign(tx *wire.MsgTx, strPrivateKey string) (string, error)

	//unspent.go
	UnlockAllUnspentTransaction() error
	LockUnspent(tx btcjson.ListUnspentResult) error
	ListUnspent() ([]btc.ListUnspentResult, error)
	ListUnspentByAccount(accountType account.AccountType) ([]btc.ListUnspentResult, []btcutil.Address, error)

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
