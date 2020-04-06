package service

//Cold wallet

import (
	"fmt"
	"strings"

	"github.com/bookerzzz/grok"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/model"
	"github.com/hiromaily/go-bitcoin/pkg/serial"
	"github.com/hiromaily/go-bitcoin/pkg/txfile"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/btc"
)

// coldwallet側から未署名トランザクションを読み込み、署名を行う

// SignatureFromFile 渡されたファイルからtransactionを読み取り、署名を行う
// TODO:いずれにせよ、入金と出金で署名もMultisigかどうかで変わってくる
func (w *Wallet) SignatureFromFile(filePath string) (string, bool, string, error) {
	if w.Type == enum.WalletTypeWatchOnly {
		return "", false, "", errors.New("it's available on ColdWallet1, ColdWallet2")
	}

	//ファイル名から、tx_receipt_idを取得する
	//payment_5_unsigned_1534466246366489473
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
	tmp := strings.Split(data, ",")
	hex = tmp[0]
	if len(tmp) > 1 {
		encodedAddrsPrevs = tmp[1]
	}

	//署名
	hexTx, isSigned, newEncodedAddrsPrevs, err := w.signatureByHex(hex, encodedAddrsPrevs, actionType)
	if err != nil {
		return "", isSigned, "", err
	}

	//ファイルに書き込むデータ
	savedata := hexTx

	//署名が完了していないとき、TxTypeUnsigned2nd
	txType := enum.TxTypeSigned
	if isSigned == false {
		txType = enum.TxTypeUnsigned2nd
		if newEncodedAddrsPrevs != "" {
			savedata = fmt.Sprintf("%s,%s", savedata, newEncodedAddrsPrevs)
		}
	}

	//ファイルに書き込む
	path := txfile.CreateFilePath(actionType, txType, txReceiptID, true)
	generatedFileName, err := txfile.WriteFile(path, savedata)
	if err != nil {
		return "", isSigned, "", err
	}

	return hexTx, isSigned, generatedFileName, nil
}

// signatureByHex 署名する
// オフラインで使うことを想定
func (w *Wallet) signatureByHex(hex, encodedAddrsPrevs string, actionType enum.ActionType) (string, bool, string, error) {
	// Hexからトランザクションを取得
	msgTx, err := w.BTC.ToMsgTx(hex)
	if err != nil {
		return "", false, "", err
	}

	// 署名
	var (
		signedTx             *wire.MsgTx
		isSigned             bool
		addrsPrevs           btc.AddrsPrevTxs
		accountKeys          []model.AccountKeyTable
		wips                 []string
		newEncodedAddrsPrevs string
	)

	if encodedAddrsPrevs == "" {
		//Bitcoin coreのバージョン17から、常に必要
		return "", false, "", errors.New("encodedAddrsPrevs must be set")
	}

	//decodeする
	serial.DecodeFromString(encodedAddrsPrevs, &addrsPrevs)

	//WIPs, RedeedScriptを取得
	//TODO:coldwallet1とcoldwallet2で挙動が違う
	//TODO:receiptの場合、wipsは不要
	//coldwallet2の場合、AccountTypeAuthorizationが必要
	if w.Type == enum.WalletTypeSignature {
		//account_key_authorizationテーブルから情報を取得
		accountKey, err := w.DB.GetOneByMaxIDOnAccountKeyTable(account.AccountTypeAuthorization)
		if err != nil {
			return "", false, "", errors.Errorf("DB.GetOneByMaxIDOnAccountKeyTable() error: %s", err)
		}
		accountKeys = append(accountKeys, *accountKey)
	} else {
		//TODO:ActionTypeが`transfer`の場合、AccountのFromから判別しないといけない。。。
		//=> addrsPrevs.SenderAccount を使うように変更
		//if val, ok := enum.ActionToAccountMap[actionType]; ok {
		//	//account_key_payment/account_key_clientテーブルから取得
		//	accountKeys, err = w.DB.GetAllAccountKeyByMultiAddrs(val, addrsPrevs.Addrs)
		//	if err != nil {
		//		return "", false, "", errors.Errorf("DB.GetWIPByMultiAddrs() error: %s", err)
		//	}
		//} else {
		//	return "", false, "", errors.New("[Fatal] actionType can not be retrieved. it should be fixed programmatically")
		//}
		//account_key_payment/account_key_clientテーブルから取得
		accountKeys, err = w.DB.GetAllAccountKeyByMultiAddrs(addrsPrevs.SenderAccount, addrsPrevs.Addrs)
		if err != nil {
			return "", false, "", errors.Errorf("DB.GetWIPByMultiAddrs() error: %s", err)
		}
	}

	//wip
	for _, val := range accountKeys {
		wips = append(wips, val.WalletImportFormat)
	}

	//multisigの場合のみの処理
	//accountType, ok := enum.ActionToAccountMap[actionType]
	if account.AccountTypeMultisig[addrsPrevs.SenderAccount] {
		if w.Type == enum.WalletTypeKeyGen {
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

			//redeemScriptセット後、シリアライズして戻す
			newEncodedAddrsPrevs, err = serial.EncodeToString(addrsPrevs)
			if err != nil {
				return "", false, "", errors.Errorf("serial.EncodeToString(): error: %s", err)
			}
		} else {
			newEncodedAddrsPrevs = encodedAddrsPrevs
		}
	}

	//署名
	//multisigかどうかで判別
	if account.AccountTypeMultisig[addrsPrevs.SenderAccount] {
		signedTx, isSigned, err = w.BTC.SignRawTransactionWithKey(msgTx, wips, addrsPrevs.PrevTxs)
	} else {
		signedTx, isSigned, err = w.BTC.SignRawTransaction(msgTx, addrsPrevs.PrevTxs)
	}

	if err != nil {
		return "", false, "", err
	}
	logger.Debugf("isSigned is %t", isSigned)

	hexTx, err := w.BTC.ToHex(signedTx)
	if err != nil {
		return "", false, "", errors.Errorf("w.BTC.ToHex(msgTx): error: %s", err)
	}

	return hexTx, isSigned, newEncodedAddrsPrevs, nil

}
