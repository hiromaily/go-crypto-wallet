package xrp

import (
	"context"

	"github.com/btcsuite/btcd/chaincfg"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	xrprpcpublic "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc/public"
)

// XRPPublicClient defines all operations provided by the public XRP client.
// Prefer this over the deprecated XRPer interface for new code.
type XRPPublicClient interface {
	XRPPublicer
	AccountInfoProvider
	BalanceChecker
	TransactionPreparer
	TransactionSubmitter
	LedgerPoller
	TransactionSigner
	SignerListPreparer
	CoinTypeProvider
	Closer
}

// XRPAdminClient defines all operations provided by the admin XRP client.
// Used by keygen wallets and standalone-mode ledger advancement.
type XRPAdminClient interface {
	XRPAdminer
	LedgerAdvancer
	CoinTypeProvider
	Closer
}

// XRPer defines the combined interface for XRP blockchain operations.
//
// Deprecated: XRPer is kept for backward compatibility only.
// Prefer XRPPublicClient and XRPAdminClient for new code.
type XRPer interface {
	XRPAdminer
	XRPPublicer
	TransactionSubmitter
	TransactionSigner

	GetAccountInfo(ctx context.Context, address string) (*dtoxrp.AccountInfo, error)
	GetBalance(ctx context.Context, addr string) (float64, error)
	GetTotalBalance(ctx context.Context, addrs []string) float64
	CreateRawTransaction(
		ctx context.Context,
		senderAccount, receiverAccount string,
		amount float64,
		instructions *dtoxrp.Instructions,
	) (*dtoxrp.TxInput, string, error)
	Close() error
	CoinTypeCode() domainCoin.CoinTypeCode
	GetChainConf() *chaincfg.Params
}

// SignerEntryInput represents a signer for creating SignerListSet transactions.
type SignerEntryInput struct {
	Account string
	Weight  uint32
}

// XRPPublicer defines the interface for XRP public node operations.
type XRPPublicer interface {
	AccountChannels(ctx context.Context, sender, receiver string) (*xrprpcpublic.ResponseAccountChannels, error)
	AccountInfo(ctx context.Context, address string) (*xrprpcpublic.ResponseAccountInfo, error)
	ServerInfo(ctx context.Context) (*xrprpcpublic.ResponseServerInfo, error)
}

// XRPAdminer defines the interface for XRP admin node operations.
// These operations typically require admin access to the XRP node.
type XRPAdminer interface {
	ValidationCreate(ctx context.Context, secret string) (*dtoxrp.ValidationKey, error)
	WalletProposeWithKey(
		ctx context.Context,
		seed string,
		keyType dtoxrp.XRPKeyType,
	) (*dtoxrp.WalletInfo, error)
	WalletPropose(ctx context.Context, passphrase string) (*dtoxrp.WalletInfo, error)
}

// AccountInfoProvider defines a minimal interface for querying XRP account information.
// Used by use cases that need account data without depending on the full XRPPublicClient.
type AccountInfoProvider interface {
	GetAccountInfo(ctx context.Context, address string) (*dtoxrp.AccountInfo, error)
	GetBalance(ctx context.Context, addr string) (float64, error)
}

// CoinTypeProvider provides coin type and chain configuration.
type CoinTypeProvider interface {
	CoinTypeCode() domainCoin.CoinTypeCode
	GetChainConf() *chaincfg.Params
}

// BalanceChecker provides balance query operations.
type BalanceChecker interface {
	GetBalance(ctx context.Context, addr string) (float64, error)
	GetTotalBalance(ctx context.Context, addrs []string) float64
}

// TransactionPreparer prepares raw transactions.
type TransactionPreparer interface {
	CreateRawTransaction(
		ctx context.Context,
		senderAccount, receiverAccount string,
		amount float64,
		instructions *dtoxrp.Instructions,
	) (*dtoxrp.TxInput, string, error)
}

// TransactionCombiner combines multiple signed transactions (for multisig).
type TransactionCombiner interface {
	CombineTransaction(ctx context.Context, signedTxs []string) (string, string, error)
}

// RegularKeyPreparer prepares SetRegularKey transactions.
type RegularKeyPreparer interface {
	PrepareSetRegularKeyTransaction(
		ctx context.Context,
		senderAccount, regularKey string,
		instructions *dtoxrp.Instructions,
	) (*dtoxrp.SetRegularKeyTxInput, string, error)
}

// SignerListPreparer prepares SignerListSet transactions.
type SignerListPreparer interface {
	PrepareSignerListSetTransaction(
		ctx context.Context,
		senderAccount string,
		signerQuorum uint32,
		signerEntries []SignerEntryInput,
		instructions *dtoxrp.Instructions,
	) (*dtoxrp.SignerListSetTxInput, string, error)
}

// KeyGenerator generates XRP keys/wallets.
type KeyGenerator interface {
	WalletPropose(ctx context.Context, passphrase string) (*dtoxrp.WalletInfo, error)
}

// Closer provides cleanup operations.
type Closer interface {
	Close() error
}

// TransactionSigner defines an interface for offline XRP transaction signing.
// No network calls — suitable for air-gapped keygen and sign wallets.
type TransactionSigner interface {
	// SignTransaction signs a transaction using the legacy gRPC-based method.
	//
	// Deprecated: Use SignTransactionNative for offline native Go signing.
	SignTransaction(ctx context.Context, txJSON *dtoxrp.TxInput, secret string) (string, string, error)

	// SignTransactionNative signs an XRP transaction using native Go (offline).
	// Supports single-signature and multi-signature transactions.
	// Pass existingSignedBlob for additional signatures in a multi-sig workflow.
	SignTransactionNative(
		ctx context.Context,
		txInput *dtoxrp.TxInput,
		secret string,
		isMultiSig bool,
		existingSignedBlob *string,
	) (string, string, error)
}

// LedgerPoller queries the current ledger index.
// Used by use cases that poll for ledger advancement after transaction submission.
type LedgerPoller interface {
	LedgerCurrent(ctx context.Context) (*dtoxrp.LedgerCurrent, error)
}

// LedgerAdvancer advances the ledger (admin operation, standalone mode only).
// May be nil when admin client is not configured.
type LedgerAdvancer interface {
	LedgerAccept(ctx context.Context) (*dtoxrp.LedgerAccept, error)
}

// TransactionSubmitter submits and retrieves XRP transactions.
// Ledger polling is handled separately via LedgerPoller and LedgerAdvancer.
type TransactionSubmitter interface {
	SubmitTransaction(ctx context.Context, signedTx string) (*dtoxrp.SentTx, uint64, error)
	GetTransaction(ctx context.Context, txID string, targetLedgerVersion uint64) (*dtoxrp.TxInfo, error)
}
