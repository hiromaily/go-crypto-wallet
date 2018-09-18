package service

//Cold wallet

import (
	"fmt"
	"strings"

	"github.com/bookerzzz/grok"
	"github.com/btcsuite/btcd/wire"
	"github.com/hiromaily/go-bitcoin/pkg/api/btc"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/model"
	"github.com/hiromaily/go-bitcoin/pkg/serial"
	"github.com/hiromaily/go-bitcoin/pkg/txfile"
	"github.com/pkg/errors"
)

// coldwallet側から未署名トランザクションを読み込み、署名を行う

// signatureByHex 署名する
// オフラインで使うことを想定
func (w *Wallet) signatureByHex(hex, encodedAddrsPrevs string, actionType enum.ActionType) (string, bool, error) {
	//first hex: 未署名トランザクションのhex
	// Hexからトランザクションを取得
	msgTx, err := w.BTC.ToMsgTx(hex)
	if err != nil {
		return "", false, err
	}

	// 署名
	var (
		signedTx    *wire.MsgTx
		isSigned    bool
		addrsPrevs  btc.AddrsPrevTxs
		accountKeys []model.AccountKeyTable
		wips        []string
	)

	if encodedAddrsPrevs != "" {
		//こちらの処理はMultisigの場合
		//decodeする
		serial.DecodeFromString(encodedAddrsPrevs, &addrsPrevs)
		grok.Value(addrsPrevs)
		//type AddrsPrevTxs struct {
		//	Addrs   []string
		//	PrevTxs []PrevTx
		//}

		//WIPs, RedeedScriptを取得
		if val, ok := enum.ActionToAccountMap[actionType]; ok {
			accountKeys, err = w.DB.GetAllAccountKeyByMultiAddrs(val, addrsPrevs.Addrs)
			if err != nil {
				return "", false, errors.Errorf("DB.GetWIPByMultiAddrs() error: %s", err)
			}
		} else {
			return "", false, errors.New("[Fatal] actionType can not be retrieved. it should be fixed programmatically")
		}

		//wip
		for _, val := range accountKeys {
			wips = append(wips, val.WalletImportFormat)
		}

		//取得したredeemScriptをPrevTxsにマッピング
		for idx, val := range addrsPrevs.Addrs {
			rs := model.GetRedeedScriptByAddress(accountKeys, val)
			if rs == "" {
				logger.Error("redeemScript can not be found")
				continue
			}
			addrsPrevs.PrevTxs[idx].RedeemScript = rs
		}
		grok.Value(addrsPrevs)

		//署名
		signedTx, isSigned, err = w.BTC.SignRawTransactionWithKey(msgTx, wips, addrsPrevs.PrevTxs)

		//panic("for now, it stops")

	} else {
		signedTx, isSigned, err = w.BTC.SignRawTransaction(msgTx)
	}
	if err != nil {
		return "", false, err
	}
	logger.Debugf("isSigned is %t", isSigned)

	hexTx, err := w.BTC.ToHex(signedTx)
	if err != nil {
		return "", false, errors.Errorf("w.BTC.ToHex(msgTx): error: %s", err)
	}

	return hexTx, isSigned, nil

}

// SignatureFromFile 渡されたファイルからtransactionを読み取り、署名を行う
// ColdWalletの機能なので、渡されたfilePathをそのまま使う?
// TODO:いずれにせよ、入金と出金で署名もMultisigかどうかで変わってくる
func (w *Wallet) SignatureFromFile(filePath string) (string, bool, string, error) {
	//ファイル名から、tx_receipt_idを取得する
	//payment_5_unsigned_1534466246366489473
	//txReceiptID, actionType, _, err := txfile.ParseFile(filePath, "unsigned")
	txReceiptID, actionType, _, err := txfile.ParseFile(filePath, []enum.TxType{enum.TxTypeUnsigned, enum.TxTypeUnsigned2nd})
	if err != nil {
		return "", false, "", err
	}

	//ファイルからhexを読み取る
	data, err := txfile.ReadFile(filePath)
	if err != nil {
		return "", false, "", err
	}

	var hex, encodedAddrsPrevs string

	//encodedPrevTxs
	//paymentの場合は、multisigのため、データが異なる
	//TODO:multisigかどうかの判別は、enum.AccountTypeMultisig[]で行う
	//TODO:ActionType/AccountTypeの相互変換が必要かも
	if val, ok := enum.ActionToAccountMap[actionType]; ok {
		if enum.AccountTypeMultisig[val] {
			//if actionType == enum.ActionTypePayment && enum.AccountTypeMultisig[enum.AccountTypePayment] {
			tmp := strings.Split(data, ",")
			if len(tmp) != 2 {
				return "", false, "", errors.New("imported tx data is wrong. encodedPrevTxs would not be found")
			}
			hex = tmp[0]
			encodedAddrsPrevs = tmp[1]
			//TODO:署名が更に必要なので、ファイル出力時にこの情報も引き継ぐ必要がある
		}
	}

	//署名
	hexTx, isSigned, err := w.signatureByHex(hex, encodedAddrsPrevs, actionType)
	if err != nil {
		return "", isSigned, "", err
	}

	//ファイルに書き込むデータ
	savedata := hexTx

	//TODO:署名が完了していないとき、TxTypeUnsigned2nd
	txType := enum.TxTypeSigned
	if isSigned == false {
		txType = enum.TxTypeUnsigned2nd
		if encodedAddrsPrevs != "" {
			savedata = fmt.Sprintf("%s,%s", savedata, encodedAddrsPrevs)
		}
	}

	//ファイルに書き込む
	//path := txfile.CreateFilePath(actionType, enum.TxTypeSigned, txReceiptID, true)
	path := txfile.CreateFilePath(actionType, txType, txReceiptID, true)
	generatedFileName, err := txfile.WriteFile(path, savedata)
	if err != nil {
		return "", isSigned, "", err
	}

	return hexTx, isSigned, generatedFileName, nil
}

// SignatureByHex Hex文字列から署名を行う
// TODO:出金/入金でフラグがほしいが、このfuncはDebug時にしか使わない
//func (w *Wallet) SignatureByHex(actionType enum.ActionType, hex string, txReceiptID int64) (string, bool, string, error) {
//	//署名
//	hexTx, isSigned, err := w.signatureByHex(hex, "")
//	if err != nil {
//		return "", isSigned, "", err
//	}
//
//	//ファイルに書き込む
//	path := txfile.CreateFilePath(actionType, enum.TxTypeSigned, txReceiptID, true)
//	generatedFileName, err := txfile.WriteFile(path, hex)
//	if err != nil {
//		return "", isSigned, "", err
//	}
//
//	return hexTx, isSigned, generatedFileName, nil
//}
