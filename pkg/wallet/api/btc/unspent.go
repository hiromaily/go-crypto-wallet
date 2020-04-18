package btc

import (
	"encoding/json"
	"sort"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/account"
)

// ListUnspentResult is response type of PRC `listunspent`
type ListUnspentResult struct {
	TxID          string  `json:"txid"`
	Vout          uint32  `json:"vout"`
	Address       string  `json:"address"`
	Label         string  `json:"label"` //to client
	RedeemScript  string  `json:"redeemScript"`
	ScriptPubKey  string  `json:"scriptPubKey"`
	Amount        float64 `json:"amount"`
	Confirmations int64   `json:"confirmations"`
	Spendable     bool    `json:"spendable"`
	Solvable      bool    `json:"solvable"`
	Desc          string  `json:"desc"`
	Safe          bool    `json:"safe"`
}

// UnlockAllUnspentTransaction unlock locked unspent tx
func (b *Bitcoin) UnlockAllUnspentTransaction() error {
	list, err := b.client.ListLockUnspent() //[]*wire.OutPoint
	if err != nil {
		return errors.Wrap(err, "fail to call client.ListLockUnspent()")
	}

	if len(list) != 0 {
		err = b.client.LockUnspent(true, list)
		if err != nil {
			return errors.Wrap(err, "fail to call client.LockUnspent()")
		}
	}

	return nil
}

// LockUnspent lock given txID
func (b *Bitcoin) LockUnspent(tx btcjson.ListUnspentResult) error {
	txIDHash, err := chainhash.NewHashFromStr(tx.TxID)
	if err != nil {
		return errors.Wrapf(err, "fail to call chainhash.NewHashFromStr(%s)", tx.TxID)
	}
	outpoint := wire.NewOutPoint(txIDHash, tx.Vout)
	err = b.client.LockUnspent(false, []*wire.OutPoint{outpoint})
	if err != nil {
		return err
	}
	return nil
}

// ListUnspent call RPC `listunspent`
func (b *Bitcoin) ListUnspent() ([]ListUnspentResult, error) {
	input, err := json.Marshal(uint64(b.confirmationBlock))
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.Marchal()")
	}
	rawResult, err := b.client.RawRequest("listunspent", []json.RawMessage{input})
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.RawRequest(listunspent)")
	}

	var listunspentResult []ListUnspentResult
	err = json.Unmarshal([]byte(rawResult), &listunspentResult)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.Unmarshal()")
	}

	if len(listunspentResult) == 0 {
		return nil, nil
	}

	return listunspentResult, nil
}

// ListUnspentByAccount get listunspent by account
func (b *Bitcoin) ListUnspentByAccount(accountType account.AccountType) ([]ListUnspentResult, []btcutil.Address, error) {
	addrs, err := b.GetAddressesByLabel(accountType.String())
	if err != nil {
		return nil, nil, errors.Wrap(err, "fail to call btc.GetAddressesByLabel()")
	}
	if len(addrs) == 0 {
		return nil, nil, errors.Errorf("address for %s can not be found", accountType)
	}

	var unspentList []ListUnspentResult

	unspentList, err = b.listUnspentByAccount(addrs)
	if err != nil {
		return nil, nil, errors.Wrap(err, "fail to call btc.listUnspentByAccount()")
	}

	// sort amount by ascending (small to big)
	sort.Slice(unspentList, func(i, j int) bool {
		//small to big
		return unspentList[i].Amount < unspentList[j].Amount
	})

	return unspentList, addrs, nil
}

func (b *Bitcoin) listUnspentByAccount(addrs []btcutil.Address) ([]ListUnspentResult, error) {
	input1, err := json.Marshal(uint64(b.confirmationBlock))
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.Marchal(confirmationBlock)")
	}

	input2, err := json.Marshal(uint64(9999999))
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.Marchal(9999999)")
	}

	//address
	strAddrs := make([]string, len(addrs))
	for idx, addr := range addrs {
		strAddrs[idx] = addr.String()
	}

	input3, err := json.Marshal(strAddrs)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.Marchal(addresses)")
	}

	rawResult, err := b.client.RawRequest("listunspent", []json.RawMessage{input1, input2, input3})
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.RawRequest(listunspent)")
	}

	var listunspentResult []ListUnspentResult
	err = json.Unmarshal([]byte(rawResult), &listunspentResult)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.Unmarshal(rawResult)")
	}

	if len(listunspentResult) == 0 {
		return nil, nil
	}

	return listunspentResult, nil
}
