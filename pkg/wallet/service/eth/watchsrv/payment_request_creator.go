package watchsrv

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/pkg/errors"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/types"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	models "github.com/hiromaily/go-bitcoin/pkg/models/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/repository/watchrepo"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/ethgrp"
)

// PaymentRequestCreate type
type PaymentRequestCreate struct {
	eth        ethgrp.Ethereumer
	logger     *zap.Logger
	dbConn     *sql.DB
	addrRepo   watchrepo.AddressRepositorier
	payReqRepo watchrepo.PaymentRequestRepositorier
	wtype      wallet.WalletType
}

// NewPaymentRequestCreate returns PaymentRequestCreate object
func NewPaymentRequestCreate(
	eth ethgrp.Ethereumer,
	logger *zap.Logger,
	dbConn *sql.DB,
	addrRepo watchrepo.AddressRepositorier,
	payReqRepo watchrepo.PaymentRequestRepositorier,
	wtype wallet.WalletType) *PaymentRequestCreate {

	return &PaymentRequestCreate{
		eth:        eth,
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
		0.001,
		0.002,
		0.0025,
		0.0015,
		0.003,
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
			Coin:            p.eth.CoinTypeCode().String(),
			PaymentID:       null.NewInt64(0, false),
			SenderAddress:   pubkeyItems[0+idx].WalletAddress,
			SenderAccount:   pubkeyItems[0+idx].Account,
			ReceiverAddress: pubkeyItems[len(amtList)+idx].WalletAddress,
			Amount:          p.floatToDecimal(amt),
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

// FloatToDecimal converts float to decimal
func (p *PaymentRequestCreate) floatToDecimal(f float64) types.Decimal {
	strAmt := fmt.Sprintf("%f", f)
	dAmt := types.Decimal{Big: new(decimal.Big)}
	dAmt.Big, _ = dAmt.SetString(strAmt)
	return dAmt
}
