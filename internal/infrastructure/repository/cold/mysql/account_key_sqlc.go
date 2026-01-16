package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainBitcoin "github.com/hiromaily/go-crypto-wallet/internal/domain/bitcoin"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
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

// convertToBtcAccountKey converts sqlcgen.BtcAccountKey to domain.BtcAccountKey entity.
// SECURITY: Handles WIF (private key) data - never log the wallet import format field.
func convertToBtcAccountKey(sqlcKey *sqlcgen.BtcAccountKey) (*domainBitcoin.BtcAccountKey, error) {
	addrStatus, err := domainAddress.AddrStatusFromInt8(sqlcKey.AddrStatus)
	if err != nil {
		return nil, fmt.Errorf("invalid addr status in database: %w", err)
	}

	coinTypeCode := domainCoin.CoinTypeCode(sqlcKey.Coin)
	if !domainCoin.IsCoinTypeCode(string(coinTypeCode)) {
		return nil, fmt.Errorf("invalid coin type from database: %s", sqlcKey.Coin)
	}

	accountType := domainAccount.AccountType(sqlcKey.Account)
	if !domainAccount.ValidateAccountType(string(accountType)) {
		return nil, fmt.Errorf("invalid account type from database: %s", sqlcKey.Account)
	}

	key := &domainBitcoin.BtcAccountKey{
		ID:                 sqlcKey.ID,
		CoinTypeCode:       coinTypeCode,
		KeyType:            sqlcKey.KeyType,
		Account:            accountType,
		P2pkhAddress:       sqlcKey.P2pkhAddress,
		P2shSegwitAddress:  sqlcKey.P2shSegwitAddress,
		Bech32Address:      sqlcKey.Bech32Address,
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
		key.AccountExtendedPrivkey = &sqlcKey.AccountExtendedPrivkey.String // NEVER log
	}
	if sqlcKey.UpdatedAt.Valid {
		key.UpdatedAt = &sqlcKey.UpdatedAt.Time
	}

	return key, nil
}

// convertFromBtcAccountKey converts domain.BtcAccountKey entity to sqlcgen.BtcAccountKey.
func convertFromBtcAccountKey(key *domainBitcoin.BtcAccountKey) *sqlcgen.BtcAccountKey {
	sqlcKey := &sqlcgen.BtcAccountKey{
		ID:                 key.ID,
		Coin:               sqlcgen.BtcAccountKeyCoin(key.CoinTypeCode.String()),
		KeyType:            key.KeyType,
		Account:            sqlcgen.BtcAccountKeyAccount(key.Account.String()),
		P2pkhAddress:       key.P2pkhAddress,
		P2shSegwitAddress:  key.P2shSegwitAddress,
		Bech32Address:      key.Bech32Address,
		FullPublicKey:      key.FullPublicKey,
		MultisigAddress:    key.MultisigAddress,
		RedeemScript:       key.RedeemScript,
		WalletImportFormat: key.WalletImportFormat,
		Idx:                key.Idx,
		AddrStatus:         key.AddrStatus.Int8(),
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
func (r *BTCAccountKeyRepositorySqlc) GetMaxIndex(accountType domainAccount.AccountType) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.GetMaxBtcAccountKeyIndex(ctx, sqlcgen.GetMaxBtcAccountKeyIndexParams{
		Coin:    sqlcgen.BtcAccountKeyCoin(r.coinTypeCode.String()),
		Account: sqlcgen.BtcAccountKeyAccount(accountType.String()),
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
) (*domainBitcoin.BtcAccountKey, error) {
	ctx := context.Background()

	accountKey, err := r.queries.GetOneBtcAccountKeyByMaxID(ctx, sqlcgen.GetOneBtcAccountKeyByMaxIDParams{
		Coin:    sqlcgen.BtcAccountKeyCoin(r.coinTypeCode.String()),
		Account: sqlcgen.BtcAccountKeyAccount(accountType.String()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetOneBtcAccountKeyByMaxID(): %w", err)
	}

	return convertToBtcAccountKey(&accountKey)
}

// GetAllAddrStatus returns all BtcAccountKey by addr_status
func (r *BTCAccountKeyRepositorySqlc) GetAllAddrStatus(
	accountType domainAccount.AccountType, addrStatus domainAddress.AddrStatus,
) ([]*domainBitcoin.BtcAccountKey, error) {
	ctx := context.Background()

	accountKeys, err := r.queries.GetBtcAccountKeysByAddrStatus(ctx, sqlcgen.GetBtcAccountKeysByAddrStatusParams{
		Coin:       sqlcgen.BtcAccountKeyCoin(r.coinTypeCode.String()),
		Account:    sqlcgen.BtcAccountKeyAccount(accountType.String()),
		AddrStatus: addrStatus.Int8(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetBtcAccountKeysByAddrStatus(): %w", err)
	}

	result := make([]*domainBitcoin.BtcAccountKey, 0, len(accountKeys))
	for i := range accountKeys {
		domainKey, err := convertToBtcAccountKey(&accountKeys[i])
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
) ([]*domainBitcoin.BtcAccountKey, error) {
	ctx := context.Background()

	accountKeys, err := r.queries.GetBtcAccountKeysByMultisigAddresses(
		ctx,
		sqlcgen.GetBtcAccountKeysByMultisigAddressesParams{
			Coin:    sqlcgen.BtcAccountKeyCoin(r.coinTypeCode.String()),
			Account: sqlcgen.BtcAccountKeyAccount(accountType.String()),
			Addrs:   addrs,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetBtcAccountKeysByMultisigAddresses(): %w", err)
	}

	result := make([]*domainBitcoin.BtcAccountKey, 0, len(accountKeys))
	for i := range accountKeys {
		domainKey, err := convertToBtcAccountKey(&accountKeys[i])
		if err != nil {
			return nil, fmt.Errorf("failed to convert account key at index %d: %w", i, err)
		}
		result = append(result, domainKey)
	}

	return result, nil
}

// InsertBulk inserts multiple records
func (r *BTCAccountKeyRepositorySqlc) InsertBulk(items []*domainBitcoin.BtcAccountKey) error {
	ctx := context.Background()

	// Begin transaction to ensure atomicity of bulk insert
	tx, err := r.dbConn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	qtx := r.queries.WithTx(tx)

	for _, item := range items {
		sqlcItem := convertFromBtcAccountKey(item)
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

	result, err := r.queries.UpdateBtcAccountKeyAddress(ctx, sqlcgen.UpdateBtcAccountKeyAddressParams{
		P2pkhAddress:      addr,
		UpdatedAt:         sql.NullTime{Time: time.Now(), Valid: true},
		Coin:              sqlcgen.BtcAccountKeyCoin(r.coinTypeCode.String()),
		Account:           sqlcgen.BtcAccountKeyAccount(accountType.String()),
		P2shSegwitAddress: keyAddress,
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
			AddrStatus:         addrStatus.Int8(),
			UpdatedAt:          sql.NullTime{Time: time.Now(), Valid: true},
			Coin:               sqlcgen.BtcAccountKeyCoin(r.coinTypeCode.String()),
			Account:            sqlcgen.BtcAccountKeyAccount(accountType.String()),
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
	accountType domainAccount.AccountType, item *domainBitcoin.BtcAccountKey,
) (int64, error) {
	ctx := context.Background()

	sqlcItem := convertFromBtcAccountKey(item)
	result, err := r.queries.UpdateBtcAccountKeyMultisigAddr(ctx, sqlcgen.UpdateBtcAccountKeyMultisigAddrParams{
		MultisigAddress: sqlcItem.MultisigAddress,
		RedeemScript:    sqlcItem.RedeemScript,
		AddrStatus:      sqlcItem.AddrStatus,
		UpdatedAt:       sql.NullTime{Time: time.Now(), Valid: true},
		Coin:            sqlcgen.BtcAccountKeyCoin(r.coinTypeCode.String()),
		Account:         sqlcgen.BtcAccountKeyAccount(accountType.String()),
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
	accountType domainAccount.AccountType, items []*domainBitcoin.BtcAccountKey,
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
		sqlcItem := convertFromBtcAccountKey(item)
		result, updateErr := qtx.UpdateBtcAccountKeyMultisigAddr(ctx, sqlcgen.UpdateBtcAccountKeyMultisigAddrParams{
			MultisigAddress: sqlcItem.MultisigAddress,
			RedeemScript:    sqlcItem.RedeemScript,
			AddrStatus:      sqlcItem.AddrStatus,
			UpdatedAt:       sql.NullTime{Time: time.Now(), Valid: true},
			Coin:            sqlcgen.BtcAccountKeyCoin(r.coinTypeCode.String()),
			Account:         sqlcgen.BtcAccountKeyAccount(accountType.String()),
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

// AccountKeyRepositorySqlc is kept for backward compatibility.
//
// Deprecated: Use BTCAccountKeyRepositorySqlc instead.
type AccountKeyRepositorySqlc = BTCAccountKeyRepositorySqlc
