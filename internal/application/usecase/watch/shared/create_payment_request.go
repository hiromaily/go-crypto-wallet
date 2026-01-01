package shared

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/quagmt/udecimal"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	sqlc "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch"
)

type createPaymentRequestUseCase struct {
	dbConn       *sql.DB
	addrRepo     watch.AddressRepositorier
	payReqRepo   watch.PaymentRequestRepositorier
	coinTypeCode domainCoin.CoinTypeCode
	wtype        domainWallet.WalletType
}

// NewCreatePaymentRequestUseCase creates a new CreatePaymentRequestUseCase for watch wallet
func NewCreatePaymentRequestUseCase(
	dbConn *sql.DB,
	addrRepo watch.AddressRepositorier,
	payReqRepo watch.PaymentRequestRepositorier,
	coinTypeCode domainCoin.CoinTypeCode,
	wtype domainWallet.WalletType,
) watchusecase.CreatePaymentRequestUseCase {
	return &createPaymentRequestUseCase{
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
	payReqItems := make([]*sqlc.PaymentRequest, 0, len(input.AmountList))
	var idx int
	for _, amt := range input.AmountList {
		// Convert float amount to string using decimal library for financial precision
		amount, err := udecimal.NewFromFloat64(amt)
		if err != nil {
			return fmt.Errorf("fail to convert amount %f to decimal: %w", amt, err)
		}
		payReqItems = append(payReqItems, &sqlc.PaymentRequest{
			Coin:            sqlc.PaymentRequestCoin(u.coinTypeCode.String()),
			PaymentID:       sql.NullInt64{},
			SenderAddress:   pubkeyItems[0+idx].WalletAddress,
			SenderAccount:   string(pubkeyItems[0+idx].Account),
			ReceiverAddress: pubkeyItems[len(input.AmountList)+idx].WalletAddress,
			Amount:          amount.String(),
			IsDone:          false,
			UpdatedAt:       sql.NullTime{Time: time.Now(), Valid: true},
		})
		idx++
	}
	if err = u.payReqRepo.InsertBulk(payReqItems); err != nil {
		return fmt.Errorf("fail to call payReqRepo.InsertBulk(): %w", err)
	}
	return nil
}
