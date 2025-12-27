package btc

import (
	"context"
	"fmt"

	signusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/config/account"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type signTransactionUseCase struct {
	btc             bitcoin.Bitcoiner
	accountKeyRepo  cold.AccountKeyRepositorier
	authKeyRepo     cold.AuthAccountKeyRepositorier
	txFileRepo      file.TransactionFileRepositorier
	multisigAccount account.MultisigAccounter
	wtype           domainWallet.WalletType
}

// NewSignTransactionUseCase creates a new SignTransactionUseCase for sign wallet
func NewSignTransactionUseCase(
	btcAPI bitcoin.Bitcoiner,
	accountKeyRepo cold.AccountKeyRepositorier,
	authKeyRepo cold.AuthAccountKeyRepositorier,
	txFileRepo file.TransactionFileRepositorier,
	multisigAccount account.MultisigAccounter,
	wtype domainWallet.WalletType,
) signusecase.SignTransactionUseCase {
	return &signTransactionUseCase{
		btc:             btcAPI,
		accountKeyRepo:  accountKeyRepo,
		authKeyRepo:     authKeyRepo,
		txFileRepo:      txFileRepo,
		multisigAccount: multisigAccount,
		wtype:           wtype,
	}
}

func (u *signTransactionUseCase) Sign(
	ctx context.Context,
	input signusecase.SignTransactionInput,
) (signusecase.SignTransactionOutput, error) {
	// Get tx_deposit_id from tx file name
	//  if payment_5_unsigned_1_1534466246366489473.psbt, 5 is target
	actionType, _, txID, signedCount, err := u.txFileRepo.ValidateFilePath(input.FilePath, domainTx.TxTypeUnsigned)
	if err != nil {
		return signusecase.SignTransactionOutput{}, err
	}

	// Read PSBT from file
	psbtBase64, err := u.txFileRepo.ReadPSBTFile(input.FilePath)
	if err != nil {
		return signusecase.SignTransactionOutput{}, fmt.Errorf("fail to read PSBT file: %w", err)
	}

	// Sign PSBT (add second signature for multisig)
	signedPSBT, isSigned, err := u.sign(psbtBase64, actionType)
	if err != nil {
		return signusecase.SignTransactionOutput{}, err
	}

	// If sign is not finished because of multisig, signedCount should be increment
	txType := domainTx.TxTypeSigned
	if !isSigned {
		txType = domainTx.TxTypeUnsigned
		signedCount++ // Increment for additional signatures needed
	}

	// Write signed PSBT file
	path := u.txFileRepo.CreateFilePath(actionType, txType, txID, signedCount)
	generatedFileName, err := u.txFileRepo.WritePSBTFile(path, signedPSBT)
	if err != nil {
		return signusecase.SignTransactionOutput{}, fmt.Errorf("fail to write signed PSBT file: %w", err)
	}

	logger.Debug("signed PSBT",
		"action", actionType.String(),
		"txID", txID,
		"signedCount", signedCount,
		"isSigned", isSigned,
		"fileName", generatedFileName,
	)

	return signusecase.SignTransactionOutput{
		SignedHex:    signedPSBT, // PSBT base64 (not hex anymore)
		IsComplete:   isSigned,
		NextFilePath: generatedFileName,
	}, nil
}

// sign signs a PSBT using offline signing with btcd library (no Bitcoin Core RPC).
// This function supports all Bitcoin address types including Taproot:
//   - Legacy (P2PKH): ECDSA signature
//   - SegWit (P2WPKH, P2SH-SegWit): ECDSA signature with witness data
//   - Taproot (P2TR): Schnorr signature (BIP340) with witness data
//
// The signature algorithm is automatically selected based on the PSBT input's scriptPubKey type.
// For Taproot addresses, Schnorr signatures (BIP340) are used.
//
// The Sign wallet acts as the second (or subsequent) signer in a multisig setup, using keys from
// the auth_account_key table. This adds signatures to a partially signed PSBT from Keygen wallet.
//
// Transaction flow:
//   - [actionType:deposit]  [from] client [to] deposit (not multisig addr)
//   - [actionType:payment]  [from] payment [to] unknown (multisig addr)
//   - [actionType:transfer] [from] account [to] account (multisig addr)
//
// Note: This operates OFFLINE - no Bitcoin Core RPC required.
func (u *signTransactionUseCase) sign(
	psbtBase64 string,
	actionType domainTx.ActionType,
) (string, bool, error) {
	// Sign wallet always signs multisig transactions
	// Add second signature to PSBT using auth key
	signedPSBT, isSigned, err := u.signMultisigPSBT(psbtBase64)
	if err != nil {
		return "", false, err
	}

	logger.Debug("PSBT signing completed",
		"action", actionType.String(),
		"wallet_type", u.wtype.String(),
		"isSigned", isSigned,
	)

	return signedPSBT, isSigned, nil
}

// signMultisigPSBT adds the second signature to a partially signed PSBT for multisig transactions.
// The Sign wallet uses the auth key from auth_account_key table to add its signature.
// This function works offline using btcd library, supporting all address types:
//   - P2SH (Legacy multisig): ECDSA signature
//   - P2SH-SegWit: ECDSA signature with witness
//   - P2WSH (Native SegWit multisig): ECDSA signature with witness
//   - P2TR (Taproot): Schnorr signature (BIP340) or script path
//
// For 2-of-2 multisig, this signature typically completes the transaction.
// For 2-of-N multisig (N>2), the transaction is complete once 2 signatures are present.
func (u *signTransactionUseCase) signMultisigPSBT(
	psbtBase64 string,
) (string, bool, error) {
	// Get auth key from auth_account_key table (Sign wallet's key)
	authKey, err := u.authKeyRepo.GetOne("")
	if err != nil {
		return "", false, fmt.Errorf("fail to get auth key: %w", err)
	}

	logger.Debug("signing PSBT with auth key",
		"wallet_type", u.wtype.String(),
	)

	// Sign PSBT with Sign wallet's private key (offline, using btcd)
	// This adds the second signature to the partially signed PSBT
	signedPSBT, isSigned, err := u.btc.SignPSBTWithKey(psbtBase64, []string{authKey.WalletImportFormat})
	if err != nil {
		return "", false, fmt.Errorf("fail to sign PSBT with auth key: %w", err)
	}

	return signedPSBT, isSigned, nil
}
