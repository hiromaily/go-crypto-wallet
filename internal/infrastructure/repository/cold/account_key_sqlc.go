package cold

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
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
) (*sqlcgen.BtcAccountKey, error) {
	ctx := context.Background()

	accountKey, err := r.queries.GetOneBtcAccountKeyByMaxID(ctx, sqlcgen.GetOneBtcAccountKeyByMaxIDParams{
		Coin:    sqlcgen.BtcAccountKeyCoin(r.coinTypeCode.String()),
		Account: sqlcgen.BtcAccountKeyAccount(accountType.String()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetOneBtcAccountKeyByMaxID(): %w", err)
	}

	return &accountKey, nil
}

// GetAllAddrStatus returns all BtcAccountKey by addr_status
func (r *BTCAccountKeyRepositorySqlc) GetAllAddrStatus(
	accountType domainAccount.AccountType, addrStatus domainAddress.AddrStatus,
) ([]*sqlcgen.BtcAccountKey, error) {
	ctx := context.Background()

	accountKeys, err := r.queries.GetBtcAccountKeysByAddrStatus(ctx, sqlcgen.GetBtcAccountKeysByAddrStatusParams{
		Coin:       sqlcgen.BtcAccountKeyCoin(r.coinTypeCode.String()),
		Account:    sqlcgen.BtcAccountKeyAccount(accountType.String()),
		AddrStatus: addrStatus.Int8(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetBtcAccountKeysByAddrStatus(): %w", err)
	}

	result := make([]*sqlcgen.BtcAccountKey, len(accountKeys))
	for i := range accountKeys {
		result[i] = &accountKeys[i]
	}

	return result, nil
}

// GetAllMultiAddr returns all BtcAccountKey by multisig_address
func (r *BTCAccountKeyRepositorySqlc) GetAllMultiAddr(
	accountType domainAccount.AccountType, addrs []string,
) ([]*sqlcgen.BtcAccountKey, error) {
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

	result := make([]*sqlcgen.BtcAccountKey, len(accountKeys))
	for i := range accountKeys {
		result[i] = &accountKeys[i]
	}

	return result, nil
}

// InsertBulk inserts multiple records
func (r *BTCAccountKeyRepositorySqlc) InsertBulk(items []*sqlcgen.BtcAccountKey) error {
	ctx := context.Background()

	for _, item := range items {
		_, err := r.queries.InsertBtcAccountKey(ctx, sqlcgen.InsertBtcAccountKeyParams{
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
	accountType domainAccount.AccountType, item *sqlcgen.BtcAccountKey,
) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateBtcAccountKeyMultisigAddr(ctx, sqlcgen.UpdateBtcAccountKeyMultisigAddrParams{
		MultisigAddress: item.MultisigAddress,
		RedeemScript:    item.RedeemScript,
		AddrStatus:      item.AddrStatus,
		UpdatedAt:       sql.NullTime{Time: time.Now(), Valid: true},
		Coin:            sqlcgen.BtcAccountKeyCoin(r.coinTypeCode.String()),
		Account:         sqlcgen.BtcAccountKeyAccount(accountType.String()),
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
	accountType domainAccount.AccountType, items []*sqlcgen.BtcAccountKey,
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
		result, updateErr := qtx.UpdateBtcAccountKeyMultisigAddr(ctx, sqlcgen.UpdateBtcAccountKeyMultisigAddrParams{
			MultisigAddress: item.MultisigAddress,
			RedeemScript:    item.RedeemScript,
			AddrStatus:      item.AddrStatus,
			UpdatedAt:       sql.NullTime{Time: time.Now(), Valid: true},
			Coin:            sqlcgen.BtcAccountKeyCoin(r.coinTypeCode.String()),
			Account:         sqlcgen.BtcAccountKeyAccount(accountType.String()),
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
