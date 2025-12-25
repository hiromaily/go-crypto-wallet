package btc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/wire"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	signusecase "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/sign"
	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/bitcoin"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/bitcoin/btc"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/repository/cold"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/storage/file"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/serial"
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
	// get tx_deposit_id from tx file name
	//  if payment_5_unsigned_0_1534466246366489473, 5 is target
	actionType, _, txID, signedCount, err := u.txFileRepo.ValidateFilePath(input.FilePath, domainTx.TxTypeUnsigned)
	if err != nil {
		return signusecase.SignTransactionOutput{}, err
	}

	// get hex tx from file
	data, err := u.txFileRepo.ReadFile(input.FilePath)
	if err != nil {
		return signusecase.SignTransactionOutput{}, err
	}

	var hex, encodedPrevsAddrs string
	tmp := strings.Split(data, ",")
	// file: hex, prev_address
	hex = tmp[0]
	if len(tmp) > 1 {
		encodedPrevsAddrs = tmp[1]
	}
	if encodedPrevsAddrs == "" {
		// it's required data since Bitcoin core ver17
		return signusecase.SignTransactionOutput{}, errors.New("encodedPrevsAddrs must be set in csv file")
	}

	// sing
	hexTx, isSigned, newEncodedPrevsAddrs, err := u.sign(hex, encodedPrevsAddrs)
	if err != nil {
		return signusecase.SignTransactionOutput{}, err
	}

	// hexTx for save data as file
	saveData := hexTx

	// if sign is not finished because of multisig, signedCount should be increment
	txType := domainTx.TxTypeSigned
	if !isSigned {
		txType = domainTx.TxTypeUnsigned
		signedCount++
		if newEncodedPrevsAddrs != "" {
			saveData = fmt.Sprintf("%s,%s", saveData, newEncodedPrevsAddrs)
		}
	}

	// write file
	path := u.txFileRepo.CreateFilePath(actionType, txType, txID, signedCount)
	generatedFileName, err := u.txFileRepo.WriteFile(path, saveData)
	if err != nil {
		return signusecase.SignTransactionOutput{}, err
	}

	return signusecase.SignTransactionOutput{
		SignedHex:    hexTx,
		IsComplete:   isSigned,
		NextFilePath: generatedFileName,
	}, nil
}

// sign
// - coin is sent from sender account to receiver account. Sender's privKey(sender account) is required
// - [actionType:deposit]  [from] client [to] deposit, (not multisig addr)
// - [actionType:payment]  [from] payment [to] unknown, (multisig addr)
// - [actionType:transfer] [from] from [to] to, (multisig addr)
func (u *signTransactionUseCase) sign(hex, encodedPrevsAddrs string) (string, bool, string, error) {
	// get tx from hex
	msgTx, err := u.btc.ToMsgTx(hex)
	if err != nil {
		return "", false, "", err
	}

	// decode encodedPrevsAddrs string to btc.AddrsPrevTxs struct
	var prevsAddrs btc.PreviousTxs
	if err = serial.DecodeFromString(encodedPrevsAddrs, &prevsAddrs); err != nil {
		return "", false, "", err
	}

	// single signature address
	var (
		signedTx             *wire.MsgTx
		isSigned             bool
		newEncodedPrevsAddrs string
	)

	// sign
	if !u.multisigAccount.IsMultisigAccount(prevsAddrs.SenderAccount) {
		signedTx, isSigned, err = u.btc.SignRawTransaction(msgTx, prevsAddrs.PrevTxs)
	} else {
		signedTx, isSigned, newEncodedPrevsAddrs, err = u.signMultisig(msgTx, &prevsAddrs)
	}

	if newEncodedPrevsAddrs == "" {
		newEncodedPrevsAddrs = encodedPrevsAddrs
	}

	// after sign
	if err != nil {
		return "", false, "", err
	}

	hexTx, err := u.btc.ToHex(signedTx)
	if err != nil {
		return "", false, "", fmt.Errorf("fail to call u.btc.ToHex(signedTx): %w", err)
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
	var wips []string
	var newEncodedPrevsAddrs string

	// get WIPs, RedeedScript
	switch u.wtype {
	case domainWallet.WalletTypeKeyGen:
		accountKeys, err := u.accountKeyRepo.GetAllMultiAddr(prevsAddrs.SenderAccount, prevsAddrs.Addrs)
		if err != nil {
			return nil, false, "", fmt.Errorf("fail to call accountKeyRepo.GetAllMultiAddr(): %w", err)
		}

		// retrieve WIPs
		for _, val := range accountKeys {
			wips = append(wips, val.WalletImportFormat)
		}

		// mapping redeemScript to PrevTxs
		for idx, val := range prevsAddrs.Addrs {
			rs := cold.GetRedeemScriptByAddress(accountKeys, val)
			if rs == "" {
				logger.Error("redeemScript can not be found")
				continue
			}
			prevsAddrs.PrevTxs[idx].RedeemScript = rs
		}

		// serialize prevsAddrs with redeemScript
		newEncodedPrevsAddrs, err = serial.EncodeToString(prevsAddrs)
		if err != nil {
			return nil, false, "", fmt.Errorf("fail to call serial.EncodeToString(): %w", err)
		}

	case domainWallet.WalletTypeSign:
		authKey, err := u.authKeyRepo.GetOne("")
		if err != nil {
			return nil, false, "", fmt.Errorf("fail to call authKeyRepo.GetOne(): %w", err)
		}
		// wip
		wips = []string{authKey.WalletImportFormat}
	case domainWallet.WalletTypeWatchOnly:
		return nil, false, "", fmt.Errorf("WalletType is invalid: %s", u.wtype.String())
	default:
		return nil, false, "", fmt.Errorf("WalletType is invalid: %s", u.wtype.String())
	}

	// sign
	signedTx, isSigned, err := u.btc.SignRawTransactionWithKey(msgTx, wips, prevsAddrs.PrevTxs)
	return signedTx, isSigned, newEncodedPrevsAddrs, err
}
