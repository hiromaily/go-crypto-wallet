package watch

import (
	domainBitcoin "github.com/hiromaily/go-crypto-wallet/internal/domain/bitcoin"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
)

// BTCTxRepositorier is BTCTxRepository interface
type BTCTxRepositorier interface {
	GetOne(id int64) (*domainBitcoin.BTCTransaction, error)
	GetCountByUnsignedHex(actionType domainTx.ActionType, hex string) (int64, error)
	GetTxIDBySentHash(actionType domainTx.ActionType, hash string) (int64, error)
	GetSentHashTx(actionType domainTx.ActionType, txType domainTx.TxType) ([]string, error)
	InsertUnsignedTx(actionType domainTx.ActionType, txItem *domainBitcoin.BTCTransaction) (int64, error)
	Update(txItem *domainBitcoin.BTCTransaction) (int64, error)
	UpdateAfterTxSent(txID int64, txType domainTx.TxType, signedHex, sentHashTx string) (int64, error)
	UpdateTxType(id int64, txType domainTx.TxType) (int64, error)
	UpdateTxTypeBySentHashTx(actionType domainTx.ActionType, txType domainTx.TxType, sentHashTx string) (int64, error)
	DeleteAll() (int64, error)
}

// TxInputRepositorier is TxInputRepository interface
type TxInputRepositorier interface {
	GetOne(id int64) (*domainBitcoin.BTCTxInput, error)
	GetAllByTxID(id int64) ([]*domainBitcoin.BTCTxInput, error)
	Insert(txItem *domainBitcoin.BTCTxInput) error
	InsertBulk(txItems []*domainBitcoin.BTCTxInput) error
}

// TxOutputRepositorier is TxOutputRepository interface
type TxOutputRepositorier interface {
	GetOne(id int64) (*domainBitcoin.BTCTxOutput, error)
	GetAllByTxID(id int64) ([]*domainBitcoin.BTCTxOutput, error)
	Insert(txItem *domainBitcoin.BTCTxOutput) error
	InsertBulk(txItems []*domainBitcoin.BTCTxOutput) error
}
