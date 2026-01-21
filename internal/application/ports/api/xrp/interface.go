// Package xrp defines interfaces for Ripple/XRP blockchain operations.
//
// This package follows the Dependency Inversion Principle of Clean Architecture
// by defining interfaces in the application layer that are implemented by the
// infrastructure layer.
package xrp

import (
	"context"

	"github.com/btcsuite/btcd/chaincfg"

	dtoRipple "github.com/hiromaily/go-crypto-wallet/internal/application/dto/ripple"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
)

// Rippler defines the main interface for Ripple/XRP blockchain operations.
// It embeds specialized interfaces for admin, public, and API operations.
type Rippler interface {
	RippleAdminer
	RipplePublicer
	RippleAPIer

	// balance
	GetBalance(ctx context.Context, addr string) (float64, error)
	GetTotalBalance(ctx context.Context, addrs []string) float64

	// transaction
	CreateRawTransaction(
		ctx context.Context,
		senderAccount, receiverAccount string,
		amount float64,
		instructions *dtoRipple.Instructions,
	) (*dtoRipple.TxInput, string, error)

	// ripple
	Close() error
	CoinTypeCode() domainCoin.CoinTypeCode
	GetChainConf() *chaincfg.Params
}

// SignerEntryInput represents a signer for creating SignerListSet transactions
type SignerEntryInput struct {
	Account string
	Weight  uint32
}

// RippleAPIer defines the interface for Ripple API operations.
// Implementations handle account management, address generation, and transaction operations.
type RippleAPIer interface {
	// RippleAccountAPI
	GetAccountInfo(ctx context.Context, address string) (*dtoRipple.ResponseGetAccountInfo, error)
	// RippleAddressAPI
	GenerateAddress(ctx context.Context) (*dtoRipple.ResponseGenerateAddress, error)
	GenerateXAddress(ctx context.Context) (*dtoRipple.ResponseGenerateXAddress, error)
	IsValidAddress(ctx context.Context, addr string) (bool, error)
	// RippleTxAPI
	PrepareTransaction(
		ctx context.Context,
		senderAccount, receiverAccount string,
		amount float64,
		instructions *dtoRipple.Instructions,
	) (*dtoRipple.TxInput, string, error)
	SignTransaction(ctx context.Context, txJSON *dtoRipple.TxInput, secret string) (string, string, error)
	CombineTransaction(ctx context.Context, signedTxs []string) (string, string, error)
	SubmitTransaction(ctx context.Context, signedTx string) (*dtoRipple.SentTx, uint64, error)
	WaitValidation(ctx context.Context, targetledgerVarsion uint64) (uint64, error)
	GetTransaction(ctx context.Context, txID string, targetLedgerVersion uint64) (*dtoRipple.TxInfo, error)

	// Multi-signature operations
	// Reference: https://xrpl.org/docs/concepts/accounts/multi-signing
	PrepareSignerListSetTransaction(
		ctx context.Context,
		senderAccount string,
		signerQuorum uint32,
		signerEntries []SignerEntryInput,
		instructions *dtoRipple.Instructions,
	) (*dtoRipple.SignerListSetTxInput, string, error)

	// Trust line operations
	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/trustset
	PrepareTrustSetTransaction(
		ctx context.Context,
		senderAccount string,
		limitAmount *dtoRipple.IssuedCurrencyAmount,
		qualityIn, qualityOut uint32,
		instructions *dtoRipple.Instructions,
	) (*dtoRipple.TrustSetTxInput, string, error)

	// Escrow operations
	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/escrowcreate
	PrepareEscrowCreateTransaction(
		ctx context.Context,
		senderAccount, destinationAccount string,
		amount float64,
		cancelAfter, finishAfter uint32,
		condition string,
		destinationTag uint32,
		instructions *dtoRipple.Instructions,
	) (*dtoRipple.EscrowCreateTxInput, string, error)

	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/escrowfinish
	PrepareEscrowFinishTransaction(
		ctx context.Context,
		senderAccount, owner string,
		offerSequence uint32,
		condition, fulfillment string,
		instructions *dtoRipple.Instructions,
	) (*dtoRipple.EscrowFinishTxInput, string, error)

	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/escrowcancel
	PrepareEscrowCancelTransaction(
		ctx context.Context,
		senderAccount, owner string,
		offerSequence uint32,
		instructions *dtoRipple.Instructions,
	) (*dtoRipple.EscrowCancelTxInput, string, error)

	// PaymentChannel operations
	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/paymentchannelcreate
	PreparePaymentChannelCreateTransaction(
		ctx context.Context,
		senderAccount, destinationAccount string,
		amount float64,
		settleDelay uint32,
		publicKey string,
		cancelAfter, destinationTag, sourceTag uint32,
		instructions *dtoRipple.Instructions,
	) (*dtoRipple.PaymentChannelCreateTxInput, string, error)

	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/paymentchannelfund
	PreparePaymentChannelFundTransaction(
		ctx context.Context,
		senderAccount, channel string,
		amount float64,
		expiration uint32,
		instructions *dtoRipple.Instructions,
	) (*dtoRipple.PaymentChannelFundTxInput, string, error)

	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/paymentchannelclaim
	PreparePaymentChannelClaimTransaction(
		ctx context.Context,
		senderAccount, channel string,
		balance string,
		amount float64,
		signature, publicKey string,
		instructions *dtoRipple.Instructions,
	) (*dtoRipple.PaymentChannelClaimTxInput, string, error)

	// NFToken operations
	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/nftokenmint
	PrepareNFTokenMintTransaction(
		ctx context.Context,
		senderAccount string,
		nfTokenTaxon uint32,
		issuer, uri string,
		transferFee uint32,
		instructions *dtoRipple.Instructions,
	) (*dtoRipple.NFTokenMintTxInput, string, error)

	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/nftokenburn
	PrepareNFTokenBurnTransaction(
		ctx context.Context,
		senderAccount, nfTokenID, owner string,
		instructions *dtoRipple.Instructions,
	) (*dtoRipple.NFTokenBurnTxInput, string, error)

	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/nftokencreateoffer
	PrepareNFTokenCreateOfferTransaction(
		ctx context.Context,
		senderAccount, nfTokenID string,
		amount float64,
		owner, destination string,
		expiration uint32,
		instructions *dtoRipple.Instructions,
	) (*dtoRipple.NFTokenCreateOfferTxInput, string, error)

	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/nftokenacceptoffer
	PrepareNFTokenAcceptOfferTransaction(
		ctx context.Context,
		senderAccount string,
		nfTokenSellOffer, nfTokenBuyOffer string,
		nfTokenBrokerFee float64,
		instructions *dtoRipple.Instructions,
	) (*dtoRipple.NFTokenAcceptOfferTxInput, string, error)

	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/nftokencanceloffer
	PrepareNFTokenCancelOfferTransaction(
		ctx context.Context,
		senderAccount string,
		nfTokenOffers []string,
		instructions *dtoRipple.Instructions,
	) (*dtoRipple.NFTokenCancelOfferTxInput, string, error)
}

// RipplePublicer defines the interface for Ripple public node operations.
// These operations query public information from the Ripple network.
type RipplePublicer interface {
	// public_account
	AccountChannels(ctx context.Context, sender, receiver string) (*dtoRipple.ResponseAccountChannels, error)
	AccountInfo(ctx context.Context, address string) (*dtoRipple.ResponseAccountInfo, error)
	// public_server_info
	ServerInfo(ctx context.Context) (*dtoRipple.ResponseServerInfo, error)
}

// RippleAdminer defines the interface for Ripple admin node operations.
// These operations typically require admin access to the Ripple node.
type RippleAdminer interface {
	// admin_keygen
	ValidationCreate(ctx context.Context, secret string) (*dtoRipple.ResponseValidationCreate, error)
	WalletProposeWithKey(
		ctx context.Context,
		seed string,
		keyType dtoRipple.XRPKeyType,
	) (*dtoRipple.ResponseWalletPropose, error)
	WalletPropose(ctx context.Context, passphrase string) (*dtoRipple.ResponseWalletPropose, error)
}
