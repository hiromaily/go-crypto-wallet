package service

import (
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/file"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/pkg/errors"
)

// coldwallet側から未署名トランザクションを読み込み、署名を行う

// signatureByHex 署名する
// オフラインで使うことを想定
func (w *Wallet) signatureByHex(hex string) (string, bool, error) {
	//first hex: 未署名トランザクションのhex
	// Hexからトランザクションを取得
	msgTx, err := w.BTC.ToMsgTx(hex)
	if err != nil {
		return "", false, err
	}

	// 署名
	signedTx, isSigned, err := w.BTC.SignRawTransaction(msgTx)
	if err != nil {
		return "", false, err
	}
	logger.Debugf("isSigned is %t", isSigned)

	hexTx, err := w.BTC.ToHex(signedTx)
	if err != nil {
		return "", false, errors.Errorf("w.BTC.ToHex(msgTx): error: %v", err)
	}

	return hexTx, isSigned, nil

}

// SignatureByHex Hex文字列から署名を行う
// TODO:出金/入金でフラグがほしいが、このfuncはDebug時にしか使わない
func (w *Wallet) SignatureByHex(actionType enum.ActionType, hex string, txReceiptID int64) (string, bool, string, error) {
	//署名
	hexTx, isSigned, err := w.signatureByHex(hex)
	if err != nil {
		return "", isSigned, "", err
	}

	//ファイルに書き込む
	path := file.CreateFilePath(actionType, enum.TxTypeSigned, txReceiptID)
	generatedFileName, err := file.WriteFile(path, hex)
	if err != nil {
		return "", isSigned, "", err
	}

	return hexTx, isSigned, generatedFileName, nil
}

// SignatureFromFile 渡されたファイルからtransactionを読み取り、署名を行う
// ColdWalletの機能なので、渡されたfilePathをそのまま使う?
// TODO:いずれにせよ、入金と出金で署名もMultisigかどうかで変わってくる
func (w *Wallet) SignatureFromFile(filePath string) (string, bool, string, error) {
	//ファイル名から、tx_receipt_idを取得する
	//payment_5_unsigned_1534466246366489473
	txReceiptID, actionType, _, err := file.ParseFile(filePath, "unsigned")
	if err != nil {
		return "", false, "", err
	}

	//ファイルからhexを読み取る
	hex, err := file.ReadFile(filePath)
	if err != nil {
		return "", false, "", err
	}

	//署名
	hexTx, isSigned, err := w.signatureByHex(hex)
	if err != nil {
		return "", isSigned, "", err
	}

	//ファイルに書き込む
	path := file.CreateFilePath(actionType, enum.TxTypeSigned, txReceiptID)
	generatedFileName, err := file.WriteFile(path, hexTx)
	if err != nil {
		return "", isSigned, "", err
	}

	return hexTx, isSigned, generatedFileName, nil
}
