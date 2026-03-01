package btc

import (
	"fmt"
)

// listlabels

// SetLabel sets account to existing imported address
func (b *Bitcoin) SetLabel(addr, label string) error {
	_, err := b.DecodeAddress(addr)
	if err != nil {
		return fmt.Errorf("fail to call btc.DecodeAddress(%s): %w", addr, err)
	}

	if err := b.pkgrpc.SetLabel(addr, label); err != nil {
		return fmt.Errorf("fail to call btcrpc.SetLabel(): %w", err)
	}

	return nil
}

// GetReceivedByLabelAndMinConf returns balance by label(account)
// FIXME: even if spent utxo is left as balance
// - please use GetBalanceByAccount()
// func (b *Bitcoin) GetReceivedByLabelAndMinConf(accountName string, minConf int) (btcutil.Amount, error) {
//	input1, err := json.Marshal(accountName)
//	if err != nil {
//		return 0, fmt.Errorf("fail to call json.Marchal(accountName): %w", err)
//	}
//	input2, err := json.Marshal(minConf)
//	if err != nil {
//		return 0, fmt.Errorf("fail to call json.Marchal(minConf): %w", err)
//	}
//
//	rawResult, err := b.Client.RawRequest("getreceivedbylabel", []json.RawMessage{input1, input2})
//	if err != nil {
//		return 0, fmt.Errorf("fail to call json.RawRequest(getreceivedbylabel): %w", err)
//	}
//
//	var receivedAmt float64
//	err = json.Unmarshal([]byte(rawResult), &receivedAmt)
//	if err != nil {
//		return 0, fmt.Errorf("fail to call json.Unmarshal(): %w", err)
//	}
//
//	//convert float to amout
//	return b.FloatToAmount(receivedAmt)
//}
