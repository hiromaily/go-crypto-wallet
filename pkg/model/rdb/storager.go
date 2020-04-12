package rdb

import (
	"github.com/jmoiron/sqlx"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb/coldrepo"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb/walletrepo"
)

// TODO: decouple interface into interface of each repo
type WalletStorager interface {
	Close() error
	MustBegin() *sqlx.Tx
	//acount_pubkey_repo
	GetAllAccountPubKeyTable(accountType account.AccountType) ([]walletrepo.AccountPublicKeyTable, error)
	GetOneUnAllocatedAccountPubKeyTable(accountType account.AccountType) (*walletrepo.AccountPublicKeyTable, error)
	InsertAccountPubKeyTable(accountType account.AccountType, accountPubKeyTables []walletrepo.AccountPublicKeyTable, tx *sqlx.Tx, isCommit bool) error
	UpdateAccountOnAccountPubKeyTable(accountType account.AccountType, accountKeyTable []walletrepo.AccountPublicKeyTable, tx *sqlx.Tx, isCommit bool) error
	UpdateIsAllocatedOnAccountPubKeyTable(accountType account.AccountType, accountKeyTable []walletrepo.AccountPublicKeyTable, tx *sqlx.Tx, isCommit bool) error
	//payment_request_repo
	GetPaymentRequestAll() ([]walletrepo.PaymentRequest, error)
	GetPaymentRequestByPaymentID(paymentID int64) ([]walletrepo.PaymentRequest, error)
	InsertPaymentRequest(paymentRequests []walletrepo.PaymentRequest, tx *sqlx.Tx, isCommit bool) error
	UpdateIsDoneOnPaymentRequest(paymentID int64, tx *sqlx.Tx, isCommit bool) (int64, error)
	UpdatePaymentIDOnPaymentRequest(paymentID int64, ids []int64, tx *sqlx.Tx, isCommit bool) (int64, error)
	ResetAnyFlagOnPaymentRequestForTestOnly(tx *sqlx.Tx, isCommit bool) (int64, error)
	//tx_input_repo
	GetTxInputByReceiptID(actionType enum.ActionType, receiptID int64) ([]walletrepo.TxInput, error)
	InsertTxInputForUnsigned(actionType enum.ActionType, txReceiptInputs []walletrepo.TxInput, tx *sqlx.Tx, isCommit bool) error
	//tx_output_repo
	GetTxOutputByReceiptID(actionType enum.ActionType, receiptID int64) ([]walletrepo.TxOutput, error)
	InsertTxOutputForUnsigned(actionType enum.ActionType, txReceiptOutputs []walletrepo.TxOutput, tx *sqlx.Tx, isCommit bool) error
	//tx_repo
	GetTxByID(actionType enum.ActionType, id int64) (*walletrepo.TxTable, error)
	GetTxCountByUnsignedHex(actionType enum.ActionType, hex string) (int64, error)
	GetTxIDBySentHash(actionType enum.ActionType, hash string) (int64, error)
	GetSentTxHashByTxTypeSent(actionType enum.ActionType) ([]string, error)
	GetSentTxHashByTxTypeDone(actionType enum.ActionType) ([]string, error)
	InsertTxForUnsigned(actionType enum.ActionType, txReceipt *walletrepo.TxTable, tx *sqlx.Tx, isCommit bool) (int64, error)
	UpdateTxAfterSent(actionType enum.ActionType, txReceipt *walletrepo.TxTable, tx *sqlx.Tx, isCommit bool) (int64, error)
	UpdateTxTypeDoneByTxHash(actionType enum.ActionType, hash string, tx *sqlx.Tx, isCommit bool) (int64, error)
	UpdateTxTypeNotifiedByTxHash(actionType enum.ActionType, hash string, tx *sqlx.Tx, isCommit bool) (int64, error)
	UpdateTxTypeDoneByID(actionType enum.ActionType, ID int64, tx *sqlx.Tx, isCommit bool) (int64, error)
	UpdateTxTypeNotifiedByID(actionType enum.ActionType, ID int64, tx *sqlx.Tx, isCommit bool) (int64, error)
}

type ColdStorager interface {
	Close() error
	//acount_pubkey_repo
	GetMaxIndexOnAccountKeyTable(accountType account.AccountType) (int64, error)
	GetOneByMaxIDOnAccountKeyTable(accountType account.AccountType) (*coldrepo.AccountKeyTable, error)
	GetAllAccountKeyByKeyStatus(accountType account.AccountType, keyStatus enum.KeyStatus) ([]coldrepo.AccountKeyTable, error)
	GetAllAccountKeyByMultiAddrs(accountType account.AccountType, addrs []string) ([]coldrepo.AccountKeyTable, error)
	InsertAccountKeyTable(accountType account.AccountType, accountKeyTables []coldrepo.AccountKeyTable, tx *sqlx.Tx, isCommit bool) error
	UpdateKeyStatusByWIF(accountType account.AccountType, keyStatus enum.KeyStatus, strWIF string, tx *sqlx.Tx, isCommit bool) (int64, error)
	UpdateKeyStatusByWIFs(accountType account.AccountType, keyStatus enum.KeyStatus, wifs []string, tx *sqlx.Tx, isCommit bool) (int64, error)
	UpdateMultisigAddrOnAccountKeyTableByFullPubKey(accountType account.AccountType, accountKeyTable []coldrepo.AccountKeyTable, tx *sqlx.Tx, isCommit bool) error
	//added_pubkey_history_repo
	GetAddedPubkeyHistoryTableByNoWalletMultisigAddress(accountType account.AccountType) ([]coldrepo.AddedPubkeyHistoryTable, error)
	GetAddedPubkeyHistoryTableByNotExported(accountType account.AccountType) ([]coldrepo.AddedPubkeyHistoryTable, error)
	InsertAddedPubkeyHistoryTable(accountType account.AccountType, addedPubkeyHistoryTables []coldrepo.AddedPubkeyHistoryTable, tx *sqlx.Tx, isCommit bool) error
	UpdateMultisigAddrOnAddedPubkeyHistoryTable(accountType account.AccountType, multiSigAddr, redeemScript, authAddr1, fullPublicKey string, tx *sqlx.Tx, isCommit bool) error
	UpdateIsExportedOnAddedPubkeyHistoryTable(accountType account.AccountType, ids []int64, tx *sqlx.Tx, isCommit bool) (int64, error)
	//seed_repo
	GetSeedAll() ([]coldrepo.Seed, error)
	GetSeedOne() (coldrepo.Seed, error)
	GetSeedCount() (int64, error)
	InsertSeed(seed string, tx *sqlx.Tx, isCommit bool) (int64, error)
}

//TODO: use only required func
//type KeygenStorager interface {
//	Close() error
//	//acount_pubkey_repo
//	GetMaxIndexOnAccountKeyTable(accountType account.AccountType) (int64, error)
//	GetOneByMaxIDOnAccountKeyTable(accountType account.AccountType) (*keygenrepo.AccountKeyTable, error)
//	GetAllAccountKeyByKeyStatus(accountType account.AccountType, keyStatus enum.KeyStatus) ([]keygenrepo.AccountKeyTable, error)
//	GetAllAccountKeyByMultiAddrs(accountType account.AccountType, addrs []string) ([]keygenrepo.AccountKeyTable, error)
//	InsertAccountKeyTable(accountType account.AccountType, accountKeyTables []keygenrepo.AccountKeyTable, tx *sqlx.Tx, isCommit bool) error
//	UpdateKeyStatusByWIF(accountType account.AccountType, keyStatus enum.KeyStatus, strWIF string, tx *sqlx.Tx, isCommit bool) (int64, error)
//	UpdateKeyStatusByWIFs(accountType account.AccountType, keyStatus enum.KeyStatus, wifs []string, tx *sqlx.Tx, isCommit bool) (int64, error)
//	UpdateMultisigAddrOnAccountKeyTableByFullPubKey(accountType account.AccountType, accountKeyTable []keygenrepo.AccountKeyTable, tx *sqlx.Tx, isCommit bool) error
//	//added_pubkey_history_repo
//	GetAddedPubkeyHistoryTableByNoWalletMultisigAddress(accountType account.AccountType) ([]keygenrepo.AddedPubkeyHistoryTable, error)
//	GetAddedPubkeyHistoryTableByNotExported(accountType account.AccountType) ([]keygenrepo.AddedPubkeyHistoryTable, error)
//	InsertAddedPubkeyHistoryTable(accountType account.AccountType, addedPubkeyHistoryTables []keygenrepo.AddedPubkeyHistoryTable, tx *sqlx.Tx, isCommit bool) error
//	UpdateMultisigAddrOnAddedPubkeyHistoryTable(accountType account.AccountType, multiSigAddr, redeemScript, authAddr1, fullPublicKey string, tx *sqlx.Tx, isCommit bool) error
//	UpdateIsExportedOnAddedPubkeyHistoryTable(accountType account.AccountType, ids []int64, tx *sqlx.Tx, isCommit bool) (int64, error)
//	//seed_repo
//	GetSeedAll() ([]keygenrepo.Seed, error)
//	GetSeedOne() (keygenrepo.Seed, error)
//	GetSeedCount() (int64, error)
//	InsertSeed(seed string, tx *sqlx.Tx, isCommit bool) (int64, error)
//}

//TODO: use only required func
//type SignatureStorager interface {
//	Close() error
//	//acount_pubkey_repo
//	GetMaxIndexOnAccountKeyTable(accountType account.AccountType) (int64, error)
//	GetOneByMaxIDOnAccountKeyTable(accountType account.AccountType) (*signaturerepo.AccountKeyTable, error)
//	GetAllAccountKeyByKeyStatus(accountType account.AccountType, keyStatus enum.KeyStatus) ([]signaturerepo.AccountKeyTable, error)
//	GetAllAccountKeyByMultiAddrs(accountType account.AccountType, addrs []string) ([]signaturerepo.AccountKeyTable, error)
//	InsertAccountKeyTable(accountType account.AccountType, accountKeyTables []signaturerepo.AccountKeyTable, tx *sqlx.Tx, isCommit bool) error
//	UpdateKeyStatusByWIF(accountType account.AccountType, keyStatus enum.KeyStatus, strWIF string, tx *sqlx.Tx, isCommit bool) (int64, error)
//	UpdateKeyStatusByWIFs(accountType account.AccountType, keyStatus enum.KeyStatus, wifs []string, tx *sqlx.Tx, isCommit bool) (int64, error)
//	UpdateMultisigAddrOnAccountKeyTableByFullPubKey(accountType account.AccountType, accountKeyTable []signaturerepo.AccountKeyTable, tx *sqlx.Tx, isCommit bool) error
//	//added_pubkey_history_repo
//	GetAddedPubkeyHistoryTableByNoWalletMultisigAddress(accountType account.AccountType) ([]signaturerepo.AddedPubkeyHistoryTable, error)
//	GetAddedPubkeyHistoryTableByNotExported(accountType account.AccountType) ([]signaturerepo.AddedPubkeyHistoryTable, error)
//	InsertAddedPubkeyHistoryTable(accountType account.AccountType, addedPubkeyHistoryTables []signaturerepo.AddedPubkeyHistoryTable, tx *sqlx.Tx, isCommit bool) error
//	UpdateMultisigAddrOnAddedPubkeyHistoryTable(accountType account.AccountType, multiSigAddr, redeemScript, authAddr1, fullPublicKey string, tx *sqlx.Tx, isCommit bool) error
//	UpdateIsExportedOnAddedPubkeyHistoryTable(accountType account.AccountType, ids []int64, tx *sqlx.Tx, isCommit bool) (int64, error)
//	//seed_repo
//	GetSeedAll() ([]signaturerepo.Seed, error)
//	GetSeedOne() (signaturerepo.Seed, error)
//	GetSeedCount() (int64, error)
//	InsertSeed(seed string, tx *sqlx.Tx, isCommit bool) (int64, error)
//}
