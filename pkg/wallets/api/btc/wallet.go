package btc

import (
	"encoding/json"

	"github.com/pkg/errors"

	ctype "github.com/hiromaily/go-bitcoin/pkg/wallets/api/types"
)

//{
//"filename": "/Applications/sample01.dat"
//}

const (
	dumpWallet   = "dumpwallet"
	importWallet = "importwallet"
)

// WalletResult dumpwallet/dumpwallet RPC のレスポンス
type WalletResult struct {
	FileName string `json:"filename"`
}

// LoadWalletResult loadwallet/createwallet　RPCのレスポンス
type LoadWalletResult struct {
	WalletName string `json:"name"`
	Warning    string `json:"warning"`
}

// BackupWallet walletファイルをunloadする
//  Safely copies current wallet file to destination, which can be a directory or a path with filename
func (b *Bitcoin) BackupWallet(fileName string) error {
	//backupwallet
	bFileName, err := json.Marshal(fileName)
	if err != nil {
		return errors.Errorf("json.Marchal(): error: %s", err)
	}
	_, err = b.client.RawRequest("backupwallet", []json.RawMessage{bFileName})
	if err != nil {
		return errors.Errorf("json.RawRequest(backupwallet): error: %s", err)
	}

	return nil
}

// DumpWallet walletを指定したファイル名でdumpする (これをすることでkeyの存在しないwalletが作られるはず)
//  fileNameは絶対パス
//  This does not allow overwriting existing files.
func (b *Bitcoin) DumpWallet(fileName string) error {
	return b.dumpImportWallet(fileName, dumpWallet)
}

// ImportWallet walletファイルをimportする
func (b *Bitcoin) ImportWallet(fileName string) error {
	return b.dumpImportWallet(fileName, importWallet)
}

func (b *Bitcoin) dumpImportWallet(fileName, method string) error {
	bFileName, err := json.Marshal(string(fileName))
	if err != nil {
		return errors.Errorf("json.Marchal(): error: %s", err)
	}

	rawResult, err := b.client.RawRequest(method, []json.RawMessage{bFileName})
	if err != nil {
		return errors.Errorf("json.RawRequest(%s): error: %s", method, err)
	}

	var walletResult WalletResult
	err = json.Unmarshal([]byte(rawResult), &walletResult)
	if err != nil {
		return errors.Errorf("json.Unmarshal(): error: %s", err)
	}

	return nil
}

// EncryptWallet Walletの初回パスフレーズを設定
// 0.16.2の場合、`encryptwallet`というmethod名になっている。これは古いのかも
func (b *Bitcoin) EncryptWallet(passphrase string) error {
	return b.client.CreateEncryptedWallet(passphrase)
}

// WalletLock walletにロックをかける(walletpassphraseによってロックを解除するまで利用できない)
func (b *Bitcoin) WalletLock() error {
	return b.client.WalletLock()
}

// WalletPassphrase ロックを解除する
func (b *Bitcoin) WalletPassphrase(passphrase string, timeoutSecs int64) error {
	return b.client.WalletPassphrase(passphrase, timeoutSecs)
}

// WalletPassphraseChange パスフレーズを変更する
func (b *Bitcoin) WalletPassphraseChange(old, new string) error {
	return b.client.WalletPassphraseChange(old, new)
}

// LoadWallet walletファイルをimportする
//  Loads a wallet from a wallet file or directory.
//  Note that all wallet command-line options used when starting bitcoind will be
//  applied to the new wallet (eg -zapwallettxes, upgradewallet, rescan, etc).
//  e.g. bitcoin-cli loadwallet "test.dat"
func (b *Bitcoin) LoadWallet(fileName string) error {
	if b.Version() < ctype.BTCVer17 {
		return errors.New("`loadwallet` is available from bitcoin version 0.17")
	}
	//loadwallet "filename"
	bFileName, err := json.Marshal(fileName)
	if err != nil {
		return errors.Errorf("json.Marchal(): error: %s", err)
	}
	rawResult, err := b.client.RawRequest("loadwallet", []json.RawMessage{bFileName})
	if err != nil {
		return errors.Errorf("json.RawRequest(loadwallet): error: %s", err)
	}

	loadWalletResult := LoadWalletResult{}
	err = json.Unmarshal([]byte(rawResult), &loadWalletResult)
	if err != nil {
		return errors.Errorf("json.Unmarshal(): error: %s", err)
	}
	if loadWalletResult.Warning != "" {
		//TODO:この扱いをどうするか
		return errors.Errorf("json.RawRequest(loadwallet): warning: %s", loadWalletResult.Warning)
	}

	return nil
}

// UnLoadWallet walletファイルをunloadする
//  Unloads the wallet referenced by the request endpoint otherwise unloads the wallet specified in the argument.
//  Specifying the wallet name on a wallet endpoint is invalid.
//  e.g. bitcoin-cli unloadwallet wallet_name
func (b *Bitcoin) UnLoadWallet(fileName string) error {
	if b.Version() < ctype.BTCVer17 {
		return errors.New("`unloadwallet` is available from bitcoin version 0.17")
	}
	//unloadwallet ( "wallet_name" )
	bFileName, err := json.Marshal(fileName)
	if err != nil {
		return errors.Errorf("json.Marchal(): error: %s", err)
	}
	_, err = b.client.RawRequest("unloadwallet", []json.RawMessage{bFileName})
	if err != nil {
		return errors.Errorf("json.RawRequest(unloadwallet): error: %s", err)
	}

	return nil
}

// CreateWallet Creates and loads a new wallet
func (b *Bitcoin) CreateWallet(fileName string, disablePrivKey bool) error {
	if b.Version() < ctype.BTCVer17 {
		return errors.New("`createwallet` is available from bitcoin version 0.17")
	}
	//createwallet "wallet_name" ( disable_private_keys )
	bFileName, err := json.Marshal(fileName)
	if err != nil {
		return errors.Errorf("json.Marchal(): error: %s", err)
	}

	bDisablePrivKey, err := json.Marshal(disablePrivKey)
	if err != nil {
		return errors.Errorf("json.Marchal(): error: %s", err)
	}

	rawResult, err := b.client.RawRequest("createwallet", []json.RawMessage{bFileName, bDisablePrivKey})
	if err != nil {
		return errors.Errorf("json.RawRequest(createwallet): error: %s", err)
	}

	loadWalletResult := LoadWalletResult{}
	err = json.Unmarshal([]byte(rawResult), &loadWalletResult)
	if err != nil {
		return errors.Errorf("json.Unmarshal(): error: %s", err)
	}
	if loadWalletResult.Warning != "" {
		//TODO:この扱いをどうするか
		return errors.Errorf("json.RawRequest(loadwallet): warning: %s", loadWalletResult.Warning)
	}

	return nil
}
