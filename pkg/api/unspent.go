package api

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
)

// UnlockAllUnspentTransaction Lockされたトランザクションの解除
func (b *Bitcoin) UnlockAllUnspentTransaction() error {
	list, err := b.client.ListLockUnspent() //[]*wire.OutPoint
	if err != nil {
		return errors.Errorf("client.ListLockUnspent(): error: %s", err)
	}

	if len(list) != 0 {
		err = b.client.LockUnspent(true, list)
		if err != nil {
			//FIXME: -8: Invalid parameter, expected unspent output たまにこのエラーが出る。。。Bitcoin Coreの再起動が必要
			// Bitcoin Coreから先のP2Pネットワークへの接続が失敗しているときに起きる
			// よって、Bitcoin Coreの再起動が必要
			// loggingコマンド, もしくは ~/Library/Application Support/Bitcoin/testnet3/debug.logのチェック??
			return errors.Errorf("client.LockUnspent(): error: %s", err)
		}
	}

	return nil
}

// LockUnspent 渡されたtxIDをロックする
func (b *Bitcoin) LockUnspent(tx btcjson.ListUnspentResult) error {
	txIDHash, err := chainhash.NewHashFromStr(tx.TxID)
	if err != nil {
		return errors.Errorf("chainhash.NewHashFromStr(): error: %s", err)
	}
	outpoint := wire.NewOutPoint(txIDHash, tx.Vout)
	err = b.client.LockUnspent(false, []*wire.OutPoint{outpoint})
	if err != nil {
		return err
	}
	return nil
}
