package btcgrp

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/volatiletech/sqlboiler/types"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp/btc"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// Bitcoiner Bitcoin/BitcoinCash Interface
type Bitcoiner interface {
	//public_account.go -> wrapper of GetAddressInfo to return account
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
	GetBalanceByListUnspent(confirmationNum uint64) (btcutil.Amount, error)
	GetBalanceByAccount(accountType account.AccountType, confirmationNum uint64) (btcutil.Amount, error)

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
	Version() btc.BTCVersion
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

	//label.go
	SetLabel(addr, label string) error
	//GetReceivedByLabelAndMinConf(accountName string, minConf int) (btcutil.Amount, error)

	//logging.go
	Logging() (*btc.LoggingResult, error)

	//multisig.go
	AddMultisigAddress(requiredSigs int, addresses []string, accountName string, addressType address.AddrType) (*btc.AddMultisigAddressResult, error)

	//network.go
	GetNetworkInfo() (*btc.GetNetworkInfoResult, error)
	GetBlockchainInfo() (*btc.GetBlockchainInfoResult, error)

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
	ListUnspent(confirmationNum uint64) ([]btc.ListUnspentResult, error)
	ListUnspentByAccount(accountType account.AccountType, confirmationNum uint64) ([]btc.ListUnspentResult, error)
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
