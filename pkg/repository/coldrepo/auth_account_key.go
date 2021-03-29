package coldrepo

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// AuthAccountKeyRepositorier is AuthAccountKeyRepository interface
type AuthAccountKeyRepositorier interface {
	GetOne(authType account.AuthType) (*models.AuthAccountKey, error)
	Insert(item *models.AuthAccountKey) error
	UpdateAddrStatus(addrStatus address.AddrStatus, strWIF string) (int64, error)
}

// AuthAccountKeyRepository is repository  for auth_account table
type AuthAccountKeyRepository struct {
	dbConn       *sql.DB
	tableName    string
	coinTypeCode coin.CoinTypeCode
	authType     account.AuthType
	logger       *zap.Logger
}

// NewAuthAccountKeyRepository returns AuthAccountKeyRepository object
func NewAuthAccountKeyRepository(
	dbConn *sql.DB,
	coinTypeCode coin.CoinTypeCode,
	authType account.AuthType,
	logger *zap.Logger) *AuthAccountKeyRepository {
	return &AuthAccountKeyRepository{
		dbConn:       dbConn,
		tableName:    "auth_account_key",
		coinTypeCode: coinTypeCode,
		authType:     authType,
		logger:       logger,
	}
}

// GetOne returns one records
func (r *AuthAccountKeyRepository) GetOne(authType account.AuthType) (*models.AuthAccountKey, error) {
	ctx := context.Background()

	if authType == "" {
		authType = r.authType
	}

	item, err := models.AuthAccountKeys(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("auth_account=?", authType.String()),
	).One(ctx, r.dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.AuthAccountKeys().One()")
	}
	return item, nil
}

// Insert inserts record
func (r *AuthAccountKeyRepository) Insert(item *models.AuthAccountKey) error {
	ctx := context.Background()
	if err := item.Insert(ctx, r.dbConn, boil.Infer()); err != nil {
		return errors.Wrap(err, "failed to call item.Insert()")
	}

	return nil
}

// UpdateAddrStatus updates addr_status
func (r *AuthAccountKeyRepository) UpdateAddrStatus(addrStatus address.AddrStatus, strWIF string) (int64, error) {
	ctx := context.Background()

	// Set updating columns
	updCols := map[string]interface{}{
		models.AuthAccountKeyColumns.AddrStatus: addrStatus.Int8(),
	}

	return models.AuthAccountKeys(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("auth_account=?", r.authType.String()),
		qm.And("wallet_import_format=?", strWIF),
	).UpdateAll(ctx, r.dbConn, updCols)
}
