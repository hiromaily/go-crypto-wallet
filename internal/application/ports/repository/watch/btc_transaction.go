package watch

import (
	domainBitcoin "github.com/hiromaily/go-crypto-wallet/internal/domain/bitcoin"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
)

// BTCTxRepositorier is BTCTxRepository interface
type BTCTxRepositorier interface {
	GetOne(id int64) (*domainBitcoin.BtcTransaction, error)
	GetCountByUnsignedHex(actionType domainTx.ActionType, hex string) (int64, error)
	GetTxIDBySentHash(actionType domainTx.ActionType, hash string) (int64, error)
	GetSentHashTx(actionType domainTx.ActionType, txType domainTx.TxType) ([]string, error)
	InsertUnsignedTx(actionType domainTx.ActionType, txItem *domainBitcoin.BtcTransaction) (int64, error)
	Update(txItem *domainBitcoin.BtcTransaction) (int64, error)
	UpdateAfterTxSent(txID int64, txType domainTx.TxType, signedHex, sentHashTx string) (int64, error)
	UpdateTxType(id int64, txType domainTx.TxType) (int64, error)
	UpdateTxTypeBySentHashTx(actionType domainTx.ActionType, txType domainTx.TxType, sentHashTx string) (int64, error)
	DeleteAll() (int64, error)
}

// TxInputRepositorier is TxInputRepository interface
type TxInputRepositorier interface {
	GetOne(id int64) (*domainBitcoin.BtcTxInput, error)
	GetAllByTxID(id int64) ([]*domainBitcoin.BtcTxInput, error)
	Insert(txItem *domainBitcoin.BtcTxInput) error
	InsertBulk(txItems []*domainBitcoin.BtcTxInput) error
}

// TxOutputRepositorier is TxOutputRepository interface
type TxOutputRepositorier interface {
	GetOne(id int64) (*domainBitcoin.BtcTxOutput, error)
	GetAllByTxID(id int64) ([]*domainBitcoin.BtcTxOutput, error)
	Insert(txItem *domainBitcoin.BtcTxOutput) error
	InsertBulk(txItems []*domainBitcoin.BtcTxOutput) error
}
