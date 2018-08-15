package service

import "github.com/hiromaily/go-bitcoin/pkg/file"

// coldwallet側で署名済みトランザクションを作成したものから、送金処理を行う

// SendFromFile 渡されたファイルから署名済transactionを読み取り、送信を行う
func (w *Wallet) SendFromFile(filePath string) (string, error) {
	//ファイルからhexを読み取る
	hex, err := file.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	//送信
	hash, err := w.Btc.SendTransactionByHex(hex)
	if err != nil {
		return "", err
	}

	//DB更新

	return hash.String(), nil
}
