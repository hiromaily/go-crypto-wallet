package coldrepo

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// XRPAccountKeyRepositorier is XRPAccountKeyRepository interface
type XRPAccountKeyRepositorier interface {
	GetAllAddrStatus(accountType account.AccountType, addrStatus address.AddrStatus) ([]*models.XRPAccountKey, error)
	InsertBulk(items []*models.XRPAccountKey) error
	UpdateAddrStatus(accountType account.AccountType, addrStatus address.AddrStatus, strWIFs []string) (int64, error)
}

// XRPAccountKeyRepository is repository for xrp_account_key table
type XRPAccountKeyRepository struct {
	dbConn       *sql.DB
	tableName    string
	coinTypeCode coin.CoinTypeCode
	logger       *zap.Logger
}

// NewXRPAccountKeyRepository returns XRPAccountKeyRepository object
func NewXRPAccountKeyRepository(dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger *zap.Logger) *XRPAccountKeyRepository {
	return &XRPAccountKeyRepository{
		dbConn:       dbConn,
		tableName:    "xrp_account_key",
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetAllAddrStatus returns all XRPAccountKey by addr_status
func (r *XRPAccountKeyRepository) GetAllAddrStatus(accountType account.AccountType, addrStatus address.AddrStatus) ([]*models.XRPAccountKey, error) {
	ctx := context.Background()

	items, err := models.XRPAccountKeys(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("account=?", accountType.String()),
		qm.And("addr_status=?", addrStatus.Int8()),
	).All(ctx, r.dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.AccountKeys().All()")
	}

	return items, nil
}

// InsertBulk inserts multiple records
func (r *XRPAccountKeyRepository) InsertBulk(items []*models.XRPAccountKey) error {
	ctx := context.Background()
	return models.XRPAccountKeySlice(items).InsertAll(ctx, r.dbConn, boil.Infer())
}

// UpdateAddrStatus updates addr_status
func (r *XRPAccountKeyRepository) UpdateAddrStatus(accountType account.AccountType, addrStatus address.AddrStatus, accountIDs []string) (int64, error) {
	//sql := `UPDATE %s SET addr_status=? WHERE wallet_import_format=?`
	ctx := context.Background()

	// Set updating columns
	updCols := map[string]interface{}{
		models.XRPAccountKeyColumns.AddrStatus: addrStatus.Int8(),
	}

	targetIDs := make([]interface{}, len(accountIDs))
	for i, v := range accountIDs {
		targetIDs[i] = v
	}

	return models.XRPAccountKeys(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("account=?", accountType.String()),
		qm.AndIn("account_id IN ?", targetIDs...),
	).UpdateAll(ctx, r.dbConn, updCols)
}
