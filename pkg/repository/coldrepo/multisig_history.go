package coldrepo

import (
	"database/sql"
	"github.com/hiromaily/go-bitcoin/pkg/address"
	models "github.com/hiromaily/go-bitcoin/pkg/models/rdb"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/coin"
)

type MultisigHistoryRepository interface {
}

type multisigHistoryRepository struct {
	dbConn       *sql.DB
	tableName    string
	coinTypeCode coin.CoinTypeCode
	logger       *zap.Logger
}

// NewAccountKeyRepository returns AccountKeyRepository interface
func NewMultisigHistoryRepository(dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger *zap.Logger) MultisigHistoryRepository {
	return &accountKeyRepository{
		dbConn:       dbConn,
		tableName:    "multisig_history",
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetAllAddrStatus returns all AccountKey by addr_status
// - replaced from GetAddedPubkeyHistoryTableByNoWalletMultisigAddress
func (r *accountKeyRepository) GetAllNoMultisig() ([]*models.MultisigHistory, error) {
	//sql := "SELECT * FROM %s WHERE addr_status=?;"
	ctx := context.Background()

	items, err := models.AccountKeys(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("addr_status=?", addrStatus.Int8()),
	).All(ctx, r.dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.AccountKeys().All()")
	}

	return items, nil
}
