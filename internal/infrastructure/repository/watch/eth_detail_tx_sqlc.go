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

// ETHDetailTXInputRepositorySqlc is repository for eth_detail_tx table using sqlc
type ETHDetailTXInputRepositorySqlc struct {
	queries      *sqlc.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewETHDetailTXInputRepositorySqlc returns ETHDetailTXInputRepositorySqlc object
func NewETHDetailTXInputRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *ETHDetailTXInputRepositorySqlc {
	return &ETHDetailTXInputRepositorySqlc{
		queries:      sqlc.New(dbConn),
		coinTypeCode: coinTypeCode,
	}
}

// GetOne get one record by ID
func (r *ETHDetailTXInputRepositorySqlc) GetOne(id int64) (*sqlc.ETHDetailTX, error) {
	ctx := context.Background()

	ethTx, err := r.queries.GetETHDetailTXByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetETHDetailTXByID(): %w", err)
	}

	return &ethTx, nil
}

// GetAllByTxID returns all records searched by tx_id
func (r *ETHDetailTXInputRepositorySqlc) GetAllByTxID(id int64) ([]*sqlc.ETHDetailTX, error) {
	ctx := context.Background()

	ethTxs, err := r.queries.GetETHDetailTXsByTxID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetETHDetailTXsByTxID(): %w", err)
	}

	result := make([]*sqlc.ETHDetailTX, len(ethTxs))
	for i := range ethTxs {
		result[i] = &ethTxs[i]
	}

	return result, nil
}

// GetSentHashTx returns list of sent_hash_tx by txType
func (r *ETHDetailTXInputRepositorySqlc) GetSentHashTx(txType domainTx.TxType) ([]string, error) {
	ctx := context.Background()

	hashes, err := r.queries.GetETHDetailTXSentHashList(ctx, sqlc.GetETHDetailTXSentHashListParams{
		Coin:          sqlc.TxCoin(r.coinTypeCode.String()),
		CurrentTxType: txType.Int8(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetETHDetailTXSentHashList(): %w", err)
	}

	return hashes, nil
}

// Insert inserts one record
func (r *ETHDetailTXInputRepositorySqlc) Insert(txItem *sqlc.ETHDetailTX) error {
	ctx := context.Background()

	_, err := r.queries.InsertETHDetailTX(ctx, sqlc.InsertETHDetailTXParams{
		TxID:              txItem.TxID,
		Uuid:              txItem.Uuid,
		CurrentTxType:     txItem.CurrentTxType,
		SenderAccount:     txItem.SenderAccount,
		SenderAddress:     txItem.SenderAddress,
		ReceiverAccount:   txItem.ReceiverAccount,
		ReceiverAddress:   txItem.ReceiverAddress,
		Amount:            txItem.Amount,
		Fee:               txItem.Fee,
		GasLimit:          txItem.GasLimit,
		Nonce:             txItem.Nonce,
		UnsignedHexTx:     txItem.UnsignedHexTx,
		SignedHexTx:       txItem.SignedHexTx,
		SentHashTx:        txItem.SentHashTx,
		UnsignedUpdatedAt: txItem.UnsignedUpdatedAt,
		SentUpdatedAt:     txItem.SentUpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to call InsertETHDetailTX(): %w", err)
	}

	return nil
}

// InsertBulk inserts multiple records
func (r *ETHDetailTXInputRepositorySqlc) InsertBulk(txItems []*sqlc.ETHDetailTX) error {
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

	result, err := r.queries.UpdateETHDetailTXAfterSent(ctx, sqlc.UpdateETHDetailTXAfterSentParams{
		CurrentTxType: txType.Int8(),
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

	result, err := r.queries.UpdateETHDetailTXType(ctx, sqlc.UpdateETHDetailTXTypeParams{
		CurrentTxType: txType.Int8(),
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

	result, err := r.queries.UpdateETHDetailTXTypeBySentHash(ctx, sqlc.UpdateETHDetailTXTypeBySentHashParams{
		CurrentTxType: txType.Int8(),
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

// WithTx returns a new repository instance that uses the provided transaction
func (r *ETHDetailTXInputRepositorySqlc) WithTx(tx *sql.Tx) portsPersistence.ETHDetailTXRepositorier {
	return &ETHDetailTXInputRepositorySqlc{
		queries:      r.queries.WithTx(tx),
		coinTypeCode: r.coinTypeCode,
	}
}
