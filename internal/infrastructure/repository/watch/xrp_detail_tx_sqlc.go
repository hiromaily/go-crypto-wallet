package watch

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	models "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/sqlc"
)

// XrpDetailTxInputRepositorySqlc is repository for xrp_detail_tx table using sqlc
type XrpDetailTxInputRepositorySqlc struct {
	queries      *sqlc.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewXrpDetailTxInputRepositorySqlc returns XrpDetailTxInputRepositorySqlc object
func NewXrpDetailTxInputRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *XrpDetailTxInputRepositorySqlc {
	return &XrpDetailTxInputRepositorySqlc{
		queries:      sqlc.New(dbConn),
		coinTypeCode: coinTypeCode,
	}
}

// GetOne get one record by ID
func (r *XrpDetailTxInputRepositorySqlc) GetOne(id int64) (*models.XRPDetailTX, error) {
	ctx := context.Background()

	xrpTx, err := r.queries.GetXrpDetailTxByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetXrpDetailTxByID(): %w", err)
	}

	return convertSqlcXrpDetailTxToModel(&xrpTx), nil
}

// GetAllByTxID returns all records searched by tx_id
func (r *XrpDetailTxInputRepositorySqlc) GetAllByTxID(id int64) ([]*models.XRPDetailTX, error) {
	ctx := context.Background()

	xrpTxs, err := r.queries.GetXrpDetailTxsByTxID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetXrpDetailTxsByTxID(): %w", err)
	}

	result := make([]*models.XRPDetailTX, len(xrpTxs))
	for i, xrpTx := range xrpTxs {
		result[i] = convertSqlcXrpDetailTxToModel(&xrpTx)
	}

	return result, nil
}

// GetSentHashTx returns list of tx_blob by txType
func (r *XrpDetailTxInputRepositorySqlc) GetSentHashTx(txType domainTx.TxType) ([]string, error) {
	ctx := context.Background()

	blobs, err := r.queries.GetXrpDetailTxBlobList(ctx, sqlc.GetXrpDetailTxBlobListParams{
		Coin:          sqlc.TxCoin(r.coinTypeCode.String()),
		CurrentTxType: txType.Int8(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetXrpDetailTxBlobList(): %w", err)
	}

	return blobs, nil
}

// Insert inserts one record
func (r *XrpDetailTxInputRepositorySqlc) Insert(txItem *models.XRPDetailTX) error {
	ctx := context.Background()

	_, err := r.queries.InsertXrpDetailTx(ctx, sqlc.InsertXrpDetailTxParams{
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
		return fmt.Errorf("failed to call InsertXrpDetailTx(): %w", err)
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
	txType domainTx.TxType,
	signedTxID,
	txBlob string,
	earlistLedgerVersion uint64,
) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateXrpDetailTxAfterSent(ctx, sqlc.UpdateXrpDetailTxAfterSentParams{
		CurrentTxType:         txType.Int8(),
		SignedTxID:            signedTxID,
		TxBlob:                txBlob,
		EarliestLedgerVersion: earlistLedgerVersion,
		SentUpdatedAt:         sql.NullTime{Time: time.Now(), Valid: true},
		Uuid:                  uuid,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call UpdateXrpDetailTxAfterSent(): %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
	}

	return rowsAffected, nil
}

// UpdateTxType updates txType
func (r *XrpDetailTxInputRepositorySqlc) UpdateTxType(id int64, txType domainTx.TxType) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateXrpDetailTxType(ctx, sqlc.UpdateXrpDetailTxTypeParams{
		CurrentTxType: txType.Int8(),
		ID:            id,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call UpdateXrpDetailTxType(): %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
	}

	return rowsAffected, nil
}

// UpdateTxTypeBySentHashTx updates txType by tx_blob
func (r *XrpDetailTxInputRepositorySqlc) UpdateTxTypeBySentHashTx(
	txType domainTx.TxType, sentHashTx string,
) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateXrpDetailTxTypeBySentHash(ctx, sqlc.UpdateXrpDetailTxTypeBySentHashParams{
		CurrentTxType: txType.Int8(),
		TxBlob:        sentHashTx, // sentHashTx is actually tx_blob for XRP
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call UpdateXrpDetailTxTypeBySentHash(): %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
	}

	return rowsAffected, nil
}

// Helper functions

func convertSqlcXrpDetailTxToModel(xrpTx *sqlc.XrpDetailTx) *models.XRPDetailTX {
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
