package watchsrv

import (
	"math/big"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/action"
	models "github.com/hiromaily/go-bitcoin/pkg/models/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/serial"
)

// CreateTransferTx create unsigned tx for transfer coin among internal account except client, authorization
// FIXME: for now, receiver account covers fee, but is should be flexible
// TODO: implement
// sender pays fee,
// TODO: maybe any internal account should have only one address in Ethereum because no utxo
func (t *TxCreate) CreateTransferTx(sender, receiver account.AccountType, floatAmount, adjustmentFee float64) (string, string, error) {
	targetAction := action.ActionTypeTransfer

	// validation account
	if receiver == account.AccountTypeClient || receiver == account.AccountTypeAuthorization {
		return "", "", errors.New("invalid receiver account. client, authorization account is not allowed as receiver")
	}
	if sender == receiver {
		return "", "", errors.New("invalid account. sender and receiver is same")
	}

	// check sernder's total balance
	addrs, err := t.addrRepo.GetAllAddress(sender)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call addrRepo.GetAll(account.AccountTypeClient)")
	}
	total, userAmounts := t.eth.GetTotalBalance(addrs)
	if total.Uint64() == 0 || len(userAmounts) == 0 {
		return "", "", errors.New("sender has no balance")
	}

	// convert float amout to big.Int
	value := t.eth.FromFloatEther(floatAmount)
	t.logger.Debug("amount",
		zap.Float64("floatAmount(Ether)", floatAmount),
		zap.Uint64("bigIntAmount(Wei)", value.Uint64()),
		zap.Uint64("total", total.Uint64()),
	)
	if floatAmount != 0 && (total.Uint64() <= value.Uint64()) {
		return "", "", errors.New("sender balance is insufficient to send")
	}

	// get receiver address
	receiverAddr, err := t.addrRepo.GetOneUnAllocated(receiver)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call addrRepo.GetOneUnAllocated(receiver)")
	}

	// create waw transaction until amount is reached to total
	targetAmount := new(big.Int).SetInt64(total.Int64())

	serializedTxs := make([]string, 0, len(userAmounts))
	txDetailItems := make([]*models.EthDetailTX, 0, len(userAmounts))
	for _, userAmount := range userAmounts {
		var value uint64
		if userAmount.Amount > targetAmount.Uint64() {
			value = targetAmount.Uint64()
		} else {
			value = userAmount.Amount
		}
		targetAmount = new(big.Int).Sub(targetAmount, new(big.Int).SetUint64(value))
		// call CreateRawTransaction
		rawTx, txDetailItem, err := t.eth.CreateRawTransaction(userAmount.Address, receiverAddr.WalletAddress, value)
		if err != nil {
			return "", "", errors.Wrapf(err, "fail to call addrRepo.CreateRawTransaction(), sender address: %s", userAmount.Address)
		}

		rawTxHex := rawTx.TxHex
		t.logger.Debug("rawTxHex", zap.String("rawTxHex", rawTxHex))
		//TODO: `rawTxHex` should be used to trace progress to update database

		serializedTx, err := serial.EncodeToString(rawTx)
		if err != nil {
			return "", "", errors.Wrap(err, "fail to call serial.EncodeToString(rawTx)")
		}
		serializedTxs = append(serializedTxs, serializedTx)

		// create insert data for　eth_detail_tx
		txDetailItem.SenderAccount = sender.String()
		txDetailItem.ReceiverAccount = receiver.String()
		txDetailItems = append(txDetailItems, txDetailItem)

		if targetAmount.Uint64() == 0 {
			break
		}
	}

	return t.afterTxCreation(targetAction, sender, serializedTxs, txDetailItems)
}
