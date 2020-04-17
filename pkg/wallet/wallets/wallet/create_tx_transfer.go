package wallet

import (
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/action"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb/walletrepo"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/btc"
)

// CreateTransferTx transfer coin among internal account except client, authorization
// - if amount=0, all coin is sent
// TODO: temporary fixed account is used `receipt to payment`
// TODO: after this func, what if `listtransactions` api is called to see result
func (w *Wallet) CreateTransferTx(sender, receiver account.AccountType, floatAmount, adjustmentFee float64) (string, string, error) {

	// validation
	if receiver == account.AccountTypeClient || receiver == account.AccountTypeAuthorization {
		return "", "", errors.New("invalid receiver account. client, authorization account is not allowed as receiver")
	}
	if sender == receiver {
		return "", "", errors.New("invalid account. sender and receiver is same")
	}
	//amount btcutil.Amount
	amount, err := w.btc.FloatBitToAmount(floatAmount)
	if err != nil {
		return "", "", err
	}

	// check balance for sender
	balance, err := w.btc.GetReceivedByLabelAndMinConf(sender.String(), w.btc.ConfirmationBlock())
	if err != nil {
		return "", "", err
	}
	if balance <= amount {
		//balance is short
		return "", "", errors.Errorf("account: %s balance is insufficient", sender)
	}

	//指定金額になるまで、utxoからinputを作成する
	// Listunspent()にてpaymentアカウント用のutxoをすべて取得する
	unspentList, _, err := w.btc.ListUnspentByAccount(sender)
	if err != nil {
		return "", "", errors.Errorf("BTC.ListUnspentByAccount(%s) error: %s", sender, err)
	}

	if len(unspentList) == 0 {
		w.logger.Info("no listunspent for from account", zap.String("account", sender.String()))
		return "", "", nil
	}

	//送金先は1つのアドレスなので、receipt.goのロジックに近い
	var (
		inputs          []btcjson.TransactionInput
		inputTotal      btcutil.Amount
		txReceiptInputs []walletrepo.TxInput
		prevTxs         []btc.PrevTx
		addresses       []string
	)

	for _, tx := range unspentList {

		// Amount
		amt, err := btcutil.NewAmount(tx.Amount)
		if err != nil {
			//このエラーは起こりえない
			w.logger.Error(
				"btcutil.NewAmount()",
				zap.Float64("amount", tx.Amount),
				zap.Error(err))
			continue
		}
		inputTotal += amt //合計

		//TODO:Ver17対応が必要
		//lockunspentによって、該当トランザクションをロックして再度ListUnspent()で出力されることを防ぐ
		//if w.BTC.LockUnspent(tx) != nil {
		//	continue
		//}

		// inputs
		inputs = append(inputs, btcjson.TransactionInput{
			Txid: tx.TxID,
			Vout: tx.Vout,
		})

		// txReceiptInputs
		txReceiptInputs = append(txReceiptInputs, walletrepo.TxInput{
			ReceiptID:          0,
			InputTxid:          tx.TxID,
			InputVout:          tx.Vout,
			InputAddress:       tx.Address,
			InputAccount:       tx.Label,
			InputAmount:        fmt.Sprintf("%f", tx.Amount),
			InputConfirmations: tx.Confirmations,
		})

		// prevTxs(walletでの署名でもversion17からは必要になる。。。fuck)
		prevTxs = append(prevTxs, btc.PrevTx{
			Txid:         tx.TxID,
			Vout:         tx.Vout,
			ScriptPubKey: tx.ScriptPubKey,
			RedeemScript: tx.RedeemScript,
			Amount:       tx.Amount,
		})

		//tx.Address
		addresses = append(addresses, tx.Address)

		//totalをチェック
		if amount == 0 {
			continue
		}
		if inputTotal > amount {
			break
		}
	}

	w.logger.Debug(
		"total coin to send (Satoshi) before fee calculated, input length: %d",
		zap.Any("amount", inputTotal),
		zap.Int("len(inputs)", len(inputs)))
	if len(inputs) == 0 {
		return "", "", nil
	}

	addrsPrevs := btc.AddrsPrevTxs{
		Addrs:         addresses,
		PrevTxs:       prevTxs,
		SenderAccount: sender,
	}

	// 一連の処理を実行
	hex, fileName, err := w.createRawTx(action.ActionTypeTransfer, receiver, 0, inputs,
		inputTotal, txReceiptInputs, &addrsPrevs)

	return hex, fileName, err
}
