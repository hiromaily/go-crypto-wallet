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
	signusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// signTxBCHClient defines the minimal interface needed for BCH sign wallet transaction signing.
// This follows the Interface Segregation Principle - depend only on methods actually used.
// BCH uses Raw TX Hex with SIGHASH_FORKID instead of PSBT.
type signTxBCHClient interface {
	apibtc.RawTransactionCreator // ToMsgTx, ToHex
	apibtc.RawTransactionSigner  // SignRawTransactionWithKey
}

// signTransactionUseCase implements BCH-specific transaction signing for Sign wallet.
// Sign wallet acts as the second signer in a multisig setup.
// BCH uses Raw Transaction Hex format with SIGHASH_FORKID (0x41).
type signTransactionUseCase struct {
	bch             signTxBCHClient
	accountKeyRepo  repocold.BTCAccountKeyRepositorier
	authKeyRepo     repocold.AuthAccountKeyRepositorier
	txFileRepo      file.TransactionFileRepositorier
	multisigAccount *domainAccount.MultisigConfig
	wtype           domainWallet.WalletType
	authType        domainAccount.AuthType
}

// NewSignTransactionUseCase creates a new BCH SignTransactionUseCase for Sign wallet.
// BCH uses Raw TX Hex with SIGHASH_FORKID instead of PSBT.
func NewSignTransactionUseCase(
	bchAPI apibtc.BCHer,
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
	// Get tx_deposit_id from tx file name
	actionType, _, txID, signedCount, err := u.txFileRepo.ValidateFilePath(input.FilePath, domainTx.TxTypeUnsigned)
	if err != nil {
		return signusecase.SignTransactionOutput{}, err
	}

	// Read partially signed TX from file (BCH format)
	txContent, err := u.txFileRepo.ReadFile(input.FilePath)
	if err != nil {
		return signusecase.SignTransactionOutput{}, fmt.Errorf("fail to read tx file: %w", err)
	}

	// Parse the BCH raw tx content (hex + prevtx metadata)
	txHex, prevTxs, err := u.parseRawTxContent(txContent)
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
	if isSigned {
		generatedFileName, err = u.txFileRepo.WriteFile(path, signedHex)
	} else {
		content := u.formatSignedTxContent(signedHex, prevTxs)
		generatedFileName, err = u.txFileRepo.WriteFile(path, content)
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
	prevTxs []signPrevTx,
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
	dtoPrevTxs := u.convertPrevTxsToDTO(prevTxs)

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

// signPrevTx represents previous transaction data for BCH signing
type signPrevTx struct {
	TxID         string
	Vout         uint32
	ScriptPubKey string
	RedeemScript string
	Amount       int64
}

// parseRawTxContent parses BCH raw transaction content
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

// formatSignedTxContent formats the signed transaction content
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
