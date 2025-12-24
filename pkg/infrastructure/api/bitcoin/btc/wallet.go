package btc

import (
	"encoding/json"
	"errors"
	"fmt"
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

// LoadWalletResult is response type of PRC `loadwallet`/`unloadwallet`
type LoadWalletResult struct {
	WalletName string `json:"name"`
	Warning    string `json:"warning"`
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

// LoadWallet import wallet dat
//
//	Loads a wallet from a wallet file or directory.
//	Note that all wallet command-line options used when starting bitcoind will be
//	applied to the new wallet (eg -zapwallettxes, upgradewallet, rescan, etc).
//	e.g. bitcoin-cli loadwallet "test.dat"
func (b *Bitcoin) LoadWallet(fileName string) error {
	if b.Version() < BTCVer17 {
		return errors.New("`loadwallet` is available from bitcoin version 0.17")
	}
	// loadwallet "filename"
	bFileName, err := json.Marshal(fileName)
	if err != nil {
		return fmt.Errorf("fail to call json.Marchal(fileName): %w", err)
	}
	rawResult, err := b.Client.RawRequest("loadwallet", []json.RawMessage{bFileName})
	if err != nil {
		return fmt.Errorf("fail to call json.RawRequest(loadwallet): %w", err)
	}

	loadWalletResult := LoadWalletResult{}
	err = json.Unmarshal(rawResult, &loadWalletResult)
	if err != nil {
		return fmt.Errorf("fail to call json.Unmarshal(rawResult): %w", err)
	}
	if loadWalletResult.Warning != "" {
		// TODO: how to handle this warning
		return fmt.Errorf("detect warning: %s", loadWalletResult.Warning)
	}

	return nil
}

// UnLoadWallet unload wallet dat
//
//	Unloads the wallet referenced by the request endpoint otherwise unloads the wallet specified in the argument.
//	Specifying the wallet name on a wallet endpoint is invalid.
//	e.g. bitcoin-cli unloadwallet wallet_name
func (b *Bitcoin) UnLoadWallet(fileName string) error {
	if b.Version() < BTCVer17 {
		return errors.New("`unloadwallet` is available from bitcoin version 0.17")
	}
	// unloadwallet ( "wallet_name" )
	bFileName, err := json.Marshal(fileName)
	if err != nil {
		return fmt.Errorf("fail to call json.Marchal(): %w", err)
	}
	_, err = b.Client.RawRequest("unloadwallet", []json.RawMessage{bFileName})
	if err != nil {
		return errors.New("fail to call json.RawRequest(unloadwallet)")
	}

	return nil
}

// CreateWallet Creates and loads a new wallet
func (b *Bitcoin) CreateWallet(fileName string, disablePrivKey bool) error {
	if b.Version() < BTCVer17 {
		return errors.New("`createwallet` is available from bitcoin version 0.17")
	}
	// createwallet "wallet_name" ( disable_private_keys )
	bFileName, err := json.Marshal(fileName)
	if err != nil {
		return fmt.Errorf("fail to call json.Marchal(fileName): %w", err)
	}

	bDisablePrivKey, err := json.Marshal(disablePrivKey)
	if err != nil {
		return fmt.Errorf("fail to call json.Marchal(bool): %w", err)
	}

	rawResult, err := b.Client.RawRequest("createwallet", []json.RawMessage{bFileName, bDisablePrivKey})
	if err != nil {
		return fmt.Errorf("fail to call json.RawRequest(createwallet): %w", err)
	}

	loadWalletResult := LoadWalletResult{}
	err = json.Unmarshal(rawResult, &loadWalletResult)
	if err != nil {
		return fmt.Errorf("fail to call json.Unmarshal(rawResult): %w", err)
	}
	if loadWalletResult.Warning != "" {
		// TODO: how to handle this warning
		return fmt.Errorf("detect warning: %s", loadWalletResult.Warning)
	}

	return nil
}
