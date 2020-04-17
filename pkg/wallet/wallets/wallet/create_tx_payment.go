package wallet

import (
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

// CreatePaymentTx create unsigned tx for user(anonymous addresses)
// sender: payment, receiver: addresses coming from user_payment table
func (w *Wallet) CreatePaymentTx(adjustmentFee float64) (string, string, error) {
	sender := account.AccountTypePayment
	//receiver := account.AccountTypeAnonymous
	//targetAction := action.ActionTypePayment

	// get payment data from payment_request
	userPayments, paymentRequestIds, err := w.createUserPayment()
	if err != nil {
		return "", "", err
	}
	if len(userPayments) == 0 {
		w.logger.Debug("no data in userPayments")
		// no data
		return "", "", nil
	}

	// calculate total amount to send from payment_request
	var userTotal btcutil.Amount
	for _, val := range userPayments {
		userTotal += val.validAmount
	}

	// get balance for payment account
	balance, err := w.btc.GetReceivedByLabelAndMinConf(account.AccountTypePayment.String(), w.btc.ConfirmationBlock())
	if err != nil {
		return "", "", err
	}
	if balance <= userTotal {
		//balance is short
		return "", "", errors.New("balance for payment account is insufficient")
	}
	w.logger.Debug("balane_userTotal",
		zap.Any("balance", balance),
		zap.Any("userTotal", userTotal))

	//FIXME: how to commonalize code from here

	// get listUnspent
	unspentList, unspentAddrs, err := w.getUnspentList(sender)
	if len(unspentList) == 0 {
		w.logger.Info("no listunspent")
		return "", "", nil
	}

	// parse listUnspent
	parsedTx, inputTotal, isDone := w.parseListUnspentTx(unspentList, userTotal)
	w.logger.Debug(
		"total coin to send (Satoshi) before fee calculated",
		zap.Any("input_amount", inputTotal),
		zap.Int("len(inputs)", len(parsedTx.txInputs)))
	if len(parsedTx.txInputs) == 0 {
		return "", "", nil
	}
	if !isDone {
		return "", "", errors.New("sender account can't meet amount to send")
	}

	addrsPrevs := btc.AddrsPrevTxs{
		Addrs:         parsedTx.addresses,
		PrevTxs:       parsedTx.prevTxs,
		SenderAccount: sender,
	}

	// payment logic exclusively
	var (
		outputs    = map[btcutil.Address]btcutil.Amount{}
		tmpOutputs = map[string]btcutil.Amount{} //mapのkeyが、btcutil.Address型だとユニークかどうかkeyから判定できないため、判定用としてこちらを作成
	)

	//6.outputは別途こちらで作成
	for _, userPayment := range userPayments {
		grok.Value(userPayment)
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
	chargeAddr := unspentAddrs[0].String() //change from w.BTC.PaymentAddress()
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

	//Debug
	grok.Value(parsedTx.txInputs)
	grok.Value(tmpOutputs)
	grok.Value(outputs)

	// 一連の処理を実行
	return w.createRawTransactionForPayment(
		adjustmentFee,
		parsedTx.txInputs,
		inputTotal,
		parsedTx.txRepoTxInputs,
		outputs,
		&addrsPrevs,
		paymentRequestIds)
}

//TODO:receipt側のロジックとほぼ同じ為、まとめたほうがいい(手数料のロジックのみ異なる)
func (w *Wallet) createRawTransactionForPayment(
	adjustmentFee float64,
	inputs []btcjson.TransactionInput,
	inputTotal btcutil.Amount,
	txPaymentInputs []walletrepo.TxInput,
	outputs map[btcutil.Address]btcutil.Amount,
	addrsPrevs *btc.AddrsPrevTxs,
	paymentRequestIds []int64) (string, string, error) {

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
		generatedFileName, err = w.generateHexFile(action.ActionTypePayment, hex, encodedAddrsPrevs, txReceiptID)
		if err != nil {
			return "", "", errors.Errorf("wallet.storeHex(): error: %s", err)
		}
	}

	// 9. 出金準備に入ったことをユーザーに通知
	// TODO:NatsのPublisherとして通知すればいいか？

	return hex, generatedFileName, nil
}

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
