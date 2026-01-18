package bitcoin

import (
	"errors"
	"time"
)

// BTCTxInput represents a Bitcoin transaction input in the domain layer
type BTCTxInput struct {
	ID                 int64
	TxID               int64
	InputTxid          string
	InputVout          uint32
	InputAddress       string
	InputAccount       string
	InputAmount        string
	InputConfirmations uint64
	UpdatedAt          *time.Time
}

// NewBTCTxInput creates a new BTCTxInput entity
func NewBTCTxInput(
	txID int64,
	inputTxid string,
	inputVout uint32,
	inputAddress string,
	inputAccount string,
	inputAmount string,
	inputConfirmations uint64,
) (*BTCTxInput, error) {
	if inputTxid == "" {
		return nil, errors.New("input txid cannot be empty")
	}
	if inputAddress == "" {
		return nil, errors.New("input address cannot be empty")
	}
	if inputAmount == "" {
		return nil, errors.New("input amount cannot be empty")
	}

	now := time.Now()
	return &BTCTxInput{
		TxID:               txID,
		InputTxid:          inputTxid,
		InputVout:          inputVout,
		InputAddress:       inputAddress,
		InputAccount:       inputAccount,
		InputAmount:        inputAmount,
		InputConfirmations: inputConfirmations,
		UpdatedAt:          &now,
	}, nil
}
