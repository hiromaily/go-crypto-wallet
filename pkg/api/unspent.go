package api

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
)

// UnlockAllUnspentTransaction Lockされたトランザクションの解除
//TODO:手動解除の場合、listlockunspentコマンドでtxidの一覧を出力し
//TODO:              lockunspent true txid一覧のjson
func (b *Bitcoin) UnlockAllUnspentTransaction() error {
	list, err := b.client.ListLockUnspent() //[]*wire.OutPoint
	if len(list) != 0 {
		err = b.client.LockUnspent(true, list)
		if err != nil {
			return errors.Errorf("LockUnspent(): error: %v", err)
		}
	}

	return nil
}

// LockUnspent 渡されたtxidをロックする
func (b *Bitcoin) LockUnspent(tx btcjson.ListUnspentResult) error {
	txIDHash, err := chainhash.NewHashFromStr(tx.TxID)
	if err != nil {
		return err
	}
	outpoint := wire.NewOutPoint(txIDHash, tx.Vout)
	err = b.client.LockUnspent(false, []*wire.OutPoint{outpoint})
	if err != nil {
		return err
	}
	return nil
}
