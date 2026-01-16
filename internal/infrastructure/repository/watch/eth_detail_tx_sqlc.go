package watch

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	portsRepository "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainEth "github.com/hiromaily/go-crypto-wallet/internal/domain/ethereum"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
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

// convertToEthDetailTx converts sqlcgen.EthDetailTx to domain.EthDetailTx entity
func convertToEthDetailTx(sqlcTx *sqlcgen.EthDetailTx) (*domainEth.EthDetailTx, error) {
	currentTxType, err := domainTx.TxTypeFromInt8(sqlcTx.CurrentTxType)
	if err != nil {
		return nil, fmt.Errorf("invalid tx type in database: %w", err)
	}

	tx := &domainEth.EthDetailTx{
		ID:              sqlcTx.ID,
		TxID:            sqlcTx.TxID,
		UUID:            sqlcTx.Uuid,
		CurrentTxType:   currentTxType,
		SenderAccount:   sqlcTx.SenderAccount,
		SenderAddress:   sqlcTx.SenderAddress,
		ReceiverAccount: sqlcTx.ReceiverAccount,
		ReceiverAddress: sqlcTx.ReceiverAddress,
		Amount:          sqlcTx.Amount,
		Fee:             sqlcTx.Fee,
		GasLimit:        sqlcTx.GasLimit,
		Nonce:           sqlcTx.Nonce,
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

// convertFromEthDetailTx converts domain.EthDetailTx entity to sqlcgen.EthDetailTx
func convertFromEthDetailTx(tx *domainEth.EthDetailTx) *sqlcgen.EthDetailTx {
	sqlcTx := &sqlcgen.EthDetailTx{
		ID:              tx.ID,
		TxID:            tx.TxID,
		Uuid:            tx.UUID,
		CurrentTxType:   tx.CurrentTxType.Int8(),
		SenderAccount:   tx.SenderAccount,
		SenderAddress:   tx.SenderAddress,
		ReceiverAccount: tx.ReceiverAccount,
		ReceiverAddress: tx.ReceiverAddress,
		Amount:          tx.Amount,
		Fee:             tx.Fee,
		GasLimit:        tx.GasLimit,
		Nonce:           tx.Nonce,
		UnsignedHexTx:   tx.UnsignedHexTx,
		SignedHexTx:     tx.SignedHexTx,
		SentHashTx:      tx.SentHashTx,
	}

	if tx.UnsignedUpdatedAt != nil {
		sqlcTx.UnsignedUpdatedAt = sql.NullTime{Time: *tx.UnsignedUpdatedAt, Valid: true}
	}
	if tx.SentUpdatedAt != nil {
		sqlcTx.SentUpdatedAt = sql.NullTime{Time: *tx.SentUpdatedAt, Valid: true}
	}

	return sqlcTx
}

// GetOne get one record by ID
func (r *ETHDetailTXInputRepositorySqlc) GetOne(id int64) (*domainEth.EthDetailTx, error) {
	ctx := context.Background()

	ethTx, err := r.queries.GetETHDetailTXByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetETHDetailTXByID(): %w", err)
	}

	return convertToEthDetailTx(&ethTx)
}

// GetAllByTxID returns all records searched by tx_id
func (r *ETHDetailTXInputRepositorySqlc) GetAllByTxID(id int64) ([]*domainEth.EthDetailTx, error) {
	ctx := context.Background()

	ethTxs, err := r.queries.GetETHDetailTXsByTxID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetETHDetailTXsByTxID(): %w", err)
	}

	result := make([]*domainEth.EthDetailTx, len(ethTxs))
	for i := range ethTxs {
		domainTx, err := convertToEthDetailTx(&ethTxs[i])
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
		Coin:          sqlcgen.TxCoin(r.coinTypeCode.String()),
		CurrentTxType: txType.Int8(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetETHDetailTXSentHashList(): %w", err)
	}

	return hashes, nil
}

// Insert inserts one record
func (r *ETHDetailTXInputRepositorySqlc) Insert(txItem *domainEth.EthDetailTx) error {
	ctx := context.Background()

	sqlcTx := convertFromEthDetailTx(txItem)
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
func (r *ETHDetailTXInputRepositorySqlc) InsertBulk(txItems []*domainEth.EthDetailTx) error {
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

	result, err := r.queries.UpdateETHDetailTXType(ctx, sqlcgen.UpdateETHDetailTXTypeParams{
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

	result, err := r.queries.UpdateETHDetailTXTypeBySentHash(ctx, sqlcgen.UpdateETHDetailTXTypeBySentHashParams{
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
func (r *ETHDetailTXInputRepositorySqlc) WithTx(tx *sql.Tx) portsRepository.ETHDetailTXRepositorier {
	return &ETHDetailTXInputRepositorySqlc{
		queries:      r.queries.WithTx(tx),
		coinTypeCode: r.coinTypeCode,
	}
}
