package model

//tx_table.go側に集約

//
//import (
//	"github.com/hiromaily/go-bitcoin/pkg/enum"
//	"github.com/jmoiron/sqlx"
//)
//
////enum.Actionに応じて、テーブルを切り替える
//
//const (
//	tableNamePayment = "tx_payment"
//)
//
//// TableNamePayment tx_paymentテーブル名を返す
//func (m *DB) TableNamePayment() string {
//	return tableNamePayment
//}
//
//// GetTxPaymentByID TxReceiptテーブルから該当するIDのレコードを返す
//func (m *DB) GetTxPaymentByID(id int64) (*TxTable, error) {
//	return m.getTxReceiptByID(m.TableNamePayment(), id)
//}
//
//// GetTxPaymentCountByUnsignedHex unsigned_hex_txをキーとしてレコード数を取得する
//func (m *DB) GetTxPaymentCountByUnsignedHex(hex string) (int64, error) {
//	return m.getTxReceiptCountByUnsignedHex(m.TableNamePayment(), hex)
//}
//
//// GetTxPaymentIDBySentHash sent_hash_txをキーとしてpayment_idを取得する
//func (m *DB) GetTxPaymentIDBySentHash(hash string) (int64, error) {
//	return m.getTxReceiptIDBySentHash(m.TableNamePayment(), hash)
//}
//
//// GetSentTxHashOnTxPaymentByTxTypeSent TxPaymentテーブルから送信済ステータスであるsent_hash_txの配列を返す
//func (m *DB) GetSentTxHashOnTxPaymentByTxTypeSent() ([]string, error) {
//	txTypeValue := enum.TxTypeValue[enum.TxTypeSent]
//	return m.getSentTxHashOnTxReceipt(m.TableNamePayment(), txTypeValue)
//}
//
//// InsertTxPaymentForUnsigned TxReceiptテーブルに未署名トランザクションレコードを作成する
//func (m *DB) InsertTxPaymentForUnsigned(txReceipt *TxTable, tx *sqlx.Tx, isCommit bool) (int64, error) {
//	return m.insertTxReceiptForUnsigned(m.TableNamePayment(), txReceipt, tx, isCommit)
//}
//
//// UpdateTxPaymentForSent TxReceiptテーブルのsigned_hex_tx, sent_hash_txを更新する
//func (m *DB) UpdateTxPaymentForSent(txReceipt *TxTable, tx *sqlx.Tx, isCommit bool) (int64, error) {
//	return m.updateTxReceiptForSent(m.TableNamePayment(), txReceipt, tx, isCommit)
//}
//
//// UpdateTxPaymentDoneByTxHash TxReceiptテーブルの該当するsent_hash_txのレコードのcurrnt_tx_typeを更新する
//func (m *DB) UpdateTxPaymentDoneByTxHash(hash string, tx *sqlx.Tx, isCommit bool) (int64, error) {
//	txTypeValue := enum.TxTypeValue[enum.TxTypeDone]
//	return m.updateTxTypeOnTxReceiptByTxHash(m.TableNamePayment(), hash, txTypeValue, tx, isCommit)
//}
//
//// UpdateTxPaymentNotifiedByTxHash TxReceiptテーブルの該当するsent_hash_txのレコードのcurrnt_tx_typeを更新する
//func (m *DB) UpdateTxPaymentNotifiedByTxHash(hash string, tx *sqlx.Tx, isCommit bool) (int64, error) {
//	txTypeValue := enum.TxTypeValue[enum.TxTypeNotified]
//	return m.updateTxTypeOnTxReceiptByTxHash(m.TableNamePayment(), hash, txTypeValue, tx, isCommit)
//}
//
//// UpdateTxPaymentDoneByID TxReceiptテーブルの該当するIDのレコードのcurrnt_tx_typeを更新する
//func (m *DB) UpdateTxPaymentDoneByID(ID int64, tx *sqlx.Tx, isCommit bool) (int64, error) {
//	txTypeValue := enum.TxTypeValue[enum.TxTypeDone]
//	return m.updateTxTypeOnTxReceiptByID(m.TableNamePayment(), ID, txTypeValue, tx, isCommit)
//}
//
//// UpdateTxPaymentNotifiedByID TxReceiptテーブルの該当するIDのレコードのcurrnt_tx_typeを更新する
//func (m *DB) UpdateTxPaymentNotifiedByID(ID int64, tx *sqlx.Tx, isCommit bool) (int64, error) {
//	txTypeValue := enum.TxTypeValue[enum.TxTypeNotified]
//	return m.updateTxTypeOnTxReceiptByID(m.TableNamePayment(), ID, txTypeValue, tx, isCommit)
//}
