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
	"github.com/hiromaily/go-bitcoin/pkg/enum"
)

// ListUnspentResult listunspentの戻り値
type ListUnspentResult struct {
	TxID          string  `json:"txid"`
	Vout          uint32  `json:"vout"`
	Address       string  `json:"address"`
	Label         string  `json:"label"` //to account
	ScriptPubKey  string  `json:"scriptPubKey"`
	Amount        float64 `json:"amount"`
	Confirmations int64   `json:"confirmations"`
	RedeemScript  string  `json:"redeemScript"`
	Spendable     bool    `json:"spendable"`
	Solvable      bool    `json:"solvable"` //new
	Safe          bool    `json:"safe"`     //new
}

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

func (b *Bitcoin) convertListUnspent(listUnspent []btcjson.ListUnspentResult) []ListUnspentResult {
	//[]btcjson.ListUnspentResult
	converted := make([]ListUnspentResult, len(listUnspent))
	for idx, val := range listUnspent {
		converted[idx].TxID = val.TxID
		converted[idx].Vout = val.Vout
		converted[idx].Address = val.Address
		converted[idx].Label = val.Account
		converted[idx].ScriptPubKey = val.ScriptPubKey
		converted[idx].Amount = val.Amount
		converted[idx].Confirmations = val.Confirmations
		converted[idx].Spendable = val.Spendable
	}

	return converted
}

// ListUnspent listunspentを呼び出す
func (b *Bitcoin) ListUnspent() ([]ListUnspentResult, error) {
	if b.Version() >= enum.BTCVer17 {
		return b.listUnspentVer17()
	}
	return b.listUnspentVer16()
}

func (b *Bitcoin) listUnspentVer16() ([]ListUnspentResult, error) {
	listUnspentResult, err := b.client.ListUnspentMin(b.ConfirmationBlock())
	if err != nil {
		return nil, errors.Errorf("client.ListUnspentMin(): error: %s", err)
	}

	if len(listUnspentResult) == 0 {
		return nil, nil
	}

	return b.convertListUnspent(listUnspentResult), nil
}

func (b *Bitcoin) listUnspentVer17() ([]ListUnspentResult, error) {
	input, err := json.Marshal(uint64(b.confirmationBlock))
	if err != nil {
		return nil, errors.Errorf("json.Marchal(): error: %s", err)
	}
	rawResult, err := b.client.RawRequest("listunspent", []json.RawMessage{input})
	if err != nil {
		return nil, errors.Errorf("json.RawRequest(listunspent): error: %s", err)
	}

	var listunspentResult []ListUnspentResult
	err = json.Unmarshal([]byte(rawResult), &listunspentResult)
	if err != nil {
		return nil, errors.Errorf("json.Unmarshal(): error: %s", err)
	}

	if len(listunspentResult) == 0 {
		return nil, nil
	}

	return listunspentResult, nil
}

// ListUnspentByAccount 指定したアカウントのlistunspentを取得する
func (b *Bitcoin) ListUnspentByAccount(accountType account.AccountType) ([]ListUnspentResult, []btcutil.Address, error) {
	addrs, err := b.GetAddressesByAccount(string(accountType))
	if err != nil {
		return nil, nil, errors.Errorf("BTC.GetAddressesByAccount(): error: %s", err)
	}
	if len(addrs) == 0 {
		return nil, nil, errors.Errorf("%s addresses could not be found", accountType)
	}

	var unspentList []ListUnspentResult

	if b.Version() >= enum.BTCVer17 {
		unspentList, err = b.listUnspentByAccountVer17(addrs)
		if err != nil {
			return nil, nil, errors.Errorf("BTC.listUnspentByAccountVer17() error: %s", err)
		}
	} else {
		unspentList, err = b.listUnspentByAccountVer16(addrs)
		if err != nil {
			return nil, nil, errors.Errorf("BTC.listUnspentByAccountVer16() error: %s", err)
		}
	}

	//送金の金額と近しいutxoでtxを作成するため、ソートしておく => 小さなutxoから利用していくのに便利だが、MUSTではない
	sort.Slice(unspentList, func(i, j int) bool {
		//small to big
		return unspentList[i].Amount < unspentList[j].Amount
	})

	return unspentList, addrs, nil
}

func (b *Bitcoin) listUnspentByAccountVer16(addrs []btcutil.Address) ([]ListUnspentResult, error) {
	listUnspentResult, err := b.client.ListUnspentMinMaxAddresses(b.ConfirmationBlock(), 9999999, addrs)
	if err != nil {
		//ListUnspentが実行できない。致命的なエラー。この場合BitcoinCoreの再起動が必要
		return nil, err
	}

	if len(listUnspentResult) == 0 {
		return nil, nil
	}

	return b.convertListUnspent(listUnspentResult), nil
}

func (b *Bitcoin) listUnspentByAccountVer17(addrs []btcutil.Address) ([]ListUnspentResult, error) {
	input1, err := json.Marshal(uint64(b.confirmationBlock))
	if err != nil {
		return nil, errors.Errorf("json.Marchal(): error: %s", err)
	}

	input2, err := json.Marshal(uint64(9999999))
	if err != nil {
		return nil, errors.Errorf("json.Marchal(): error: %s", err)
	}

	//address
	strAddrs := make([]string, len(addrs))
	for idx, addr := range addrs {
		strAddrs[idx] = addr.String()
	}

	input3, err := json.Marshal(strAddrs)
	if err != nil {
		return nil, errors.Errorf("json.Marchal(): error: %s", err)
	}

	rawResult, err := b.client.RawRequest("listunspent", []json.RawMessage{input1, input2, input3})
	if err != nil {
		return nil, errors.Errorf("json.RawRequest(listunspent): error: %s", err)
	}

	var listunspentResult []ListUnspentResult
	err = json.Unmarshal([]byte(rawResult), &listunspentResult)
	if err != nil {
		return nil, errors.Errorf("json.Unmarshal(): error: %s", err)
	}

	if len(listunspentResult) == 0 {
		return nil, nil
	}

	return listunspentResult, nil
}
