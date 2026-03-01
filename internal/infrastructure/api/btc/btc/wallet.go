package btc

import (
	"encoding/json"
	"fmt"

	btcrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/rpc"
)

//{
//  "filename": "/Applications/sample01.dat"
//}

const (
	dumpWallet   = "dumpwallet"
	importWallet = "importwallet"
)

// WalletResult is response type of PRC `dumpwallet`/`dumpwallet`
type WalletResult struct {
	FileName string `json:"filename"`
}

// BackupWallet unload wallet.dat
//
//	Safely copies current wallet file to destination, which can be a directory or a path with filename
func (b *Bitcoin) BackupWallet(fileName string) error {
	// backupwallet
	bFileName, err := json.Marshal(fileName)
	if err != nil {
		return fmt.Errorf("fail to call json.Marchal(filename): %w", err)
	}
	_, err = b.Client.RawRequest("backupwallet", []json.RawMessage{bFileName})
	if err != nil {
		return fmt.Errorf("fail to call json.RawRequest(backupwallet): %w", err)
	}

	return nil
}

// DumpWallet dump wallet.dat
//   - fileName: full path
//   - This does not allow overwriting existing files.
func (b *Bitcoin) DumpWallet(fileName string) error {
	return b.dumpImportWallet(fileName, dumpWallet)
}

// ImportWallet import wallet.dat
func (b *Bitcoin) ImportWallet(fileName string) error {
	return b.dumpImportWallet(fileName, importWallet)
}

func (b *Bitcoin) dumpImportWallet(fileName, method string) error {
	bFileName, err := json.Marshal(fileName)
	if err != nil {
		return fmt.Errorf("fail to call json.Marchal(filename): %w", err)
	}

	rawResult, err := b.Client.RawRequest(method, []json.RawMessage{bFileName})
	if err != nil {
		return fmt.Errorf("fail to call json.RawRequest(%s): %w", method, err)
	}

	var walletResult WalletResult
	err = json.Unmarshal(rawResult, &walletResult)
	if err != nil {
		return fmt.Errorf("json.Unmarshal(rawResult): %w", err)
	}

	return nil
}

// EncryptWallet encrypt wallet by pass phrase
// https://bitcoincore.org/en/doc/0.19.0/rpc/wallet/encryptwallet/
// Note: For new wallets, consider using the passphrase parameter in createwallet
// instead of encrypting after creation.
func (b *Bitcoin) EncryptWallet(passphrase string) error {
	// backupwallet
	input1, err := json.Marshal(passphrase)
	if err != nil {
		return fmt.Errorf("fail to call json.Marchal(passphrase): %w", err)
	}
	_, err = b.Client.RawRequest("encryptwallet", []json.RawMessage{input1})
	if err != nil {
		return fmt.Errorf("fail to call json.RawRequest(encryptwallet): %w", err)
	}

	return nil
}

// WalletLock lock wallet
func (b *Bitcoin) WalletLock() error {
	return b.Client.WalletLock()
}

// WalletPassphrase unlock wallet
func (b *Bitcoin) WalletPassphrase(passphrase string, timeoutSecs int64) error {
	return b.Client.WalletPassphrase(passphrase, timeoutSecs)
}

// WalletPassphraseChange change pass phrase
func (b *Bitcoin) WalletPassphraseChange(old, newPass string) error {
	return b.Client.WalletPassphraseChange(old, newPass)
}

// LoadWallet loads a wallet from a wallet file or directory
func (b *Bitcoin) LoadWallet(fileName string) error {
	return b.pkgrpc.LoadWallet(fileName)
}

// UnLoadWallet unload wallet dat
//
//	Unloads the wallet referenced by the request endpoint otherwise unloads the wallet specified in the argument.
//	Specifying the wallet name on a wallet endpoint is invalid.
//	e.g. bitcoin-cli unloadwallet wallet_name
func (b *Bitcoin) UnLoadWallet(fileName string) error {
	if err := b.pkgrpc.UnloadWallet(fileName); err != nil {
		return fmt.Errorf("fail to call btcrpc.UnloadWallet(): %w", err)
	}

	return nil
}

// CreateWalletOptions configures wallet creation parameters
type CreateWalletOptions struct {
	DisablePrivateKeys bool   // If true, wallet will not have private keys (watch-only)
	Blank              bool   // If true, create a blank wallet
	Passphrase         string // Encrypt wallet with this passphrase
	AvoidReuse         bool   // If true, wallet will avoid address reuse
	Descriptors        bool   // If true, create a descriptor wallet (Bitcoin Core v0.21.0+)
	LoadOnStartup      bool   // If true, load wallet on startup
}

// CreateWallet Creates and loads a new wallet
// For descriptor wallets, set opts.Descriptors = true (Bitcoin Core v0.21.0+)
func (b *Bitcoin) CreateWallet(fileName string, disablePrivKey bool) error {
	return b.CreateWalletWithOptions(fileName, &CreateWalletOptions{
		DisablePrivateKeys: disablePrivKey,
	})
}

// CreateWalletWithOptions creates and loads a new wallet with advanced options
// Bitcoin Core RPC: createwallet "wallet_name"
//
//	( disable_private_keys blank "passphrase" avoid_reuse descriptors load_on_startup )
func (b *Bitcoin) CreateWalletWithOptions(fileName string, opts *CreateWalletOptions) error {
	if opts == nil {
		opts = &CreateWalletOptions{}
	}

	pkgOpts := &btcrpc.CreateWalletOptions{
		DisablePrivateKeys: opts.DisablePrivateKeys,
		Blank:              opts.Blank,
		Passphrase:         opts.Passphrase,
		AvoidReuse:         opts.AvoidReuse,
		Descriptors:        opts.Descriptors,
		LoadOnStartup:      opts.LoadOnStartup,
	}

	if err := b.pkgrpc.CreateWallet(fileName, pkgOpts); err != nil {
		return fmt.Errorf("fail to call btcrpc.CreateWallet(): %w", err)
	}

	return nil
}

// ===== Bitcoin Core v28.0+ New RPC Methods =====

// HDKeyResult is response type of RPC `gethdkeys` (Bitcoin Core v28.0+)
type HDKeyResult struct {
	XPub string `json:"xpub"`
}

// GetHDKeys lists all BIP32 HD keys in use by wallet descriptors
// Bitcoin Core v28.0+ only
// https://bitcoincore.org/en/doc/28.0/rpc/wallet/gethdkeys/
func (b *Bitcoin) GetHDKeys(active *bool) ([]HDKeyResult, error) {
	var params []json.RawMessage

	// Add active parameter if specified
	if active != nil {
		bActive, err := json.Marshal(map[string]bool{"active": *active})
		if err != nil {
			return nil, fmt.Errorf("fail to call json.Marshal(active): %w", err)
		}
		params = append(params, bActive)
	}

	rawResult, err := b.Client.RawRequest("gethdkeys", params)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.RawRequest(gethdkeys): %w", err)
	}

	var hdKeys []HDKeyResult
	err = json.Unmarshal(rawResult, &hdKeys)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(rawResult): %w", err)
	}

	return hdKeys, nil
}

// CreateWalletDescriptorResult is response type of RPC `createwalletdescriptor` (Bitcoin Core v28.0+)
type CreateWalletDescriptorResult struct {
	Descriptors []string `json:"descriptors"` // Array of descriptors that were added
	Warnings    []string `json:"warnings"`    // Any warnings that occurred
}

// DescriptorType represents the type of descriptor to create
type DescriptorType string

const (
	DescriptorTypeBech32  DescriptorType = "bech32"  // Native segwit (P2WPKH)
	DescriptorTypeBech32m DescriptorType = "bech32m" // Taproot (P2TR)
)

// CreateWalletDescriptor adds new automatically generated descriptors to the wallet
// Bitcoin Core v28.0+ only
// This allows upgrading wallets to support new standards like taproot
// https://bitcoincore.org/en/doc/28.0/rpc/wallet/createwalletdescriptor/
func (b *Bitcoin) CreateWalletDescriptor(descriptorType DescriptorType) (*CreateWalletDescriptorResult, error) {
	// Validate descriptor type
	if descriptorType != DescriptorTypeBech32 && descriptorType != DescriptorTypeBech32m {
		return nil, fmt.Errorf("invalid descriptor type: %s (must be 'bech32' or 'bech32m')", descriptorType)
	}

	bType, err := json.Marshal(string(descriptorType))
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(type): %w", err)
	}

	rawResult, err := b.Client.RawRequest("createwalletdescriptor", []json.RawMessage{bType})
	if err != nil {
		return nil, fmt.Errorf("fail to call json.RawRequest(createwalletdescriptor): %w", err)
	}

	var result CreateWalletDescriptorResult
	err = json.Unmarshal(rawResult, &result)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(rawResult): %w", err)
	}

	return &result, nil
}
