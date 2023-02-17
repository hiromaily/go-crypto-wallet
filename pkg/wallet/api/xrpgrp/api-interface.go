package xrpgrp

import (
	"github.com/btcsuite/btcd/chaincfg"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// Rippler Ripple Interface
type Rippler interface {
	RippleAdminer
	RipplePublicer
	RippleAPIer

	// balance
	GetBalance(addr string) (float64, error)
	GetTotalBalance(addrs []string) float64

	// transaction
	CreateRawTransaction(senderAccount, receiverAccount string, amount float64, instructions *xrp.Instructions) (*xrp.TxInput, string, error)

	// ripple
	Close() error
	CoinTypeCode() coin.CoinTypeCode
	GetChainConf() *chaincfg.Params
}

// RippleAPIer is RippleAPI interface
type RippleAPIer interface {
	// RippleAccountAPI
	GetAccountInfo(address string) (*xrp.ResponseGetAccountInfo, error)
	// RippleAddressAPI
	GenerateAddress() (*xrp.ResponseGenerateAddress, error)
	GenerateXAddress() (*xrp.ResponseGenerateXAddress, error)
	IsValidAddress(addr string) (bool, error)
	// RippleTxAPI
	PrepareTransaction(senderAccount, receiverAccount string, amount float64, instructions *xrp.Instructions) (*xrp.TxInput, string, error)
	SignTransaction(txJSON *xrp.TxInput, secret string) (string, string, error)
	CombineTransaction(signedTxs []string) (string, string, error)
	SubmitTransaction(signedTx string) (*xrp.SentTx, uint64, error)
	WaitValidation(targetledgerVarsion uint64) (uint64, error)
	GetTransaction(txID string, targetLedgerVersion uint64) (*xrp.TxInfo, error)
}

// RipplePublicer is RipplePublic interface
type RipplePublicer interface {
	// public_account
	AccountChannels(sender, receiver string) (*xrp.ResponseAccountChannels, error)
	AccountInfo(address string) (*xrp.ResponseAccountInfo, error)
	// public_server_info
	ServerInfo() (*xrp.ResponseServerInfo, error)
}

// RippleAdminer is RippleAdmin interface
type RippleAdminer interface {
	// admin_keygen
	ValidationCreate(secret string) (*xrp.ResponseValidationCreate, error)
	WalletProposeWithKey(seed string, keyType xrp.XRPKeyType) (*xrp.ResponseWalletPropose, error)
	WalletPropose(passphrase string) (*xrp.ResponseWalletPropose, error)
}
