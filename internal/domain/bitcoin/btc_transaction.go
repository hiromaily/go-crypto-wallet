package bitcoin

import (
	"time"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
)

// BTCTransaction represents a Bitcoin transaction in the domain layer
type BTCTransaction struct {
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

// NewBTCTransaction creates a new BtcTransaction entity
func NewBTCTransaction(
	coinTypeCode domainCoin.CoinTypeCode,
	actionType domainTx.ActionType,
	currentTxType domainTx.TxType,
) *BTCTransaction {
	return &BTCTransaction{
		CoinTypeCode:  coinTypeCode,
		ActionType:    actionType,
		CurrentTxType: currentTxType,
	}
}

// SetUnsignedTx sets the unsigned transaction data
func (t *BTCTransaction) SetUnsignedTx(hexTx, totalInput, totalOutput, fee string) {
	t.UnsignedHexTx = hexTx
	t.TotalInputAmount = totalInput
	t.TotalOutputAmount = totalOutput
	t.Fee = fee
	now := time.Now()
	t.UnsignedUpdatedAt = &now
}

// SetSignedTx sets the signed transaction data
func (t *BTCTransaction) SetSignedTx(signedHex, sentHash string) {
	t.SignedHexTx = signedHex
	t.SentHashTx = sentHash
	now := time.Now()
	t.SentUpdatedAt = &now
}

// UpdateTxType updates the transaction type
func (t *BTCTransaction) UpdateTxType(txType domainTx.TxType) {
	t.CurrentTxType = txType
}
