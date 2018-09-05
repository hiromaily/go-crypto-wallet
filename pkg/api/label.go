package api

import (
	"encoding/json"

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

	//TODO:第二パラメータはnilで問題ないか？
	err = json.Unmarshal([]byte(rawResult), nil)
	if err != nil {
		return errors.Errorf("json.Unmarshal(): error: %s", err)
	}

	return nil
}
