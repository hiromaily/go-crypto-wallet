package xrpwallet

import (
	"context"
	"database/sql"

	watchusecase "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/watch"
	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ripple"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// XRPWatch watch only wallet object
type XRPWatch struct {
	XRP                     ripple.Rippler
	dbConn                  *sql.DB
	wtype                   domainWallet.WalletType
	createTxUseCase         watchusecase.CreateTransactionUseCase
	monitorTxUseCase        watchusecase.MonitorTransactionUseCase
	sendTxUseCase           watchusecase.SendTransactionUseCase
	importAddrUseCase       watchusecase.ImportAddressUseCase
	createPaymentReqUseCase watchusecase.CreatePaymentRequestUseCase
}

// NewXRPWatch returns XRPWatch object
func NewXRPWatch(
	xrp ripple.Rippler,
	dbConn *sql.DB,
	createTxUseCase watchusecase.CreateTransactionUseCase,
	monitorTxUseCase watchusecase.MonitorTransactionUseCase,
	sendTxUseCase watchusecase.SendTransactionUseCase,
	importAddrUseCase watchusecase.ImportAddressUseCase,
	createPaymentReqUseCase watchusecase.CreatePaymentRequestUseCase,
	walletType domainWallet.WalletType,
) *XRPWatch {
	return &XRPWatch{
		XRP:                     xrp,
		dbConn:                  dbConn,
		wtype:                   walletType,
		createTxUseCase:         createTxUseCase,
		monitorTxUseCase:        monitorTxUseCase,
		sendTxUseCase:           sendTxUseCase,
		importAddrUseCase:       importAddrUseCase,
		createPaymentReqUseCase: createPaymentReqUseCase,
	}
}

// ImportAddress imports address
func (w *XRPWatch) ImportAddress(fileName string, _ bool) error {
	return w.importAddrUseCase.Execute(context.Background(), watchusecase.ImportAddressInput{
		FileName: fileName,
		Rescan:   false, // XRP doesn't support rescan
	})
}

// createTx is a helper method to reduce code duplication across transaction creation methods
func (w *XRPWatch) createTx(input watchusecase.CreateTransactionInput) (string, string, error) {
	output, err := w.createTxUseCase.Execute(context.Background(), input)
	if err != nil {
		return "", "", err
	}
	return output.TransactionHex, output.FileName, nil
}

// CreateDepositTx creates deposit unsigned transaction
func (w *XRPWatch) CreateDepositTx(_ float64) (string, string, error) {
	return w.createTx(watchusecase.CreateTransactionInput{
		ActionType: domainTx.ActionTypeDeposit.String(),
	})
}

// CreatePaymentTx creates payment unsigned transaction
func (w *XRPWatch) CreatePaymentTx(_ float64) (string, string, error) {
	return w.createTx(watchusecase.CreateTransactionInput{
		ActionType: domainTx.ActionTypePayment.String(),
	})
}

// CreateTransferTx creates transfer unsigned transaction
func (w *XRPWatch) CreateTransferTx(
	sender, receiver domainAccount.AccountType, floatAmount, _ float64,
) (string, string, error) {
	return w.createTx(watchusecase.CreateTransactionInput{
		ActionType:      domainTx.ActionTypeTransfer.String(),
		SenderAccount:   sender,
		ReceiverAccount: receiver,
		Amount:          floatAmount,
	})
}

// UpdateTxStatus updates transaction status
func (*XRPWatch) UpdateTxStatus() error {
	logger.Info("no functionality for XRP")
	return nil
}

// MonitorBalance monitors balance
func (w *XRPWatch) MonitorBalance(confirmationNum uint64) error {
	return w.monitorTxUseCase.MonitorBalance(context.Background(), watchusecase.MonitorBalanceInput{
		ConfirmationNum: confirmationNum,
	})
}

// SendTx sends signed transaction
func (w *XRPWatch) SendTx(filePath string) (string, error) {
	output, err := w.sendTxUseCase.Execute(context.Background(), watchusecase.SendTransactionInput{
		FilePath: filePath,
	})
	if err != nil {
		return "", err
	}
	return output.TxID, nil
}

// CreatePaymentRequest creates payment_request dummy data for development
func (w *XRPWatch) CreatePaymentRequest() error {
	amtList := []float64{
		50,
		100,
		120,
		130,
		150,
	}
	return w.createPaymentReqUseCase.Execute(context.Background(), watchusecase.CreatePaymentRequestInput{
		AmountList: amtList,
	})
}

// Done should be called before exit
func (w *XRPWatch) Done() {
	_ = w.dbConn.Close() // Best effort cleanup

	_ = w.XRP.Close() // Best effort cleanup
}

// CoinTypeCode returns domainCoin.CoinTypeCode
func (w *XRPWatch) CoinTypeCode() domainCoin.CoinTypeCode {
	return w.XRP.CoinTypeCode()
}
