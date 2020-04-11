package btc

import (
	"encoding/json"

	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

//GetLabel()は存在しない
//GetAddressInfo()を呼び出すこと

// SetLabel 既存のimport済のアドレスにラベル名をセットする
// version0.18より、setaccountは呼び出せなくなるので、SetLabel()をcallすること
func (b *Bitcoin) SetLabel(addr, label string) error {
	_, err := b.DecodeAddress(addr)
	if err != nil {
		return errors.Errorf("DecodeAddress(%s): error: %s", addr, err)
	}

	input1, err := json.Marshal(string(addr))
	if err != nil {
		return errors.Errorf("json.Marchal(): error: %s", err)
	}
	input2, err := json.Marshal(string(label))
	if err != nil {
		return errors.Errorf("json.Marchal(): error: %s", err)
	}

	rawResult, err := b.client.RawRequest("setlabel", []json.RawMessage{input1, input2})
	if err != nil {
		return errors.Errorf("json.RawRequest(setlabel): error: %s", err)
	}

	var tmp interface{}
	err = json.Unmarshal([]byte(rawResult), &tmp)
	if err != nil {
		return errors.Errorf("json.Unmarshal(): error: %s", err)
	}

	return nil
}

// GetReceivedByLabelAndMinConf ラベルに対してのBalanceを取得する
func (b *Bitcoin) GetReceivedByLabelAndMinConf(accountName string, minConf int) (btcutil.Amount, error) {
	input1, err := json.Marshal(accountName)
	if err != nil {
		return 0, errors.Errorf("json.Marchal(accountName): error: %s", err)
	}
	input2, err := json.Marshal(minConf)
	if err != nil {
		return 0, errors.Errorf("json.Marchal(minConf): error: %s", err)
	}

	rawResult, err := b.client.RawRequest("getreceivedbylabel", []json.RawMessage{input1, input2})
	if err != nil {
		return 0, errors.Errorf("json.RawRequest(getreceivedbylabel): error: %s", err)
	}

	var receivedAmt float64
	err = json.Unmarshal([]byte(rawResult), &receivedAmt)
	if err != nil {
		return 0, errors.Errorf("json.Unmarshal(): error: %s", err)
	}

	//変換
	return b.FloatBitToAmount(receivedAmt)
}
