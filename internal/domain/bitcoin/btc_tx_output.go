package bitcoin

import (
	"errors"
	"time"
)

// BtcTxOutput represents a Bitcoin transaction output in the domain layer
type BtcTxOutput struct {
	ID            int64
	TxID          int64
	OutputAddress string
	OutputAccount string
	OutputAmount  string
	IsChange      bool
	UpdatedAt     *time.Time
}

// NewBtcTxOutput creates a new BtcTxOutput entity
func NewBtcTxOutput(
	txID int64,
	outputAddress string,
	outputAccount string,
	outputAmount string,
	isChange bool,
) (*BtcTxOutput, error) {
	if outputAddress == "" {
		return nil, errors.New("output address cannot be empty")
	}
	if outputAmount == "" {
		return nil, errors.New("output amount cannot be empty")
	}

	now := time.Now()
	return &BtcTxOutput{
		TxID:          txID,
		OutputAddress: outputAddress,
		OutputAccount: outputAccount,
		OutputAmount:  outputAmount,
		IsChange:      isChange,
		UpdatedAt:     &now,
	}, nil
}
