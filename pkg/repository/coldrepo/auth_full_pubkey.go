package coldrepo

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	models "github.com/hiromaily/go-bitcoin/pkg/models/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/coin"
)

// AuthFullPubkeyRepositorier is AuthFullPubkeyRepository interface
type AuthFullPubkeyRepositorier interface {
	GetOne(authType account.AuthType) (*models.AuthFullpubkey, error)
	Insert(authType account.AuthType, fullPubKey string) error
	InsertBulk(items []*models.AuthFullpubkey) error
}

// AuthFullPubkeyRepository is repository  for auth_fullpubkey table
type AuthFullPubkeyRepository struct {
	dbConn       *sql.DB
	tableName    string
	coinTypeCode coin.CoinTypeCode
	logger       *zap.Logger
}

// NewAuthFullPubkeyRepository returns AuthFullPubkeyRepository object
func NewAuthFullPubkeyRepository(dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger *zap.Logger) *AuthFullPubkeyRepository {
	return &AuthFullPubkeyRepository{
		dbConn:       dbConn,
		tableName:    "auth_pubkey",
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetOne returns one records
func (r *AuthFullPubkeyRepository) GetOne(authType account.AuthType) (*models.AuthFullpubkey, error) {
	ctx := context.Background()

	item, err := models.AuthFullpubkeys(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("auth_account=?", authType.String()),
	).One(ctx, r.dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.AuthFullpubkeys().One()")
	}
	return item, nil
}

// Insert inserts record
func (r *AuthFullPubkeyRepository) Insert(authType account.AuthType, fullPubKey string) error {
	ctx := context.Background()
	item := &models.AuthFullpubkey{
		Coin:          r.coinTypeCode.String(),
		AuthAccount:   authType.String(),
		FullPublicKey: fullPubKey,
	}

	if err := item.Insert(ctx, r.dbConn, boil.Infer()); err != nil {
		return errors.Wrap(err, "failed to call item.Insert()")
	}

	return nil
}

// InsertBulk inserts multiple records
func (r *AuthFullPubkeyRepository) InsertBulk(items []*models.AuthFullpubkey) error {
	ctx := context.Background()
	return models.AuthFullpubkeySlice(items).InsertAll(ctx, r.dbConn, boil.Infer())
}
