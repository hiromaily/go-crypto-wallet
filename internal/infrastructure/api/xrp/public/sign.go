package xrppublic

import (
	"context"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	xrpsigner "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/xrp/signer"
)

// SignTransaction signs a payment transaction offline using the native Go implementation.
func (*PublicXRP) SignTransaction(
	ctx context.Context, txInput *dtoxrp.TxInput, secret string,
) (string, string, error) {
	s := xrpsigner.NewPeersystSigner()
	return s.SignTransactionNative(ctx, txInput, secret, false, nil)
}

// SignTransactionNative signs a transaction using the native Go implementation (Peersyst/xrpl-go).
// Supports both single-signature and multi-signature transactions.
func (*PublicXRP) SignTransactionNative(
	ctx context.Context,
	txInput *dtoxrp.TxInput,
	secret string,
	isMultiSig bool,
	existingSignedBlob *string,
) (string, string, error) {
	s := xrpsigner.NewPeersystSigner()
	return s.SignTransactionNative(ctx, txInput, secret, isMultiSig, existingSignedBlob)
}
