package btc

import (
	"encoding/json"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/address"
)

// AddMultisigAddressResult addmultisigaddressをcallしたresponseの型
type AddMultisigAddressResult struct {
	Address      string `json:"address"`
	RedeemScript string `json:"redeemScript"`
}

// AddMultisigAddress マルチシグを Rawトランザクション用に作成する
//  - requiredSigs: 取引成立に必要なサイン数
//  - addresses:    自分のアドレス+承認者のアドレスxN をいれていく
// clientのアドレスは不要。payment/receiptのアドレスはmultisig対応しておくこと
func (b *Bitcoin) AddMultisigAddress(requiredSigs int, addresses []string, accountName string, addressType address.AddrType) (*AddMultisigAddressResult, error) {

	if requiredSigs > len(addresses) {
		return nil, errors.New("number of given address should be at least same to requiredSigs or more")
	}

	//requiredSigs
	bRequiredSigs, err := json.Marshal(requiredSigs)
	if err != nil {
		return nil, errors.Errorf("json.Marchal(requiredSigs): error: %v", err)
	}

	//addresses
	bAddresses, err := json.Marshal(addresses)
	if err != nil {
		return nil, errors.Errorf("json.Marchal(addresses): error: %v", err)
	}

	//accountName
	bAccount, err := json.Marshal(accountName)
	if err != nil {
		b.logger.Error(
			"fail to json.Marshal(accountName)",
			zap.String("accountName", accountName),
			zap.Error(err))
		bAccount = nil
	}

	//addressType
	bAddrType, err := json.Marshal(string(addressType))
	if err != nil {
		b.logger.Error(
			"fail to json.Marchal(addressType)",
			zap.String("addressType", addressType.String()),
			zap.Error(err))
		bAddrType = nil
	}

	jsonRawMsg := []json.RawMessage{bRequiredSigs, bAddresses, bAccount, bAddrType}

	//call addmultisigaddress
	rawResult, err := b.client.RawRequest("addmultisigaddress", jsonRawMsg)
	if err != nil {
		return nil, errors.Errorf("client.RawRequest(addmultisigaddress): error: %v", err)
	}

	multisigAddrResult := AddMultisigAddressResult{}
	err = json.Unmarshal([]byte(rawResult), &multisigAddrResult)
	if err != nil {
		return nil, errors.Errorf("json.Unmarshal(multisigAddrResult): error: %v", err)
	}

	return &multisigAddrResult, nil
}
