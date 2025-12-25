package btcwallet

import (
	"context"
	"database/sql"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/address"
)

// BTCWatch watch only wallet object
type BTCWatch struct {
	BTC                     bitcoin.Bitcoiner
	dbConn                  *sql.DB
	addrType                address.AddrType
	wtype                   domainWallet.WalletType
	createTxUseCase         watchusecase.CreateTransactionUseCase
	monitorTxUseCase        watchusecase.MonitorTransactionUseCase
	sendTxUseCase           watchusecase.SendTransactionUseCase
	importAddrUseCase       watchusecase.ImportAddressUseCase
	createPaymentReqUseCase watchusecase.CreatePaymentRequestUseCase
}

// NewBTCWatch returns Watch object
func NewBTCWatch(
	btc bitcoin.Bitcoiner,
	dbConn *sql.DB,
	addrType address.AddrType,
	createTxUseCase watchusecase.CreateTransactionUseCase,
	monitorTxUseCase watchusecase.MonitorTransactionUseCase,
	sendTxUseCase watchusecase.SendTransactionUseCase,
	importAddrUseCase watchusecase.ImportAddressUseCase,
	createPaymentReqUseCase watchusecase.CreatePaymentRequestUseCase,
	walletType domainWallet.WalletType,
) *BTCWatch {
	return &BTCWatch{
		BTC:                     btc,
		dbConn:                  dbConn,
		addrType:                addrType,
		wtype:                   walletType,
		createTxUseCase:         createTxUseCase,
		monitorTxUseCase:        monitorTxUseCase,
		sendTxUseCase:           sendTxUseCase,
		importAddrUseCase:       importAddrUseCase,
		createPaymentReqUseCase: createPaymentReqUseCase,
	}
}

// ImportAddress imports address
func (w *BTCWatch) ImportAddress(fileName string, isRescan bool) error {
	return w.importAddrUseCase.Execute(context.Background(), watchusecase.ImportAddressInput{
		FileName: fileName,
		Rescan:   isRescan,
	})
}

// createTx is a helper method to reduce code duplication across transaction creation methods
func (w *BTCWatch) createTx(input watchusecase.CreateTransactionInput) (string, string, error) {
	output, err := w.createTxUseCase.Execute(context.Background(), input)
	if err != nil {
		return "", "", err
	}
	return output.TransactionHex, output.FileName, nil
}

// CreateDepositTx creates deposit unsigned transaction
func (w *BTCWatch) CreateDepositTx(adjustmentFee float64) (string, string, error) {
	return w.createTx(watchusecase.CreateTransactionInput{
		ActionType:    domainTx.ActionTypeDeposit.String(),
		AdjustmentFee: adjustmentFee,
	})
}

// CreatePaymentTx creates payment unsigned transaction
func (w *BTCWatch) CreatePaymentTx(adjustmentFee float64) (string, string, error) {
	return w.createTx(watchusecase.CreateTransactionInput{
		ActionType:    domainTx.ActionTypePayment.String(),
		AdjustmentFee: adjustmentFee,
	})
}

// CreateTransferTx creates transfer unsigned transaction
func (w *BTCWatch) CreateTransferTx(
	sender, receiver domainAccount.AccountType, floatAmount, adjustmentFee float64,
) (string, string, error) {
	return w.createTx(watchusecase.CreateTransactionInput{
		ActionType:      domainTx.ActionTypeTransfer.String(),
		SenderAccount:   sender,
		ReceiverAccount: receiver,
		Amount:          floatAmount,
		AdjustmentFee:   adjustmentFee,
	})
}

// UpdateTxStatus updates transaction status
func (w *BTCWatch) UpdateTxStatus() error {
	return w.monitorTxUseCase.UpdateTxStatus(context.Background())
}

// MonitorBalance monitors balance
func (w *BTCWatch) MonitorBalance(confirmationNum uint64) error {
	return w.monitorTxUseCase.MonitorBalance(context.Background(), watchusecase.MonitorBalanceInput{
		ConfirmationNum: confirmationNum,
	})
}

// SendTx sends signed transaction
func (w *BTCWatch) SendTx(filePath string) (string, error) {
	output, err := w.sendTxUseCase.Execute(context.Background(), watchusecase.SendTransactionInput{
		FilePath: filePath,
	})
	if err != nil {
		return "", err
	}
	return output.TxID, nil
}

// CreatePaymentRequest creates payment_request dummy data for development
func (w *BTCWatch) CreatePaymentRequest() error {
	amtList := []float64{
		0.00001,
		0.00002,
		0.000025,
		0.000015,
		0.00003,
	}

	return w.createPaymentReqUseCase.Execute(context.Background(), watchusecase.CreatePaymentRequestInput{
		AmountList: amtList,
	})
}

// Done should be called before exit
func (w *BTCWatch) Done() {
	_ = w.dbConn.Close() // Best effort cleanup
	w.BTC.Close()
}

// CoinTypeCode returns domainCoin.CoinTypeCode
func (w *BTCWatch) CoinTypeCode() domainCoin.CoinTypeCode {
	return w.BTC.CoinTypeCode()
}

// GetBTC gets btc
// func (w *BTCWatch) GetBTC() bitcoin.Bitcoiner {
//	return w.BTC
//}
