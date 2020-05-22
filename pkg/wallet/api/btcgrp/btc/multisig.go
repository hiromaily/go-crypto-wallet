package btc

import (
	"encoding/json"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// AddMultisigAddressResult is response type of PRC `addmultisigaddress`
type AddMultisigAddressResult struct {
	Address      string `json:"address"`
	RedeemScript string `json:"redeemScript"`
}

// AddMultisigAddress create multisig address
//  - requiredSigs: required number of signature for transaction
//  - addresses:    list of addresses(e.g. client, auth1, auth2, auth3
//  - [N:M] e.g. 2:5 => requiredSigs=2, addresses[addr1, addr2, addr3, addr4, addr5]
func (b *Bitcoin) AddMultisigAddress(
	requiredSigs int,
	addresses []string,
	accountName string,
	addressType address.AddrType) (*AddMultisigAddressResult, error) {

	if requiredSigs > len(addresses) {
		return nil, errors.Errorf("number of given address doesn't meet number of requiredSigs: requiredSigs:%d, len(addresses):%d", requiredSigs, len(addresses))
	}

	//requiredSigs
	bRequiredSigs, err := json.Marshal(requiredSigs)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.Marchal(requiredSigs)")
	}

	//addresses
	bAddresses, err := json.Marshal(addresses)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.Marchal(addresses)")
	}

	//accountName
	bAccount, err := json.Marshal(accountName)
	if err != nil {
		b.logger.Warn(
			"fail to json.Marshal(accountName)",
			zap.String("accountName", accountName),
			zap.Error(err))
		bAccount = nil
	}

	//addressType for only BTC
	var jsonRawMsg []json.RawMessage
	switch b.coinTypeCode {
	case coin.BTC:
		bAddrType, err := json.Marshal(addressType.String())
		if err != nil {
			b.logger.Warn(
				"fail to json.Marchal(addressType)",
				zap.String("addressType", addressType.String()),
				zap.Error(err))
			bAddrType = nil
		}

		jsonRawMsg = []json.RawMessage{bRequiredSigs, bAddresses, bAccount, bAddrType}
	case coin.BCH:
		jsonRawMsg = []json.RawMessage{bRequiredSigs, bAddresses, bAccount}
	default:
		return nil, errors.Errorf("not implemented for %s in AddMultisigAddress()", b.coinTypeCode.String())
	}

	//call addmultisigaddress
	rawResult, err := b.Client.RawRequest("addmultisigaddress", jsonRawMsg)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call client.RawRequest(addmultisigaddress)")
	}

	multisigAddrResult := AddMultisigAddressResult{}
	err = json.Unmarshal([]byte(rawResult), &multisigAddrResult)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.Unmarshal(rawResult)")
	}

	return &multisigAddrResult, nil
}
