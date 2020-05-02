package coldrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/address"
	models "github.com/hiromaily/go-bitcoin/pkg/models/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/coin"
)

// AccountKeyRepositorier is AccountKeyRepository interface
type AccountKeyRepositorier interface {
	GetMaxIndex(accountType account.AccountType) (int64, error)
	GetOneMaxID(accountType account.AccountType) (*models.AccountKey, error)
	GetAllAddrStatus(accountType account.AccountType, addrStatus address.AddrStatus) ([]*models.AccountKey, error)
	GetAllMultiAddr(accountType account.AccountType, addrs []string) ([]*models.AccountKey, error)
	InsertBulk(items []*models.AccountKey) error
	UpdateAddrStatus(accountType account.AccountType, addrStatus address.AddrStatus, strWIFs []string) (int64, error)
	UpdateMultisigAddr(accountType account.AccountType, item *models.AccountKey) (int64, error)
	UpdateMultisigAddrs(accountType account.AccountType, items []*models.AccountKey) (int64, error)
}

// AccountKeyRepository is repository for account_key table
type AccountKeyRepository struct {
	dbConn       *sql.DB
	tableName    string
	coinTypeCode coin.CoinTypeCode
	logger       *zap.Logger
}

// NewAccountKeyRepository returns AccountKeyRepository object
func NewAccountKeyRepository(dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger *zap.Logger) *AccountKeyRepository {
	return &AccountKeyRepository{
		dbConn:       dbConn,
		tableName:    "account_key",
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetMaxIndex returns max id
func (r *AccountKeyRepository) GetMaxIndex(accountType account.AccountType) (int64, error) {
	//sql := "SELECT MAX(idx) from %s;"
	ctx := context.Background()

	type Response struct {
		MaxCount int64 `boil:"max_count"`
	}
	var maxCount Response
	err := models.AccountKeys(
		qm.Select("MAX(idx) as max_count"),
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("account=?", accountType.String()),
	).Bind(ctx, r.dbConn, &maxCount)
	if err != nil {
		return 0, errors.Wrap(err, "failed to call models.AccountKeys().Bind()")
	}
	return maxCount.MaxCount, nil
}

// GetOneMaxID returns one records by max id
func (r *AccountKeyRepository) GetOneMaxID(accountType account.AccountType) (*models.AccountKey, error) {
	//sql := "SELECT * FROM %s ORDER BY ID DESC LIMIT 1;"
	ctx := context.Background()

	item, err := models.AccountKeys(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("account=?", accountType.String()),
		qm.OrderBy("id DESC"),
	).One(ctx, r.dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.AccountKeys().One()")
	}
	return item, nil
}

// GetAllAddrStatus returns all AccountKey by addr_status
func (r *AccountKeyRepository) GetAllAddrStatus(accountType account.AccountType, addrStatus address.AddrStatus) ([]*models.AccountKey, error) {
	//sql := "SELECT * FROM %s WHERE addr_status=?;"
	ctx := context.Background()

	items, err := models.AccountKeys(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("account=?", accountType.String()),
		qm.And("addr_status=?", addrStatus.Int8()),
	).All(ctx, r.dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.AccountKeys().All()")
	}

	return items, nil
}

// GetAllMultiAddr returns all AccountKey by multisig_address
func (r *AccountKeyRepository) GetAllMultiAddr(accountType account.AccountType, addrs []string) ([]*models.AccountKey, error) {
	//sql := "SELECT * FROM %s WHERE wallet_multisig_address IN (?);"
	ctx := context.Background()

	targetAddrs := make([]interface{}, len(addrs))
	for i, v := range addrs {
		targetAddrs[i] = v
	}

	items, err := models.AccountKeys(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("account=?", accountType.String()),
		qm.AndIn("multisig_address IN ?", targetAddrs...),
	).All(ctx, r.dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.AccountKeys().All()")
	}

	return items, nil
}

// InsertBulk inserts multiple records
func (r *AccountKeyRepository) InsertBulk(items []*models.AccountKey) error {
	ctx := context.Background()
	return models.AccountKeySlice(items).InsertAll(ctx, r.dbConn, boil.Infer())
}

// UpdateAddrStatus updates addr_status
func (r *AccountKeyRepository) UpdateAddrStatus(accountType account.AccountType, addrStatus address.AddrStatus, strWIFs []string) (int64, error) {
	//sql := `UPDATE %s SET addr_status=? WHERE wallet_import_format=?`
	ctx := context.Background()

	// Set updating columns
	updCols := map[string]interface{}{
		models.AccountKeyColumns.AddrStatus: addrStatus.Int8(),
	}

	targetWIFs := make([]interface{}, len(strWIFs))
	for i, v := range strWIFs {
		targetWIFs[i] = v
	}

	return models.AccountKeys(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("account=?", accountType.String()),
		qm.AndIn("wallet_import_format IN ?", targetWIFs...),
	).UpdateAll(ctx, r.dbConn, updCols)
}

// UpdateMultisigAddr updates all multisig_address
// TODO: how to update multiple records
func (r *AccountKeyRepository) UpdateMultisigAddr(accountType account.AccountType, item *models.AccountKey) (int64, error) {
	ctx := context.Background()

	// Set updating columns
	updCols := map[string]interface{}{
		models.AccountKeyColumns.MultisigAddress: item.MultisigAddress,
		models.AccountKeyColumns.RedeemScript:    item.RedeemScript,
		models.AccountKeyColumns.AddrStatus:      item.AddrStatus,
		models.AccountKeyColumns.UpdatedAt:       null.TimeFrom(time.Now()),
	}
	_, err := models.AccountKeys(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("account=?", accountType.String()),
		qm.And("full_public_key=?", item.FullPublicKey),
	).UpdateAll(ctx, r.dbConn, updCols)
	if err != nil {
		return 0, err
	}

	return 0, nil
}

// UpdateMultisigAddrs updates all multisig_address
// TODO: how to update multiple records
// TODO: maybe this func can be deleted
func (r *AccountKeyRepository) UpdateMultisigAddrs(accountType account.AccountType, items []*models.AccountKey) (int64, error) {
	//	sql := `
	//UPDATE %s SET wallet_multisig_address=:wallet_multisig_address, redeem_script=:redeem_script, addr_status=:addr_status, updated_at=:updated_at
	//WHERE full_public_key=:full_public_key`
	ctx := context.Background()

	// transaction
	dtx, err := r.dbConn.Begin()
	if err != nil {
		return 0, errors.Wrap(err, "failed to call db.Begin()")
	}
	defer func() {
		if err != nil {
			dtx.Rollback()
		} else {
			dtx.Commit()
		}
	}()

	for _, item := range items {
		// Set updating columns
		updCols := map[string]interface{}{
			models.AccountKeyColumns.MultisigAddress: item.MultisigAddress,
			models.AccountKeyColumns.RedeemScript:    item.RedeemScript,
			models.AccountKeyColumns.AddrStatus:      item.AddrStatus,
			models.AccountKeyColumns.UpdatedAt:       null.TimeFrom(time.Now()),
		}
		_, err := models.AccountKeys(
			qm.Where("coin=?", r.coinTypeCode.String()),
			qm.And("account=?", accountType.String()),
			qm.And("full_public_key=?", item.FullPublicKey),
		).UpdateAll(ctx, r.dbConn, updCols)
		if err != nil {
			return 0, err
		}
	}
	return 0, nil
}

// GetRedeedScriptByAddress returns RedeemScript from given multisig address
func GetRedeedScriptByAddress(accountKeys []*models.AccountKey, addr string) string {
	for _, val := range accountKeys {
		if val.MultisigAddress == addr {
			return val.RedeemScript
		}
	}
	return ""
}
