package watchrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-crypto-wallet/pkg/db/rdb/sqlcgen"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// XrpDetailTxInputRepositorySqlc is repository for xrp_detail_tx table using sqlc
type XrpDetailTxInputRepositorySqlc struct {
	queries      *sqlcgen.Queries
	coinTypeCode coin.CoinTypeCode
	logger       logger.Logger
}

// NewXrpDetailTxInputRepositorySqlc returns XrpDetailTxInputRepositorySqlc object
func NewXrpDetailTxInputRepositorySqlc(
	dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger logger.Logger,
) *XrpDetailTxInputRepositorySqlc {
	return &XrpDetailTxInputRepositorySqlc{
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetOne get one record by ID
func (r *XrpDetailTxInputRepositorySqlc) GetOne(id int64) (*models.XRPDetailTX, error) {
	ctx := context.Background()

	xrpTx, err := r.queries.GetXrpDetailTxByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call GetXrpDetailTxByID()")
	}

	return convertSqlcXrpDetailTxToModel(&xrpTx), nil
}

// GetAllByTxID returns all records searched by tx_id
func (r *XrpDetailTxInputRepositorySqlc) GetAllByTxID(id int64) ([]*models.XRPDetailTX, error) {
	ctx := context.Background()

	xrpTxs, err := r.queries.GetXrpDetailTxsByTxID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call GetXrpDetailTxsByTxID()")
	}

	result := make([]*models.XRPDetailTX, len(xrpTxs))
	for i, xrpTx := range xrpTxs {
		result[i] = convertSqlcXrpDetailTxToModel(&xrpTx)
	}

	return result, nil
}

// GetSentHashTx returns list of tx_blob by txType
func (r *XrpDetailTxInputRepositorySqlc) GetSentHashTx(txType tx.TxType) ([]string, error) {
	ctx := context.Background()

	blobs, err := r.queries.GetXrpDetailTxBlobList(ctx, sqlcgen.GetXrpDetailTxBlobListParams{
		Coin:          sqlcgen.TxCoin(r.coinTypeCode.String()),
		CurrentTxType: txType.Int8(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to call GetXrpDetailTxBlobList()")
	}

	return blobs, nil
}

// Insert inserts one record
func (r *XrpDetailTxInputRepositorySqlc) Insert(txItem *models.XRPDetailTX) error {
	ctx := context.Background()

	_, err := r.queries.InsertXrpDetailTx(ctx, sqlcgen.InsertXrpDetailTxParams{
		TxID:                  txItem.TXID,
		Uuid:                  txItem.UUID,
		CurrentTxType:         txItem.CurrentTXType,
		SenderAccount:         txItem.SenderAccount,
		SenderAddress:         txItem.SenderAddress,
		ReceiverAccount:       txItem.ReceiverAccount,
		ReceiverAddress:       txItem.ReceiverAddress,
		Amount:                txItem.Amount,
		XrpTxType:             txItem.XRPTXType,
		Fee:                   txItem.Fee,
		Flags:                 txItem.Flags,
		LastLedgerSequence:    txItem.LastLedgerSequence,
		Sequence:              txItem.Sequence,
		SigningPubkey:         txItem.SigningPubkey,
		TxnSignature:          txItem.TXNSignature,
		Hash:                  txItem.Hash,
		EarliestLedgerVersion: txItem.EarliestLedgerVersion,
		SignedTxID:            txItem.SignedTXID,
		TxBlob:                txItem.TXBlob,
		SentUpdatedAt:         convertNullTimeToSQLNullTime(txItem.SentUpdatedAt),
	})
	if err != nil {
		return errors.Wrap(err, "failed to call InsertXrpDetailTx()")
	}

	return nil
}

// InsertBulk inserts multiple records
func (r *XrpDetailTxInputRepositorySqlc) InsertBulk(txItems []*models.XRPDetailTX) error {
	for _, item := range txItems {
		if err := r.Insert(item); err != nil {
			return err
		}
	}
	return nil
}

// UpdateAfterTxSent updates when tx sent
func (r *XrpDetailTxInputRepositorySqlc) UpdateAfterTxSent(
	uuid string,
	txType tx.TxType,
	signedTxID,
	txBlob string,
	earlistLedgerVersion uint64,
) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateXrpDetailTxAfterSent(ctx, sqlcgen.UpdateXrpDetailTxAfterSentParams{
		CurrentTxType:         txType.Int8(),
		SignedTxID:            signedTxID,
		TxBlob:                txBlob,
		EarliestLedgerVersion: earlistLedgerVersion,
		SentUpdatedAt:         sql.NullTime{Time: time.Now(), Valid: true},
		Uuid:                  uuid,
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to call UpdateXrpDetailTxAfterSent()")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get RowsAffected()")
	}

	return rowsAffected, nil
}

// UpdateTxType updates txType
func (r *XrpDetailTxInputRepositorySqlc) UpdateTxType(id int64, txType tx.TxType) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateXrpDetailTxType(ctx, sqlcgen.UpdateXrpDetailTxTypeParams{
		CurrentTxType: txType.Int8(),
		ID:            id,
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to call UpdateXrpDetailTxType()")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get RowsAffected()")
	}

	return rowsAffected, nil
}

// UpdateTxTypeBySentHashTx updates txType by tx_blob
func (r *XrpDetailTxInputRepositorySqlc) UpdateTxTypeBySentHashTx(
	txType tx.TxType, sentHashTx string,
) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateXrpDetailTxTypeBySentHash(ctx, sqlcgen.UpdateXrpDetailTxTypeBySentHashParams{
		CurrentTxType: txType.Int8(),
		TxBlob:        sentHashTx, // sentHashTx is actually tx_blob for XRP
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to call UpdateXrpDetailTxTypeBySentHash()")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get RowsAffected()")
	}

	return rowsAffected, nil
}

// Helper functions

func convertSqlcXrpDetailTxToModel(xrpTx *sqlcgen.XrpDetailTx) *models.XRPDetailTX {
	return &models.XRPDetailTX{
		ID:                    xrpTx.ID,
		TXID:                  xrpTx.TxID,
		UUID:                  xrpTx.Uuid,
		CurrentTXType:         xrpTx.CurrentTxType,
		SenderAccount:         xrpTx.SenderAccount,
		SenderAddress:         xrpTx.SenderAddress,
		ReceiverAccount:       xrpTx.ReceiverAccount,
		ReceiverAddress:       xrpTx.ReceiverAddress,
		Amount:                xrpTx.Amount,
		XRPTXType:             xrpTx.XrpTxType,
		Fee:                   xrpTx.Fee,
		Flags:                 xrpTx.Flags,
		LastLedgerSequence:    xrpTx.LastLedgerSequence,
		Sequence:              xrpTx.Sequence,
		SigningPubkey:         xrpTx.SigningPubkey,
		TXNSignature:          xrpTx.TxnSignature,
		Hash:                  xrpTx.Hash,
		EarliestLedgerVersion: xrpTx.EarliestLedgerVersion,
		SignedTXID:            xrpTx.SignedTxID,
		TXBlob:                xrpTx.TxBlob,
		SentUpdatedAt:         convertSQLNullTimeToNullTime(xrpTx.SentUpdatedAt),
	}
}
