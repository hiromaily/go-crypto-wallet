package service

//Cold wallet2

import (
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
)

// AddMultisigAddressByAuthorization account_key_authorizationテーブルのwallet_addressを認証者として、
// added_pubkey_history_paymentテーブル内のwalletアドレスのmultisigアドレスを生成する
func (w *Wallet) AddMultisigAddressByAuthorization(accountType enum.AccountType) error {
	//accountチェック
	if accountType != enum.AccountTypeReceipt && accountType != enum.AccountTypePayment {
		logger.Info("AccountType should be AccountTypeReceipt or AccountTypePayment")
		return nil
	}

	//account_key_authorizationテーブルからAuthorizationのwallet_addressを取得

	//added_pubkey_history_xxxテーブルからwallet_addressを取得

	//addmultisigaddress APIをcall

	//レスポンスをadded_pubkey_history_xxxテーブルに保存
	//resAddr, err := w.BTC.CreateMultiSig(2, []string{"2N7ZwUXpo841GZDpxLGFqrhr1xwMzTba7ZP", "2NAm558FWpiaJQLz838vbzBPpqmKxyeyxsu"}, "multi01")

	return nil
}
