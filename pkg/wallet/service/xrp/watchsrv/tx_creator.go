package watchsrv

import (
	"database/sql"
	"fmt"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/repository/watch"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/storage/file"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/uuid"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ripple"
)

// TxCreator is TxCreator interface
type TxCreator interface {
	CreateDepositTx() (string, string, error)
	CreatePaymentTx() (string, string, error)
	CreateTransferTx(sender, receiver domainAccount.AccountType, floatAmount float64) (string, string, error)
}

// TxCreate type
type TxCreate struct {
	xrp             ripple.Rippler
	dbConn          *sql.DB
	uuidHandler     uuid.UUIDHandler
	addrRepo        watch.AddressRepositorier
	txRepo          watch.TxRepositorier
	txDetailRepo    watch.XrpDetailTxRepositorier
	payReqRepo      watch.PaymentRequestRepositorier
	txFileRepo      file.TransactionFileRepositorier
	depositReceiver domainAccount.AccountType
	paymentSender   domainAccount.AccountType
	wtype           domainWallet.WalletType
}

// NewTxCreate returns TxCreate object
func NewTxCreate(
	xrp ripple.Rippler,
	dbConn *sql.DB,
	uuidHandler uuid.UUIDHandler,
	addrRepo watch.AddressRepositorier,
	txRepo watch.TxRepositorier,
	txDetailRepo watch.XrpDetailTxRepositorier,
	payReqRepo watch.PaymentRequestRepositorier,
	txFileRepo file.TransactionFileRepositorier,
	depositReceiver domainAccount.AccountType,
	paymentSender domainAccount.AccountType,
	wtype domainWallet.WalletType,
) *TxCreate {
	return &TxCreate{
		xrp:             xrp,
		dbConn:          dbConn,
		uuidHandler:     uuidHandler,
		addrRepo:        addrRepo,
		txRepo:          txRepo,
		txDetailRepo:    txDetailRepo,
		payReqRepo:      payReqRepo,
		txFileRepo:      txFileRepo,
		depositReceiver: depositReceiver,
		paymentSender:   paymentSender,
		wtype:           wtype,
	}
}

// TODO: it can be commonized to ./pkg/wallet/service/eth/watchsrv/tx_creator.go
func (t *TxCreate) updateDB(
	targetAction domainTx.ActionType,
	txDetailItems []*models.XRPDetailTX,
	paymentRequestIds []int64,
) (int64, error) {
	// start transaction
	dtx, err := t.dbConn.Begin()
	if err != nil {
		return 0, fmt.Errorf("fail to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			dtx.Rollback()
		} else {
			dtx.Commit()
		}
	}()

	// Insert tx
	txID, err := t.txRepo.InsertUnsignedTx(targetAction)
	if err != nil {
		return 0, fmt.Errorf("fail to call txRepo.InsertUnsignedTx(): %w", err)
	}
	// Insert to eth_detail_tx
	for idx := range txDetailItems {
		txDetailItems[idx].TXID = txID
	}
	if err = t.txDetailRepo.InsertBulk(txDetailItems); err != nil {
		return 0, fmt.Errorf("fail to call txDetailRepo.InsertBulk(): %w", err)
	}

	if targetAction == domainTx.ActionTypePayment {
		_, err = t.payReqRepo.UpdatePaymentID(txID, paymentRequestIds)
		if err != nil {
			return 0, fmt.Errorf("fail to call repo.PayReq().UpdatePaymentID(txID, paymentRequestIds): %w", err)
		}
	}
	return txID, nil
}

// generateHexFile generate file for hex txID and encoded previous addresses
// TODO: it can be commonized to ./pkg/wallet/service/eth/watchsrv/tx_creator.go
func (t *TxCreate) generateHexFile(
	actionType domainTx.ActionType, senderAccount domainAccount.AccountType, txID int64, serializedTxs []string,
) (string, error) {
	// add senderAccount to first line
	serializedTxs = append([]string{senderAccount.String()}, serializedTxs...)

	// create file
	path := t.txFileRepo.CreateFilePath(actionType, domainTx.TxTypeUnsigned, txID, 0)
	generatedFileName, err := t.txFileRepo.WriteFileSlice(path, serializedTxs)
	if err != nil {
		return "", fmt.Errorf("fail to call txFileRepo.WriteFile(): %w", err)
	}

	return generatedFileName, nil
}
