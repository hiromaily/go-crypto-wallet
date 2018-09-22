package btc

import (
	"encoding/json"

	"github.com/pkg/errors"
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

//backupwallet

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
