package ripple

import (
	"context"

	"github.com/btcsuite/btcd/chaincfg"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/ripple/xrp"
)

// Rippler Ripple Interface
type Rippler interface {
	RippleAdminer
	RipplePublicer
	RippleAPIer

	// balance
	GetBalance(ctx context.Context, addr string) (float64, error)
	GetTotalBalance(ctx context.Context, addrs []string) float64

	// transaction
	CreateRawTransaction(
		ctx context.Context, senderAccount, receiverAccount string, amount float64, instructions *xrp.Instructions,
	) (*xrp.TxInput, string, error)

	// ripple
	Close() error
	CoinTypeCode() domainCoin.CoinTypeCode
	GetChainConf() *chaincfg.Params
}

// RippleAPIer is RippleAPI interface
type RippleAPIer interface {
	// RippleAccountAPI
	GetAccountInfo(ctx context.Context, address string) (*xrp.ResponseGetAccountInfo, error)
	// RippleAddressAPI
	GenerateAddress(ctx context.Context) (*xrp.ResponseGenerateAddress, error)
	GenerateXAddress(ctx context.Context) (*xrp.ResponseGenerateXAddress, error)
	IsValidAddress(ctx context.Context, addr string) (bool, error)
	// RippleTxAPI
	PrepareTransaction(
		ctx context.Context, senderAccount, receiverAccount string, amount float64, instructions *xrp.Instructions,
	) (*xrp.TxInput, string, error)
	SignTransaction(ctx context.Context, txJSON *xrp.TxInput, secret string) (string, string, error)
	CombineTransaction(ctx context.Context, signedTxs []string) (string, string, error)
	SubmitTransaction(ctx context.Context, signedTx string) (*xrp.SentTx, uint64, error)
	WaitValidation(ctx context.Context, targetledgerVarsion uint64) (uint64, error)
	GetTransaction(ctx context.Context, txID string, targetLedgerVersion uint64) (*xrp.TxInfo, error)
}

// RipplePublicer is RipplePublic interface
type RipplePublicer interface {
	// public_account
	AccountChannels(ctx context.Context, sender, receiver string) (*xrp.ResponseAccountChannels, error)
	AccountInfo(ctx context.Context, address string) (*xrp.ResponseAccountInfo, error)
	// public_server_info
	ServerInfo(ctx context.Context) (*xrp.ResponseServerInfo, error)
}

// RippleAdminer is RippleAdmin interface
type RippleAdminer interface {
	// admin_keygen
	ValidationCreate(ctx context.Context, secret string) (*xrp.ResponseValidationCreate, error)
	WalletProposeWithKey(ctx context.Context, seed string, keyType xrp.XRPKeyType) (*xrp.ResponseWalletPropose, error)
	WalletPropose(ctx context.Context, passphrase string) (*xrp.ResponseWalletPropose, error)
}
