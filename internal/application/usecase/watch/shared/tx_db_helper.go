package shared

import (
	"database/sql"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"

	repowatch "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/watch"
	domainBTC "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/btc"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	btcpkg "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc"
)

// TxAmountConverter provides the coin type needed for DB operations.
type TxAmountConverter interface {
	CoinTypeCode() domainCoin.CoinTypeCode
}

// TxDBHelper provides helper methods for transaction database operations.
// This struct is shared between BTC and BCH use cases.
type TxDBHelper struct {
	dbConn       *sql.DB
	txRepo       repowatch.BTCTxRepositorier
	txInputRepo  repowatch.TxInputRepositorier
	txOutputRepo repowatch.TxOutputRepositorier
	payReqRepo   repowatch.PaymentRequestRepositorier
	converter    TxAmountConverter
}

// NewTxDBHelper creates a new TxDBHelper.
func NewTxDBHelper(
	dbConn *sql.DB,
	txRepo repowatch.BTCTxRepositorier,
	txInputRepo repowatch.TxInputRepositorier,
	txOutputRepo repowatch.TxOutputRepositorier,
	payReqRepo repowatch.PaymentRequestRepositorier,
	converter TxAmountConverter,
) *TxDBHelper {
	return &TxDBHelper{
		dbConn:       dbConn,
		txRepo:       txRepo,
		txInputRepo:  txInputRepo,
		txOutputRepo: txOutputRepo,
		payReqRepo:   payReqRepo,
		converter:    converter,
	}
}

// InsertTxTableForUnsigned inserts unsigned transaction data into the database.
// This method is shared between BTC and BCH use cases.
//
// Returns:
//   - txID: The ID of the inserted transaction (0 if already exists)
//   - error: Any error that occurred
func (h *TxDBHelper) InsertTxTableForUnsigned(
	actionType domainTx.ActionType,
	hex string,
	inputTotal,
	outputTotal,
	fee btcutil.Amount,
	txInputs []*domainBTC.BTCTxInput,
	txOutputs []*domainBTC.BTCTxOutput,
	paymentRequestIds []int64,
) (int64, error) {
	// skip if same hex is already stored
	count, err := h.txRepo.GetCountByUnsignedHex(actionType, hex)
	if err != nil {
		return 0, fmt.Errorf("fail to call repo.Tx().GetCountByUnsignedHex(): %w", err)
	}
	if count != 0 {
		// skip
		return 0, nil
	}

	// TxReceipt table
	totalInputAmt, err := btcpkg.AmountToDecimal(inputTotal)
	if err != nil {
		return 0, fmt.Errorf("fail to convert total input amount to decimal: %w", err)
	}
	totalOutputAmt, err := btcpkg.AmountToDecimal(outputTotal)
	if err != nil {
		return 0, fmt.Errorf("fail to convert total output amount to decimal: %w", err)
	}
	feeAmt, err := btcpkg.AmountToDecimal(fee)
	if err != nil {
		return 0, fmt.Errorf("fail to convert fee amount to decimal: %w", err)
	}
	txItem := domainBTC.NewBTCTransaction(
		h.converter.CoinTypeCode(),
		actionType,
		domainTx.TxTypeUnsigned,
	)
	txItem.SetUnsignedTx(hex, totalInputAmt.String(), totalOutputAmt.String(), feeAmt.String())

	// start database transaction
	dtx, err := h.dbConn.Begin()
	if err != nil {
		return 0, fmt.Errorf("fail to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = dtx.Rollback() // Error already being handled
		} else {
			_ = dtx.Commit() // Error already being handled
		}
	}()

	txID, err := h.txRepo.InsertUnsignedTx(actionType, txItem)
	if err != nil {
		return 0, fmt.Errorf("fail to call repo.Tx().InsertUnsignedTx(): %w", err)
	}

	// TxReceiptInput table
	//  update txID
	for idx := range txInputs {
		txInputs[idx].TxID = txID
	}
	err = h.txInputRepo.InsertBulk(txInputs)
	if err != nil {
		return 0, fmt.Errorf("fail to call txInRepo.InsertBulk(): %w", err)
	}

	// TxReceiptOutput table
	//  update txID
	for idx := range txOutputs {
		txOutputs[idx].TxID = txID
	}
	err = h.txOutputRepo.InsertBulk(txOutputs)
	if err != nil {
		return 0, fmt.Errorf("fail to call repo.TxOutput().InsertBulk(): %w", err)
	}

	// update payment_id in payment_request table for only domainTx.ActionTypePayment
	if actionType == domainTx.ActionTypePayment {
		_, err = h.payReqRepo.UpdatePaymentID(txID, paymentRequestIds)
		if err != nil {
			return 0, fmt.Errorf("fail to call repo.PayReq().UpdatePaymentID(txID, paymentRequestIds): %w", err)
		}
	}

	return txID, nil
}
