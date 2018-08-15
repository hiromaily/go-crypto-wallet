package service

import (
	"github.com/pkg/errors"
	"log"
)

//TODO: こちらは不要になるはず。service/signature.go側にて開発

// MultiSigByHex 一連のマルチシグアドレスに対しての署名処理
// オフラインで使うことを想定
func (w *Wallet) MultiSigByHex(hex string) (string, error) {
	//first hex: 未署名トランザクションのhex
	// Hexからトランザクションを取得
	msgTx, err := w.Btc.ToMsgTx(hex)
	if err != nil {
		return "", err
	}

	// TODO:署名処理はLoopのほうがいいか？
	// 署名1
	signedTx, isSigned, err := w.Btc.SignRawTransaction(msgTx)
	if err != nil {
		return "", err
	}
	//TODO:multisigでも1回でisSignedがtrueになった。。。
	//つまり,実行環境化にて両方の秘密鍵が保持されている場合に、1回でOKになるのでは？
	log.Printf("[Debug] isSigned01 is %t, false is expected.", isSigned)

	// 署名2
	//signedTx2, isSigned, err := w.Btc.SignRawTransaction(signedTx1)
	//if err != nil {
	//	return "", err
	//}
	//log.Printf("[Debug] isSigned02 is %t, true is expected.", isSigned)

	hexTx, err := w.Btc.ToHex(signedTx)
	//hexTx, err := w.Btc.ToHex(signedTx2)
	if err != nil {
		return "", errors.Errorf("w.Btc.ToHex(msgTx): error: %v", err)
	}

	return hexTx, nil
}
