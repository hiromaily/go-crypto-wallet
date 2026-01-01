package watch

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	portsPersistence "github.com/hiromaily/go-crypto-wallet/internal/application/ports/persistence"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/sqlc"
)

// XRPDetailTxInputRepositorySqlc is repository for xrp_detail_tx table using sqlc
type XRPDetailTxInputRepositorySqlc struct {
	queries      *sqlc.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewXRPDetailTxInputRepositorySqlc returns XRPDetailTxInputRepositorySqlc object
func NewXRPDetailTxInputRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *XRPDetailTxInputRepositorySqlc {
	return &XRPDetailTxInputRepositorySqlc{
		queries:      sqlc.New(dbConn),
		coinTypeCode: coinTypeCode,
	}
}

// GetOne get one record by ID
func (r *XRPDetailTxInputRepositorySqlc) GetOne(id int64) (*sqlc.XrpDetailTx, error) {
	ctx := context.Background()

	xrpTx, err := r.queries.GetXrpDetailTxByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetXrpDetailTxByID(): %w", err)
	}

	return &xrpTx, nil
}

// GetAllByTxID returns all records searched by tx_id
func (r *XRPDetailTxInputRepositorySqlc) GetAllByTxID(id int64) ([]*sqlc.XrpDetailTx, error) {
	ctx := context.Background()

	xrpTxs, err := r.queries.GetXrpDetailTxsByTxID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetXrpDetailTxsByTxID(): %w", err)
	}

	result := make([]*sqlc.XrpDetailTx, len(xrpTxs))
	for i := range xrpTxs {
		result[i] = &xrpTxs[i]
	}

	return result, nil
}

// GetSentHashTx returns list of tx_blob by txType
func (r *XRPDetailTxInputRepositorySqlc) GetSentHashTx(txType domainTx.TxType) ([]string, error) {
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
func (r *XRPDetailTxInputRepositorySqlc) Insert(txItem *sqlc.XrpDetailTx) error {
	ctx := context.Background()

	_, err := r.queries.InsertXrpDetailTx(ctx, sqlc.InsertXrpDetailTxParams{
		TxID:                  txItem.TxID,
		Uuid:                  txItem.Uuid,
		CurrentTxType:         txItem.CurrentTxType,
		SenderAccount:         txItem.SenderAccount,
		SenderAddress:         txItem.SenderAddress,
		ReceiverAccount:       txItem.ReceiverAccount,
		ReceiverAddress:       txItem.ReceiverAddress,
		Amount:                txItem.Amount,
		XrpTxType:             txItem.XrpTxType,
		Fee:                   txItem.Fee,
		Flags:                 txItem.Flags,
		LastLedgerSequence:    txItem.LastLedgerSequence,
		Sequence:              txItem.Sequence,
		SigningPubkey:         txItem.SigningPubkey,
		TxnSignature:          txItem.TxnSignature,
		Hash:                  txItem.Hash,
		EarliestLedgerVersion: txItem.EarliestLedgerVersion,
		SignedTxID:            txItem.SignedTxID,
		TxBlob:                txItem.TxBlob,
		SentUpdatedAt:         txItem.SentUpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to call InsertXrpDetailTx(): %w", err)
	}

	return nil
}

// InsertBulk inserts multiple records
func (r *XRPDetailTxInputRepositorySqlc) InsertBulk(txItems []*sqlc.XrpDetailTx) error {
	for _, item := range txItems {
		if err := r.Insert(item); err != nil {
			return err
		}
	}
	return nil
}

// UpdateAfterTxSent updates when tx sent
func (r *XRPDetailTxInputRepositorySqlc) UpdateAfterTxSent(
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
func (r *XRPDetailTxInputRepositorySqlc) UpdateTxType(id int64, txType domainTx.TxType) (int64, error) {
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
func (r *XRPDetailTxInputRepositorySqlc) UpdateTxTypeBySentHashTx(
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

// WithTx returns a new repository instance that uses the provided transaction
func (r *XRPDetailTxInputRepositorySqlc) WithTx(tx *sql.Tx) portsPersistence.XRPDetailTxRepositorier {
	return &XRPDetailTxInputRepositorySqlc{
		queries:      r.queries.WithTx(tx),
		coinTypeCode: r.coinTypeCode,
	}
}
