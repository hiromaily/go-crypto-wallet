package bch

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcutil"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	apibtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
	file "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	usecaseshared "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/shared"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// signTransactionUseCase implements BCH-specific transaction signing.
// Unlike BTC, BCH uses Raw Transaction Hex format with SIGHASH_FORKID (0x41).
type signTransactionUseCase struct {
	bch             apibtc.BCHer
	accountKeyRepo  repocold.BTCAccountKeyRepositorier
	txFileRepo      file.TransactionFileRepositorier
	multisigAccount *domainAccount.MultisigConfig
}

// NewSignTransactionUseCase creates a new BCH SignTransactionUseCase.
// BCH uses Raw TX Hex with SIGHASH_FORKID instead of PSBT.
func NewSignTransactionUseCase(
	bch apibtc.BCHer,
	accountKeyRepo repocold.BTCAccountKeyRepositorier,
	txFileRepo file.TransactionFileRepositorier,
	multisigAccount *domainAccount.MultisigConfig,
) keygenusecase.SignTransactionUseCase {
	return &signTransactionUseCase{
		bch:             bch,
		accountKeyRepo:  accountKeyRepo,
		txFileRepo:      txFileRepo,
		multisigAccount: multisigAccount,
	}
}

func (u *signTransactionUseCase) Sign(
	ctx context.Context,
	input keygenusecase.SignTransactionInput,
) (keygenusecase.SignTransactionOutput, error) {
	// Get tx_deposit_id from tx file name
	// BCH uses .hex extension, so validate accordingly
	actionType, _, txID, signedCount, err := u.txFileRepo.ValidateFilePath(input.FilePath, domainTx.TxTypeUnsigned)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, err
	}

	// Read Raw TX Hex from file (BCH format)
	txContent, err := u.txFileRepo.ReadFile(input.FilePath)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, fmt.Errorf("fail to read tx file: %w", err)
	}

	// Parse the BCH raw tx content (hex + prevtx metadata)
	txHex, prevTxs, err := u.parseRawTxContent(txContent)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, fmt.Errorf("fail to parse raw tx content: %w", err)
	}

	// Sign raw transaction using SignRawTransactionWithKey
	// BCH signing uses SIGHASH_FORKID (0x41) internally
	signedHex, isSigned, err := u.sign(txHex, prevTxs, actionType)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, err
	}

	// Determine output file type based on signature completion
	txType := domainTx.TxTypeSigned
	if !isSigned {
		txType = domainTx.TxTypeUnsigned
		signedCount++
	}

	// Write signed transaction file
	path := u.txFileRepo.CreateFilePath(actionType, txType, txID, signedCount)
	// Use .hex extension for BCH
	if strings.HasSuffix(path, ".psbt") {
		path = path[:len(path)-5] + ".hex"
	}

	var generatedFileName string
	if isSigned {
		// For fully signed transactions, just write the hex
		generatedFileName, err = u.txFileRepo.WriteFile(path, signedHex)
	} else {
		// For partially signed, include prevTx metadata for next signer
		content := u.formatSignedTxContent(signedHex, prevTxs)
		generatedFileName, err = u.txFileRepo.WriteFile(path, content)
	}
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, fmt.Errorf("fail to write signed tx file: %w", err)
	}

	logger.Debug("signed BCH transaction",
		"action", actionType.String(),
		"txID", txID,
		"signedCount", signedCount,
		"isSigned", isSigned,
		"fileName", generatedFileName,
	)

	return keygenusecase.SignTransactionOutput{
		FilePath:      generatedFileName,
		IsDone:        isSigned,
		SignedCount:   1,
		UnsignedCount: 0,
	}, nil
}

// sign signs a BCH raw transaction using SignRawTransactionWithKey.
// BCH uses SIGHASH_FORKID (0x41) which is handled internally by the BCH node.
func (u *signTransactionUseCase) sign(
	txHex string,
	prevTxs []signPrevTx,
	actionType domainTx.ActionType,
) (string, bool, error) {
	// Infer sender account from action type
	senderAccount, err := usecaseshared.InferSenderAccount(actionType)
	if err != nil {
		return "", false, err
	}

	// Get WIF keys for this account
	wifs, err := u.getWIFsForAccount(senderAccount)
	if err != nil {
		return "", false, err
	}
	if len(wifs) == 0 {
		return "", false, fmt.Errorf("no WIF keys found for account: %s", senderAccount)
	}

	// Convert txHex to MsgTx
	msgTx, err := u.bch.ToMsgTx(txHex)
	if err != nil {
		return "", false, fmt.Errorf("fail to convert hex to MsgTx: %w", err)
	}

	// Convert prevTxs to DTO format
	dtoPrevTxs := u.convertPrevTxsToDTO(prevTxs)

	// Sign raw transaction with WIF keys
	// BCH nodes handle SIGHASH_FORKID internally
	signedTx, isSigned, err := u.bch.SignRawTransactionWithKey(msgTx, wifs, dtoPrevTxs)
	if err != nil {
		return "", false, fmt.Errorf("fail to sign raw transaction: %w", err)
	}

	// Convert signed transaction back to hex
	signedHex, err := u.bch.ToHex(signedTx)
	if err != nil {
		return "", false, fmt.Errorf("fail to convert signed tx to hex: %w", err)
	}

	logger.Debug("BCH transaction signing completed",
		"action", actionType.String(),
		"sender_account", senderAccount.String(),
		"isSigned", isSigned,
	)

	return signedHex, isSigned, nil
}

// getWIFsForAccount gets WIF private keys for the specified account
func (u *signTransactionUseCase) getWIFsForAccount(
	senderAccount domainAccount.AccountType,
) ([]string, error) {
	// Get keys with various statuses
	statuses := []domainAddress.AddrStatus{
		domainAddress.AddrStatusHDKeyGenerated,
		domainAddress.AddrStatusPrivKeyImported,
	}

	wifs := make([]string, 0)
	for _, status := range statuses {
		accountKeys, err := u.accountKeyRepo.GetAllAddrStatus(senderAccount, status)
		if err != nil {
			return nil, fmt.Errorf("fail to get account keys for %s: %w", senderAccount.String(), err)
		}
		for _, key := range accountKeys {
			if key.WalletImportFormat != "" {
				wifs = append(wifs, key.WalletImportFormat)
			}
		}
	}

	return wifs, nil
}

// signPrevTx represents previous transaction data for BCH signing
type signPrevTx struct {
	TxID         string
	Vout         uint32
	ScriptPubKey string
	RedeemScript string
	Amount       int64
}

// parseRawTxContent parses BCH raw transaction content
// Format: txHex on first line, followed by prevtx metadata lines
func (*signTransactionUseCase) parseRawTxContent(content string) (string, []signPrevTx, error) {
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) == 0 {
		return "", nil, errors.New("empty transaction content")
	}

	txHex := lines[0]
	prevTxs := make([]signPrevTx, 0, len(lines)-1)

	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" || !strings.HasPrefix(line, "prevtx:") {
			continue
		}

		// Parse prevtx:txid:vout:scriptPubKey:redeemScript:amount
		parts := strings.Split(strings.TrimPrefix(line, "prevtx:"), ":")
		if len(parts) < 5 {
			logger.Warn("invalid prevtx format, skipping", "line", line)
			continue
		}

		vout64, err := strconv.ParseUint(parts[1], 10, 32)
		if err != nil {
			logger.Warn("invalid vout in prevtx", "vout", parts[1], "error", err)
			continue
		}
		vout := uint32(vout64)

		amount, err := strconv.ParseInt(parts[4], 10, 64)
		if err != nil {
			logger.Warn("invalid amount in prevtx", "amount", parts[4], "error", err)
			continue
		}

		prevTxs = append(prevTxs, signPrevTx{
			TxID:         parts[0],
			Vout:         vout,
			ScriptPubKey: parts[2],
			RedeemScript: parts[3],
			Amount:       amount,
		})
	}

	return txHex, prevTxs, nil
}

// convertPrevTxsToDTO converts signPrevTx to DTO format for signing
func (*signTransactionUseCase) convertPrevTxsToDTO(prevTxs []signPrevTx) []dtobtc.PreviousTx {
	result := make([]dtobtc.PreviousTx, len(prevTxs))
	for i, prev := range prevTxs {
		result[i] = dtobtc.PreviousTx{
			TxID:         prev.TxID,
			Vout:         prev.Vout,
			ScriptPubKey: prev.ScriptPubKey,
			RedeemScript: prev.RedeemScript,
			Amount:       btcutil.Amount(prev.Amount),
		}
	}
	return result
}

// formatSignedTxContent formats the signed transaction content for BCH
func (*signTransactionUseCase) formatSignedTxContent(txHex string, prevTxs []signPrevTx) string {
	var builder strings.Builder
	builder.WriteString(txHex)
	builder.WriteString("\n")

	for _, prev := range prevTxs {
		fmt.Fprintf(&builder, "prevtx:%s:%d:%s:%s:%d\n",
			prev.TxID,
			prev.Vout,
			prev.ScriptPubKey,
			prev.RedeemScript,
			prev.Amount,
		)
	}
	return builder.String()
}
