package watchsrv

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/guregu/null/v6"

	"github.com/hiromaily/go-crypto-wallet/pkg/converter"
	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/watchrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// PaymentRequestCreate type
type PaymentRequestCreate struct {
	converter    converter.Converter
	dbConn       *sql.DB
	addrRepo     watchrepo.AddressRepositorier
	payReqRepo   watchrepo.PaymentRequestRepositorier
	coinTypeCode coin.CoinTypeCode
	wtype        domainWallet.WalletType
}

// NewPaymentRequestCreate returns PaymentRequestCreate object
func NewPaymentRequestCreate(
	conv converter.Converter,
	dbConn *sql.DB,
	addrRepo watchrepo.AddressRepositorier,
	payReqRepo watchrepo.PaymentRequestRepositorier,
	coinTypeCode coin.CoinTypeCode,
	wtype domainWallet.WalletType,
) *PaymentRequestCreate {
	return &PaymentRequestCreate{
		converter:    conv,
		dbConn:       dbConn,
		addrRepo:     addrRepo,
		payReqRepo:   payReqRepo,
		coinTypeCode: coinTypeCode,
		wtype:        wtype,
	}
}

// CreatePaymentRequest creates payment_request dummy data for development
func (p *PaymentRequestCreate) CreatePaymentRequest(amtList []float64) error {
	// get client pubkeys
	pubkeyItems, err := p.addrRepo.GetAll(domainAccount.AccountTypeClient)
	if err != nil {
		return fmt.Errorf("fail to call addrRepo.GetAll(): %w", err)
	}
	if len(pubkeyItems) < len(amtList)*2 {
		return fmt.Errorf("pubkey for client should be created at least %d", len(amtList))
	}
	// start transaction
	dtx, err := p.dbConn.Begin()
	if err != nil {
		return fmt.Errorf("fail to start transaction: %w", err)
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
		return fmt.Errorf("fail to call payReqRepo.DeleteAll(): %w", err)
	}
	// insert payment_request
	payReqItems := make([]*models.PaymentRequest, 0, len(amtList))
	var idx int
	for _, amt := range amtList {
		amount, err := p.converter.FloatToDecimal(amt)
		if err != nil {
			return fmt.Errorf("fail to convert amount %f to decimal: %w", amt, err)
		}
		payReqItems = append(payReqItems, &models.PaymentRequest{
			Coin:            p.coinTypeCode.String(),
			PaymentID:       null.Int64{},
			SenderAddress:   pubkeyItems[0+idx].WalletAddress,
			SenderAccount:   pubkeyItems[0+idx].Account,
			ReceiverAddress: pubkeyItems[len(amtList)+idx].WalletAddress,
			Amount:          amount,
			IsDone:          false,
			UpdatedAt:       null.TimeFrom(time.Now()),
		})
		idx++
	}
	if err = p.payReqRepo.InsertBulk(payReqItems); err != nil {
		return fmt.Errorf("fail to call payReqRepo.InsertBulk(): %w", err)
	}
	return nil
}
