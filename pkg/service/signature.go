package service

import (
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/file"
	"github.com/pkg/errors"
	"log"
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
	log.Printf("[Debug] isSigned01 is %t", isSigned)

	hexTx, err := w.BTC.ToHex(signedTx)
	if err != nil {
		return "", false, errors.Errorf("w.BTC.ToHex(msgTx): error: %v", err)
	}

	return hexTx, isSigned, nil

}

// SignatureByHex Hex文字列から署名を行う
// TODO:出金/入金でフラグがほしいが、これはDebug時にしか使わない
func (w *Wallet) SignatureByHex(hex string, txReceiptID int64) (string, bool, string, error) {
	//署名
	hexTx, isSigned, err := w.signatureByHex(hex)
	if err != nil {
		return "", isSigned, "", err
	}
	//log.Println("hex:", hexTx)

	//ファイルに書き込む
	//TODO:暫定で1を使っている
	path := file.CreateFilePath(enum.ActionReceipt, enum.TxTypeSigned, txReceiptID)
	generatedFileName, err := file.WriteFile(path, hex)
	//generatedFileName := file.WriteFileForSigned(txReceiptID, "inside/", hexTx)
	if err != nil {
		return "", isSigned, "", err
	}

	return hexTx, isSigned, generatedFileName, nil
}

// SignatureFromFile 渡されたファイルからtransactionを読み取り、署名を行う
// ColdWalletの機能なので、渡されたfilePathをそのまま使う?
// TODO:いずれにせよ、入金と出金で署名もMultisigかどうかで変わってくる
func (w *Wallet) SignatureFromFile(filePath string, actionFlg enum.Action) (string, bool, string, error) {
	//ファイル名から、tx_receipt_idを取得する
	//5_unsigned_1534466246366489473
	txReceiptID, _, err := file.ParseFile(filePath, "unsigned")
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
	//log.Println("hex:", hexTx)

	//ファイルに書き込む
	path := file.CreateFilePath(actionFlg, enum.TxTypeSigned, txReceiptID)
	generatedFileName, err := file.WriteFile(path, hexTx)
	if err != nil {
		return "", isSigned, "", err
	}

	return hexTx, isSigned, generatedFileName, nil
}
