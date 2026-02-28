package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ─── Fee ─────────────────────────────────────────────────────────────────────

// EstimateSmartFeeResult is the wire-format response of the estimatesmartfee command.
type EstimateSmartFeeResult struct {
	FeeRate float64  `json:"feerate"`
	Errors  []string `json:"errors"`
	Blocks  uint64   `json:"blocks"`
}

// EstimateSmartFee calls estimatesmartfee with the given confirmation target
// and returns the fee rate in BTC/kB.
func EstimateSmartFee(caller RPCCaller, confirmationBlock int) (float64, error) {
	input, err := json.Marshal(confirmationBlock)
	if err != nil {
		return 0, fmt.Errorf("fail to call json.Marshal(confirmationBlock): %w", err)
	}
	rawResult, err := caller.RawRequest("estimatesmartfee", []json.RawMessage{input})
	if err != nil {
		return 0, fmt.Errorf("fail to call RawRequest(estimatesmartfee): %w", err)
	}
	var result EstimateSmartFeeResult
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return 0, fmt.Errorf("fail to call json.Unmarshal(estimatesmartfee): %w", err)
	}
	if len(result.Errors) != 0 {
		return 0, fmt.Errorf("estimatesmartfee response includes error: %s", result.Errors[0])
	}
	return result.FeeRate, nil
}

// ─── Balance ─────────────────────────────────────────────────────────────────

// GetBalance calls getbalance with wildcard account and the given minimum confirmations.
// Returns the balance in BTC as a float64.
func GetBalance(caller RPCCaller, confirmationBlock int) (float64, error) {
	bWildcard, err := json.Marshal("*")
	if err != nil {
		return 0, fmt.Errorf("fail to call json.Marshal(wildcard): %w", err)
	}
	bConf, err := json.Marshal(confirmationBlock)
	if err != nil {
		return 0, fmt.Errorf("fail to call json.Marshal(confirmationBlock): %w", err)
	}
	rawResult, err := caller.RawRequest("getbalance", []json.RawMessage{bWildcard, bConf})
	if err != nil {
		return 0, fmt.Errorf("fail to call RawRequest(getbalance): %w", err)
	}
	var amount float64
	if err := json.Unmarshal(rawResult, &amount); err != nil {
		return 0, fmt.Errorf("fail to call json.Unmarshal(getbalance): %w", err)
	}
	return amount, nil
}

// ─── Label ───────────────────────────────────────────────────────────────────

// SetLabel calls setlabel to assign a label to an address.
func SetLabel(caller RPCCaller, addr, label string) error {
	bAddr, err := json.Marshal(addr)
	if err != nil {
		return fmt.Errorf("fail to call json.Marshal(addr): %w", err)
	}
	bLabel, err := json.Marshal(label)
	if err != nil {
		return fmt.Errorf("fail to call json.Marshal(label): %w", err)
	}
	_, err = caller.RawRequest("setlabel", []json.RawMessage{bAddr, bLabel})
	if err != nil {
		return fmt.Errorf("fail to call RawRequest(setlabel): %w", err)
	}
	return nil
}

// ─── Multisig ────────────────────────────────────────────────────────────────

// AddMultisigAddressResult is the wire-format response of addmultisigaddress.
type AddMultisigAddressResult struct {
	Address      string `json:"address"`
	RedeemScript string `json:"redeemScript"`
}

// AddMultisigAddress calls addmultisigaddress.
// addressType is optional (pass "" to omit the parameter, which is required for BCH).
func AddMultisigAddress(
	caller RPCCaller, requiredSigs int, addresses []string, accountName, addressType string,
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

	rawResult, err := caller.RawRequest("addmultisigaddress", params)
	if err != nil {
		return nil, fmt.Errorf("fail to call RawRequest(addmultisigaddress): %w", err)
	}
	var result AddMultisigAddressResult
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(addmultisigaddress): %w", err)
	}
	return &result, nil
}

// ─── Logging ─────────────────────────────────────────────────────────────────

// LoggingResult is the wire-format response of the logging command.
type LoggingResult struct {
	Net         bool `json:"net"`
	Tor         bool `json:"tor"`
	Mempool     bool `json:"mempool"`
	HTTP        bool `json:"http"`
	Bench       bool `json:"bench"`
	Zmq         bool `json:"zmq"`
	Walletdb    bool `json:"walletdb"`
	RPC         bool `json:"rpc"`
	Estimatefee bool `json:"estimatefee"`
	Addrman     bool `json:"addrman"`
	Selectcoins bool `json:"selectcoins"`
	Reindex     bool `json:"reindex"`
	Cmpctblock  bool `json:"cmpctblock"`
	Rand        bool `json:"rand"`
	Prune       bool `json:"prune"`
	Proxy       bool `json:"proxy"`
	Mempoolrej  bool `json:"mempoolrej"`
	Libevent    bool `json:"libevent"`
	Coindb      bool `json:"coindb"`
	Qt          bool `json:"qt"`
	Leveldb     bool `json:"leveldb"`
	Validation  bool `json:"validation"`
}

// Logging calls the logging RPC and returns the current logging configuration.
func Logging(caller RPCCaller) (*LoggingResult, error) {
	rawResult, err := caller.RawRequest("logging", []json.RawMessage{})
	if err != nil {
		return nil, fmt.Errorf("fail to call RawRequest(logging): %w", err)
	}
	var result LoggingResult
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(logging): %w", err)
	}
	return &result, nil
}

// ─── Import ───────────────────────────────────────────────────────────────────

// ImportAddress calls importaddress to import an address or script.
func ImportAddress(caller RPCCaller, address, label string, rescan bool) error {
	bAddress, err := json.Marshal(address)
	if err != nil {
		return fmt.Errorf("fail to call json.Marshal(address): %w", err)
	}
	bLabel, err := json.Marshal(label)
	if err != nil {
		return fmt.Errorf("fail to call json.Marshal(label): %w", err)
	}
	bRescan, err := json.Marshal(rescan)
	if err != nil {
		return fmt.Errorf("fail to call json.Marshal(rescan): %w", err)
	}
	_, err = caller.RawRequest("importaddress", []json.RawMessage{bAddress, bLabel, bRescan})
	if err != nil {
		return fmt.Errorf("fail to call RawRequest(importaddress): %w", err)
	}
	return nil
}

// ImportPrivKey calls importprivkey to import a WIF-encoded private key.
func ImportPrivKey(caller RPCCaller, wifKey, label string, rescan bool) error {
	bKey, err := json.Marshal(wifKey)
	if err != nil {
		return fmt.Errorf("fail to call json.Marshal(wifKey): %w", err)
	}
	bLabel, err := json.Marshal(label)
	if err != nil {
		return fmt.Errorf("fail to call json.Marshal(label): %w", err)
	}
	bRescan, err := json.Marshal(rescan)
	if err != nil {
		return fmt.Errorf("fail to call json.Marshal(rescan): %w", err)
	}
	_, err = caller.RawRequest("importprivkey", []json.RawMessage{bKey, bLabel, bRescan})
	if err != nil {
		return fmt.Errorf("fail to call RawRequest(importprivkey): %w", err)
	}
	return nil
}

// ─── Descriptor ──────────────────────────────────────────────────────────────

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
func ImportDescriptors(caller RPCCaller, requests []ImportDescriptorsRequest) ([]ImportDescriptorsResponse, error) {
	if len(requests) == 0 {
		return nil, errors.New("no descriptors to import")
	}
	bRequests, err := json.Marshal(requests)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(requests): %w", err)
	}
	rawResult, err := caller.RawRequest("importdescriptors", []json.RawMessage{bRequests})
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
func ImportMulti(
	caller RPCCaller, requests []ImportMultiRequest, options *ImportMultiOptions,
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
	rawResult, err := caller.RawRequest("importmulti", params)
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
func GetDescriptorInfo(caller RPCCaller, descriptor string) (*DescriptorInfo, error) {
	bDescriptor, err := json.Marshal(descriptor)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(descriptor): %w", err)
	}
	rawResult, err := caller.RawRequest("getdescriptorinfo", []json.RawMessage{bDescriptor})
	if err != nil {
		return nil, fmt.Errorf("fail to call RawRequest(getdescriptorinfo): %w", err)
	}
	var result DescriptorInfo
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(getdescriptorinfo): %w", err)
	}
	return &result, nil
}

// ListDescriptors calls listdescriptors to list all imported descriptors.
func ListDescriptors(caller RPCCaller, privateDescriptors bool) (*ListDescriptorsResult, error) {
	bPrivate, err := json.Marshal(privateDescriptors)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(privateDescriptors): %w", err)
	}
	rawResult, err := caller.RawRequest("listdescriptors", []json.RawMessage{bPrivate})
	if err != nil {
		return nil, fmt.Errorf("fail to call RawRequest(listdescriptors): %w", err)
	}
	var result ListDescriptorsResult
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(listdescriptors): %w", err)
	}
	return &result, nil
}

// ─── Wallet lifecycle ────────────────────────────────────────────────────────

// LoadWalletResult is the wire-format response of loadwallet/createwallet.
type LoadWalletResult struct {
	WalletName string `json:"name"`
	Warning    string `json:"warning"`
}

// CreateWalletOptions configures wallet creation.
type CreateWalletOptions struct {
	DisablePrivateKeys bool
	Blank              bool
	Passphrase         string
	AvoidReuse         bool
	Descriptors        bool
	LoadOnStartup      bool
}

// CreateWallet calls createwallet.
func CreateWallet(caller RPCCaller, fileName string, opts *CreateWalletOptions) error {
	if opts == nil {
		opts = &CreateWalletOptions{}
	}
	params := make([]json.RawMessage, 0, 7)
	for _, v := range []any{
		fileName, opts.DisablePrivateKeys, opts.Blank,
		opts.Passphrase, opts.AvoidReuse, opts.Descriptors, opts.LoadOnStartup,
	} {
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("fail to call json.Marshal(createwallet param): %w", err)
		}
		params = append(params, b)
	}
	rawResult, err := caller.RawRequest("createwallet", params)
	if err != nil {
		return fmt.Errorf("fail to call RawRequest(createwallet): %w", err)
	}
	var result LoadWalletResult
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return fmt.Errorf("fail to call json.Unmarshal(createwallet): %w", err)
	}
	if result.Warning != "" {
		return fmt.Errorf("createwallet warning: %s", result.Warning)
	}
	return nil
}

// LoadWallet calls loadwallet.
func LoadWallet(caller RPCCaller, fileName string) error {
	bFileName, err := json.Marshal(fileName)
	if err != nil {
		return fmt.Errorf("fail to call json.Marshal(fileName): %w", err)
	}
	rawResult, err := caller.RawRequest("loadwallet", []json.RawMessage{bFileName})
	if err != nil {
		return fmt.Errorf("fail to call RawRequest(loadwallet): %w", err)
	}
	var result LoadWalletResult
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return fmt.Errorf("fail to call json.Unmarshal(loadwallet): %w", err)
	}
	if result.Warning != "" {
		return fmt.Errorf("loadwallet warning: %s", result.Warning)
	}
	return nil
}

// UnloadWallet calls unloadwallet.
func UnloadWallet(caller RPCCaller, fileName string) error {
	bFileName, err := json.Marshal(fileName)
	if err != nil {
		return fmt.Errorf("fail to call json.Marshal(fileName): %w", err)
	}
	_, err = caller.RawRequest("unloadwallet", []json.RawMessage{bFileName})
	if err != nil {
		return fmt.Errorf("fail to call RawRequest(unloadwallet): %w", err)
	}
	return nil
}
