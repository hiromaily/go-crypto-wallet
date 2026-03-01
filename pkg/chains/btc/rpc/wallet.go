package rpc

import (
	"encoding/json"
	"fmt"
)

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
func (c *Client) CreateWallet(fileName string, opts *CreateWalletOptions) error {
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
	rawResult, err := c.client.RawRequest("createwallet", params)
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
func (c *Client) LoadWallet(fileName string) error {
	bFileName, err := json.Marshal(fileName)
	if err != nil {
		return fmt.Errorf("fail to call json.Marshal(fileName): %w", err)
	}
	rawResult, err := c.client.RawRequest("loadwallet", []json.RawMessage{bFileName})
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
func (c *Client) UnloadWallet(fileName string) error {
	bFileName, err := json.Marshal(fileName)
	if err != nil {
		return fmt.Errorf("fail to call json.Marshal(fileName): %w", err)
	}
	_, err = c.client.RawRequest("unloadwallet", []json.RawMessage{bFileName})
	if err != nil {
		return fmt.Errorf("fail to call RawRequest(unloadwallet): %w", err)
	}
	return nil
}
