package watchsrv

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/volatiletech/null"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	models "github.com/hiromaily/go-bitcoin/pkg/models/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/repository/watchrepo"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
)

// PaymentRequestCreator is PaymentRequestCreate interface
type PaymentRequestCreator interface {
	CreatePaymentRequest() error
}

// PaymentRequestCreate type
type PaymentRequestCreate struct {
	btc        api.Bitcoiner
	logger     *zap.Logger
	dbConn     *sql.DB
	addrRepo   watchrepo.AddressRepositorier
	payReqRepo watchrepo.PaymentRequestRepositorier
	wtype      wallet.WalletType
}

// NewPaymentRequestCreate returns PaymentRequestCreate object
func NewPaymentRequestCreate(
	btc api.Bitcoiner,
	logger *zap.Logger,
	dbConn *sql.DB,
	addrRepo watchrepo.AddressRepositorier,
	payReqRepo watchrepo.PaymentRequestRepositorier,
	wtype wallet.WalletType) *PaymentRequestCreate {

	return &PaymentRequestCreate{
		btc:        btc,
		logger:     logger,
		dbConn:     dbConn,
		addrRepo:   addrRepo,
		payReqRepo: payReqRepo,
		wtype:      wtype,
	}
}

// CreatePaymentRequest creates payment_request dummy data for development
func (p *PaymentRequestCreate) CreatePaymentRequest() error {
	// create payment_request table
	amtList := []float64{
		0.00001,
		0.00002,
		0.000025,
		0.000015,
		0.00003,
	}

	// get client pubkeys
	pubkeyItems, err := p.addrRepo.GetAll(account.AccountTypeClient)
	if err != nil {
		return errors.Wrap(err, "fail to call addrRepo.GetAll()")
	}
	if len(pubkeyItems) < len(amtList)*2 {
		return errors.Errorf("pubkey for client should be created at least %d", len(amtList))
	}
	// start transaction
	dtx, err := p.dbConn.Begin()
	if err != nil {
		return errors.Wrap(err, "fail to start transaction")
	}
	defer func() {
		if err != nil {
			dtx.Rollback()
		} else {
			dtx.Commit()
		}
	}()

	// delete payment request
	_, err = p.payReqRepo.DeleteAll()
	if err != nil {
		return errors.Wrap(err, "fail to call payReqRepo.DeleteAll()")
	}
	// insert payment_request
	payReqItems := make([]*models.PaymentRequest, 0, len(amtList))
	var idx int
	for _, amt := range amtList {
		payReqItems = append(payReqItems, &models.PaymentRequest{
			Coin:            p.btc.CoinTypeCode().String(),
			PaymentID:       null.NewInt64(0, false),
			SenderAddress:   pubkeyItems[0+idx].WalletAddress,
			SenderAccount:   pubkeyItems[0+idx].Account,
			ReceiverAddress: pubkeyItems[len(amtList)+idx].WalletAddress,
			Amount:          p.btc.FloatToDecimal(amt),
			IsDone:          false,
			UpdatedAt:       null.TimeFrom(time.Now()),
		})
		idx++
	}
	if err = p.payReqRepo.InsertBulk(payReqItems); err != nil {
		return errors.Wrap(err, "fail to call payReqRepo.InsertBulk()")
	}
	return nil
}
