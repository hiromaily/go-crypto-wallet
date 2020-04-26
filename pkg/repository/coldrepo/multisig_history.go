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

type MultisigHistoryRepository interface {
	GetAllNoMultisig(accountType account.AccountType) ([]*models.MultisigHistory, error)
	GetAllNotExported(accountType account.AccountType) ([]*models.MultisigHistory, error)
	InsertBulk(items []*models.MultisigHistory) error
	UpdateMultisigAddr(accountType account.AccountType, multiSigAddr, redeemScript, authAddr1, fullPublicKey string) (int64, error)
	UpdateIsExported(accountType account.AccountType, ids []int64) (int64, error)
}

type multisigHistoryRepository struct {
	dbConn       *sql.DB
	tableName    string
	coinTypeCode coin.CoinTypeCode
	logger       *zap.Logger
}

// NewAccountKeyRepository returns AccountKeyRepository interface
func NewMultisigHistoryRepository(dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger *zap.Logger) MultisigHistoryRepository {
	return &multisigHistoryRepository{
		dbConn:       dbConn,
		tableName:    "multisig_history",
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetAllNoMultisig returns all MultisigHistory by blank mulstisig_address
// - replaced from GetAddedPubkeyHistoryTableByNoWalletMultisigAddress
func (r *multisigHistoryRepository) GetAllNoMultisig(accountType account.AccountType) ([]*models.MultisigHistory, error) {
	//sql := "SELECT * FROM %s WHERE wallet_multisig_address = '';"
	ctx := context.Background()

	items, err := models.MultisigHistories(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("account=?", accountType.String()),
		qm.And("wallet_multisig_address=?", ""),
	).All(ctx, r.dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.MultisigHistories().All()")
	}
	return items, nil
}

// GetAllNotExported returns all MultisigHistory by not exported
// - replaced from GetAddedPubkeyHistoryTableByNotExported
func (r *multisigHistoryRepository) GetAllNotExported(accountType account.AccountType) ([]*models.MultisigHistory, error) {
	//sql := "SELECT * FROM %s WHERE wallet_multisig_address != '' AND is_exported=false;"
	ctx := context.Background()

	items, err := models.MultisigHistories(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("account=?", accountType.String()),
		qm.And("wallet_multisig_address!=?", ""),
		qm.And("is_exported=?", false),
	).All(ctx, r.dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.MultisigHistories().All()")
	}
	return items, nil
}

// Insert inserts multiple records
// - replaced from InsertAddedPubkeyHistoryTable()
func (r *multisigHistoryRepository) InsertBulk(items []*models.MultisigHistory) error {
	ctx := context.Background()
	return models.MultisigHistorySlice(items).InsertAll(ctx, r.dbConn, boil.Infer())
}

// UpdateMultisigAddr updates multisig_address
// - replaced from UpdateMultisigAddrOnAddedPubkeyHistoryTable
func (r *multisigHistoryRepository) UpdateMultisigAddr(accountType account.AccountType, multiSigAddr, redeemScript, authAddr1, fullPublicKey string) (int64, error) {
	//sql := `UPDATE %s SET wallet_multisig_address=?, redeem_script=?, auth_address1=? WHERE full_public_key=?`
	ctx := context.Background()

	// Set updating columns
	updCols := map[string]interface{}{
		models.MultisigHistoryColumns.WalletMultisigAddress: multiSigAddr,
		models.MultisigHistoryColumns.RedeemScript:          redeemScript,
		models.MultisigHistoryColumns.AuthAddress1:          authAddr1,
	}
	return models.AccountKeys(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("account=?", accountType.String()),
		qm.And("full_public_key=?", fullPublicKey),
	).UpdateAll(ctx, r.dbConn, updCols)
}

// UpdateIsExported updates is_exported
// - replaced from UpdateIsExportedOnAddedPubkeyHistoryTable
func (r *multisigHistoryRepository) UpdateIsExported(accountType account.AccountType, ids []int64) (int64, error) {
	//sql = "UPDATE %s SET is_exported=true WHERE id IN (?);"
	ctx := context.Background()

	// Set updating columns
	updCols := map[string]interface{}{
		models.MultisigHistoryColumns.IsExported: true,
	}

	targetIDs := make([]interface{}, len(ids))
	for i, v := range ids {
		targetIDs[i] = v
	}

	return models.AccountKeys(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.AndIn("id IN ?", targetIDs...),
	).UpdateAll(ctx, r.dbConn, updCols)
}
