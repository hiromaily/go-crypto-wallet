package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainAuth "github.com/hiromaily/go-crypto-wallet/internal/domain/auth"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/postgres/sqlcgen"
)

// AuthAccountKeyRepositorySqlc is repository for auth_account_key table using sqlc
type AuthAccountKeyRepositorySqlc struct {
	queries      *sqlcgen.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewAuthAccountKeyRepositorySqlc returns AuthAccountKeyRepositorySqlc object
func NewAuthAccountKeyRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *AuthAccountKeyRepositorySqlc {
	return &AuthAccountKeyRepositorySqlc{
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
	}
}

// authAccountKeyRow defines the common interface for auth account key query results
type authAccountKeyRow interface {
	*sqlcgen.GetAuthAccountKeyRow | *sqlcgen.GetAuthAccountKeyByAccountRow
}

// convertAuthAccountKeyRow is a helper to convert query row types to domain.AuthAccountKey
// This reduces code duplication between different query result converters
func convertAuthAccountKeyRow[T authAccountKeyRow](row T) (*domainAuth.AuthAccountKey, error) {
	// Extract fields using type switch to handle different row types
	var (
		id                     int16
		coin                   string
		keyType                string
		authAccount            string
		account                string
		p2pkhAddress           string
		p2shSegwitAddress      string
		bech32Address          string
		taprootAddress         sql.NullString
		fullPublicKey          string
		multisigAddress        string
		redeemScript           string
		walletImportFormat     string
		accountExtendedPrivkey sql.NullString
		idx                    int64
		addrStatus             int16
		updatedAt              sql.NullTime
	)

	switch v := any(row).(type) {
	case *sqlcgen.GetAuthAccountKeyRow:
		id = v.ID
		coin = v.Coin
		keyType = v.KeyType
		authAccount = v.AuthAccount
		account = v.Account
		p2pkhAddress = v.P2pkhAddress
		p2shSegwitAddress = v.P2shSegwitAddress
		bech32Address = v.Bech32Address
		taprootAddress = v.TaprootAddress
		fullPublicKey = v.FullPublicKey
		multisigAddress = v.MultisigAddress
		redeemScript = v.RedeemScript
		walletImportFormat = v.WalletImportFormat
		accountExtendedPrivkey = v.AccountExtendedPrivkey
		idx = v.Idx
		addrStatus = v.AddrStatus
		updatedAt = v.UpdatedAt
	case *sqlcgen.GetAuthAccountKeyByAccountRow:
		id = v.ID
		coin = v.Coin
		keyType = v.KeyType
		authAccount = v.AuthAccount
		account = v.Account
		p2pkhAddress = v.P2pkhAddress
		p2shSegwitAddress = v.P2shSegwitAddress
		bech32Address = v.Bech32Address
		taprootAddress = v.TaprootAddress
		fullPublicKey = v.FullPublicKey
		multisigAddress = v.MultisigAddress
		redeemScript = v.RedeemScript
		walletImportFormat = v.WalletImportFormat
		accountExtendedPrivkey = v.AccountExtendedPrivkey
		idx = v.Idx
		addrStatus = v.AddrStatus
		updatedAt = v.UpdatedAt
	}

	status, err := domainAddress.AddrStatusFromInt8(int8(addrStatus))
	if err != nil {
		return nil, fmt.Errorf("invalid addr status in database: %w", err)
	}

	key := &domainAuth.AuthAccountKey{
		ID:                 id,
		CoinTypeCode:       domainCoin.CoinTypeCode(coin),
		KeyType:            keyType,
		AuthAccount:        domainAccount.AuthType(authAccount),
		Account:            domainAccount.AccountType(account),
		P2pkhAddress:       p2pkhAddress,
		P2shSegwitAddress:  p2shSegwitAddress,
		Bech32Address:      bech32Address,
		FullPublicKey:      fullPublicKey,
		MultisigAddress:    multisigAddress,
		RedeemScript:       redeemScript,
		WalletImportFormat: walletImportFormat,
		Idx:                idx,
		AddrStatus:         status,
	}

	if taprootAddress.Valid {
		key.TaprootAddress = &taprootAddress.String
	}
	if accountExtendedPrivkey.Valid {
		key.AccountExtendedPrivkey = &accountExtendedPrivkey.String
	}
	if updatedAt.Valid {
		key.UpdatedAt = &updatedAt.Time
	}

	return key, nil
}

// convertGetAuthAccountKeyRow converts sqlcgen.GetAuthAccountKeyRow to domain.AuthAccountKey entity.
func convertGetAuthAccountKeyRow(row *sqlcgen.GetAuthAccountKeyRow) (*domainAuth.AuthAccountKey, error) {
	return convertAuthAccountKeyRow(row)
}

// convertGetAuthAccountKeyByAccountRow converts sqlcgen.GetAuthAccountKeyByAccountRow to domain.AuthAccountKey entity.
func convertGetAuthAccountKeyByAccountRow(
	row *sqlcgen.GetAuthAccountKeyByAccountRow,
) (*domainAuth.AuthAccountKey, error) {
	return convertAuthAccountKeyRow(row)
}

// convertFromAuthAccountKey converts domain.AuthAccountKey entity to sqlcgen.AuthAccountKey.
func convertFromAuthAccountKey(key *domainAuth.AuthAccountKey) *sqlcgen.AuthAccountKey {
	sqlcKey := &sqlcgen.AuthAccountKey{
		ID:                 key.ID,
		Coin:               key.CoinTypeCode.String(),
		KeyType:            key.KeyType,
		AuthAccount:        key.AuthAccount.String(),
		Account:            key.Account.String(),
		P2pkhAddress:       key.P2pkhAddress,
		FullPublicKey:      key.FullPublicKey,
		MultisigAddress:    key.MultisigAddress,
		RedeemScript:       key.RedeemScript,
		WalletImportFormat: key.WalletImportFormat,
		Idx:                key.Idx,
		AddrStatus:         int16(key.AddrStatus.Int8()),
	}

	// Handle nullable string fields (string -> sql.NullString)
	// Empty strings are converted to NULL to avoid UNIQUE constraint violations for BCH
	if key.P2shSegwitAddress != "" {
		sqlcKey.P2shSegwitAddress = sql.NullString{String: key.P2shSegwitAddress, Valid: true}
	}
	if key.Bech32Address != "" {
		sqlcKey.Bech32Address = sql.NullString{String: key.Bech32Address, Valid: true}
	}
	if key.TaprootAddress != nil {
		sqlcKey.TaprootAddress = sql.NullString{String: *key.TaprootAddress, Valid: true}
	}
	if key.AccountExtendedPrivkey != nil {
		sqlcKey.AccountExtendedPrivkey = sql.NullString{String: *key.AccountExtendedPrivkey, Valid: true}
	}
	if key.UpdatedAt != nil {
		sqlcKey.UpdatedAt = sql.NullTime{Time: *key.UpdatedAt, Valid: true}
	}

	return sqlcKey
}

// GetOne returns one record by authType
func (r *AuthAccountKeyRepositorySqlc) GetOne(
	ctx context.Context, authType domainAccount.AuthType,
) (*domainAuth.AuthAccountKey, error) {
	authKey, err := r.queries.GetAuthAccountKey(ctx, sqlcgen.GetAuthAccountKeyParams{
		Coin:        r.coinTypeCode.String(),
		AuthAccount: authType.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetAuthAccountKey(): %w", err)
	}

	return convertGetAuthAccountKeyRow(&authKey)
}

// GetByAccount returns one record by authType and accountType
func (r *AuthAccountKeyRepositorySqlc) GetByAccount(
	authType domainAccount.AuthType, accountType domainAccount.AccountType,
) (*domainAuth.AuthAccountKey, error) {
	ctx := context.Background()

	authKey, err := r.queries.GetAuthAccountKeyByAccount(ctx, sqlcgen.GetAuthAccountKeyByAccountParams{
		Coin:        r.coinTypeCode.String(),
		AuthAccount: authType.String(),
		Account:     accountType.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetAuthAccountKeyByAccount(): %w", err)
	}

	return convertGetAuthAccountKeyByAccountRow(&authKey)
}

// Insert inserts record
func (r *AuthAccountKeyRepositorySqlc) Insert(item *domainAuth.AuthAccountKey) error {
	ctx := context.Background()

	sqlcItem := convertFromAuthAccountKey(item)
	_, err := r.queries.InsertAuthAccountKey(ctx, sqlcgen.InsertAuthAccountKeyParams{
		Coin:                   sqlcItem.Coin,
		KeyType:                sqlcItem.KeyType,
		AuthAccount:            sqlcItem.AuthAccount,
		Account:                sqlcItem.Account,
		P2pkhAddress:           sqlcItem.P2pkhAddress,
		P2shSegwitAddress:      sqlcItem.P2shSegwitAddress,
		Bech32Address:          sqlcItem.Bech32Address,
		TaprootAddress:         sqlcItem.TaprootAddress,
		FullPublicKey:          sqlcItem.FullPublicKey,
		MultisigAddress:        sqlcItem.MultisigAddress,
		RedeemScript:           sqlcItem.RedeemScript,
		WalletImportFormat:     sqlcItem.WalletImportFormat,
		AccountExtendedPrivkey: sqlcItem.AccountExtendedPrivkey,
		Idx:                    sqlcItem.Idx,
		AddrStatus:             sqlcItem.AddrStatus,
	})
	if err != nil {
		return fmt.Errorf("failed to call InsertAuthAccountKey(): %w", err)
	}

	return nil
}

// UpdateAddrStatus updates addr_status
func (r *AuthAccountKeyRepositorySqlc) UpdateAddrStatus(
	addrStatus domainAddress.AddrStatus, strWIF string,
) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateAuthAccountKeyAddrStatus(ctx, sqlcgen.UpdateAuthAccountKeyAddrStatusParams{
		AddrStatus:         int16(addrStatus.Int8()),
		UpdatedAt:          sql.NullTime{Time: time.Now(), Valid: true},
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
