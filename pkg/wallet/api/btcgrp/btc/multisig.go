package btc

import (
	"encoding/json"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// AddMultisigAddressResult is response type of PRC `addmultisigaddress`
type AddMultisigAddressResult struct {
	Address      string `json:"address"`
	RedeemScript string `json:"redeemScript"`
}

// AddMultisigAddress create multisig address
//   - requiredSigs: required number of signature for transaction
//   - addresses:    list of addresses(e.g. client, auth1, auth2, auth3
//   - [N:M] e.g. 2:5 => requiredSigs=2, addresses[addr1, addr2, addr3, addr4, addr5]
func (b *Bitcoin) AddMultisigAddress(
	requiredSigs int,
	addresses []string,
	accountName string,
	addressType address.AddrType,
) (*AddMultisigAddressResult, error) {
	if requiredSigs > len(addresses) {
		return nil, fmt.Errorf(
			"number of given address doesn't meet number of requiredSigs: requiredSigs:%d, len(addresses):%d",
			requiredSigs, len(addresses))
	}

	// requiredSigs
	bRequiredSigs, err := json.Marshal(requiredSigs)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marchal(requiredSigs): %w", err)
	}

	// addresses
	bAddresses, err := json.Marshal(addresses)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marchal(addresses): %w", err)
	}

	// accountName
	bAccount, err := json.Marshal(accountName)
	if err != nil {
		logger.Warn(
			"fail to json.Marshal(accountName)",
			"accountName", accountName,
			"error", err)
		bAccount = nil
	}

	// addressType for only BTC
	var jsonRawMsg []json.RawMessage
	switch b.coinTypeCode {
	case coin.BTC:
		var bAddrType []byte
		bAddrType, err = json.Marshal(addressType.String())
		if err != nil {
			logger.Warn(
				"fail to json.Marchal(addressType)",
				"addressType", addressType.String(),
				"error", err)
			bAddrType = nil
		}

		jsonRawMsg = []json.RawMessage{bRequiredSigs, bAddresses, bAccount, bAddrType}
	case coin.BCH:
		jsonRawMsg = []json.RawMessage{bRequiredSigs, bAddresses, bAccount}
	case coin.LTC, coin.ETH, coin.XRP, coin.ERC20, coin.HYC:
		return nil, fmt.Errorf("not implemented for %s in AddMultisigAddress()", b.coinTypeCode.String())
	default:
		return nil, fmt.Errorf("not implemented for %s in AddMultisigAddress()", b.coinTypeCode.String())
	}

	// call addmultisigaddress
	rawResult, err := b.Client.RawRequest("addmultisigaddress", jsonRawMsg)
	if err != nil {
		return nil, fmt.Errorf("fail to call client.RawRequest(addmultisigaddress): %w", err)
	}

	multisigAddrResult := AddMultisigAddressResult{}
	err = json.Unmarshal(rawResult, &multisigAddrResult)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(rawResult): %w", err)
	}

	return &multisigAddrResult, nil
}
