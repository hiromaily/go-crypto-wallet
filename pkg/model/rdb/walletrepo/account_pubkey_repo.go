package walletrepo

import (
	"time"

	"github.com/hiromaily/go-bitcoin/pkg/account"
)

// AccountPublicKeyTable account_key_clientテーブル
type AccountPublicKeyTable struct {
	ID            int64      `db:"id"`
	WalletAddress string     `db:"wallet_address"`
	Account       string     `db:"account"`
	IsAllocated   bool       `db:"is_allocated"`
	UpdatedAt     *time.Time `db:"updated_at"`
}

var accountPubKeyTableName = map[account.AccountType]string{
	account.AccountTypeClient:  "account_pubkey_client",
	account.AccountTypeReceipt: "account_pubkey_receipt",
	account.AccountTypePayment: "account_pubkey_payment",
	account.AccountTypeFee:     "account_pubkey_fee",
	account.AccountTypeStored:  "account_pubkey_stored",
}
