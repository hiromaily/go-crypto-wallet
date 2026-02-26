package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainBTC "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/btc"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/postgres/sqlcgen"
)

// BTCAccountKeyRepositorySqlc is repository for btc_account_key table using sqlc
type BTCAccountKeyRepositorySqlc struct {
	queries      *sqlcgen.Queries
	dbConn       *sql.DB
	coinTypeCode domainCoin.CoinTypeCode
}

// NewBTCAccountKeyRepositorySqlc returns BTCAccountKeyRepositorySqlc object
func NewBTCAccountKeyRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *BTCAccountKeyRepositorySqlc {
	return &BTCAccountKeyRepositorySqlc{
		queries:      sqlcgen.New(dbConn),
		dbConn:       dbConn,
		coinTypeCode: coinTypeCode,
	}
}

// btcAccountKeyRow defines the common interface for BTC account key query results
type btcAccountKeyRow interface {
	*sqlcgen.GetOneBtcAccountKeyByMaxIDRow |
		*sqlcgen.GetBtcAccountKeysByAddrStatusRow |
		*sqlcgen.GetBtcAccountKeysByMultisigAddressesRow
}

// convertBtcAccountKeyRow is a helper to convert query row types to domain.BTCAccountKey
// This reduces code duplication between different query result converters
func convertBtcAccountKeyRow[T btcAccountKeyRow](row T) (*domainBTC.BTCAccountKey, error) {
	// Extract fields using type switch to handle different row types
	var (
		id                     int64
		coin                   string
		keyType                string
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
	case *sqlcgen.GetOneBtcAccountKeyByMaxIDRow:
		id = v.ID
		coin = v.Coin
		keyType = v.KeyType
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
	case *sqlcgen.GetBtcAccountKeysByAddrStatusRow:
		id = v.ID
		coin = v.Coin
		keyType = v.KeyType
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
	case *sqlcgen.GetBtcAccountKeysByMultisigAddressesRow:
		id = v.ID
		coin = v.Coin
		keyType = v.KeyType
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

	key := &domainBTC.BTCAccountKey{
		ID:                 id,
		CoinTypeCode:       domainCoin.CoinTypeCode(coin),
		KeyType:            keyType,
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

// convertGetOneBtcAccountKeyByMaxIDRow converts row to domain.BtcAccountKey entity.
func convertGetOneBtcAccountKeyByMaxIDRow(
	row *sqlcgen.GetOneBtcAccountKeyByMaxIDRow,
) (*domainBTC.BTCAccountKey, error) {
	return convertBtcAccountKeyRow(row)
}

// convertGetBtcAccountKeysByAddrStatusRow converts row to domain.BtcAccountKey entity.
func convertGetBtcAccountKeysByAddrStatusRow(
	row *sqlcgen.GetBtcAccountKeysByAddrStatusRow,
) (*domainBTC.BTCAccountKey, error) {
	return convertBtcAccountKeyRow(row)
}

// convertGetBtcAccountKeysByMultisigAddressesRow converts row to domain.BtcAccountKey entity.
func convertGetBtcAccountKeysByMultisigAddressesRow(
	row *sqlcgen.GetBtcAccountKeysByMultisigAddressesRow,
) (*domainBTC.BTCAccountKey, error) {
	return convertBtcAccountKeyRow(row)
}

// convertFromBTCAccountKey converts domain.BTCAccountKey entity to sqlcgen.BtcAccountKey.
func convertFromBTCAccountKey(key *domainBTC.BTCAccountKey) *sqlcgen.BtcAccountKey {
	sqlcKey := &sqlcgen.BtcAccountKey{
		ID:                 key.ID,
		Coin:               key.CoinTypeCode.String(),
		KeyType:            key.KeyType,
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

// GetMaxIndex returns max idx
func (r *BTCAccountKeyRepositorySqlc) GetMaxIndex(
	ctx context.Context, accountType domainAccount.AccountType,
) (int64, error) {
	result, err := r.queries.GetMaxBtcAccountKeyIndex(ctx, sqlcgen.GetMaxBtcAccountKeyIndexParams{
		Coin:    r.coinTypeCode.String(),
		Account: accountType.String(),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call GetMaxBtcAccountKeyIndex(): %w", err)
	}

	// Type assert interface{} to int64
	if maxIdx, ok := result.(int64); ok {
		return maxIdx, nil
	}

	return 0, nil
}

// GetOneMaxID returns one record by max id
func (r *BTCAccountKeyRepositorySqlc) GetOneMaxID(accountType domainAccount.AccountType,
) (*domainBTC.BTCAccountKey, error) {
	ctx := context.Background()

	accountKey, err := r.queries.GetOneBtcAccountKeyByMaxID(ctx, sqlcgen.GetOneBtcAccountKeyByMaxIDParams{
		Coin:    r.coinTypeCode.String(),
		Account: accountType.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetOneBtcAccountKeyByMaxID(): %w", err)
	}

	return convertGetOneBtcAccountKeyByMaxIDRow(&accountKey)
}

// GetAllAddrStatus returns all BtcAccountKey by addr_status
func (r *BTCAccountKeyRepositorySqlc) GetAllAddrStatus(
	accountType domainAccount.AccountType, addrStatus domainAddress.AddrStatus,
) ([]*domainBTC.BTCAccountKey, error) {
	ctx := context.Background()

	accountKeys, err := r.queries.GetBtcAccountKeysByAddrStatus(ctx, sqlcgen.GetBtcAccountKeysByAddrStatusParams{
		Coin:       r.coinTypeCode.String(),
		Account:    accountType.String(),
		AddrStatus: int16(addrStatus.Int8()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetBtcAccountKeysByAddrStatus(): %w", err)
	}

	result := make([]*domainBTC.BTCAccountKey, 0, len(accountKeys))
	for i := range accountKeys {
		domainKey, err := convertGetBtcAccountKeysByAddrStatusRow(&accountKeys[i])
		if err != nil {
			return nil, fmt.Errorf("failed to convert account key at index %d: %w", i, err)
		}
		result = append(result, domainKey)
	}

	return result, nil
}

// GetAllMultiAddr returns all BtcAccountKey by multisig_address
func (r *BTCAccountKeyRepositorySqlc) GetAllMultiAddr(
	accountType domainAccount.AccountType, addrs []string,
) ([]*domainBTC.BTCAccountKey, error) {
	ctx := context.Background()

	accountKeys, err := r.queries.GetBtcAccountKeysByMultisigAddresses(
		ctx,
		sqlcgen.GetBtcAccountKeysByMultisigAddressesParams{
			Coin:    r.coinTypeCode.String(),
			Account: accountType.String(),
			Addrs:   addrs,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetBtcAccountKeysByMultisigAddresses(): %w", err)
	}

	result := make([]*domainBTC.BTCAccountKey, 0, len(accountKeys))
	for i := range accountKeys {
		domainKey, err := convertGetBtcAccountKeysByMultisigAddressesRow(&accountKeys[i])
		if err != nil {
			return nil, fmt.Errorf("failed to convert account key at index %d: %w", i, err)
		}
		result = append(result, domainKey)
	}

	return result, nil
}

// InsertBulk inserts multiple records
func (r *BTCAccountKeyRepositorySqlc) InsertBulk(items []*domainBTC.BTCAccountKey) error {
	ctx := context.Background()

	// Begin transaction to ensure atomicity of bulk insert
	tx, err := r.dbConn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	qtx := r.queries.WithTx(tx)

	for _, item := range items {
		sqlcItem := convertFromBTCAccountKey(item)
		_, err := qtx.InsertBtcAccountKey(ctx, sqlcgen.InsertBtcAccountKeyParams{
			Coin:                   sqlcItem.Coin,
			KeyType:                sqlcItem.KeyType,
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
			_ = tx.Rollback()
			return fmt.Errorf("failed to call InsertBtcAccountKey(): %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// UpdateAddr updates address by P2SHSegWitAddr
func (r *BTCAccountKeyRepositorySqlc) UpdateAddr(
	accountType domainAccount.AccountType, addr, keyAddress string,
) (int64, error) {
	ctx := context.Background()

	// Convert keyAddress to sql.NullString (it might be empty for BCH)
	var p2shSegwitAddr sql.NullString
	if keyAddress != "" {
		p2shSegwitAddr = sql.NullString{String: keyAddress, Valid: true}
	}

	result, err := r.queries.UpdateBtcAccountKeyAddress(ctx, sqlcgen.UpdateBtcAccountKeyAddressParams{
		P2pkhAddress:      addr,
		UpdatedAt:         sql.NullTime{Time: time.Now(), Valid: true},
		Coin:              r.coinTypeCode.String(),
		Account:           accountType.String(),
		P2shSegwitAddress: p2shSegwitAddr,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call UpdateBtcAccountKeyAddress(): %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
	}

	return rowsAffected, nil
}

// UpdateAddrStatus updates addr_status
func (r *BTCAccountKeyRepositorySqlc) UpdateAddrStatus(
	accountType domainAccount.AccountType, addrStatus domainAddress.AddrStatus, strWIFs []string,
) (int64, error) {
	ctx := context.Background()
	var totalAffected int64

	// sqlc doesn't support IN clauses with variable arguments, so update one at a time
	for _, wif := range strWIFs {
		result, err := r.queries.UpdateBtcAccountKeyAddrStatus(ctx, sqlcgen.UpdateBtcAccountKeyAddrStatusParams{
			AddrStatus:         int16(addrStatus.Int8()),
			UpdatedAt:          sql.NullTime{Time: time.Now(), Valid: true},
			Coin:               r.coinTypeCode.String(),
			Account:            accountType.String(),
			WalletImportFormat: wif,
		})
		if err != nil {
			return 0, fmt.Errorf("failed to call UpdateBtcAccountKeyAddrStatus(): %w", err)
		}

		affected, err := result.RowsAffected()
		if err != nil {
			return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
		}
		totalAffected += affected
	}

	return totalAffected, nil
}

// UpdateMultisigAddr updates multisig_address
func (r *BTCAccountKeyRepositorySqlc) UpdateMultisigAddr(
	accountType domainAccount.AccountType, item *domainBTC.BTCAccountKey,
) (int64, error) {
	ctx := context.Background()

	sqlcItem := convertFromBTCAccountKey(item)
	result, err := r.queries.UpdateBtcAccountKeyMultisigAddr(ctx, sqlcgen.UpdateBtcAccountKeyMultisigAddrParams{
		MultisigAddress: sqlcItem.MultisigAddress,
		RedeemScript:    sqlcItem.RedeemScript,
		AddrStatus:      sqlcItem.AddrStatus,
		UpdatedAt:       sql.NullTime{Time: time.Now(), Valid: true},
		Coin:            r.coinTypeCode.String(),
		Account:         accountType.String(),
		FullPublicKey:   sqlcItem.FullPublicKey,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call UpdateBtcAccountKeyMultisigAddr(): %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
	}

	return rowsAffected, nil
}

// UpdateMultisigAddrs updates all multisig_address with transaction
func (r *BTCAccountKeyRepositorySqlc) UpdateMultisigAddrs(
	accountType domainAccount.AccountType, items []*domainBTC.BTCAccountKey,
) (int64, error) {
	ctx := context.Background()

	// transaction
	dtx, err := r.dbConn.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to call db.Begin(): %w", err)
	}
	defer func() {
		if err != nil {
			_ = dtx.Rollback() // Error already being handled
		} else {
			_ = dtx.Commit() // Error already being handled
		}
	}()

	qtx := r.queries.WithTx(dtx)
	var totalAffected int64

	for _, item := range items {
		sqlcItem := convertFromBTCAccountKey(item)
		result, updateErr := qtx.UpdateBtcAccountKeyMultisigAddr(ctx, sqlcgen.UpdateBtcAccountKeyMultisigAddrParams{
			MultisigAddress: sqlcItem.MultisigAddress,
			RedeemScript:    sqlcItem.RedeemScript,
			AddrStatus:      sqlcItem.AddrStatus,
			UpdatedAt:       sql.NullTime{Time: time.Now(), Valid: true},
			Coin:            r.coinTypeCode.String(),
			Account:         accountType.String(),
			FullPublicKey:   sqlcItem.FullPublicKey,
		})
		if updateErr != nil {
			return 0, fmt.Errorf("failed to call UpdateBtcAccountKeyMultisigAddr(): %w", updateErr)
		}

		affected, affectedErr := result.RowsAffected()
		if affectedErr != nil {
			return 0, fmt.Errorf("failed to get RowsAffected(): %w", affectedErr)
		}
		totalAffected += affected
	}

	return totalAffected, nil
}

// NewAccountKeyRepositorySqlc is kept for backward compatibility.
//
// Deprecated: Use NewBTCAccountKeyRepositorySqlc instead.
func NewAccountKeyRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *BTCAccountKeyRepositorySqlc {
	return NewBTCAccountKeyRepositorySqlc(dbConn, coinTypeCode)
}
