package cold

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	sqlc "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/address"
)

// BTCAccountKeyRepositorySqlc is repository for btc_account_key table using sqlc
type BTCAccountKeyRepositorySqlc struct {
	queries      *sqlc.Queries
	dbConn       *sql.DB
	coinTypeCode domainCoin.CoinTypeCode
}

// NewBTCAccountKeyRepositorySqlc returns BTCAccountKeyRepositorySqlc object
func NewBTCAccountKeyRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *BTCAccountKeyRepositorySqlc {
	return &BTCAccountKeyRepositorySqlc{
		queries:      sqlc.New(dbConn),
		dbConn:       dbConn,
		coinTypeCode: coinTypeCode,
	}
}

// GetMaxIndex returns max idx
func (r *BTCAccountKeyRepositorySqlc) GetMaxIndex(accountType domainAccount.AccountType) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.GetMaxBtcAccountKeyIndex(ctx, sqlc.GetMaxBtcAccountKeyIndexParams{
		Coin:    sqlc.BtcAccountKeyCoin(r.coinTypeCode.String()),
		Account: sqlc.BtcAccountKeyAccount(accountType.String()),
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
func (r *BTCAccountKeyRepositorySqlc) GetOneMaxID(accountType domainAccount.AccountType) (*sqlc.BtcAccountKey, error) {
	ctx := context.Background()

	accountKey, err := r.queries.GetOneBtcAccountKeyByMaxID(ctx, sqlc.GetOneBtcAccountKeyByMaxIDParams{
		Coin:    sqlc.BtcAccountKeyCoin(r.coinTypeCode.String()),
		Account: sqlc.BtcAccountKeyAccount(accountType.String()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetOneBtcAccountKeyByMaxID(): %w", err)
	}

	return &accountKey, nil
}

// GetAllAddrStatus returns all BtcAccountKey by addr_status
func (r *BTCAccountKeyRepositorySqlc) GetAllAddrStatus(
	accountType domainAccount.AccountType, addrStatus address.AddrStatus,
) ([]*sqlc.BtcAccountKey, error) {
	ctx := context.Background()

	accountKeys, err := r.queries.GetBtcAccountKeysByAddrStatus(ctx, sqlc.GetBtcAccountKeysByAddrStatusParams{
		Coin:       sqlc.BtcAccountKeyCoin(r.coinTypeCode.String()),
		Account:    sqlc.BtcAccountKeyAccount(accountType.String()),
		AddrStatus: addrStatus.Int8(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetBtcAccountKeysByAddrStatus(): %w", err)
	}

	result := make([]*sqlc.BtcAccountKey, len(accountKeys))
	for i := range accountKeys {
		result[i] = &accountKeys[i]
	}

	return result, nil
}

// GetAllMultiAddr returns all BtcAccountKey by multisig_address
func (r *BTCAccountKeyRepositorySqlc) GetAllMultiAddr(
	accountType domainAccount.AccountType, addrs []string,
) ([]*sqlc.BtcAccountKey, error) {
	ctx := context.Background()

	accountKeys, err := r.queries.GetBtcAccountKeysByMultisigAddresses(
		ctx,
		sqlc.GetBtcAccountKeysByMultisigAddressesParams{
			Coin:    sqlc.BtcAccountKeyCoin(r.coinTypeCode.String()),
			Account: sqlc.BtcAccountKeyAccount(accountType.String()),
			Addrs:   addrs,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetBtcAccountKeysByMultisigAddresses(): %w", err)
	}

	result := make([]*sqlc.BtcAccountKey, len(accountKeys))
	for i := range accountKeys {
		result[i] = &accountKeys[i]
	}

	return result, nil
}

// InsertBulk inserts multiple records
func (r *BTCAccountKeyRepositorySqlc) InsertBulk(items []*sqlc.BtcAccountKey) error {
	ctx := context.Background()

	for _, item := range items {
		_, err := r.queries.InsertBtcAccountKey(ctx, sqlc.InsertBtcAccountKeyParams{
			Coin:               item.Coin,
			KeyType:            item.KeyType,
			Account:            item.Account,
			P2pkhAddress:       item.P2pkhAddress,
			P2shSegwitAddress:  item.P2shSegwitAddress,
			Bech32Address:      item.Bech32Address,
			TaprootAddress:     item.TaprootAddress,
			FullPublicKey:      item.FullPublicKey,
			MultisigAddress:    item.MultisigAddress,
			RedeemScript:       item.RedeemScript,
			WalletImportFormat: item.WalletImportFormat,
			Idx:                item.Idx,
			AddrStatus:         item.AddrStatus,
		})
		if err != nil {
			return fmt.Errorf("failed to call InsertBtcAccountKey(): %w", err)
		}
	}

	return nil
}

// UpdateAddr updates address by P2SHSegWitAddr
func (r *BTCAccountKeyRepositorySqlc) UpdateAddr(
	accountType domainAccount.AccountType, addr, keyAddress string,
) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateBtcAccountKeyAddress(ctx, sqlc.UpdateBtcAccountKeyAddressParams{
		P2pkhAddress:      addr,
		UpdatedAt:         sql.NullTime{Time: time.Now(), Valid: true},
		Coin:              sqlc.BtcAccountKeyCoin(r.coinTypeCode.String()),
		Account:           sqlc.BtcAccountKeyAccount(accountType.String()),
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
	accountType domainAccount.AccountType, addrStatus address.AddrStatus, strWIFs []string,
) (int64, error) {
	ctx := context.Background()
	var totalAffected int64

	// sqlc doesn't support IN clauses with variable arguments, so update one at a time
	for _, wif := range strWIFs {
		result, err := r.queries.UpdateBtcAccountKeyAddrStatus(ctx, sqlc.UpdateBtcAccountKeyAddrStatusParams{
			AddrStatus:         addrStatus.Int8(),
			UpdatedAt:          sql.NullTime{Time: time.Now(), Valid: true},
			Coin:               sqlc.BtcAccountKeyCoin(r.coinTypeCode.String()),
			Account:            sqlc.BtcAccountKeyAccount(accountType.String()),
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
	accountType domainAccount.AccountType, item *sqlc.BtcAccountKey,
) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateBtcAccountKeyMultisigAddr(ctx, sqlc.UpdateBtcAccountKeyMultisigAddrParams{
		MultisigAddress: item.MultisigAddress,
		RedeemScript:    item.RedeemScript,
		AddrStatus:      item.AddrStatus,
		UpdatedAt:       sql.NullTime{Time: time.Now(), Valid: true},
		Coin:            sqlc.BtcAccountKeyCoin(r.coinTypeCode.String()),
		Account:         sqlc.BtcAccountKeyAccount(accountType.String()),
		FullPublicKey:   item.FullPublicKey,
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
	accountType domainAccount.AccountType, items []*sqlc.BtcAccountKey,
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
		result, updateErr := qtx.UpdateBtcAccountKeyMultisigAddr(ctx, sqlc.UpdateBtcAccountKeyMultisigAddrParams{
			MultisigAddress: item.MultisigAddress,
			RedeemScript:    item.RedeemScript,
			AddrStatus:      item.AddrStatus,
			UpdatedAt:       sql.NullTime{Time: time.Now(), Valid: true},
			Coin:            sqlc.BtcAccountKeyCoin(r.coinTypeCode.String()),
			Account:         sqlc.BtcAccountKeyAccount(accountType.String()),
			FullPublicKey:   item.FullPublicKey,
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
