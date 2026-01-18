package bitcoin

import (
	"errors"
	"time"
)

// BTCTxOutput represents a Bitcoin transaction output in the domain layer
type BTCTxOutput struct {
	ID            int64
	TxID          int64
	OutputAddress string
	OutputAccount string
	OutputAmount  string
	IsChange      bool
	UpdatedAt     *time.Time
}

// NewBTCTxOutput creates a new BTCTxOutput entity
func NewBTCTxOutput(
	txID int64,
	outputAddress string,
	outputAccount string,
	outputAmount string,
	isChange bool,
) (*BTCTxOutput, error) {
	if outputAddress == "" {
		return nil, errors.New("output address cannot be empty")
	}
	if outputAmount == "" {
		return nil, errors.New("output amount cannot be empty")
	}

	now := time.Now()
	return &BTCTxOutput{
		TxID:          txID,
		OutputAddress: outputAddress,
		OutputAccount: outputAccount,
		OutputAmount:  outputAmount,
		IsChange:      isChange,
		UpdatedAt:     &now,
	}, nil
}
