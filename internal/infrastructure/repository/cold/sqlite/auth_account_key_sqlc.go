package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainAuth "github.com/hiromaily/go-crypto-wallet/internal/domain/auth"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/sqlite/sqlcgen"
)

// AuthAccountKeyRepositorySqlc is repository for auth_account_key table using sqlc
type AuthAccountKeyRepositorySqlc struct {
	db           *sql.DB
	queries      *sqlcgen.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewAuthAccountKeyRepositorySqlc returns AuthAccountKeyRepositorySqlc object
func NewAuthAccountKeyRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *AuthAccountKeyRepositorySqlc {
	return &AuthAccountKeyRepositorySqlc{
		db:           dbConn,
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
	}
}

// GetOne returns one record by authType
// Note: Uses raw SQL with COALESCE to handle NULL values in BCH (no SegWit support)
func (r *AuthAccountKeyRepositorySqlc) GetOne(
	ctx context.Context, authType domainAccount.AuthType,
) (*domainAuth.AuthAccountKey, error) {
	query := `SELECT
		id, coin, key_type, auth_account, account, p2pkh_address,
		COALESCE(p2sh_segwit_address, '') as p2sh_segwit_address,
		COALESCE(bech32_address, '') as bech32_address,
		taproot_address, full_public_key, multisig_address, redeem_script,
		wallet_import_format, account_extended_privkey, idx, addr_status, updated_at
	FROM auth_account_key WHERE coin = ? AND auth_account = ? LIMIT 1`

	var key domainAuth.AuthAccountKey
	var taprootAddress, accountExtendedPrivkey, updatedAt sql.NullString
	var addrStatus int64

	err := r.db.QueryRowContext(ctx, query, r.coinTypeCode.String(), authType.String()).Scan(
		&key.ID,
		&key.CoinTypeCode,
		&key.KeyType,
		&key.AuthAccount,
		&key.Account,
		&key.P2pkhAddress,
		&key.P2shSegwitAddress,
		&key.Bech32Address,
		&taprootAddress,
		&key.FullPublicKey,
		&key.MultisigAddress,
		&key.RedeemScript,
		&key.WalletImportFormat,
		&accountExtendedPrivkey,
		&key.Idx,
		&addrStatus,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query auth account key: %w", err)
	}

	// Convert scanned values
	status, err := domainAddress.AddrStatusFromInt8(int8(addrStatus))
	if err != nil {
		return nil, fmt.Errorf("invalid addr status: %w", err)
	}
	key.AddrStatus = status
	if taprootAddress.Valid {
		key.TaprootAddress = &taprootAddress.String
	}
	if accountExtendedPrivkey.Valid {
		key.AccountExtendedPrivkey = &accountExtendedPrivkey.String
	}
	if updatedAt.Valid {
		t, err := time.Parse("2006-01-02 15:04:05", updatedAt.String)
		if err != nil {
			return nil, fmt.Errorf("invalid timestamp format in database: %w", err)
		}
		key.UpdatedAt = &t
	}

	return &key, nil
}

// GetByAccount returns one record by authType and accountType
// Note: Uses raw SQL with COALESCE to handle NULL values in BCH (no SegWit support)
func (r *AuthAccountKeyRepositorySqlc) GetByAccount(
	authType domainAccount.AuthType, accountType domainAccount.AccountType,
) (*domainAuth.AuthAccountKey, error) {
	ctx := context.Background()

	query := `SELECT
		id, coin, key_type, auth_account, account, p2pkh_address,
		COALESCE(p2sh_segwit_address, '') as p2sh_segwit_address,
		COALESCE(bech32_address, '') as bech32_address,
		taproot_address, full_public_key, multisig_address, redeem_script,
		wallet_import_format, account_extended_privkey, idx, addr_status, updated_at
	FROM auth_account_key WHERE coin = ? AND auth_account = ? AND account = ? LIMIT 1`

	var key domainAuth.AuthAccountKey
	var taprootAddress, accountExtendedPrivkey, updatedAt sql.NullString
	var addrStatus int64

	err := r.db.QueryRowContext(ctx, query,
		r.coinTypeCode.String(), authType.String(), accountType.String()).Scan(
		&key.ID,
		&key.CoinTypeCode,
		&key.KeyType,
		&key.AuthAccount,
		&key.Account,
		&key.P2pkhAddress,
		&key.P2shSegwitAddress,
		&key.Bech32Address,
		&taprootAddress,
		&key.FullPublicKey,
		&key.MultisigAddress,
		&key.RedeemScript,
		&key.WalletImportFormat,
		&accountExtendedPrivkey,
		&key.Idx,
		&addrStatus,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query auth account key: %w", err)
	}

	// Convert scanned values
	status, err := domainAddress.AddrStatusFromInt8(int8(addrStatus))
	if err != nil {
		return nil, fmt.Errorf("invalid addr status: %w", err)
	}
	key.AddrStatus = status
	if taprootAddress.Valid {
		key.TaprootAddress = &taprootAddress.String
	}
	if accountExtendedPrivkey.Valid {
		key.AccountExtendedPrivkey = &accountExtendedPrivkey.String
	}
	if updatedAt.Valid {
		t, err := time.Parse("2006-01-02 15:04:05", updatedAt.String)
		if err != nil {
			return nil, fmt.Errorf("invalid timestamp format in database: %w", err)
		}
		key.UpdatedAt = &t
	}

	return &key, nil
}

// Insert inserts record
// For SQLite, we use raw SQL with NULLIF to convert empty strings to NULL
// to avoid UNIQUE constraint violations on bech32_address and p2sh_segwit_address
func (r *AuthAccountKeyRepositorySqlc) Insert(item *domainAuth.AuthAccountKey) error {
	ctx := context.Background()

	// Use raw SQL with NULLIF to handle empty strings for BCH compatibility
	// BCH doesn't support SegWit/bech32, so these fields are empty and must be NULL
	query := `INSERT INTO auth_account_key (
		coin, key_type, auth_account, account, p2pkh_address,
		p2sh_segwit_address, bech32_address, taproot_address,
		full_public_key, multisig_address, redeem_script,
		wallet_import_format, account_extended_privkey, idx, addr_status
	) VALUES (?, ?, ?, ?, ?, NULLIF(?, ''), NULLIF(?, ''), ?, ?, ?, ?, ?, ?, ?, ?)`

	// Handle nullable pointer fields
	var taprootAddr interface{} = nil
	if item.TaprootAddress != nil {
		taprootAddr = *item.TaprootAddress
	}
	var accountExtPrivkey interface{} = nil
	if item.AccountExtendedPrivkey != nil {
		accountExtPrivkey = *item.AccountExtendedPrivkey
	}

	_, err := r.db.ExecContext(ctx, query,
		item.CoinTypeCode.String(),
		item.KeyType,
		item.AuthAccount.String(),
		item.Account.String(),
		item.P2pkhAddress,
		item.P2shSegwitAddress,
		item.Bech32Address,
		taprootAddr,
		item.FullPublicKey,
		item.MultisigAddress,
		item.RedeemScript,
		item.WalletImportFormat,
		accountExtPrivkey,
		item.Idx,
		item.AddrStatus.Int8(),
	)
	if err != nil {
		return fmt.Errorf("failed to insert auth account key: %w", err)
	}

	return nil
}

// UpdateAddrStatus updates addr_status
func (r *AuthAccountKeyRepositorySqlc) UpdateAddrStatus(
	addrStatus domainAddress.AddrStatus, strWIF string,
) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateAuthAccountKeyAddrStatus(ctx, sqlcgen.UpdateAuthAccountKeyAddrStatusParams{
		AddrStatus:         int64(addrStatus.Int8()),
		UpdatedAt:          sql.NullString{String: time.Now().Format("2006-01-02 15:04:05"), Valid: true},
		Coin:               r.coinTypeCode.String(),
		WalletImportFormat: strWIF,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call UpdateAuthAccountKeyAddrStatus(): %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
	}

	return rowsAffected, nil
}
