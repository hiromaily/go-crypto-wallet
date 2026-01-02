package payment

import (
	"errors"
	"time"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
)

// PaymentRequest represents a payment request in the domain layer
type PaymentRequest struct {
	ID              int64
	CoinTypeCode    domainCoin.CoinTypeCode
	PaymentID       *int64
	SenderAddress   string
	SenderAccount   string
	ReceiverAddress string
	Amount          string
	IsDone          bool
	UpdatedAt       *time.Time
}

// NewPaymentRequest creates a new PaymentRequest entity
func NewPaymentRequest(
	coinTypeCode domainCoin.CoinTypeCode,
	senderAddress string,
	senderAccount string,
	receiverAddress string,
	amount string,
) (*PaymentRequest, error) {
	if senderAddress == "" {
		return nil, errors.New("sender address cannot be empty")
	}
	if receiverAddress == "" {
		return nil, errors.New("receiver address cannot be empty")
	}
	if amount == "" {
		return nil, errors.New("amount cannot be empty")
	}

	now := time.Now()
	return &PaymentRequest{
		CoinTypeCode:    coinTypeCode,
		SenderAddress:   senderAddress,
		SenderAccount:   senderAccount,
		ReceiverAddress: receiverAddress,
		Amount:          amount,
		IsDone:          false,
		UpdatedAt:       &now,
	}, nil
}

// SetPaymentID sets the payment ID
func (p *PaymentRequest) SetPaymentID(paymentID int64) {
	p.PaymentID = &paymentID
}

// MarkAsDone marks the payment request as done
func (p *PaymentRequest) MarkAsDone() {
	p.IsDone = true
	now := time.Now()
	p.UpdatedAt = &now
}
