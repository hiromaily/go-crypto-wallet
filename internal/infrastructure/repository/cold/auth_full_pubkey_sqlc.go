package cold

import (
	"context"
	"database/sql"
	"fmt"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAuth "github.com/hiromaily/go-crypto-wallet/internal/domain/auth"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
)

// AuthFullPubkeyRepositorySqlc is repository for auth_fullpubkey table using sqlc
type AuthFullPubkeyRepositorySqlc struct {
	queries      *sqlcgen.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewAuthFullPubkeyRepositorySqlc returns AuthFullPubkeyRepositorySqlc object
func NewAuthFullPubkeyRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *AuthFullPubkeyRepositorySqlc {
	return &AuthFullPubkeyRepositorySqlc{
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
	}
}

// convertToAuthFullPubkey converts sqlcgen.AuthFullpubkey to domain.AuthFullPubkey entity.
func convertToAuthFullPubkey(sqlcKey *sqlcgen.AuthFullpubkey) (*domainAuth.AuthFullPubkey, error) {
	key := &domainAuth.AuthFullPubkey{
		ID:            sqlcKey.ID,
		CoinTypeCode:  domainCoin.CoinTypeCode(sqlcKey.Coin),
		AuthAccount:   domainAccount.AuthType(sqlcKey.AuthAccount),
		FullPublicKey: sqlcKey.FullPublicKey,
	}

	if sqlcKey.ExtendedPubkey.Valid {
		key.ExtendedPubKey = sqlcKey.ExtendedPubkey.String
	}
	if sqlcKey.Fingerprint.Valid {
		fingerprint, err := domainKey.NewFingerprint(sqlcKey.Fingerprint.String)
		if err != nil {
			return nil, fmt.Errorf("invalid fingerprint in database: %w", err)
		}
		key.Fingerprint = &fingerprint
	}
	if sqlcKey.DerivationPath.Valid {
		key.DerivationPath = sqlcKey.DerivationPath.String
	}
	if sqlcKey.UpdatedAt.Valid {
		key.UpdatedAt = &sqlcKey.UpdatedAt.Time
	}

	return key, nil
}

// convertFromAuthFullPubkey converts domain.AuthFullPubkey entity to sqlcgen.AuthFullpubkey.
func convertFromAuthFullPubkey(key *domainAuth.AuthFullPubkey) *sqlcgen.AuthFullpubkey {
	sqlcKey := &sqlcgen.AuthFullpubkey{
		ID:            key.ID,
		Coin:          sqlcgen.AuthFullpubkeyCoin(key.CoinTypeCode.String()),
		AuthAccount:   key.AuthAccount.String(),
		FullPublicKey: key.FullPublicKey,
	}

	if key.ExtendedPubKey != "" {
		sqlcKey.ExtendedPubkey = sql.NullString{String: key.ExtendedPubKey, Valid: true}
	}
	if key.Fingerprint != nil {
		sqlcKey.Fingerprint = sql.NullString{String: key.Fingerprint.String(), Valid: true}
	}
	if key.DerivationPath != "" {
		sqlcKey.DerivationPath = sql.NullString{String: key.DerivationPath, Valid: true}
	}
	if key.UpdatedAt != nil {
		sqlcKey.UpdatedAt = sql.NullTime{Time: *key.UpdatedAt, Valid: true}
	}

	return sqlcKey
}

// GetOne returns one record by authType
func (r *AuthFullPubkeyRepositorySqlc) GetOne(authType domainAccount.AuthType) (*domainAuth.AuthFullPubkey, error) {
	ctx := context.Background()

	authPubkey, err := r.queries.GetAuthFullPubkey(ctx, sqlcgen.GetAuthFullPubkeyParams{
		Coin:        sqlcgen.AuthFullpubkeyCoin(r.coinTypeCode.String()),
		AuthAccount: authType.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetAuthFullPubkey(): %w", err)
	}

	return convertToAuthFullPubkey(&authPubkey)
}

// Insert inserts record
func (r *AuthFullPubkeyRepositorySqlc) Insert(authType domainAccount.AuthType, fullPubKey string) error {
	ctx := context.Background()

	_, err := r.queries.InsertAuthFullPubkey(ctx, sqlcgen.InsertAuthFullPubkeyParams{
		Coin:          sqlcgen.AuthFullpubkeyCoin(r.coinTypeCode.String()),
		AuthAccount:   authType.String(),
		FullPublicKey: fullPubKey,
	})
	if err != nil {
		return fmt.Errorf("failed to call InsertAuthFullPubkey(): %w", err)
	}

	return nil
}

// InsertBulk inserts multiple records
func (r *AuthFullPubkeyRepositorySqlc) InsertBulk(items []*domainAuth.AuthFullPubkey) error {
	ctx := context.Background()

	for _, item := range items {
		sqlcItem := convertFromAuthFullPubkey(item)
		_, err := r.queries.InsertAuthFullPubkey(ctx, sqlcgen.InsertAuthFullPubkeyParams{
			Coin:           sqlcItem.Coin,
			AuthAccount:    sqlcItem.AuthAccount,
			FullPublicKey:  sqlcItem.FullPublicKey,
			ExtendedPubkey: sqlcItem.ExtendedPubkey,
			Fingerprint:    sqlcItem.Fingerprint,
			DerivationPath: sqlcItem.DerivationPath,
		})
		if err != nil {
			return fmt.Errorf("failed to call InsertAuthFullPubkey(): %w", err)
		}
	}

	return nil
}
