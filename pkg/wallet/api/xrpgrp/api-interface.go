package xrpgrp

import (
	"github.com/btcsuite/btcd/chaincfg"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
	pb "github.com/hiromaily/ripple-lib-proto/pb/go/rippleapi"
)

// Rippler Ripple Interface
type Rippler interface {
	RippleAdminer
	RipplePublicer
	RippleAPIer

	// ripple
	Close()
	CoinTypeCode() coin.CoinTypeCode
	GetChainConf() *chaincfg.Params
}

// RippleAPIer is RippleAPI interface
type RippleAPIer interface {
	// RippleAccountAPI
	GetAccountInfo(address string) (*pb.ResponseGetAccountInfo, error)
	// RippleAddressAPI
	GenerateAddress() (*pb.ResponseGenerateAddress, error)
	GenerateXAddress() (*pb.ResponseGenerateXAddress, error)
	IsValidAddress(addr string) (bool, error)
	// RippleTxAPI
	PrepareTransaction(senderAccount, receiverAccount string, amount float64) (*xrp.TxInput, error)
	SignTransaction(txJSON *xrp.TxInput, secret string) (string, string, error)
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
