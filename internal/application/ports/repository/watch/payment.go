package watch

import (
	"database/sql"

	domainPayment "github.com/hiromaily/go-crypto-wallet/internal/domain/payment"
)

// PaymentRequestRepositorier is PaymentRequestRepository interface
type PaymentRequestRepositorier interface {
	GetAll() ([]*domainPayment.PaymentRequest, error)
	GetAllByPaymentID(paymentID int64) ([]*domainPayment.PaymentRequest, error)
	InsertBulk(items []*domainPayment.PaymentRequest) error
	UpdatePaymentID(paymentID int64, ids []int64) (int64, error)
	UpdateIsDone(paymentID int64) (int64, error)
	DeleteAll() (int64, error)
	WithTx(tx *sql.Tx) PaymentRequestRepositorier
}
