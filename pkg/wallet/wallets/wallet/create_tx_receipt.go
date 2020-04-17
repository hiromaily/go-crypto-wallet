package wallet

import (
	"fmt"

	"github.com/bookerzzz/grok"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/action"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb/walletrepo"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/btc"
)

// - check unspentTx for client
// - get utxo and create unsigned tx
// - fee is fluctuating, so outdated unsigned tx must not be used, re-create unsigned tx
// - after signed tx sent, utxo could not be retrieved by listUnspent()

type parsedTx struct {
	txInputs       []btcjson.TransactionInput
	txRepoTxInputs []walletrepo.TxInput
	prevTxs        []btc.PrevTx
	addresses      []string
}

// CreateReceiptTx create unsigned tx if client accounts have coins
// - sender account: client, receiver account: receipt
// - receiver account covers fee
func (w *Wallet) CreateReceiptTx(adjustmentFee float64) (string, string, error) {
	fixedAccount := account.AccountTypeClient

	// get listUnspent
	unspentList, err := w.getUnspentList(fixedAccount)
	if len(unspentList) == 0 {
		w.logger.Info("no listunspent")
		return "", "", nil
	}

	// parse listUnspent
	parsedTx, inputTotal := w.parseListUnspentTx(unspentList)
	w.logger.Debug(
		"total coin to send (Satoshi) before fee calculated",
		zap.Any("input_amount", inputTotal),
		zap.Int("len(inputs)", len(parsedTx.txInputs)))
	if len(parsedTx.txInputs) == 0 {
		return "", "", nil
	}

	addrsPrevs := btc.AddrsPrevTxs{
		Addrs:         parsedTx.addresses,
		PrevTxs:       parsedTx.prevTxs,
		SenderAccount: fixedAccount,
	}

	//TODO: how this code can be integrated with CreateTransferTx ??

	// create raw tx
	hex, fileName, err := w.createRawTx(
		action.ActionTypeReceipt,
		account.AccountTypeReceipt,
		adjustmentFee,
		parsedTx.txInputs,
		inputTotal,
		parsedTx.txRepoTxInputs,
		&addrsPrevs)

	//TODO: notify what unsigned tx is created
	// => in command pkg

	return hex, fileName, err
}

// call API `unspentlist`
// no result and no error is possible, so caller should check both returned value
func (w *Wallet) getUnspentList(accountType account.AccountType) ([]btc.ListUnspentResult, error) {
	// unlock locked UnspentTransaction
	//if err := w.BTC.UnlockAllUnspentTransaction(); err != nil {
	//	return "", "", err
	//}

	// get listUnspent
	unspentList, _, err := w.btc.ListUnspentByAccount(accountType)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call btc.Client().ListUnspent()")
	}
	grok.Value(unspentList)
	if len(unspentList) == 0 {
		w.logger.Info("no listunspent")
		return nil, nil
	}
	return unspentList, nil
}

// parse result of listUnspent
// retured *parsedTx could be nil
func (w *Wallet) parseListUnspentTx(unspentList []btc.ListUnspentResult) (*parsedTx, btcutil.Amount) {
	var inputTotal btcutil.Amount
	txInputs := make([]btcjson.TransactionInput, 0, len(unspentList))
	txRepoTxInputs := make([]walletrepo.TxInput, 0, len(unspentList))
	prevTxs := make([]btc.PrevTx, 0, len(unspentList))
	addresses := make([]string, 0, len(unspentList))

	for _, tx := range unspentList {
		// Amount
		amt, err := btcutil.NewAmount(tx.Amount)
		if err != nil {
			//this error is not expected
			w.logger.Error(
				"btcutil.NewAmount()",
				zap.Float64("tx amount", tx.Amount),
				zap.Error(err))
			continue
		}
		inputTotal += amt

		txInputs = append(txInputs, btcjson.TransactionInput{
			Txid: tx.TxID,
			Vout: tx.Vout,
		})

		txRepoTxInputs = append(txRepoTxInputs, walletrepo.TxInput{
			ReceiptID:          0,
			InputTxid:          tx.TxID,
			InputVout:          tx.Vout,
			InputAddress:       tx.Address,
			InputAccount:       tx.Label,
			InputAmount:        fmt.Sprintf("%f", tx.Amount),
			InputConfirmations: tx.Confirmations,
		})

		prevTxs = append(prevTxs, btc.PrevTx{
			Txid:         tx.TxID,
			Vout:         tx.Vout,
			ScriptPubKey: tx.ScriptPubKey,
			RedeemScript: "", //required if target account is multsig address
			Amount:       tx.Amount,
		})

		addresses = append(addresses, tx.Address)
	}
	return &parsedTx{
		txInputs:       txInputs,
		txRepoTxInputs: txRepoTxInputs,
		prevTxs:        prevTxs,
		addresses:      addresses,
	}, inputTotal
}
