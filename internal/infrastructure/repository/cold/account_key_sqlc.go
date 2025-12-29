package cold

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/sqlc"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/address"
)

// AccountKeyRepositorySqlc is repository for account_key table using sqlc
type AccountKeyRepositorySqlc struct {
	queries      *sqlc.Queries
	dbConn       *sql.DB
	coinTypeCode domainCoin.CoinTypeCode
}

// NewAccountKeyRepositorySqlc returns AccountKeyRepositorySqlc object
func NewAccountKeyRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *AccountKeyRepositorySqlc {
	return &AccountKeyRepositorySqlc{
		queries:      sqlc.New(dbConn),
		dbConn:       dbConn,
		coinTypeCode: coinTypeCode,
	}
}

// GetMaxIndex returns max idx
func (r *AccountKeyRepositorySqlc) GetMaxIndex(accountType domainAccount.AccountType) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.GetMaxAccountKeyIndex(ctx, sqlc.GetMaxAccountKeyIndexParams{
		Coin:    sqlc.AccountKeyCoin(r.coinTypeCode.String()),
		Account: sqlc.AccountKeyAccount(accountType.String()),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call GetMaxAccountKeyIndex(): %w", err)
	}

	// Type assert interface{} to int64
	if maxIdx, ok := result.(int64); ok {
		return maxIdx, nil
	}

	return 0, nil
}

// GetOneMaxID returns one record by max id
func (r *AccountKeyRepositorySqlc) GetOneMaxID(accountType domainAccount.AccountType) (*sqlc.AccountKey, error) {
	ctx := context.Background()

	accountKey, err := r.queries.GetOneAccountKeyByMaxID(ctx, sqlc.GetOneAccountKeyByMaxIDParams{
		Coin:    sqlc.AccountKeyCoin(r.coinTypeCode.String()),
		Account: sqlc.AccountKeyAccount(accountType.String()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetOneAccountKeyByMaxID(): %w", err)
	}

	return &accountKey, nil
}

// GetAllAddrStatus returns all AccountKey by addr_status
func (r *AccountKeyRepositorySqlc) GetAllAddrStatus(
	accountType domainAccount.AccountType, addrStatus address.AddrStatus,
) ([]*sqlc.AccountKey, error) {
	ctx := context.Background()

	accountKeys, err := r.queries.GetAccountKeysByAddrStatus(ctx, sqlc.GetAccountKeysByAddrStatusParams{
		Coin:       sqlc.AccountKeyCoin(r.coinTypeCode.String()),
		Account:    sqlc.AccountKeyAccount(accountType.String()),
		AddrStatus: addrStatus.Int8(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetAccountKeysByAddrStatus(): %w", err)
	}

	result := make([]*sqlc.AccountKey, len(accountKeys))
	for i := range accountKeys {
		result[i] = &accountKeys[i]
	}

	return result, nil
}

// GetAllMultiAddr returns all AccountKey by multisig_address
func (r *AccountKeyRepositorySqlc) GetAllMultiAddr(
	accountType domainAccount.AccountType, addrs []string,
) ([]*sqlc.AccountKey, error) {
	ctx := context.Background()

	accountKeys, err := r.queries.GetAccountKeysByMultisigAddresses(
		ctx,
		sqlc.GetAccountKeysByMultisigAddressesParams{
			Coin:    sqlc.AccountKeyCoin(r.coinTypeCode.String()),
			Account: sqlc.AccountKeyAccount(accountType.String()),
			Addrs:   addrs,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetAccountKeysByMultisigAddresses(): %w", err)
	}

	result := make([]*sqlc.AccountKey, len(accountKeys))
	for i := range accountKeys {
		result[i] = &accountKeys[i]
	}

	return result, nil
}

// InsertBulk inserts multiple records
func (r *AccountKeyRepositorySqlc) InsertBulk(items []*sqlc.AccountKey) error {
	ctx := context.Background()

	for _, item := range items {
		_, err := r.queries.InsertAccountKey(ctx, sqlc.InsertAccountKeyParams{
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
			return fmt.Errorf("failed to call InsertAccountKey(): %w", err)
		}
	}

	return nil
}

// UpdateAddr updates address by P2SHSegWitAddr
func (r *AccountKeyRepositorySqlc) UpdateAddr(
	accountType domainAccount.AccountType, addr, keyAddress string,
) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateAccountKeyAddress(ctx, sqlc.UpdateAccountKeyAddressParams{
		P2pkhAddress:      addr,
		UpdatedAt:         sql.NullTime{Time: time.Now(), Valid: true},
		Coin:              sqlc.AccountKeyCoin(r.coinTypeCode.String()),
		Account:           sqlc.AccountKeyAccount(accountType.String()),
		P2shSegwitAddress: keyAddress,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call UpdateAccountKeyAddress(): %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
	}

	return rowsAffected, nil
}

// UpdateAddrStatus updates addr_status
func (r *AccountKeyRepositorySqlc) UpdateAddrStatus(
	accountType domainAccount.AccountType, addrStatus address.AddrStatus, strWIFs []string,
) (int64, error) {
	ctx := context.Background()
	var totalAffected int64

	// sqlc doesn't support IN clauses with variable arguments, so update one at a time
	for _, wif := range strWIFs {
		result, err := r.queries.UpdateAccountKeyAddrStatus(ctx, sqlc.UpdateAccountKeyAddrStatusParams{
			AddrStatus:         addrStatus.Int8(),
			UpdatedAt:          sql.NullTime{Time: time.Now(), Valid: true},
			Coin:               sqlc.AccountKeyCoin(r.coinTypeCode.String()),
			Account:            sqlc.AccountKeyAccount(accountType.String()),
			WalletImportFormat: wif,
		})
		if err != nil {
			return 0, fmt.Errorf("failed to call UpdateAccountKeyAddrStatus(): %w", err)
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
func (r *AccountKeyRepositorySqlc) UpdateMultisigAddr(
	accountType domainAccount.AccountType, item *sqlc.AccountKey,
) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateAccountKeyMultisigAddr(ctx, sqlc.UpdateAccountKeyMultisigAddrParams{
		MultisigAddress: item.MultisigAddress,
		RedeemScript:    item.RedeemScript,
		AddrStatus:      item.AddrStatus,
		UpdatedAt:       sql.NullTime{Time: time.Now(), Valid: true},
		Coin:            sqlc.AccountKeyCoin(r.coinTypeCode.String()),
		Account:         sqlc.AccountKeyAccount(accountType.String()),
		FullPublicKey:   item.FullPublicKey,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call UpdateAccountKeyMultisigAddr(): %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
	}

	return rowsAffected, nil
}

// UpdateMultisigAddrs updates all multisig_address with transaction
func (r *AccountKeyRepositorySqlc) UpdateMultisigAddrs(
	accountType domainAccount.AccountType, items []*sqlc.AccountKey,
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
		result, updateErr := qtx.UpdateAccountKeyMultisigAddr(ctx, sqlc.UpdateAccountKeyMultisigAddrParams{
			MultisigAddress: item.MultisigAddress,
			RedeemScript:    item.RedeemScript,
			AddrStatus:      item.AddrStatus,
			UpdatedAt:       sql.NullTime{Time: time.Now(), Valid: true},
			Coin:            sqlc.AccountKeyCoin(r.coinTypeCode.String()),
			Account:         sqlc.AccountKeyAccount(accountType.String()),
			FullPublicKey:   item.FullPublicKey,
		})
		if updateErr != nil {
			return 0, fmt.Errorf("failed to call UpdateAccountKeyMultisigAddr(): %w", updateErr)
		}

		affected, affectedErr := result.RowsAffected()
		if affectedErr != nil {
			return 0, fmt.Errorf("failed to get RowsAffected(): %w", affectedErr)
		}
		totalAffected += affected
	}

	return totalAffected, nil
}
