package watchrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/hiromaily/go-crypto-wallet/pkg/db/rdb/sqlcgen"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
	"github.com/pkg/errors"
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
		return nil, errors.Wrap(err, "failed to call GetEthDetailTxByID()")
	}

	return convertSqlcEthDetailTxToModel(&ethTx), nil
}

// GetAllByTxID returns all records searched by tx_id
func (r *EthDetailTxInputRepositorySqlc) GetAllByTxID(id int64) ([]*models.EthDetailTX, error) {
	ctx := context.Background()

	ethTxs, err := r.queries.GetEthDetailTxsByTxID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call GetEthDetailTxsByTxID()")
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
		return nil, errors.Wrap(err, "failed to call GetEthDetailTxSentHashList()")
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
		Amount:            uint64(txItem.Amount),
		Fee:               uint64(txItem.Fee),
		GasLimit:          txItem.GasLimit,
		Nonce:             uint64(txItem.Nonce),
		UnsignedHexTx:     txItem.UnsignedHexTX,
		SignedHexTx:       txItem.SignedHexTX,
		SentHashTx:        txItem.SentHashTX,
		UnsignedUpdatedAt: convertNullTimeToSqlNullTime(txItem.UnsignedUpdatedAt),
		SentUpdatedAt:     convertNullTimeToSqlNullTime(txItem.SentUpdatedAt),
	})
	if err != nil {
		return errors.Wrap(err, "failed to call InsertEthDetailTx()")
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
		return 0, errors.Wrap(err, "failed to call UpdateEthDetailTxAfterSent()")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get RowsAffected()")
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
		return 0, errors.Wrap(err, "failed to call UpdateEthDetailTxType()")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get RowsAffected()")
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
		return 0, errors.Wrap(err, "failed to call UpdateEthDetailTxTypeBySentHash()")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get RowsAffected()")
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
		UnsignedUpdatedAt: convertSqlNullTimeToNullTime(ethTx.UnsignedUpdatedAt),
		SentUpdatedAt:     convertSqlNullTimeToNullTime(ethTx.SentUpdatedAt),
	}
}
