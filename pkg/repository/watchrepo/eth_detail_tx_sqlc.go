package watchrepo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/hiromaily/go-crypto-wallet/pkg/db/rdb/sqlcgen"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// EthDetailTxInputRepositorySqlc is repository for eth_detail_tx table using sqlc
type EthDetailTxInputRepositorySqlc struct {
	queries      *sqlcgen.Queries
	coinTypeCode coin.CoinTypeCode
	logger       logger.Logger
}

// NewEthDetailTxInputRepositorySqlc returns EthDetailTxInputRepositorySqlc object
func NewEthDetailTxInputRepositorySqlc(
	dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger logger.Logger,
) *EthDetailTxInputRepositorySqlc {
	return &EthDetailTxInputRepositorySqlc{
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetOne get one record by ID
func (r *EthDetailTxInputRepositorySqlc) GetOne(id int64) (*models.EthDetailTX, error) {
	ctx := context.Background()

	ethTx, err := r.queries.GetEthDetailTxByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetEthDetailTxByID(): %w", err)
	}

	return convertSqlcEthDetailTxToModel(&ethTx), nil
}

// GetAllByTxID returns all records searched by tx_id
func (r *EthDetailTxInputRepositorySqlc) GetAllByTxID(id int64) ([]*models.EthDetailTX, error) {
	ctx := context.Background()

	ethTxs, err := r.queries.GetEthDetailTxsByTxID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetEthDetailTxsByTxID(): %w", err)
	}

	result := make([]*models.EthDetailTX, len(ethTxs))
	for i, ethTx := range ethTxs {
		result[i] = convertSqlcEthDetailTxToModel(&ethTx)
	}

	return result, nil
}

// GetSentHashTx returns list of sent_hash_tx by txType
func (r *EthDetailTxInputRepositorySqlc) GetSentHashTx(txType tx.TxType) ([]string, error) {
	ctx := context.Background()

	hashes, err := r.queries.GetEthDetailTxSentHashList(ctx, sqlcgen.GetEthDetailTxSentHashListParams{
		Coin:          sqlcgen.TxCoin(r.coinTypeCode.String()),
		CurrentTxType: txType.Int8(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetEthDetailTxSentHashList(): %w", err)
	}

	return hashes, nil
}

// Insert inserts one record
func (r *EthDetailTxInputRepositorySqlc) Insert(txItem *models.EthDetailTX) error {
	ctx := context.Background()

	_, err := r.queries.InsertEthDetailTx(ctx, sqlcgen.InsertEthDetailTxParams{
		TxID:              txItem.TXID,
		Uuid:              txItem.UUID,
		CurrentTxType:     txItem.CurrentTXType,
		SenderAccount:     txItem.SenderAccount,
		SenderAddress:     txItem.SenderAddress,
		ReceiverAccount:   txItem.ReceiverAccount,
		ReceiverAddress:   txItem.ReceiverAddress,
		Amount:            txItem.Amount,
		Fee:               txItem.Fee,
		GasLimit:          txItem.GasLimit,
		Nonce:             txItem.Nonce,
		UnsignedHexTx:     txItem.UnsignedHexTX,
		SignedHexTx:       txItem.SignedHexTX,
		SentHashTx:        txItem.SentHashTX,
		UnsignedUpdatedAt: convertNullTimeToSQLNullTime(txItem.UnsignedUpdatedAt),
		SentUpdatedAt:     convertNullTimeToSQLNullTime(txItem.SentUpdatedAt),
	})
	if err != nil {
		return fmt.Errorf("failed to call InsertEthDetailTx(): %w", err)
	}

	return nil
}

// InsertBulk inserts multiple records
func (r *EthDetailTxInputRepositorySqlc) InsertBulk(txItems []*models.EthDetailTX) error {
	for _, item := range txItems {
		if err := r.Insert(item); err != nil {
			return err
		}
	}
	return nil
}

// UpdateAfterTxSent updates when tx sent
func (r *EthDetailTxInputRepositorySqlc) UpdateAfterTxSent(
	uuid string,
	txType tx.TxType,
	signedHex,
	sentHashTx string,
) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateEthDetailTxAfterSent(ctx, sqlcgen.UpdateEthDetailTxAfterSentParams{
		CurrentTxType: txType.Int8(),
		SignedHexTx:   signedHex,
		SentHashTx:    sentHashTx,
		SentUpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Uuid:          uuid,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call UpdateEthDetailTxAfterSent(): %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
	}

	return rowsAffected, nil
}

// UpdateTxType updates txType
func (r *EthDetailTxInputRepositorySqlc) UpdateTxType(id int64, txType tx.TxType) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateEthDetailTxType(ctx, sqlcgen.UpdateEthDetailTxTypeParams{
		CurrentTxType: txType.Int8(),
		ID:            id,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call UpdateEthDetailTxType(): %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
	}

	return rowsAffected, nil
}

// UpdateTxTypeBySentHashTx updates txType
func (r *EthDetailTxInputRepositorySqlc) UpdateTxTypeBySentHashTx(txType tx.TxType, sentHashTx string) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateEthDetailTxTypeBySentHash(ctx, sqlcgen.UpdateEthDetailTxTypeBySentHashParams{
		CurrentTxType: txType.Int8(),
		SentHashTx:    sentHashTx,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call UpdateEthDetailTxTypeBySentHash(): %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
	}

	return rowsAffected, nil
}

// Helper functions

func convertSqlcEthDetailTxToModel(ethTx *sqlcgen.EthDetailTx) *models.EthDetailTX {
	return &models.EthDetailTX{
		ID:                ethTx.ID,
		TXID:              ethTx.TxID,
		UUID:              ethTx.Uuid,
		CurrentTXType:     ethTx.CurrentTxType,
		SenderAccount:     ethTx.SenderAccount,
		SenderAddress:     ethTx.SenderAddress,
		ReceiverAccount:   ethTx.ReceiverAccount,
		ReceiverAddress:   ethTx.ReceiverAddress,
		Amount:            ethTx.Amount,
		Fee:               ethTx.Fee,
		GasLimit:          ethTx.GasLimit,
		Nonce:             ethTx.Nonce,
		UnsignedHexTX:     ethTx.UnsignedHexTx,
		SignedHexTX:       ethTx.SignedHexTx,
		SentHashTX:        ethTx.SentHashTx,
		UnsignedUpdatedAt: convertSQLNullTimeToNullTime(ethTx.UnsignedUpdatedAt),
		SentUpdatedAt:     convertSQLNullTimeToNullTime(ethTx.SentUpdatedAt),
	}
}
