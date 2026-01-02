package bitcoin

import (
	"time"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
)

// BtcTransaction represents a Bitcoin transaction in the domain layer
type BtcTransaction struct {
	ID                int64
	CoinTypeCode      domainCoin.CoinTypeCode
	ActionType        domainTx.ActionType
	UnsignedHexTx     string
	SignedHexTx       string
	SentHashTx        string
	TotalInputAmount  string
	TotalOutputAmount string
	Fee               string
	CurrentTxType     domainTx.TxType
	UnsignedUpdatedAt *time.Time
	SentUpdatedAt     *time.Time
}

// NewBtcTransaction creates a new BtcTransaction entity
func NewBtcTransaction(
	coinTypeCode domainCoin.CoinTypeCode,
	actionType domainTx.ActionType,
	currentTxType domainTx.TxType,
) *BtcTransaction {
	return &BtcTransaction{
		CoinTypeCode:  coinTypeCode,
		ActionType:    actionType,
		CurrentTxType: currentTxType,
	}
}

// SetUnsignedTx sets the unsigned transaction data
func (t *BtcTransaction) SetUnsignedTx(hexTx, totalInput, totalOutput, fee string) {
	t.UnsignedHexTx = hexTx
	t.TotalInputAmount = totalInput
	t.TotalOutputAmount = totalOutput
	t.Fee = fee
	now := time.Now()
	t.UnsignedUpdatedAt = &now
}

// SetSignedTx sets the signed transaction data
func (t *BtcTransaction) SetSignedTx(signedHex, sentHash string) {
	t.SignedHexTx = signedHex
	t.SentHashTx = sentHash
	now := time.Now()
	t.SentUpdatedAt = &now
}

// UpdateTxType updates the transaction type
func (t *BtcTransaction) UpdateTxType(txType domainTx.TxType) {
	t.CurrentTxType = txType
}
