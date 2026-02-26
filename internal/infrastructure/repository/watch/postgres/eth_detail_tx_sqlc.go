package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/hiromaily/go-crypto-wallet/internal/application/ports/persistence"
	repowatch "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/watch"
	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/postgres/sqlcgen"
)

// ETHDetailTXInputRepositorySqlc is repository for eth_detail_tx table using sqlc
type ETHDetailTXInputRepositorySqlc struct {
	queries      *sqlcgen.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewETHDetailTXInputRepositorySqlc returns ETHDetailTXInputRepositorySqlc object
func NewETHDetailTXInputRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *ETHDetailTXInputRepositorySqlc {
	return &ETHDetailTXInputRepositorySqlc{
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

	if sqlcTx.UnsignedUpdatedAt.Valid {
		tx.UnsignedUpdatedAt = &sqlcTx.UnsignedUpdatedAt.Time
	}
	if sqlcTx.SentUpdatedAt.Valid {
		tx.SentUpdatedAt = &sqlcTx.SentUpdatedAt.Time
	}

	return tx, nil
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

	result := make([]*domainETH.ETHDetailTx, len(ethTxs))
	for i := range ethTxs {
		domainTx, err := convertToETHDetailTx(&ethTxs[i])
		if err != nil {
			return nil, fmt.Errorf("failed to convert eth detail tx: %w", err)
		}
		result[i] = domainTx
	}

	return result, nil
}

// GetSentHashTx returns list of sent_hash_tx by txType
func (r *ETHDetailTXInputRepositorySqlc) GetSentHashTx(txType domainTx.TxType) ([]string, error) {
	ctx := context.Background()

	hashes, err := r.queries.GetETHDetailTXSentHashList(ctx, sqlcgen.GetETHDetailTXSentHashListParams{
		Coin:          r.coinTypeCode.String(),
		CurrentTxType: int16(txType.Int8()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetETHDetailTXSentHashList(): %w", err)
	}

	return hashes, nil
}

// Insert inserts one record
func (r *ETHDetailTXInputRepositorySqlc) Insert(txItem *domainETH.ETHDetailTx) error {
	ctx := context.Background()

	var unsignedUpdatedAt, sentUpdatedAt sql.NullTime
	if txItem.UnsignedUpdatedAt != nil {
		unsignedUpdatedAt = sql.NullTime{Time: *txItem.UnsignedUpdatedAt, Valid: true}
	}
	if txItem.SentUpdatedAt != nil {
		sentUpdatedAt = sql.NullTime{Time: *txItem.SentUpdatedAt, Valid: true}
	}

	_, err := r.queries.InsertETHDetailTX(ctx, sqlcgen.InsertETHDetailTXParams{
		TxID:              txItem.TxID,
		Uuid:              txItem.UUID,
		CurrentTxType:     int16(txItem.CurrentTxType.Int8()),
		SenderAccount:     txItem.SenderAccount,
		SenderAddress:     txItem.SenderAddress,
		ReceiverAccount:   txItem.ReceiverAccount,
		ReceiverAddress:   txItem.ReceiverAddress,
		Amount:            int64(txItem.Amount),
		Fee:               int64(txItem.Fee),
		GasLimit:          int32(txItem.GasLimit),
		Nonce:             int64(txItem.Nonce),
		UnsignedHexTx:     txItem.UnsignedHexTx,
		SignedHexTx:       txItem.SignedHexTx,
		SentHashTx:        txItem.SentHashTx,
		UnsignedUpdatedAt: unsignedUpdatedAt,
		SentUpdatedAt:     sentUpdatedAt,
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
		CurrentTxType: int16(txType.Int8()),
		SignedHexTx:   signedHex,
		SentHashTx:    sentHashTx,
		SentUpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
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
		CurrentTxType: int16(txType.Int8()),
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
		CurrentTxType: int16(txType.Int8()),
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
	tx persistence.Transaction,
) (repowatch.ETHDetailTXRepositorier, error) {
	sqlTx := database.UnwrapSQLTx(tx)
	if sqlTx == nil {
		return nil, database.ErrUnsupportedTransaction
	}
	return &ETHDetailTXInputRepositorySqlc{
		queries:      r.queries.WithTx(sqlTx),
		coinTypeCode: r.coinTypeCode,
	}, nil
}
