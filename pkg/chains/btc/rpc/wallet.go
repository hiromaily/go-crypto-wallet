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
func (c *rpcClient) CreateWallet(fileName string, opts *CreateWalletOptions) error {
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
func (c *rpcClient) LoadWallet(fileName string) error {
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
func (c *rpcClient) UnloadWallet(fileName string) error {
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

// walletResult is the wire-format response of dumpwallet/importwallet.
type walletResult struct {
	FileName string `json:"filename"`
}

// HDKeyResult is the response type of RPC `gethdkeys` (Bitcoin Core v28.0+).
type HDKeyResult struct {
	XPub string `json:"xpub"`
}

// CreateWalletDescriptorResult is the response type of RPC `createwalletdescriptor` (Bitcoin Core v28.0+).
type CreateWalletDescriptorResult struct {
	Descriptors []string `json:"descriptors"`
	Warnings    []string `json:"warnings"`
}

// WalletOutputType represents the output type for the createwalletdescriptor RPC.
// Valid values are "bech32" (P2WPKH) and "bech32m" (P2TR).
type WalletOutputType string

const (
	WalletOutputTypeBech32  WalletOutputType = "bech32"  // Native segwit (P2WPKH)
	WalletOutputTypeBech32m WalletOutputType = "bech32m" // Taproot (P2TR)
)

// BackupWallet calls backupwallet.
func (c *rpcClient) BackupWallet(fileName string) error {
	bFileName, err := json.Marshal(fileName)
	if err != nil {
		return fmt.Errorf("fail to call json.Marshal(fileName): %w", err)
	}
	_, err = c.client.RawRequest("backupwallet", []json.RawMessage{bFileName})
	if err != nil {
		return fmt.Errorf("fail to call RawRequest(backupwallet): %w", err)
	}
	return nil
}

// DumpWallet calls dumpwallet.
func (c *rpcClient) DumpWallet(fileName string) error {
	return c.dumpImportWallet(fileName, "dumpwallet")
}

// ImportWallet calls importwallet.
func (c *rpcClient) ImportWallet(fileName string) error {
	return c.dumpImportWallet(fileName, "importwallet")
}

func (c *rpcClient) dumpImportWallet(fileName, method string) error {
	bFileName, err := json.Marshal(fileName)
	if err != nil {
		return fmt.Errorf("fail to call json.Marshal(fileName): %w", err)
	}
	rawResult, err := c.client.RawRequest(method, []json.RawMessage{bFileName})
	if err != nil {
		return fmt.Errorf("fail to call RawRequest(%s): %w", method, err)
	}
	var result walletResult
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return fmt.Errorf("fail to call json.Unmarshal for %s response: %w", method, err)
	}
	return nil
}

// EncryptWallet calls encryptwallet.
func (c *rpcClient) EncryptWallet(passphrase string) error {
	bPassphrase, err := json.Marshal(passphrase)
	if err != nil {
		return fmt.Errorf("fail to call json.Marshal(passphrase): %w", err)
	}
	_, err = c.client.RawRequest("encryptwallet", []json.RawMessage{bPassphrase})
	if err != nil {
		return fmt.Errorf("fail to call RawRequest(encryptwallet): %w", err)
	}
	return nil
}

// WalletLock calls walletlock.
func (c *rpcClient) WalletLock() error {
	return c.client.WalletLock()
}

// WalletPassphrase calls walletpassphrase.
func (c *rpcClient) WalletPassphrase(passphrase string, timeoutSecs int64) error {
	return c.client.WalletPassphrase(passphrase, timeoutSecs)
}

// WalletPassphraseChange calls walletpassphrasechange.
func (c *rpcClient) WalletPassphraseChange(old, newPass string) error {
	return c.client.WalletPassphraseChange(old, newPass)
}

// GetHDKeys calls gethdkeys (Bitcoin Core v28.0+).
func (c *rpcClient) GetHDKeys(active *bool) ([]HDKeyResult, error) {
	var params []json.RawMessage
	if active != nil {
		bActive, err := json.Marshal(map[string]bool{"active": *active})
		if err != nil {
			return nil, fmt.Errorf("fail to call json.Marshal(active): %w", err)
		}
		params = append(params, bActive)
	}
	rawResult, err := c.client.RawRequest("gethdkeys", params)
	if err != nil {
		return nil, fmt.Errorf("fail to call RawRequest(gethdkeys): %w", err)
	}
	var hdKeys []HDKeyResult
	if err := json.Unmarshal(rawResult, &hdKeys); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(gethdkeys): %w", err)
	}
	return hdKeys, nil
}

// CreateWalletDescriptor calls createwalletdescriptor (Bitcoin Core v28.0+).
func (c *rpcClient) CreateWalletDescriptor(outputType WalletOutputType) (*CreateWalletDescriptorResult, error) {
	if outputType != WalletOutputTypeBech32 && outputType != WalletOutputTypeBech32m {
		return nil, fmt.Errorf("invalid output type: %s (must be 'bech32' or 'bech32m')", outputType)
	}
	bType, err := json.Marshal(string(outputType))
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(type): %w", err)
	}
	rawResult, err := c.client.RawRequest("createwalletdescriptor", []json.RawMessage{bType})
	if err != nil {
		return nil, fmt.Errorf("fail to call RawRequest(createwalletdescriptor): %w", err)
	}
	var result CreateWalletDescriptorResult
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(createwalletdescriptor): %w", err)
	}
	return &result, nil
}
