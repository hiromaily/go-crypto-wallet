package ethwallet

import (
	"context"
	"database/sql"

	watchusecase "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/watch"
	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum"
)

// ETHWatch watch only wallet object
type ETHWatch struct {
	ETH                     ethereum.Ethereumer
	dbConn                  *sql.DB
	wtype                   domainWallet.WalletType
	createTxUseCase         watchusecase.CreateTransactionUseCase
	monitorTxUseCase        watchusecase.MonitorTransactionUseCase
	sendTxUseCase           watchusecase.SendTransactionUseCase
	importAddrUseCase       watchusecase.ImportAddressUseCase
	createPaymentReqUseCase watchusecase.CreatePaymentRequestUseCase
}

// NewETHWatch returns ETHWatch object
func NewETHWatch(
	eth ethereum.Ethereumer,
	dbConn *sql.DB,
	createTxUseCase watchusecase.CreateTransactionUseCase,
	monitorTxUseCase watchusecase.MonitorTransactionUseCase,
	sendTxUseCase watchusecase.SendTransactionUseCase,
	importAddrUseCase watchusecase.ImportAddressUseCase,
	createPaymentReqUseCase watchusecase.CreatePaymentRequestUseCase,
	walletType domainWallet.WalletType,
) *ETHWatch {
	return &ETHWatch{
		ETH:                     eth,
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
func (w *ETHWatch) ImportAddress(fileName string, _ bool) error {
	return w.importAddrUseCase.Execute(context.Background(), watchusecase.ImportAddressInput{
		FileName: fileName,
		Rescan:   false, // ETH doesn't support rescan
	})
}

// createTx is a helper method to reduce code duplication across transaction creation methods
func (w *ETHWatch) createTx(input watchusecase.CreateTransactionInput) (string, string, error) {
	output, err := w.createTxUseCase.Execute(context.Background(), input)
	if err != nil {
		return "", "", err
	}
	return output.TransactionHex, output.FileName, nil
}

// CreateDepositTx creates deposit unsigned transaction
func (w *ETHWatch) CreateDepositTx(_ float64) (string, string, error) {
	return w.createTx(watchusecase.CreateTransactionInput{
		ActionType: domainTx.ActionTypeDeposit.String(),
	})
}

// CreatePaymentTx creates payment unsigned transaction
func (w *ETHWatch) CreatePaymentTx(_ float64) (string, string, error) {
	return w.createTx(watchusecase.CreateTransactionInput{
		ActionType: domainTx.ActionTypePayment.String(),
	})
}

// CreateTransferTx creates transfer unsigned transaction
func (w *ETHWatch) CreateTransferTx(
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
func (w *ETHWatch) UpdateTxStatus() error {
	return w.monitorTxUseCase.UpdateTxStatus(context.Background())
}

// MonitorBalance monitors balance
func (w *ETHWatch) MonitorBalance(confirmationNum uint64) error {
	return w.monitorTxUseCase.MonitorBalance(context.Background(), watchusecase.MonitorBalanceInput{
		ConfirmationNum: confirmationNum,
	})
}

// SendTx sends signed transaction
func (w *ETHWatch) SendTx(filePath string) (string, error) {
	output, err := w.sendTxUseCase.Execute(context.Background(), watchusecase.SendTransactionInput{
		FilePath: filePath,
	})
	if err != nil {
		return "", err
	}
	return output.TxID, nil
}

// CreatePaymentRequest creates payment_request dummy data for development
func (w *ETHWatch) CreatePaymentRequest() error {
	amtList := []float64{
		0.001,
		0.002,
		0.0025,
		0.0015,
		0.003,
	}
	return w.createPaymentReqUseCase.Execute(context.Background(), watchusecase.CreatePaymentRequestInput{
		AmountList: amtList,
	})
}

// Done should be called before exit
func (w *ETHWatch) Done() {
	_ = w.dbConn.Close() // Best effort cleanup
	w.ETH.Close()
}

// CoinTypeCode returns domainCoin.CoinTypeCode
func (w *ETHWatch) CoinTypeCode() domainCoin.CoinTypeCode {
	return w.ETH.CoinTypeCode()
}
