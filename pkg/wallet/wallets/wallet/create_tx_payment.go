package wallet

//Watch only wallet

import (
	"fmt"
	"strconv"

	"github.com/bookerzzz/grok"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/action"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb/walletrepo"
	"github.com/hiromaily/go-bitcoin/pkg/serial"
	"github.com/hiromaily/go-bitcoin/pkg/tx"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/btc"
)

// ユーザーからの出金依頼事に出金する処理
// Cayenneユーザーへの送金は内部振替が可能なため、こちらのWalletには依頼はこない

// UserPayment ユーザーの支払先アドレスと金額
type UserPayment struct {
	senderAddr   string          //送信者のアドレス (履歴を追うためだけに保持)
	receiverAddr string          //受信者のアドレス
	validRecAddr btcutil.Address //受診者のアドレス(変換後)
	amount       float64         //送信金額
	validAmount  btcutil.Amount  //送金金額(変換後)
}

// createUserPayment 出金依頼テーブルから処理するためのデータを取得する
func (w *Wallet) createUserPayment() ([]UserPayment, []int64, error) {
	paymentRequests, err := w.storager.GetPaymentRequestAll()
	if err != nil {
		return nil, nil, errors.Errorf("DB.GetPaymentRequestAll() error: %s", err)
	}
	if len(paymentRequests) == 0 {
		//処理するデータが存在しない。(エラーではない)
		return nil, nil, nil
	}

	userPayments := make([]UserPayment, len(paymentRequests))
	paymentRequestIds := make([]int64, len(paymentRequests))

	//TODO:更新用にidの配列も保持しておくこと
	for idx, val := range paymentRequests {
		paymentRequestIds[idx] = val.ID

		userPayments[idx].senderAddr = val.AddressFrom
		userPayments[idx].receiverAddr = val.AddressTo
		amt, err := strconv.ParseFloat(val.Amount, 64)
		if err != nil {
			//致命的なエラー、起きるのであればプログラムがおかしい
			w.logger.Error("payment_request table includes invalid amount field")
			return nil, nil, errors.New("payment_request table includes invalid amount field")
		}
		userPayments[idx].amount = amt

		//Address TODO:このタイミングでaddressは不要かもしれない
		userPayments[idx].validRecAddr, err = w.btc.DecodeAddress(userPayments[idx].receiverAddr)
		if err != nil {
			//致命的なエラー: これは本来事前にチェックされるため、ありえないはず
			w.logger.Error("unexpected error occurred converting string type receiverAddr to address type")
			return nil, nil, errors.New("unexpected error occurred converting string type receiverAddr to address type")
		}
		//grok.Value(userPayments[idx].validRecAddr)

		//Amount
		userPayments[idx].validAmount, err = w.btc.FloatBitToAmount(userPayments[idx].amount)
		if err != nil {
			//致命的なエラー: これは本来事前にチェックされるため、ありえないはず
			w.logger.Error("unexpected error occurred converting float64 type amount to Amount type")
			return nil, nil, errors.New("unexpected error occurred converting float64 type amount to Amount type")
		}
	}

	//Addressが重複する場合、合算したほうがいいかも => transactionのoutputを作成するタイミングで合算可能
	return userPayments, paymentRequestIds, nil
}

// isFoundTxIDAndVout 指定したtxIDとvoutに紐づくinputが既に存在しているか確認する
// TODO:utilityとして有用なので、transaction.goに移動するか
func (w *Wallet) isFoundTxIDAndVout(txID string, vout uint32, inputs []btcjson.TransactionInput) bool {
	for _, val := range inputs {
		if val.Txid == txID && val.Vout == vout {
			return true
		}
	}
	return false
}

// CreateUnsignedPaymentTx 支払いのための未署名トランザクションを作成する
func (w *Wallet) CreateUnsignedPaymentTx(adjustmentFee float64) (string, string, error) {

	//DBから情報を取得、もしくはNatsからのリクエストをトリガーにするか、今はまだ未定
	//とりあえず、テストデータで実装
	//実際には、こちらもamountでソートしておく=> ソートは不要

	//1.出金データを取得
	userPayments, paymentRequestIds, err := w.createUserPayment()
	if err != nil {
		return "", "", err
	}
	if len(userPayments) == 0 {
		//処理するデータがない
		return "", "", nil
	}

	//2.送信合計金額を算出
	var userTotal btcutil.Amount
	for _, val := range userPayments {
		userTotal += val.validAmount
	}

	//3.payment用アドレスの残高を確認し、金額が不足しているのであればエラー
	balance, err := w.btc.GetReceivedByLabelAndMinConf(string(account.AccountTypePayment), w.btc.ConfirmationBlock())
	if err != nil {
		return "", "", err
	}
	if balance <= userTotal {
		//残高が不足している
		return "", "", errors.New("Payment account balance is insufficient")
	}

	//4. Listunspent()にてpaymentアカウント用のutxoをすべて取得する
	unspentList, addrs, err := w.btc.ListUnspentByAccount(account.AccountTypePayment)
	if err != nil {
		return "", "", errors.Errorf("BTC.ListUnspentByAccount() error: %s", err)
	}

	var (
		inputs          []btcjson.TransactionInput
		inputTotal      btcutil.Amount
		txPaymentInputs []walletrepo.TxInput
		outputs         = map[btcutil.Address]btcutil.Amount{}
		tmpOutputs      = map[string]btcutil.Amount{} //mapのkeyが、btcutil.Address型だとユニークかどうかkeyから判定できないため、判定用としてこちらを作成
		isDone          bool
		prevTxs         []btc.PrevTx
		addresses       []string
	)

	//5.合計金額を超えるまで、listunspentからinputsを作成する
	for _, tx := range unspentList {
		amt, err := btcutil.NewAmount(tx.Amount)
		if err != nil {
			w.logger.Error(
				"btcutil.NewAmount()",
				zap.Any("amount", tx.Amount),
				zap.Error(err))
			continue
		}
		inputTotal += amt

		//utxoの新規利用のため、inputを作成
		inputs = append(inputs, btcjson.TransactionInput{
			Txid: tx.TxID,
			Vout: tx.Vout,
		})
		// txReceiptInputs
		txPaymentInputs = append(txPaymentInputs, walletrepo.TxInput{
			ReceiptID:          0,
			InputTxid:          tx.TxID,
			InputVout:          tx.Vout,
			InputAddress:       tx.Address,
			InputAccount:       tx.Label,
			InputAmount:        fmt.Sprintf("%f", tx.Amount),
			InputConfirmations: tx.Confirmations,
		})

		// prevTxs(multisigアドレスからの送金時、未署名トランザクションへの署名に必要となる)
		prevTxs = append(prevTxs, btc.PrevTx{
			Txid:         tx.TxID,
			Vout:         tx.Vout,
			ScriptPubKey: tx.ScriptPubKey,
			RedeemScript: tx.RedeemScript,
			Amount:       tx.Amount,
		})
		//tx.Address
		addresses = append(addresses, tx.Address)

		if inputTotal > userTotal {
			isDone = true
			break
		}
	}

	if !isDone {
		//これは金額が足りなかったことを意味するので致命的
		//上記のGetBalanceByAccountAndMinConf()のチェックでここは通らないはず
		return "", "", errors.New("[fatal error] bitcoin is insufficient in trading account of ours")
	}

	//6.outputは別途こちらで作成
	for _, userPayment := range userPayments {
		if _, ok := tmpOutputs[userPayment.receiverAddr]; ok {
			//加算する
			tmpOutputs[userPayment.receiverAddr] += userPayment.validAmount
		} else {
			//新規送信先を作成
			tmpOutputs[userPayment.receiverAddr] = userPayment.validAmount
		}
	}

	//7.差分でお釣り用のoutputを作成する (paymentのおつりaddressに返金する処理)
	//おつり受け取りのpaymentaddressが複数ある場合、何を使うか。とりあえずaddrs[0]
	//TODO:ユーザーの出金先にpaymentのアドレスが指定された場合のためのロジック(イレギュラーだが運用上可能)
	//FIXME: paymentが複数addressであることを考慮
	//FIXME:BIP44でaddressを生成したが、おつり用としては分けていない。これが問題をおこなさないか？
	chargeAddr := addrs[0].String() //change from w.BTC.PaymentAddress()
	change := inputTotal - userTotal
	if _, ok := tmpOutputs[chargeAddr]; ok {
		//加算
		//TODO:実運用ではありえない
		tmpOutputs[chargeAddr] += change
	} else {
		//新規
		tmpOutputs[chargeAddr] = change
	}

	//8.tmpOutputsをoutputsとして変換する
	for key, val := range tmpOutputs {
		addr, err := w.btc.DecodeAddress(key)
		if err != nil {
			//これは本来事前にチェックされるため、ありえないはず
			w.logger.Error("unexpected error converting string to address")
			continue
		}
		outputs[addr] = val
	}

	//TODO:inputsをすべてlockする必要がある？？

	addrsPrevs := btc.AddrsPrevTxs{
		Addrs:         addresses,
		PrevTxs:       prevTxs,
		SenderAccount: account.AccountTypePayment,
	}

	//Debug
	grok.Value(inputs)
	grok.Value(tmpOutputs)
	grok.Value(outputs)

	// 一連の処理を実行
	return w.createRawTransactionForPayment(adjustmentFee, inputs, inputTotal, txPaymentInputs, outputs,
		paymentRequestIds, &addrsPrevs)
}

//TODO:receipt側のロジックとほぼ同じ為、まとめたほうがいい(手数料のロジックのみ異なる)
func (w *Wallet) createRawTransactionForPayment(adjustmentFee float64, inputs []btcjson.TransactionInput,
	inputTotal btcutil.Amount, txPaymentInputs []walletrepo.TxInput, outputs map[btcutil.Address]btcutil.Amount,
	paymentRequestIds []int64, addrsPrevs *btc.AddrsPrevTxs) (string, string, error) {

	var (
		outputTotal      btcutil.Amount
		txPaymentOutputs []walletrepo.TxOutput
	)

	// 1.CreateRawTransactionWithOutput(仮で作成し、この後サイズから手数料を算出する)
	msgTx, err := w.btc.CreateRawTransactionWithOutput(inputs, outputs)
	if err != nil {
		return "", "", errors.Errorf("BTC.CreateRawTransactionWithOutput(): error: %s", err)
	}
	w.logger.Debug("CreateRawTransactionWithOutput",
		zap.Any("CreateRawTransactionWithOutput", msgTx))

	// 2.fee算出
	fee, err := w.btc.GetFee(msgTx, adjustmentFee)
	if err != nil {
		return "", "", errors.Errorf("BTC.GetFee(): error: %s", err)
	}

	// 3.お釣り用のoutputのトランザクションから、手数料を差し引く
	// FIXME: これが足りない場合がめんどくさい。。。これをどう回避すべきか
	for addr, amt := range outputs {
		//if addr.String() == w.BTC.PaymentAddress() {
		if acnt, _ := w.btc.GetAccount(addr.String()); acnt == string(account.AccountTypePayment) {
			outputs[addr] -= fee
			//break
			txPaymentOutputs = append(txPaymentOutputs, walletrepo.TxOutput{
				ReceiptID:     0,
				OutputAddress: addr.String(),
				OutputAccount: string(account.AccountTypePayment),
				OutputAmount:  w.btc.AmountString(amt - fee),
				IsChange:      true,
			})
		} else {
			txPaymentOutputs = append(txPaymentOutputs, walletrepo.TxOutput{
				ReceiptID:     0,
				OutputAddress: addr.String(),
				OutputAccount: "",
				OutputAmount:  w.btc.AmountString(amt),
				IsChange:      false,
			})
		}
		//total
		outputTotal += amt

	}
	//total
	outputTotal -= fee

	// 4.再度 CreateRawTransaction
	msgTx, err = w.btc.CreateRawTransactionWithOutput(inputs, outputs)
	if err != nil {
		return "", "", errors.Errorf("BTC.CreateRawTransactionWithOutput(): error: %s", err)
	}

	// 5.出力用にHexに変換する
	hex, err := w.btc.ToHex(msgTx)
	if err != nil {
		return "", "", errors.Errorf("BTC.ToHex(msgTx): error: %s", err)
	}

	// 6. Databaseに必要な情報を保存
	//  txType //1.未署名
	txReceiptID, err := w.insertTxTableForUnsigned(action.ActionTypePayment, hex, inputTotal, outputTotal, fee, tx.TxTypeValue[tx.TxTypeUnsigned], txPaymentInputs, txPaymentOutputs, paymentRequestIds)
	if err != nil {
		return "", "", errors.Errorf("insertTxTableForUnsigned(): error: %s", err)
	}

	// 7. serialize previous txs for multisig signature
	encodedAddrsPrevs, err := serial.EncodeToString(*addrsPrevs)
	if err != nil {
		return "", "", errors.Errorf("serial.EncodeToString(): error: %s", err)
	}
	w.logger.Debug(
		"encodedAddrsPrevs",
		zap.String("encodedAddrsPrevs", encodedAddrsPrevs))

	// 8. GCSにトランザクションファイルを作成
	var generatedFileName string
	if txReceiptID != 0 {
		generatedFileName, err = w.storeHex(hex, encodedAddrsPrevs, txReceiptID, action.ActionTypePayment)
		if err != nil {
			return "", "", errors.Errorf("wallet.storeHex(): error: %s", err)
		}
	}

	// 9. 出金準備に入ったことをユーザーに通知
	// TODO:NatsのPublisherとして通知すればいいか？

	return hex, generatedFileName, nil
}
