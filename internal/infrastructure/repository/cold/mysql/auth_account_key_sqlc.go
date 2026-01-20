package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainAuth "github.com/hiromaily/go-crypto-wallet/internal/domain/auth"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
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

// convertGetAuthAccountKeyRow converts sqlcgen.GetAuthAccountKeyRow to domain.AuthAccountKey entity.
func convertGetAuthAccountKeyRow(row *sqlcgen.GetAuthAccountKeyRow) (*domainAuth.AuthAccountKey, error) {
	addrStatus, err := domainAddress.AddrStatusFromInt8(row.AddrStatus)
	if err != nil {
		return nil, fmt.Errorf("invalid addr status in database: %w", err)
	}

	key := &domainAuth.AuthAccountKey{
		ID:                 row.ID,
		CoinTypeCode:       domainCoin.CoinTypeCode(row.Coin),
		KeyType:            row.KeyType,
		AuthAccount:        domainAccount.AuthType(row.AuthAccount),
		Account:            domainAccount.AccountType(row.Account),
		P2pkhAddress:       row.P2pkhAddress,
		P2shSegwitAddress:  row.P2shSegwitAddress,
		Bech32Address:      row.Bech32Address,
		FullPublicKey:      row.FullPublicKey,
		MultisigAddress:    row.MultisigAddress,
		RedeemScript:       row.RedeemScript,
		WalletImportFormat: row.WalletImportFormat,
		Idx:                row.Idx,
		AddrStatus:         addrStatus,
	}

	if row.TaprootAddress.Valid {
		key.TaprootAddress = &row.TaprootAddress.String
	}
	if row.AccountExtendedPrivkey.Valid {
		key.AccountExtendedPrivkey = &row.AccountExtendedPrivkey.String
	}
	if row.UpdatedAt.Valid {
		key.UpdatedAt = &row.UpdatedAt.Time
	}

	return key, nil
}

// convertGetAuthAccountKeyByAccountRow converts sqlcgen.GetAuthAccountKeyByAccountRow to domain.AuthAccountKey entity.
func convertGetAuthAccountKeyByAccountRow(row *sqlcgen.GetAuthAccountKeyByAccountRow) (*domainAuth.AuthAccountKey, error) {
	addrStatus, err := domainAddress.AddrStatusFromInt8(row.AddrStatus)
	if err != nil {
		return nil, fmt.Errorf("invalid addr status in database: %w", err)
	}

	key := &domainAuth.AuthAccountKey{
		ID:                 row.ID,
		CoinTypeCode:       domainCoin.CoinTypeCode(row.Coin),
		KeyType:            row.KeyType,
		AuthAccount:        domainAccount.AuthType(row.AuthAccount),
		Account:            domainAccount.AccountType(row.Account),
		P2pkhAddress:       row.P2pkhAddress,
		P2shSegwitAddress:  row.P2shSegwitAddress,
		Bech32Address:      row.Bech32Address,
		FullPublicKey:      row.FullPublicKey,
		MultisigAddress:    row.MultisigAddress,
		RedeemScript:       row.RedeemScript,
		WalletImportFormat: row.WalletImportFormat,
		Idx:                row.Idx,
		AddrStatus:         addrStatus,
	}

	if row.TaprootAddress.Valid {
		key.TaprootAddress = &row.TaprootAddress.String
	}
	if row.AccountExtendedPrivkey.Valid {
		key.AccountExtendedPrivkey = &row.AccountExtendedPrivkey.String
	}
	if row.UpdatedAt.Valid {
		key.UpdatedAt = &row.UpdatedAt.Time
	}

	return key, nil
}

// convertToAuthAccountKey converts sqlcgen.AuthAccountKey to domain.AuthAccountKey entity.
// SECURITY: Handles WIF (private key) data - never log the wallet import format field.
func convertToAuthAccountKey(sqlcKey *sqlcgen.AuthAccountKey) (*domainAuth.AuthAccountKey, error) {
	addrStatus, err := domainAddress.AddrStatusFromInt8(sqlcKey.AddrStatus)
	if err != nil {
		return nil, fmt.Errorf("invalid addr status in database: %w", err)
	}

	// Handle nullable string fields (sql.NullString -> string)
	p2shSegwitAddr := ""
	if sqlcKey.P2shSegwitAddress.Valid {
		p2shSegwitAddr = sqlcKey.P2shSegwitAddress.String
	}
	bech32Addr := ""
	if sqlcKey.Bech32Address.Valid {
		bech32Addr = sqlcKey.Bech32Address.String
	}

	key := &domainAuth.AuthAccountKey{
		ID:                 sqlcKey.ID,
		CoinTypeCode:       domainCoin.CoinTypeCode(sqlcKey.Coin),
		KeyType:            sqlcKey.KeyType,
		AuthAccount:        domainAccount.AuthType(sqlcKey.AuthAccount),
		Account:            domainAccount.AccountType(sqlcKey.Account),
		P2pkhAddress:       sqlcKey.P2pkhAddress,
		P2shSegwitAddress:  p2shSegwitAddr,
		Bech32Address:      bech32Addr,
		FullPublicKey:      sqlcKey.FullPublicKey,
		MultisigAddress:    sqlcKey.MultisigAddress,
		RedeemScript:       sqlcKey.RedeemScript,
		WalletImportFormat: sqlcKey.WalletImportFormat, // WIF - NEVER log
		Idx:                sqlcKey.Idx,
		AddrStatus:         addrStatus,
	}

	if sqlcKey.TaprootAddress.Valid {
		key.TaprootAddress = &sqlcKey.TaprootAddress.String
	}
	if sqlcKey.AccountExtendedPrivkey.Valid {
		key.AccountExtendedPrivkey = &sqlcKey.AccountExtendedPrivkey.String
	}
	if sqlcKey.UpdatedAt.Valid {
		key.UpdatedAt = &sqlcKey.UpdatedAt.Time
	}

	return key, nil
}

// convertFromAuthAccountKey converts domain.AuthAccountKey entity to sqlcgen.AuthAccountKey.
func convertFromAuthAccountKey(key *domainAuth.AuthAccountKey) *sqlcgen.AuthAccountKey {
	sqlcKey := &sqlcgen.AuthAccountKey{
		ID:                 key.ID,
		Coin:               sqlcgen.AuthAccountKeyCoin(key.CoinTypeCode.String()),
		KeyType:            key.KeyType,
		AuthAccount:        key.AuthAccount.String(),
		Account:            key.Account.String(),
		P2pkhAddress:       key.P2pkhAddress,
		FullPublicKey:      key.FullPublicKey,
		MultisigAddress:    key.MultisigAddress,
		RedeemScript:       key.RedeemScript,
		WalletImportFormat: key.WalletImportFormat,
		Idx:                key.Idx,
		AddrStatus:         key.AddrStatus.Int8(),
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
		Coin:        sqlcgen.AuthAccountKeyCoin(r.coinTypeCode.String()),
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
		Coin:        sqlcgen.AuthAccountKeyCoin(r.coinTypeCode.String()),
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
		AddrStatus:         addrStatus.Int8(),
		UpdatedAt:          sql.NullTime{Time: time.Now(), Valid: true},
		Coin:               sqlcgen.AuthAccountKeyCoin(r.coinTypeCode.String()),
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
