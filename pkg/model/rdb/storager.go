package rdb

import (
	"github.com/jmoiron/sqlx"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb/coldrepo"
)

// WalletStorager is wallet storager interface
// TODO: decouple interface into interface of each repo
type WalletStorager interface {
	Close() error
	MustBegin() *sqlx.Tx
	//account_pubkey_repo
	//GetOneUnAllocatedAccountPubKeyTable(accountType account.AccountType) (*walletrepo.AccountPublicKeyTable, error)
	//InsertAccountPubKeyTable(accountType account.AccountType, accountPubKeyTables []walletrepo.AccountPublicKeyTable, tx *sqlx.Tx, isCommit bool) error
	//UpdateIsAllocatedOnAccountPubKeyTable(accountType account.AccountType, accountKeyTable []walletrepo.AccountPublicKeyTable, tx *sqlx.Tx, isCommit bool) error
}

// ColdStorager is coldwalet storager interface
type ColdStorager interface {
	Close() error
	//account_pubkey_repo
	GetMaxIndexOnAccountKeyTable(accountType account.AccountType) (int64, error)
	GetOneByMaxIDOnAccountKeyTable(accountType account.AccountType) (*coldrepo.AccountKeyTable, error)
	GetAllAccountKeyByAddrStatus(accountType account.AccountType, addrStatus address.AddrStatus) ([]coldrepo.AccountKeyTable, error)
	GetAllAccountKeyByMultiAddrs(accountType account.AccountType, addrs []string) ([]coldrepo.AccountKeyTable, error)
	InsertAccountKeyTable(accountType account.AccountType, accountKeyTables []coldrepo.AccountKeyTable, tx *sqlx.Tx, isCommit bool) error
	UpdateAddrStatusByWIF(accountType account.AccountType, addrStatus address.AddrStatus, strWIF string, tx *sqlx.Tx, isCommit bool) (int64, error)
	UpdateAddrStatusByWIFs(accountType account.AccountType, addrStatus address.AddrStatus, wifs []string, tx *sqlx.Tx, isCommit bool) (int64, error)
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
