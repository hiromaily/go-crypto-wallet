package eth

import (
	"errors"
	"time"

	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
)

// ETHDetailTx represents an Ethereum transaction detail in the domain layer
type ETHDetailTx struct {
	ID                int64
	TxID              int64
	UUID              string
	CurrentTxType     domainTx.TxType
	SenderAccount     string
	SenderAddress     string
	ReceiverAccount   string
	ReceiverAddress   string
	Amount            uint64
	Fee               uint64
	GasLimit          uint32
	Nonce             uint64
	UnsignedHexTx     string
	SignedHexTx       string
	SentHashTx        string
	UnsignedUpdatedAt *time.Time
	SentUpdatedAt     *time.Time
}

// NewETHDetailTx creates a new ETHDetailTx entity
func NewETHDetailTx(
	txID int64,
	uuid string,
	currentTxType domainTx.TxType,
	senderAccount string,
	senderAddress string,
	receiverAccount string,
	receiverAddress string,
	amount uint64,
	fee uint64,
	gasLimit uint32,
	nonce uint64,
) (*ETHDetailTx, error) {
	if uuid == "" {
		return nil, errors.New("uuid cannot be empty")
	}
	if senderAddress == "" {
		return nil, errors.New("sender address cannot be empty")
	}
	if receiverAddress == "" {
		return nil, errors.New("receiver address cannot be empty")
	}
	// Note: amount can be zero for ERC20 token transactions where only the token is transferred

	return &ETHDetailTx{
		TxID:            txID,
		UUID:            uuid,
		CurrentTxType:   currentTxType,
		SenderAccount:   senderAccount,
		SenderAddress:   senderAddress,
		ReceiverAccount: receiverAccount,
		ReceiverAddress: receiverAddress,
		Amount:          amount,
		Fee:             fee,
		GasLimit:        gasLimit,
		Nonce:           nonce,
	}, nil
}

// SetUnsignedTx sets unsigned transaction data
func (e *ETHDetailTx) SetUnsignedTx(hexTx string) {
	e.UnsignedHexTx = hexTx
	now := time.Now()
	e.UnsignedUpdatedAt = &now
}

// SetSignedTx sets signed transaction data
func (e *ETHDetailTx) SetSignedTx(signedHex, sentHash string) {
	e.SignedHexTx = signedHex
	e.SentHashTx = sentHash
	now := time.Now()
	e.SentUpdatedAt = &now
}

// UpdateTxType updates the transaction type
func (e *ETHDetailTx) UpdateTxType(txType domainTx.TxType) {
	e.CurrentTxType = txType
}
