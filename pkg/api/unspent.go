package api

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
	"log"
)

// UnlockAllUnspentTransaction Lockされたトランザクションの解除
func (b *Bitcoin) UnlockAllUnspentTransaction() error {
	list, err := b.client.ListLockUnspent() //[]*wire.OutPoint
	if err != nil {
		return errors.Errorf("ListLockUnspent(): error: %v", err)
	}

	if len(list) != 0 {
		log.Printf("b.client.ListLockUnspent() %v", list)
		err = b.client.LockUnspent(true, list)
		if err != nil {
			//FIXME: -8: Invalid parameter, expected unspent output たまにこのエラーが出る。。。
			// 以前はBitcoin Coreを再起動したら直ったが。。。
			// loggingコマンド, もしくは ./testnet3/wallets/db.logのチェック??
			return errors.Errorf("LockUnspent(): error: %v", err)
		}
	}

	return nil
}

// LockUnspent 渡されたtxIDをロックする
func (b *Bitcoin) LockUnspent(tx btcjson.ListUnspentResult) error {
	txIDHash, err := chainhash.NewHashFromStr(tx.TxID)
	if err != nil {
		return errors.Errorf("chainhash.NewHashFromStr(): error: %v", err)
	}
	outpoint := wire.NewOutPoint(txIDHash, tx.Vout)
	err = b.client.LockUnspent(false, []*wire.OutPoint{outpoint})
	if err != nil {
		return err
	}
	return nil
}
