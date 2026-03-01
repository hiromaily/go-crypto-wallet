package rpc

import (
	"encoding/json"
	"fmt"
)

// AddMultisigAddressResult is the wire-format response of addmultisigaddress.
type AddMultisigAddressResult struct {
	Address      string `json:"address"`
	RedeemScript string `json:"redeemScript"`
}

// AddMultisigAddress calls addmultisigaddress.
// addressType is optional (pass "" to omit the parameter, which is required for BCH).
func (c *Client) AddMultisigAddress(
	requiredSigs int, addresses []string, accountName, addressType string,
) (*AddMultisigAddressResult, error) {
	bReqSigs, err := json.Marshal(requiredSigs)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(requiredSigs): %w", err)
	}
	bAddresses, err := json.Marshal(addresses)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(addresses): %w", err)
	}
	bAccount, err := json.Marshal(accountName)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(accountName): %w", err)
	}

	params := []json.RawMessage{bReqSigs, bAddresses, bAccount}
	if addressType != "" {
		bAddrType, err := json.Marshal(addressType)
		if err != nil {
			return nil, fmt.Errorf("fail to call json.Marshal(addressType): %w", err)
		}
		params = append(params, bAddrType)
	}

	rawResult, err := c.client.RawRequest("addmultisigaddress", params)
	if err != nil {
		return nil, fmt.Errorf("fail to call RawRequest(addmultisigaddress): %w", err)
	}
	var result AddMultisigAddressResult
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(addmultisigaddress): %w", err)
	}
	return &result, nil
}
