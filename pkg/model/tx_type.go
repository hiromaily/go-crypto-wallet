package model

//import (
//	"time"
//)
//
////TxType tx_typeテーブル
//type TxType struct {
//	ID        uint8      `db:"id"`
//	Type      string     `db:"type"`
//	UpdatedAt *time.Time `db:"updated_at"`
//}
//
//// GetTxTypeAll tx_typeテーブル全体を返す
//func (m *DB) GetTxTypeAll() ([]TxType, error) {
//	var txTypes []TxType
//	err := m.RDB.Select(&txTypes, "SELECT * FROM tx_type")
//
//	return txTypes, err
//}
