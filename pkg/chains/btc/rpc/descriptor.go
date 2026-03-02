package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ImportDescriptorsRequest is a single descriptor import request for importdescriptors.
type ImportDescriptorsRequest struct {
	Descriptor string `json:"desc"`
	Timestamp  any    `json:"timestamp"`
	Active     bool   `json:"active"`
	Range      any    `json:"range,omitempty"`
	NextIndex  *int   `json:"next_index,omitempty"`
	Label      string `json:"label,omitempty"`
	Internal   bool   `json:"internal,omitempty"`
	Watchonly  bool   `json:"watchonly,omitempty"`
}

// ImportDescriptorsResponse is the response for a single importdescriptors request.
type ImportDescriptorsResponse struct {
	Success  bool                   `json:"success"`
	Warnings []string               `json:"warnings,omitempty"`
	Error    *ImportDescriptorError `json:"error,omitempty"`
}

// ImportDescriptorError contains error details for a failed descriptor import.
type ImportDescriptorError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ImportMultiRequest is a single request for the importmulti RPC.
type ImportMultiRequest struct {
	ScriptPubKey any      `json:"scriptPubKey"`
	Timestamp    any      `json:"timestamp"`
	RedeemScript string   `json:"redeemscript,omitempty"`
	PubKeys      []string `json:"pubkeys,omitempty"`
	Keys         []string `json:"keys,omitempty"`
	Internal     bool     `json:"internal,omitempty"`
	WatchOnly    bool     `json:"watchonly,omitempty"`
	Label        string   `json:"label,omitempty"`
}

// ImportMultiOptions are options for the importmulti RPC.
type ImportMultiOptions struct {
	Rescan bool `json:"rescan"`
}

// ImportMultiResponse is the response for a single importmulti request.
type ImportMultiResponse struct {
	Success  bool              `json:"success"`
	Warnings []string          `json:"warnings,omitempty"`
	Error    *ImportMultiError `json:"error,omitempty"`
}

// ImportMultiError contains error details for a failed importmulti request.
type ImportMultiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// DescriptorInfo is the wire-format response of getdescriptorinfo.
type DescriptorInfo struct {
	Descriptor     string `json:"descriptor"`
	Checksum       string `json:"checksum"`
	IsRange        bool   `json:"isrange"`
	IsSolvable     bool   `json:"issolvable"`
	HasPrivateKeys bool   `json:"hasprivatekeys"`
}

// ListDescriptorsResult is the wire-format response of listdescriptors.
type ListDescriptorsResult struct {
	WalletName  string            `json:"wallet_name"`
	Descriptors []DescriptorEntry `json:"descriptors"`
}

// DescriptorEntry is a single descriptor in ListDescriptorsResult.
type DescriptorEntry struct {
	Desc      string  `json:"desc"`
	Timestamp int64   `json:"timestamp"`
	Active    bool    `json:"active"`
	Internal  *bool   `json:"internal,omitempty"`
	Range     *[2]int `json:"range,omitempty"`
	Next      *int    `json:"next,omitempty"`
}

// ImportDescriptors calls importdescriptors (Bitcoin Core v0.21.0+).
func (c *rpcClient) ImportDescriptors(
	requests []ImportDescriptorsRequest,
) ([]ImportDescriptorsResponse, error) {
	if len(requests) == 0 {
		return nil, errors.New("no descriptors to import")
	}
	bRequests, err := json.Marshal(requests)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(requests): %w", err)
	}
	rawResult, err := c.btcdClient.RawRequest("importdescriptors", []json.RawMessage{bRequests})
	if err != nil {
		return nil, fmt.Errorf("fail to call RawRequest(importdescriptors): %w", err)
	}
	var responses []ImportDescriptorsResponse
	if err := json.Unmarshal(rawResult, &responses); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(importdescriptors): %w", err)
	}
	return responses, nil
}

// ImportMulti calls importmulti (legacy wallet method).
func (c *rpcClient) ImportMulti(
	requests []ImportMultiRequest, options *ImportMultiOptions,
) ([]ImportMultiResponse, error) {
	if len(requests) == 0 {
		return nil, errors.New("no addresses to import")
	}
	bRequests, err := json.Marshal(requests)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(requests): %w", err)
	}
	params := []json.RawMessage{bRequests}
	if options != nil {
		bOptions, err := json.Marshal(options)
		if err != nil {
			return nil, fmt.Errorf("fail to call json.Marshal(options): %w", err)
		}
		params = append(params, bOptions)
	}
	rawResult, err := c.btcdClient.RawRequest("importmulti", params)
	if err != nil {
		return nil, fmt.Errorf("fail to call RawRequest(importmulti): %w", err)
	}
	var responses []ImportMultiResponse
	if err := json.Unmarshal(rawResult, &responses); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(importmulti): %w", err)
	}
	return responses, nil
}

// GetDescriptorInfo calls getdescriptorinfo to analyze a descriptor and calculate its checksum.
func (c *rpcClient) GetDescriptorInfo(descriptor string) (*DescriptorInfo, error) {
	bDescriptor, err := json.Marshal(descriptor)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(descriptor): %w", err)
	}
	rawResult, err := c.btcdClient.RawRequest("getdescriptorinfo", []json.RawMessage{bDescriptor})
	if err != nil {
		return nil, fmt.Errorf("fail to call RawRequest(getdescriptorinfo): %w", err)
	}
	var result DescriptorInfo
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(getdescriptorinfo): %w", err)
	}
	return &result, nil
}

// DeriveAddresses calls the deriveaddresses RPC to derive addresses from a descriptor range.
func (c *rpcClient) DeriveAddresses(descriptor string, startIdx, endIdx uint32) ([]string, error) {
	rangeParam := fmt.Sprintf("[%d,%d]", startIdx, endIdx)
	params := []json.RawMessage{
		json.RawMessage(fmt.Sprintf(`"%s"`, descriptor)),
		json.RawMessage(rangeParam),
	}
	rawResult, err := c.btcdClient.RawRequest("deriveaddresses", params)
	if err != nil {
		return nil, fmt.Errorf("fail to call RawRequest(deriveaddresses): %w", err)
	}
	var addresses []string
	if err := json.Unmarshal(rawResult, &addresses); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(deriveaddresses): %w", err)
	}
	return addresses, nil
}

// ListDescriptors calls listdescriptors to list all imported descriptors.
func (c *rpcClient) ListDescriptors(privateDescriptors bool) (*ListDescriptorsResult, error) {
	bPrivate, err := json.Marshal(privateDescriptors)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(privateDescriptors): %w", err)
	}
	rawResult, err := c.btcdClient.RawRequest("listdescriptors", []json.RawMessage{bPrivate})
	if err != nil {
		return nil, fmt.Errorf("fail to call RawRequest(listdescriptors): %w", err)
	}
	var result ListDescriptorsResult
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(listdescriptors): %w", err)
	}
	return &result, nil
}
