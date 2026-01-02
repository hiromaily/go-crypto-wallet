package watch

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	portsPersistence "github.com/hiromaily/go-crypto-wallet/internal/application/ports/persistence"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	domainXrp "github.com/hiromaily/go-crypto-wallet/internal/domain/xrp"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
)

// XRPDetailTxInputRepositorySqlc is repository for xrp_detail_tx table using sqlc
type XRPDetailTxInputRepositorySqlc struct {
	queries      *sqlcgen.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewXRPDetailTxInputRepositorySqlc returns XRPDetailTxInputRepositorySqlc object
func NewXRPDetailTxInputRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *XRPDetailTxInputRepositorySqlc {
	return &XRPDetailTxInputRepositorySqlc{
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
	}
}

// convertToXrpDetailTx converts sqlcgen.XrpDetailTx to domain.XrpDetailTx entity
func convertToXrpDetailTx(sqlcTx *sqlcgen.XrpDetailTx) (*domainXrp.XrpDetailTx, error) {
	tx := &domainXrp.XrpDetailTx{
		ID:                    sqlcTx.ID,
		TxID:                  sqlcTx.TxID,
		UUID:                  sqlcTx.Uuid,
		CurrentTxType:         domainTx.TxTypeFromInt8(sqlcTx.CurrentTxType),
		SenderAccount:         sqlcTx.SenderAccount,
		SenderAddress:         sqlcTx.SenderAddress,
		ReceiverAccount:       sqlcTx.ReceiverAccount,
		ReceiverAddress:       sqlcTx.ReceiverAddress,
		Amount:                sqlcTx.Amount,
		XrpTxType:             sqlcTx.XrpTxType,
		Fee:                   sqlcTx.Fee,
		Flags:                 sqlcTx.Flags,
		LastLedgerSequence:    sqlcTx.LastLedgerSequence,
		Sequence:              sqlcTx.Sequence,
		SigningPubkey:         sqlcTx.SigningPubkey,
		TxnSignature:          sqlcTx.TxnSignature,
		Hash:                  sqlcTx.Hash,
		EarliestLedgerVersion: sqlcTx.EarliestLedgerVersion,
		SignedTxID:            sqlcTx.SignedTxID,
		TxBlob:                sqlcTx.TxBlob,
	}

	if sqlcTx.SentUpdatedAt.Valid {
		tx.SentUpdatedAt = &sqlcTx.SentUpdatedAt.Time
	}

	return tx, nil
}

// convertFromXrpDetailTx converts domain.XrpDetailTx entity to sqlcgen.XrpDetailTx
func convertFromXrpDetailTx(tx *domainXrp.XrpDetailTx) *sqlcgen.XrpDetailTx {
	sqlcTx := &sqlcgen.XrpDetailTx{
		ID:                    tx.ID,
		TxID:                  tx.TxID,
		Uuid:                  tx.UUID,
		CurrentTxType:         tx.CurrentTxType.Int8(),
		SenderAccount:         tx.SenderAccount,
		SenderAddress:         tx.SenderAddress,
		ReceiverAccount:       tx.ReceiverAccount,
		ReceiverAddress:       tx.ReceiverAddress,
		Amount:                tx.Amount,
		XrpTxType:             tx.XrpTxType,
		Fee:                   tx.Fee,
		Flags:                 tx.Flags,
		LastLedgerSequence:    tx.LastLedgerSequence,
		Sequence:              tx.Sequence,
		SigningPubkey:         tx.SigningPubkey,
		TxnSignature:          tx.TxnSignature,
		Hash:                  tx.Hash,
		EarliestLedgerVersion: tx.EarliestLedgerVersion,
		SignedTxID:            tx.SignedTxID,
		TxBlob:                tx.TxBlob,
	}

	if tx.SentUpdatedAt != nil {
		sqlcTx.SentUpdatedAt = sql.NullTime{Time: *tx.SentUpdatedAt, Valid: true}
	}

	return sqlcTx
}

// GetOne get one record by ID
func (r *XRPDetailTxInputRepositorySqlc) GetOne(id int64) (*domainXrp.XrpDetailTx, error) {
	ctx := context.Background()

	xrpTx, err := r.queries.GetXrpDetailTxByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetXrpDetailTxByID(): %w", err)
	}

	return convertToXrpDetailTx(&xrpTx)
}

// GetAllByTxID returns all records searched by tx_id
func (r *XRPDetailTxInputRepositorySqlc) GetAllByTxID(id int64) ([]*domainXrp.XrpDetailTx, error) {
	ctx := context.Background()

	xrpTxs, err := r.queries.GetXrpDetailTxsByTxID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetXrpDetailTxsByTxID(): %w", err)
	}

	result := make([]*domainXrp.XrpDetailTx, len(xrpTxs))
	for i := range xrpTxs {
		domainTx, err := convertToXrpDetailTx(&xrpTxs[i])
		if err != nil {
			return nil, fmt.Errorf("failed to convert xrp detail tx: %w", err)
		}
		result[i] = domainTx
	}

	return result, nil
}

// GetSentHashTx returns list of tx_blob by txType
func (r *XRPDetailTxInputRepositorySqlc) GetSentHashTx(txType domainTx.TxType) ([]string, error) {
	ctx := context.Background()

	blobs, err := r.queries.GetXrpDetailTxBlobList(ctx, sqlcgen.GetXrpDetailTxBlobListParams{
		Coin:          sqlcgen.TxCoin(r.coinTypeCode.String()),
		CurrentTxType: txType.Int8(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetXrpDetailTxBlobList(): %w", err)
	}

	return blobs, nil
}

// Insert inserts one record
func (r *XRPDetailTxInputRepositorySqlc) Insert(txItem *domainXrp.XrpDetailTx) error {
	ctx := context.Background()

	sqlcTx := convertFromXrpDetailTx(txItem)
	_, err := r.queries.InsertXrpDetailTx(ctx, sqlcgen.InsertXrpDetailTxParams{
		TxID:                  sqlcTx.TxID,
		Uuid:                  sqlcTx.Uuid,
		CurrentTxType:         sqlcTx.CurrentTxType,
		SenderAccount:         sqlcTx.SenderAccount,
		SenderAddress:         sqlcTx.SenderAddress,
		ReceiverAccount:       sqlcTx.ReceiverAccount,
		ReceiverAddress:       sqlcTx.ReceiverAddress,
		Amount:                sqlcTx.Amount,
		XrpTxType:             sqlcTx.XrpTxType,
		Fee:                   sqlcTx.Fee,
		Flags:                 sqlcTx.Flags,
		LastLedgerSequence:    sqlcTx.LastLedgerSequence,
		Sequence:              sqlcTx.Sequence,
		SigningPubkey:         sqlcTx.SigningPubkey,
		TxnSignature:          sqlcTx.TxnSignature,
		Hash:                  sqlcTx.Hash,
		EarliestLedgerVersion: sqlcTx.EarliestLedgerVersion,
		SignedTxID:            sqlcTx.SignedTxID,
		TxBlob:                sqlcTx.TxBlob,
		SentUpdatedAt:         sqlcTx.SentUpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to call InsertXrpDetailTx(): %w", err)
	}

	return nil
}

// InsertBulk inserts multiple records
func (r *XRPDetailTxInputRepositorySqlc) InsertBulk(txItems []*domainXrp.XrpDetailTx) error {
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

	result, err := r.queries.UpdateXrpDetailTxAfterSent(ctx, sqlcgen.UpdateXrpDetailTxAfterSentParams{
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

	result, err := r.queries.UpdateXrpDetailTxType(ctx, sqlcgen.UpdateXrpDetailTxTypeParams{
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

	result, err := r.queries.UpdateXrpDetailTxTypeBySentHash(ctx, sqlcgen.UpdateXrpDetailTxTypeBySentHashParams{
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
func (r *XRPDetailTxInputRepositorySqlc) WithTx(tx *sql.Tx) portsPersistence.XRPDetailTXRepositorier {
	return &XRPDetailTxInputRepositorySqlc{
		queries:      r.queries.WithTx(tx),
		coinTypeCode: r.coinTypeCode,
	}
}
