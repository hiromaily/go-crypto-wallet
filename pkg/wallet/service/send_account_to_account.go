package service

import (
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/model"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/btc"
)

//SendToAccount 内部アカウント間での送金
// amountが0のとき、全額送金する
// TODO:実行後、`listtransactions`methodで確認できるかも(未チェック)
func (w *Wallet) SendToAccount(from, to account.AccountType, amount btcutil.Amount) (string, string, error) {
	if w.Type != enum.WalletTypeWatchOnly {
		return "", "", errors.New("it's available on WatchOnlyWallet")
	}

	//Validation
	//とりあえず、receipt to paymentで実装
	//AccountTypeClient, AccountTypeAuthorizationは除外する
	if to == account.AccountTypeClient || to == account.AccountTypeAuthorization {
		return "", "", errors.New("Client, Authorization account can not receive coin")
	}
	if from == to {
		return "", "", errors.New("Validation error. `from` and `to` accountType should be different")
	}

	//残高確認
	balance, err := w.BTC.GetReceivedByAccountAndMinConf(string(from), w.BTC.ConfirmationBlock())
	if err != nil {
		return "", "", err
	}
	if balance <= amount {
		//残高が不足している
		return "", "", errors.Errorf("%s account balance is insufficient", from)
	}

	//指定金額になるまで、utxoからinputを作成する
	// Listunspent()にてpaymentアカウント用のutxoをすべて取得する
	unspentList, _, err := w.BTC.ListUnspentByAccount(from)
	if err != nil {
		return "", "", errors.Errorf("BTC.ListUnspentByAccount(%s) error: %s", from, err)
	}

	if len(unspentList) == 0 {
		logger.Infof("no listunspent for %s", from)
		return "", "", nil
	}

	//送金先は1つのアドレスなので、receipt.goのロジックに近い
	var (
		inputs          []btcjson.TransactionInput
		inputTotal      btcutil.Amount
		txReceiptInputs []model.TxInput
		prevTxs         []btc.PrevTx
		addresses       []string
	)

	for _, tx := range unspentList {

		// Amount
		amt, err := btcutil.NewAmount(tx.Amount)
		if err != nil {
			//このエラーは起こりえない
			logger.Errorf("btcutil.NewAmount(%f): error: %s", tx.Amount, err)
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
		txReceiptInputs = append(txReceiptInputs, model.TxInput{
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

	logger.Debugf("total coin to send:%d(Satoshi) before fee calculated, input length: %d", inputTotal, len(inputs))
	if len(inputs) == 0 {
		return "", "", nil
	}

	addrsPrevs := btc.AddrsPrevTxs{
		Addrs:         addresses,
		PrevTxs:       prevTxs,
		SenderAccount: from,
	}

	// 一連の処理を実行
	hex, fileName, err := w.createRawTransactionAndFee(enum.ActionTypeTransfer, to, 0, inputs,
		inputTotal, txReceiptInputs, &addrsPrevs)

	return hex, fileName, err
}
