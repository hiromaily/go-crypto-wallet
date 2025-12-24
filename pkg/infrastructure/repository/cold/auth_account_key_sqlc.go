package cold

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/database/sqlc"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
)

// AuthAccountKeyRepositorySqlc is repository for auth_account_key table using sqlc
type AuthAccountKeyRepositorySqlc struct {
	queries      *sqlc.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewAuthAccountKeyRepositorySqlc returns AuthAccountKeyRepositorySqlc object
func NewAuthAccountKeyRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *AuthAccountKeyRepositorySqlc {
	return &AuthAccountKeyRepositorySqlc{
		queries:      sqlc.New(dbConn),
		coinTypeCode: coinTypeCode,
	}
}

// GetOne returns one record by authType
func (r *AuthAccountKeyRepositorySqlc) GetOne(authType domainAccount.AuthType) (*models.AuthAccountKey, error) {
	ctx := context.Background()

	authKey, err := r.queries.GetAuthAccountKey(ctx, sqlc.GetAuthAccountKeyParams{
		Coin:        sqlc.AuthAccountKeyCoin(r.coinTypeCode.String()),
		AuthAccount: authType.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetAuthAccountKey(): %w", err)
	}

	return convertSqlcAuthAccountKeyToModel(&authKey), nil
}

// Insert inserts record
func (r *AuthAccountKeyRepositorySqlc) Insert(item *models.AuthAccountKey) error {
	ctx := context.Background()

	_, err := r.queries.InsertAuthAccountKey(ctx, sqlc.InsertAuthAccountKeyParams{
		Coin:               sqlc.AuthAccountKeyCoin(item.Coin),
		AuthAccount:        item.AuthAccount,
		P2pkhAddress:       item.P2PKHAddress,
		P2shSegwitAddress:  item.P2SHSegwitAddress,
		Bech32Address:      item.Bech32Address,
		FullPublicKey:      item.FullPublicKey,
		MultisigAddress:    item.MultisigAddress,
		RedeemScript:       item.RedeemScript,
		WalletImportFormat: item.WalletImportFormat,
		Idx:                item.Idx,
		AddrStatus:         item.AddrStatus,
	})
	if err != nil {
		return fmt.Errorf("failed to call InsertAuthAccountKey(): %w", err)
	}

	return nil
}

// UpdateAddrStatus updates addr_status
func (r *AuthAccountKeyRepositorySqlc) UpdateAddrStatus(addrStatus address.AddrStatus, strWIF string) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateAuthAccountKeyAddrStatus(ctx, sqlc.UpdateAuthAccountKeyAddrStatusParams{
		AddrStatus:         addrStatus.Int8(),
		UpdatedAt:          sql.NullTime{Time: time.Now(), Valid: true},
		Coin:               sqlc.AuthAccountKeyCoin(r.coinTypeCode.String()),
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

// Helper functions

func convertSqlcAuthAccountKeyToModel(authKey *sqlc.AuthAccountKey) *models.AuthAccountKey {
	return &models.AuthAccountKey{
		ID:                 authKey.ID,
		Coin:               string(authKey.Coin),
		AuthAccount:        authKey.AuthAccount,
		P2PKHAddress:       authKey.P2pkhAddress,
		P2SHSegwitAddress:  authKey.P2shSegwitAddress,
		Bech32Address:      authKey.Bech32Address,
		FullPublicKey:      authKey.FullPublicKey,
		MultisigAddress:    authKey.MultisigAddress,
		RedeemScript:       authKey.RedeemScript,
		WalletImportFormat: authKey.WalletImportFormat,
		Idx:                authKey.Idx,
		AddrStatus:         authKey.AddrStatus,
		UpdatedAt:          convertSQLNullTimeToNullTime(authKey.UpdatedAt),
	}
}
