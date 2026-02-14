// Package xrp defines interfaces for XRP blockchain operations.
//
// This package follows the Dependency Inversion Principle of Clean Architecture
// by defining interfaces in the application layer that are implemented by the
// infrastructure layer.
package xrp

import (
	"context"

	"github.com/btcsuite/btcd/chaincfg"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
)

// XRPer defines the main interface for XRP blockchain operations.
// It embeds specialized interfaces for admin, public, and API operations.
type XRPer interface {
	XRPAdminer
	XRPPublicer
	XRPAPIer

	// balance
	GetBalance(ctx context.Context, addr string) (float64, error)
	GetTotalBalance(ctx context.Context, addrs []string) float64

	// transaction
	CreateRawTransaction(
		ctx context.Context,
		senderAccount, receiverAccount string,
		amount float64,
		instructions *dtoxrp.Instructions,
	) (*dtoxrp.TxInput, string, error)

	// xrp
	Close() error
	CoinTypeCode() domainCoin.CoinTypeCode
	GetChainConf() *chaincfg.Params
}

// SignerEntryInput represents a signer for creating SignerListSet transactions
type SignerEntryInput struct {
	Account string
	Weight  uint32
}

// XRPAPIer defines the interface for XRP API operations.
// Implementations handle account management, address generation, and transaction operations.
type XRPAPIer interface {
	// XRPAccountAPI
	GetAccountInfo(ctx context.Context, address string) (*dtoxrp.ResponseGetAccountInfo, error)
	// XRPAddressAPI
	GenerateAddress(ctx context.Context) (*dtoxrp.ResponseGenerateAddress, error)
	GenerateXAddress(ctx context.Context) (*dtoxrp.ResponseGenerateXAddress, error)
	IsValidAddress(ctx context.Context, addr string) (bool, error)
	// XRPTxAPI
	PrepareTransaction(
		ctx context.Context,
		senderAccount, receiverAccount string,
		amount float64,
		instructions *dtoxrp.Instructions,
	) (*dtoxrp.TxInput, string, error)
	SignTransaction(ctx context.Context, txJSON *dtoxrp.TxInput, secret string) (string, string, error)
	SignTransactionNative(
		ctx context.Context, txInput *dtoxrp.TxInput, secret string, isMultiSig bool,
	) (string, string, error)
	CombineTransaction(ctx context.Context, signedTxs []string) (string, string, error)
	SubmitTransaction(ctx context.Context, signedTx string) (*dtoxrp.SentTx, uint64, error)
	WaitValidation(ctx context.Context, targetledgerVarsion uint64) (uint64, error)
	GetTransaction(ctx context.Context, txID string, targetLedgerVersion uint64) (*dtoxrp.TxInfo, error)

	// Regular Key operations
	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/setregularkey
	PrepareSetRegularKeyTransaction(
		ctx context.Context,
		senderAccount, regularKey string,
		instructions *dtoxrp.Instructions,
	) (*dtoxrp.SetRegularKeyTxInput, string, error)

	// Account settings operations
	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/accountset
	PrepareAccountSetTransaction(
		ctx context.Context,
		senderAccount string,
		setFlag, clearFlag uint32,
		instructions *dtoxrp.Instructions,
	) (*dtoxrp.AccountSetTxInput, string, error)

	// Multi-signature operations
	// Reference: https://xrpl.org/docs/concepts/accounts/multi-signing
	PrepareSignerListSetTransaction(
		ctx context.Context,
		senderAccount string,
		signerQuorum uint32,
		signerEntries []SignerEntryInput,
		instructions *dtoxrp.Instructions,
	) (*dtoxrp.SignerListSetTxInput, string, error)

	// Trust line operations
	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/trustset
	PrepareTrustSetTransaction(
		ctx context.Context,
		senderAccount string,
		limitAmount *dtoxrp.IssuedCurrencyAmount,
		qualityIn, qualityOut uint32,
		instructions *dtoxrp.Instructions,
	) (*dtoxrp.TrustSetTxInput, string, error)

	// Escrow operations
	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/escrowcreate
	PrepareEscrowCreateTransaction(
		ctx context.Context,
		senderAccount, destinationAccount string,
		amount float64,
		cancelAfter, finishAfter uint32,
		condition string,
		destinationTag uint32,
		instructions *dtoxrp.Instructions,
	) (*dtoxrp.EscrowCreateTxInput, string, error)

	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/escrowfinish
	PrepareEscrowFinishTransaction(
		ctx context.Context,
		senderAccount, owner string,
		offerSequence uint32,
		condition, fulfillment string,
		instructions *dtoxrp.Instructions,
	) (*dtoxrp.EscrowFinishTxInput, string, error)

	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/escrowcancel
	PrepareEscrowCancelTransaction(
		ctx context.Context,
		senderAccount, owner string,
		offerSequence uint32,
		instructions *dtoxrp.Instructions,
	) (*dtoxrp.EscrowCancelTxInput, string, error)

	// PaymentChannel operations
	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/paymentchannelcreate
	PreparePaymentChannelCreateTransaction(
		ctx context.Context,
		senderAccount, destinationAccount string,
		amount float64,
		settleDelay uint32,
		publicKey string,
		cancelAfter, destinationTag, sourceTag uint32,
		instructions *dtoxrp.Instructions,
	) (*dtoxrp.PaymentChannelCreateTxInput, string, error)

	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/paymentchannelfund
	PreparePaymentChannelFundTransaction(
		ctx context.Context,
		senderAccount, channel string,
		amount float64,
		expiration uint32,
		instructions *dtoxrp.Instructions,
	) (*dtoxrp.PaymentChannelFundTxInput, string, error)

	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/paymentchannelclaim
	PreparePaymentChannelClaimTransaction(
		ctx context.Context,
		senderAccount, channel string,
		balance string,
		amount float64,
		signature, publicKey string,
		instructions *dtoxrp.Instructions,
	) (*dtoxrp.PaymentChannelClaimTxInput, string, error)

	// NFToken operations
	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/nftokenmint
	PrepareNFTokenMintTransaction(
		ctx context.Context,
		senderAccount string,
		nfTokenTaxon uint32,
		issuer, uri string,
		transferFee uint32,
		instructions *dtoxrp.Instructions,
	) (*dtoxrp.NFTokenMintTxInput, string, error)

	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/nftokenburn
	PrepareNFTokenBurnTransaction(
		ctx context.Context,
		senderAccount, nfTokenID, owner string,
		instructions *dtoxrp.Instructions,
	) (*dtoxrp.NFTokenBurnTxInput, string, error)

	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/nftokencreateoffer
	PrepareNFTokenCreateOfferTransaction(
		ctx context.Context,
		senderAccount, nfTokenID string,
		amount float64,
		owner, destination string,
		expiration uint32,
		instructions *dtoxrp.Instructions,
	) (*dtoxrp.NFTokenCreateOfferTxInput, string, error)

	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/nftokenacceptoffer
	PrepareNFTokenAcceptOfferTransaction(
		ctx context.Context,
		senderAccount string,
		nfTokenSellOffer, nfTokenBuyOffer string,
		nfTokenBrokerFee float64,
		instructions *dtoxrp.Instructions,
	) (*dtoxrp.NFTokenAcceptOfferTxInput, string, error)

	// Reference: https://xrpl.org/docs/references/protocol/transactions/types/nftokencanceloffer
	PrepareNFTokenCancelOfferTransaction(
		ctx context.Context,
		senderAccount string,
		nfTokenOffers []string,
		instructions *dtoxrp.Instructions,
	) (*dtoxrp.NFTokenCancelOfferTxInput, string, error)
}

// XRPPublicer defines the interface for XRP public node operations.
// These operations query public information from the XRP network.
type XRPPublicer interface {
	// public_account
	AccountChannels(ctx context.Context, sender, receiver string) (*dtoxrp.ResponseAccountChannels, error)
	AccountInfo(ctx context.Context, address string) (*dtoxrp.ResponseAccountInfo, error)
	// public_server_info
	ServerInfo(ctx context.Context) (*dtoxrp.ResponseServerInfo, error)
}

// XRPAdminer defines the interface for XRP admin node operations.
// These operations typically require admin access to the XRP node.
type XRPAdminer interface {
	// admin_keygen
	ValidationCreate(ctx context.Context, secret string) (*dtoxrp.ResponseValidationCreate, error)
	WalletProposeWithKey(
		ctx context.Context,
		seed string,
		keyType dtoxrp.XRPKeyType,
	) (*dtoxrp.ResponseWalletPropose, error)
	WalletPropose(ctx context.Context, passphrase string) (*dtoxrp.ResponseWalletPropose, error)
}
