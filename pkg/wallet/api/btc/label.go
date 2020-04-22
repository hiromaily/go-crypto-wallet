package btc

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// listlabels

// SetLabel sets account to existing imported address
func (b *Bitcoin) SetLabel(addr, label string) error {
	_, err := b.DecodeAddress(addr)
	if err != nil {
		return errors.Wrapf(err, "fail to call btc.DecodeAddress(%s)", addr)
	}

	input1, err := json.Marshal(string(addr))
	if err != nil {
		return errors.Wrap(err, "fail to call json.Marchal(addr)")
	}
	input2, err := json.Marshal(string(label))
	if err != nil {
		return errors.Wrap(err, "fail to call json.Marchal(label)")
	}

	rawResult, err := b.client.RawRequest("setlabel", []json.RawMessage{input1, input2})
	if err != nil {
		return errors.Wrap(err, "fail to call json.RawRequest(setlabel)")
	}

	var tmp interface{}
	err = json.Unmarshal([]byte(rawResult), &tmp)
	if err != nil {
		return errors.Wrap(err, "fail to call json.Unmarshal(rawResult)")
	}

	return nil
}

// GetReceivedByLabelAndMinConf returns balance by label(account)
// FIXME: even if spent utxo is left as balance
// - please use GetBalanceByAccount()
//func (b *Bitcoin) GetReceivedByLabelAndMinConf(accountName string, minConf int) (btcutil.Amount, error) {
//	input1, err := json.Marshal(accountName)
//	if err != nil {
//		return 0, errors.Wrap(err, "fail to call json.Marchal(accountName)")
//	}
//	input2, err := json.Marshal(minConf)
//	if err != nil {
//		return 0, errors.Wrap(err, "fail to call json.Marchal(minConf)")
//	}
//
//	rawResult, err := b.client.RawRequest("getreceivedbylabel", []json.RawMessage{input1, input2})
//	if err != nil {
//		return 0, errors.Wrap(err, "fail to call json.RawRequest(getreceivedbylabel)")
//	}
//
//	var receivedAmt float64
//	err = json.Unmarshal([]byte(rawResult), &receivedAmt)
//	if err != nil {
//		return 0, errors.Wrap(err, "fail to call json.Unmarshal()")
//	}
//
//	//convert float to amout
//	return b.FloatToAmount(receivedAmt)
//}