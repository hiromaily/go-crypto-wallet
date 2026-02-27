package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	repowatch "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/watch"
	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/sqlite/sqlcgen"
	dbtx "github.com/hiromaily/go-crypto-wallet/pkg/db/tx"
)

// ETHDetailTXInputRepositorySqlc is repository for eth_detail_tx table using sqlc (SQLite)
type ETHDetailTXInputRepositorySqlc struct {
	db           *sql.DB
	queries      *sqlcgen.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewETHDetailTXInputRepositorySqlc returns ETHDetailTXInputRepositorySqlc object
func NewETHDetailTXInputRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *ETHDetailTXInputRepositorySqlc {
	return &ETHDetailTXInputRepositorySqlc{
		db:           dbConn,
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
	}
}

// convertToETHDetailTx converts sqlcgen.EthDetailTx to domain.ETHDetailTx entity
func convertToETHDetailTx(sqlcTx *sqlcgen.EthDetailTx) (*domainETH.ETHDetailTx, error) {
	currentTxType, err := domainTx.TxTypeFromInt8(int8(sqlcTx.CurrentTxType))
	if err != nil {
		return nil, fmt.Errorf("invalid tx type in database: %w", err)
	}

	tx := &domainETH.ETHDetailTx{
		ID:              sqlcTx.ID,
		TxID:            sqlcTx.TxID,
		UUID:            sqlcTx.Uuid,
		CurrentTxType:   currentTxType,
		SenderAccount:   sqlcTx.SenderAccount,
		SenderAddress:   sqlcTx.SenderAddress,
		ReceiverAccount: sqlcTx.ReceiverAccount,
		ReceiverAddress: sqlcTx.ReceiverAddress,
		Amount:          uint64(sqlcTx.Amount),
		Fee:             uint64(sqlcTx.Fee),
		GasLimit:        uint32(sqlcTx.GasLimit),
		Nonce:           uint64(sqlcTx.Nonce),
		UnsignedHexTx:   sqlcTx.UnsignedHexTx,
		SignedHexTx:     sqlcTx.SignedHexTx,
		SentHashTx:      sqlcTx.SentHashTx,
	}

	// Parse TEXT timestamps
	if sqlcTx.UnsignedUpdatedAt.Valid {
		t, err := time.Parse("2006-01-02 15:04:05", sqlcTx.UnsignedUpdatedAt.String)
		if err == nil {
			tx.UnsignedUpdatedAt = &t
		}
	}
	if sqlcTx.SentUpdatedAt.Valid {
		t, err := time.Parse("2006-01-02 15:04:05", sqlcTx.SentUpdatedAt.String)
		if err == nil {
			tx.SentUpdatedAt = &t
		}
	}

	return tx, nil
}

// convertFromETHDetailTx converts domain.ETHDetailTx entity to sqlcgen.EthDetailTx
func convertFromETHDetailTx(tx *domainETH.ETHDetailTx) *sqlcgen.EthDetailTx {
	sqlcTx := &sqlcgen.EthDetailTx{
		ID:              tx.ID,
		TxID:            tx.TxID,
		Uuid:            tx.UUID,
		CurrentTxType:   int64(tx.CurrentTxType.Int8()),
		SenderAccount:   tx.SenderAccount,
		SenderAddress:   tx.SenderAddress,
		ReceiverAccount: tx.ReceiverAccount,
		ReceiverAddress: tx.ReceiverAddress,
		Amount:          int64(tx.Amount),
		Fee:             int64(tx.Fee),
		GasLimit:        int64(tx.GasLimit),
		Nonce:           int64(tx.Nonce),
		UnsignedHexTx:   tx.UnsignedHexTx,
		SignedHexTx:     tx.SignedHexTx,
		SentHashTx:      tx.SentHashTx,
	}

	// Convert time.Time to TEXT format
	if tx.UnsignedUpdatedAt != nil {
		sqlcTx.UnsignedUpdatedAt = sql.NullString{
			String: tx.UnsignedUpdatedAt.Format("2006-01-02 15:04:05"),
			Valid:  true,
		}
	}
	if tx.SentUpdatedAt != nil {
		sqlcTx.SentUpdatedAt = sql.NullString{
			String: tx.SentUpdatedAt.Format("2006-01-02 15:04:05"),
			Valid:  true,
		}
	}

	return sqlcTx
}

// GetOne get one record by ID
func (r *ETHDetailTXInputRepositorySqlc) GetOne(id int64) (*domainETH.ETHDetailTx, error) {
	ctx := context.Background()

	ethTx, err := r.queries.GetETHDetailTXByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetETHDetailTXByID(): %w", err)
	}

	return convertToETHDetailTx(&ethTx)
}

// GetAllByTxID returns all records searched by tx_id
func (r *ETHDetailTXInputRepositorySqlc) GetAllByTxID(id int64) ([]*domainETH.ETHDetailTx, error) {
	ctx := context.Background()

	ethTxs, err := r.queries.GetETHDetailTXsByTxID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetETHDetailTXsByTxID(): %w", err)
	}

	result := make([]*domainETH.ETHDetailTx, 0, len(ethTxs))
	for i := range ethTxs {
		domainTx, err := convertToETHDetailTx(&ethTxs[i])
		if err != nil {
			return nil, fmt.Errorf("failed to convert eth detail tx: %w", err)
		}
		result = append(result, domainTx)
	}

	return result, nil
}

// GetSentHashTx returns list of sent_hash_tx by txType
func (r *ETHDetailTXInputRepositorySqlc) GetSentHashTx(txType domainTx.TxType) ([]string, error) {
	ctx := context.Background()

	hashes, err := r.queries.GetETHDetailTXSentHashList(ctx, sqlcgen.GetETHDetailTXSentHashListParams{
		Coin:          r.coinTypeCode.String(),
		CurrentTxType: int64(txType.Int8()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetETHDetailTXSentHashList(): %w", err)
	}

	return hashes, nil
}

// Insert inserts one record
func (r *ETHDetailTXInputRepositorySqlc) Insert(txItem *domainETH.ETHDetailTx) error {
	ctx := context.Background()

	sqlcTx := convertFromETHDetailTx(txItem)
	_, err := r.queries.InsertETHDetailTX(ctx, sqlcgen.InsertETHDetailTXParams{
		TxID:              sqlcTx.TxID,
		Uuid:              sqlcTx.Uuid,
		CurrentTxType:     sqlcTx.CurrentTxType,
		SenderAccount:     sqlcTx.SenderAccount,
		SenderAddress:     sqlcTx.SenderAddress,
		ReceiverAccount:   sqlcTx.ReceiverAccount,
		ReceiverAddress:   sqlcTx.ReceiverAddress,
		Amount:            sqlcTx.Amount,
		Fee:               sqlcTx.Fee,
		GasLimit:          sqlcTx.GasLimit,
		Nonce:             sqlcTx.Nonce,
		UnsignedHexTx:     sqlcTx.UnsignedHexTx,
		SignedHexTx:       sqlcTx.SignedHexTx,
		SentHashTx:        sqlcTx.SentHashTx,
		UnsignedUpdatedAt: sqlcTx.UnsignedUpdatedAt,
		SentUpdatedAt:     sqlcTx.SentUpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to call InsertETHDetailTX(): %w", err)
	}

	return nil
}

// InsertBulk inserts multiple records
func (r *ETHDetailTXInputRepositorySqlc) InsertBulk(txItems []*domainETH.ETHDetailTx) error {
	for _, item := range txItems {
		if err := r.Insert(item); err != nil {
			return err
		}
	}
	return nil
}

// UpdateAfterTxSent updates when tx sent
func (r *ETHDetailTXInputRepositorySqlc) UpdateAfterTxSent(
	uuid string,
	txType domainTx.TxType,
	signedHex,
	sentHashTx string,
) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateETHDetailTXAfterSent(ctx, sqlcgen.UpdateETHDetailTXAfterSentParams{
		CurrentTxType: int64(txType.Int8()),
		SignedHexTx:   signedHex,
		SentHashTx:    sentHashTx,
		SentUpdatedAt: sql.NullString{String: time.Now().Format("2006-01-02 15:04:05"), Valid: true},
		Uuid:          uuid,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call UpdateETHDetailTXAfterSent(): %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
	}

	return rowsAffected, nil
}

// UpdateTxType updates txType
func (r *ETHDetailTXInputRepositorySqlc) UpdateTxType(id int64, txType domainTx.TxType) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateETHDetailTXType(ctx, sqlcgen.UpdateETHDetailTXTypeParams{
		CurrentTxType: int64(txType.Int8()),
		ID:            id,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call UpdateETHDetailTXType(): %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
	}

	return rowsAffected, nil
}

// UpdateTxTypeBySentHashTx updates txType
func (r *ETHDetailTXInputRepositorySqlc) UpdateTxTypeBySentHashTx(
	txType domainTx.TxType, sentHashTx string,
) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateETHDetailTXTypeBySentHash(ctx, sqlcgen.UpdateETHDetailTXTypeBySentHashParams{
		CurrentTxType: int64(txType.Int8()),
		SentHashTx:    sentHashTx,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call UpdateETHDetailTXTypeBySentHash(): %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
	}

	return rowsAffected, nil
}

// WithTransaction returns a new repository instance that uses the provided transaction
func (r *ETHDetailTXInputRepositorySqlc) WithTransaction(
	tx dbtx.Transaction,
) (repowatch.ETHDetailTXRepositorier, error) {
	sqlTx := dbtx.UnwrapSQLTx(tx)
	if sqlTx == nil {
		return nil, dbtx.ErrUnsupportedTransaction
	}
	return &ETHDetailTXInputRepositorySqlc{
		db:           r.db,
		queries:      r.queries.WithTx(sqlTx),
		coinTypeCode: r.coinTypeCode,
	}, nil
}
