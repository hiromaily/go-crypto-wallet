package coldrepo

import (
	"database/sql"

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
