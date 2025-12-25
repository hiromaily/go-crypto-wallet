package shared

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/guregu/null/v6"

	watchusecase "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/watch"
	"github.com/hiromaily/go-crypto-wallet/pkg/converter"
	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/repository/watch"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
)

type createPaymentRequestUseCase struct {
	converter    converter.Converter
	dbConn       *sql.DB
	addrRepo     watch.AddressRepositorier
	payReqRepo   watch.PaymentRequestRepositorier
	coinTypeCode domainCoin.CoinTypeCode
	wtype        domainWallet.WalletType
}

// NewCreatePaymentRequestUseCase creates a new CreatePaymentRequestUseCase for watch wallet
func NewCreatePaymentRequestUseCase(
	conv converter.Converter,
	dbConn *sql.DB,
	addrRepo watch.AddressRepositorier,
	payReqRepo watch.PaymentRequestRepositorier,
	coinTypeCode domainCoin.CoinTypeCode,
	wtype domainWallet.WalletType,
) watchusecase.CreatePaymentRequestUseCase {
	return &createPaymentRequestUseCase{
		converter:    conv,
		dbConn:       dbConn,
		addrRepo:     addrRepo,
		payReqRepo:   payReqRepo,
		coinTypeCode: coinTypeCode,
		wtype:        wtype,
	}
}

func (u *createPaymentRequestUseCase) Execute(ctx context.Context, input watchusecase.CreatePaymentRequestInput) error {
	// get client pubkeys
	pubkeyItems, err := u.addrRepo.GetAll(domainAccount.AccountTypeClient)
	if err != nil {
		return fmt.Errorf("fail to call addrRepo.GetAll(): %w", err)
	}
	if len(pubkeyItems) < len(input.AmountList)*2 {
		return fmt.Errorf("pubkey for client should be created at least %d", len(input.AmountList))
	}

	// start transaction
	dtx, err := u.dbConn.Begin()
	if err != nil {
		return fmt.Errorf("fail to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = dtx.Rollback() // Error already being handled
		} else {
			_ = dtx.Commit() // Error already being handled
		}
	}()

	// delete payment request
	_, err = u.payReqRepo.DeleteAll()
	if err != nil {
		return fmt.Errorf("fail to call payReqRepo.DeleteAll(): %w", err)
	}

	// insert payment_request
	payReqItems := make([]*models.PaymentRequest, 0, len(input.AmountList))
	var idx int
	for _, amt := range input.AmountList {
		amount, err := u.converter.FloatToDecimal(amt)
		if err != nil {
			return fmt.Errorf("fail to convert amount %f to decimal: %w", amt, err)
		}
		payReqItems = append(payReqItems, &models.PaymentRequest{
			Coin:            u.coinTypeCode.String(),
			PaymentID:       null.Int64{},
			SenderAddress:   pubkeyItems[0+idx].WalletAddress,
			SenderAccount:   pubkeyItems[0+idx].Account,
			ReceiverAddress: pubkeyItems[len(input.AmountList)+idx].WalletAddress,
			Amount:          amount,
			IsDone:          false,
			UpdatedAt:       null.TimeFrom(time.Now()),
		})
		idx++
	}
	if err = u.payReqRepo.InsertBulk(payReqItems); err != nil {
		return fmt.Errorf("fail to call payReqRepo.InsertBulk(): %w", err)
	}
	return nil
}
