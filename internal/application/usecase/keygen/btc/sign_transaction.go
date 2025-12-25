package btc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/wire"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin/btc"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file"
	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/serial"
)

type signTransactionUseCase struct {
	btc             bitcoin.Bitcoiner
	accountKeyRepo  cold.AccountKeyRepositorier
	txFileRepo      file.TransactionFileRepositorier
	multisigAccount account.MultisigAccounter
}

// NewSignTransactionUseCase creates a new SignTransactionUseCase for BTC keygen
func NewSignTransactionUseCase(
	btc bitcoin.Bitcoiner,
	accountKeyRepo cold.AccountKeyRepositorier,
	txFileRepo file.TransactionFileRepositorier,
	multisigAccount account.MultisigAccounter,
) keygenusecase.SignTransactionUseCase {
	return &signTransactionUseCase{
		btc:             btc,
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
	//  if payment_5_unsigned_0_1534466246366489473, 5 is target
	actionType, _, txID, signedCount, err := u.txFileRepo.ValidateFilePath(input.FilePath, domainTx.TxTypeUnsigned)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, err
	}

	// Get hex tx from file
	data, err := u.txFileRepo.ReadFile(input.FilePath)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, err
	}

	var hex, encodedPrevsAddrs string
	tmp := strings.Split(data, ",")
	// file: hex, prev_address
	hex = tmp[0]
	if len(tmp) > 1 {
		encodedPrevsAddrs = tmp[1]
	}
	if encodedPrevsAddrs == "" {
		// It's required data since Bitcoin core ver17
		return keygenusecase.SignTransactionOutput{}, errors.New("encodedPrevsAddrs must be set in csv file")
	}

	// Sign
	hexTx, isSigned, newEncodedPrevsAddrs, err := u.sign(hex, encodedPrevsAddrs)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, err
	}

	// hexTx for save data as file
	saveData := hexTx

	// If sign is not finished because of multisig, signedCount should be increment
	txType := domainTx.TxTypeSigned
	if !isSigned {
		txType = domainTx.TxTypeUnsigned
		signedCount++
		if newEncodedPrevsAddrs != "" {
			saveData = fmt.Sprintf("%s,%s", saveData, newEncodedPrevsAddrs)
		}
	}

	// Write file
	path := u.txFileRepo.CreateFilePath(actionType, txType, txID, signedCount)
	generatedFileName, err := u.txFileRepo.WriteFile(path, saveData)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, err
	}

	return keygenusecase.SignTransactionOutput{
		FilePath:      generatedFileName,
		IsDone:        isSigned,
		SignedCount:   1, // BTC signs one transaction at a time
		UnsignedCount: 0, // BTC doesn't track unsigned separately in this interface
	}, nil
}

// sign
// - coin is sent from sender account to receiver account. Sender's privKey(sender account) is required
// - [actionType:deposit]  [from] client [to] deposit, (not multisig addr)
// - [actionType:payment]  [from] payment [to] unknown, (multisig addr)
// - [actionType:transfer] [from] from [to] to, (multisig addr)
func (u *signTransactionUseCase) sign(hex, encodedPrevsAddrs string) (string, bool, string, error) {
	// Get tx from hex
	msgTx, err := u.btc.ToMsgTx(hex)
	if err != nil {
		return "", false, "", err
	}

	// Decode encodedPrevsAddrs string to btc.PreviousTxs struct
	var prevsAddrs btc.PreviousTxs
	if err = serial.DecodeFromString(encodedPrevsAddrs, &prevsAddrs); err != nil {
		return "", false, "", err
	}

	// Single signature address or multisig
	var (
		signedTx             *wire.MsgTx
		isSigned             bool
		newEncodedPrevsAddrs string
	)

	// Sign
	if !u.multisigAccount.IsMultisigAccount(prevsAddrs.SenderAccount) {
		signedTx, isSigned, err = u.btc.SignRawTransaction(msgTx, prevsAddrs.PrevTxs)
	} else {
		signedTx, isSigned, newEncodedPrevsAddrs, err = u.signMultisig(msgTx, &prevsAddrs)
	}

	if newEncodedPrevsAddrs == "" {
		newEncodedPrevsAddrs = encodedPrevsAddrs
	}

	// After sign
	if err != nil {
		return "", false, "", err
	}

	hexTx, err := u.btc.ToHex(signedTx)
	if err != nil {
		return "", false, "", fmt.Errorf("fail to call btc.ToHex(signedTx): %w", err)
	}
	logger.Debug(
		"call btc.SignRawTransaction()",
		"hexTx", hexTx,
		"isSigned", isSigned)

	return hexTx, isSigned, newEncodedPrevsAddrs, nil
}

func (u *signTransactionUseCase) signMultisig(
	msgTx *wire.MsgTx,
	prevsAddrs *btc.PreviousTxs,
) (*wire.MsgTx, bool, string, error) {
	var newEncodedPrevsAddrs string

	// Get WIPs and RedeemScript for keygen wallet
	accountKeys, err := u.accountKeyRepo.GetAllMultiAddr(prevsAddrs.SenderAccount, prevsAddrs.Addrs)
	if err != nil {
		return nil, false, "", fmt.Errorf("fail to call accountKeyRepo.GetAllMultiAddr(): %w", err)
	}

	// Retrieve WIPs - pre-allocate slice
	wips := make([]string, 0, len(accountKeys))
	for _, val := range accountKeys {
		wips = append(wips, val.WalletImportFormat)
	}

	// Mapping redeemScript to PrevTxs
	for idx, val := range prevsAddrs.Addrs {
		rs := cold.GetRedeemScriptByAddress(accountKeys, val)
		if rs == "" {
			logger.Error("redeemScript can not be found")
			continue
		}
		prevsAddrs.PrevTxs[idx].RedeemScript = rs
	}

	// Serialize prevsAddrs with redeemScript
	newEncodedPrevsAddrs, err = serial.EncodeToString(prevsAddrs)
	if err != nil {
		return nil, false, "", fmt.Errorf("fail to call serial.EncodeToString(): %w", err)
	}

	// Sign
	signedTx, isSigned, err := u.btc.SignRawTransactionWithKey(msgTx, wips, prevsAddrs.PrevTxs)
	return signedTx, isSigned, newEncodedPrevsAddrs, err
}
