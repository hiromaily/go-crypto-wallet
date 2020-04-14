package rdb

import (
	"github.com/jmoiron/sqlx"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/action"
	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb/coldrepo"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb/walletrepo"
)

// TODO: decouple interface into interface of each repo
type WalletStorager interface {
	Close() error
	MustBegin() *sqlx.Tx
	//account_pubkey_repo
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
	GetTxInputByReceiptID(actionType action.ActionType, receiptID int64) ([]walletrepo.TxInput, error)
	InsertTxInputForUnsigned(actionType action.ActionType, txReceiptInputs []walletrepo.TxInput, tx *sqlx.Tx, isCommit bool) error
	//tx_output_repo
	GetTxOutputByReceiptID(actionType action.ActionType, receiptID int64) ([]walletrepo.TxOutput, error)
	InsertTxOutputForUnsigned(actionType action.ActionType, txReceiptOutputs []walletrepo.TxOutput, tx *sqlx.Tx, isCommit bool) error
	//tx_repo
	GetTxByID(actionType action.ActionType, id int64) (*walletrepo.TxTable, error)
	GetTxCountByUnsignedHex(actionType action.ActionType, hex string) (int64, error)
	GetTxIDBySentHash(actionType action.ActionType, hash string) (int64, error)
	GetSentTxHashByTxTypeSent(actionType action.ActionType) ([]string, error)
	GetSentTxHashByTxTypeDone(actionType action.ActionType) ([]string, error)
	InsertTxForUnsigned(actionType action.ActionType, txReceipt *walletrepo.TxTable, tx *sqlx.Tx, isCommit bool) (int64, error)
	UpdateTxAfterSent(actionType action.ActionType, txReceipt *walletrepo.TxTable, tx *sqlx.Tx, isCommit bool) (int64, error)
	UpdateTxTypeDoneByTxHash(actionType action.ActionType, hash string, tx *sqlx.Tx, isCommit bool) (int64, error)
	UpdateTxTypeNotifiedByTxHash(actionType action.ActionType, hash string, tx *sqlx.Tx, isCommit bool) (int64, error)
	UpdateTxTypeDoneByID(actionType action.ActionType, ID int64, tx *sqlx.Tx, isCommit bool) (int64, error)
	UpdateTxTypeNotifiedByID(actionType action.ActionType, ID int64, tx *sqlx.Tx, isCommit bool) (int64, error)
}

type ColdStorager interface {
	Close() error
	//account_pubkey_repo
	GetMaxIndexOnAccountKeyTable(accountType account.AccountType) (int64, error)
	GetOneByMaxIDOnAccountKeyTable(accountType account.AccountType) (*coldrepo.AccountKeyTable, error)
	GetAllAccountKeyByAddressStatus(accountType account.AccountType, keyStatus address.AddressStatus) ([]coldrepo.AccountKeyTable, error)
	GetAllAccountKeyByMultiAddrs(accountType account.AccountType, addrs []string) ([]coldrepo.AccountKeyTable, error)
	InsertAccountKeyTable(accountType account.AccountType, accountKeyTables []coldrepo.AccountKeyTable, tx *sqlx.Tx, isCommit bool) error
	UpdateAddressStatusByWIF(accountType account.AccountType, keyStatus address.AddressStatus, strWIF string, tx *sqlx.Tx, isCommit bool) (int64, error)
	UpdateAddressStatusByWIFs(accountType account.AccountType, keyStatus address.AddressStatus, wifs []string, tx *sqlx.Tx, isCommit bool) (int64, error)
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
