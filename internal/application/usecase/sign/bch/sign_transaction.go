package bch

import (
	"context"
	"fmt"
	"strings"

	apibtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
	file "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	bchshared "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/shared/bch"
	signusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// signTransactionUseCase implements BCH-specific transaction signing for Sign wallet.
// Sign wallet acts as the second signer in a multisig setup.
// BCH uses Raw Transaction Hex format with SIGHASH_FORKID (0x41).
type signTransactionUseCase struct {
	bch             apibtc.BCHTransactionSigner
	accountKeyRepo  repocold.BTCAccountKeyRepositorier
	authKeyRepo     repocold.AuthAccountKeyRepositorier
	txFileRepo      file.TransactionFileRepositorier
	multisigAccount *domainAccount.MultisigConfig
	wtype           domainWallet.WalletType
	authType        domainAccount.AuthType
}

// NewBCHSignTransactionUseCase creates a new BCH SignTransactionUseCase for Sign wallet.
// BCH uses Raw TX Hex with SIGHASH_FORKID instead of PSBT.
func NewBCHSignTransactionUseCase(
	bchAPI apibtc.BCHTransactionSigner,
	accountKeyRepo repocold.BTCAccountKeyRepositorier,
	authKeyRepo repocold.AuthAccountKeyRepositorier,
	txFileRepo file.TransactionFileRepositorier,
	multisigAccount *domainAccount.MultisigConfig,
	wtype domainWallet.WalletType,
	authType domainAccount.AuthType,
) signusecase.SignTransactionUseCase {
	return &signTransactionUseCase{
		bch:             bchAPI,
		accountKeyRepo:  accountKeyRepo,
		authKeyRepo:     authKeyRepo,
		txFileRepo:      txFileRepo,
		multisigAccount: multisigAccount,
		wtype:           wtype,
		authType:        authType,
	}
}

func (u *signTransactionUseCase) Sign(
	ctx context.Context,
	input signusecase.SignTransactionInput,
) (signusecase.SignTransactionOutput, error) {
	// Get tx metadata from file name
	//
	// Workaround for BCH Node RPC bug: Accept both "unsigned" and "signed" transaction types.
	//
	// Issue: BCH Node's signrawtransactionwithkey incorrectly reports complete=true after
	// adding the 2nd signature to a 3-of-3 multisig transaction. This causes the keygen
	// wallet to mark the transaction file as "signed" when it should remain "unsigned"
	// for subsequent signers.
	//
	// Impact: Sign1 wallet receives a file marked as "signed" when adding the 2nd signature
	// in a 3-of-3 multisig. Sign2 wallet also receives "signed" files. Without this
	// workaround, both would reject the file with "invalid txType: signed" error.
	//
	// Solution: Accept both transaction types to handle BCH Node's incorrect complete flag.
	//
	// Reference: docs/chains/bch/README.md - "Known Issues and Workarounds"
	// Related: Issue #485, Issue #433
	fileType, err := u.txFileRepo.GetFileNameType(input.FilePath)
	if err != nil {
		return signusecase.SignTransactionOutput{}, fmt.Errorf("fail to parse file name: %w", err)
	}

	// Validate transaction type - accept both unsigned and signed for BCH multisig workaround
	if fileType.TxType != domainTx.TxTypeUnsigned && fileType.TxType != domainTx.TxTypeSigned {
		return signusecase.SignTransactionOutput{}, fmt.Errorf(
			"invalid txType: %s (expected unsigned or signed)", fileType.TxType)
	}

	actionType := fileType.ActionType
	txID := fileType.TxID
	signedCount := fileType.SignedCount

	// Read partially signed TX from file (BCH format)
	txContent, err := u.txFileRepo.ReadFile(input.FilePath)
	if err != nil {
		return signusecase.SignTransactionOutput{}, fmt.Errorf("fail to read tx file: %w", err)
	}

	// Parse the BCH raw tx content (hex + prevtx metadata)
	txHex, prevTxs, err := bchshared.ParseRawTxContent(txContent)
	if err != nil {
		return signusecase.SignTransactionOutput{}, fmt.Errorf("fail to parse raw tx content: %w", err)
	}

	// Sign raw transaction (add second signature for multisig)
	signedHex, isSigned, err := u.sign(txHex, prevTxs, actionType)
	if err != nil {
		return signusecase.SignTransactionOutput{}, err
	}

	// Determine output file type
	txType := domainTx.TxTypeSigned
	if !isSigned {
		txType = domainTx.TxTypeUnsigned
		signedCount++
	}

	// Write signed transaction file
	path := u.txFileRepo.CreateFilePath(actionType, txType, txID, signedCount)
	if strings.HasSuffix(path, ".psbt") {
		path = path[:len(path)-5] + ".hex"
	}

	var generatedFileName string
	// Workaround for BCH Node RPC bug in signrawtransactionwithkey.
	//
	// Issue: BCH Node incorrectly reports complete=true after adding the 2nd signature
	// to a 3-of-3 multisig transaction (should require all 3 signatures).
	//
	// Impact: When Sign1 wallet adds the 2nd signature and BCH reports complete=true,
	// the transaction file would be marked as "signed" and prevTx metadata would be
	// omitted. Without this metadata, Sign2 wallet cannot add the 3rd signature,
	// resulting in error: "Input not found or already spent"
	//
	// Solution: For multisig transactions, always include prevTx metadata regardless of
	// the complete flag. This ensures Sign2 wallet (and any subsequent signer) can add
	// their signatures.
	//
	// Detection: Use non-nil multisigAccount to determine if transaction is multisig.
	//
	// Reference: docs/chains/bch/README.md - "Known Issues and Workarounds"
	// Related: Issue #485, Issue #433
	isMultisig := u.multisigAccount != nil
	if bchshared.ShouldIncludePrevTxMetadata(isSigned, isMultisig) {
		// Include prevTx metadata for next signer (multisig or not yet fully signed)
		content := bchshared.FormatSignedTxContent(signedHex, prevTxs)
		generatedFileName, err = u.txFileRepo.WriteFile(path, content)
	} else {
		// For fully signed single-sig transactions, just write the hex
		generatedFileName, err = u.txFileRepo.WriteFile(path, signedHex)
	}
	if err != nil {
		return signusecase.SignTransactionOutput{}, fmt.Errorf("fail to write signed tx file: %w", err)
	}

	logger.Debug("signed BCH transaction (Sign wallet)",
		"action", actionType.String(),
		"txID", txID,
		"signedCount", signedCount,
		"isSigned", isSigned,
		"fileName", generatedFileName,
	)

	return signusecase.SignTransactionOutput{
		SignedData:   signedHex,
		IsComplete:   isSigned,
		NextFilePath: generatedFileName,
	}, nil
}

// sign adds the second signature to a BCH raw transaction for multisig.
// Sign wallet uses auth key from auth_account_key table.
func (u *signTransactionUseCase) sign(
	txHex string,
	prevTxs []bchshared.PrevTx,
	actionType domainTx.ActionType,
) (string, bool, error) {
	// Get account type from action
	accountType := actionType.ToAccountType()

	// Get auth key for this account
	authKey, err := u.authKeyRepo.GetByAccount(u.authType, accountType)
	if err != nil {
		return "", false, fmt.Errorf("fail to get auth key for authType %s, account %s: %w",
			u.authType, accountType, err)
	}

	logger.Debug("signing BCH transaction with auth key",
		"wallet_type", u.wtype.String(),
		"auth_type", u.authType.String(),
		"account_type", accountType.String(),
	)

	// Convert txHex to MsgTx
	msgTx, err := u.bch.ToMsgTx(txHex)
	if err != nil {
		return "", false, fmt.Errorf("fail to convert hex to MsgTx: %w", err)
	}

	// Convert prevTxs to DTO format
	dtoPrevTxs := bchshared.ConvertPrevTxsToDTO(prevTxs)

	// Sign with auth key
	wif := authKey.WalletImportFormat
	signedTx, isSigned, err := u.bch.SignRawTransactionWithKey(msgTx, []string{wif}, dtoPrevTxs)
	if err != nil {
		return "", false, fmt.Errorf("fail to sign raw transaction with auth key: %w", err)
	}

	// Convert to hex
	signedHex, err := u.bch.ToHex(signedTx)
	if err != nil {
		return "", false, fmt.Errorf("fail to convert signed tx to hex: %w", err)
	}

	logger.Debug("BCH transaction signing completed (Sign wallet)",
		"action", actionType.String(),
		"wallet_type", u.wtype.String(),
		"isSigned", isSigned,
	)

	return signedHex, isSigned, nil
}
