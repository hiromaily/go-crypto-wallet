package service

//Cold wallet2

import (
	"fmt"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/pkg/errors"
)

// AddMultisigAddressByAuthorization account_key_authorizationテーブルのwallet_addressを認証者として、
// added_pubkey_history_paymentテーブル内のwalletアドレスのmultisigアドレスを生成する
// TODO:第4パラメータに、address_typeを追加する。Bitcoinの場合は、p2sh-segwit とする
func (w *Wallet) AddMultisigAddressByAuthorization(accountType enum.AccountType, addressType enum.AddressType) error {
	//accountチェック
	if accountType != enum.AccountTypeReceipt && accountType != enum.AccountTypePayment {
		logger.Info("AccountType should be AccountTypeReceipt or AccountTypePayment")
		return nil
	}

	//account_key_authorizationテーブルからAuthorizationのwallet_addressを取得
	authKeyTable, err := w.DB.GetOneByMaxIDOnAccountKeyTable(enum.AccountTypeAuthorization)
	if err != nil {
		return errors.Errorf("DB.GetOneByMaxIDOnAccountKeyTable(enum.AccountTypeAuthorization) error: %s", err)
	}
	//grok.Value(authKeyTable)

	//added_pubkey_history_xxxテーブルからwallet_address(full-pubkeyである必要がある)を取得
	addedPubkeyHistoryTable, err := w.DB.GetAddedPubkeyHistoryTableByNoWalletMultisigAddress(accountType)
	if err != nil {
		return errors.Errorf("DB.GetAddedPubkeyHistoryTableByNoWalletMultisigAddress(%s) error: %s", accountType, err)
	}
	//grok.Value(addedPubkeyHistoryTable)

	//addmultisigaddress APIをcall
	for _, val := range addedPubkeyHistoryTable {
		resAddr, err := w.BTC.CreateMultiSig(
			2,
			[]string{
				val.FullPublicKey, // receipt or payment address
				//authKeyTable.WalletAddress, // authorization address TODO:P2shSegwitAddressじゃなくていいのか？これが原因でmultisigの送信に失敗した？？
				authKeyTable.P2shSegwitAddress,
			},
			fmt.Sprintf("multi_%s", accountType), //TODO:ここのアカウント名はどうすべきか
			addressType,
		)
		if err != nil {
			//[Error] -5: no full public key for address mkPmdpo59gpU7ZioGYwwoMTQJjh7MiqUvd
			logger.Errorf("BTC.CreateMultiSig(2,%s,%s) error: %s", val.FullPublicKey, authKeyTable.WalletAddress, err)
			continue
		}

		//レスポンスをadded_pubkey_history_xxxテーブルに保存
		//err = w.DB.UpdateMultisigAddrOnAddedPubkeyHistoryTable(accountType, resAddr.Address,
		//	resAddr.RedeemScript, authKeyTable.WalletAddress, val.FullPublicKey, nil, true)
		err = w.DB.UpdateMultisigAddrOnAddedPubkeyHistoryTable(accountType, resAddr.Address,
			resAddr.RedeemScript, authKeyTable.P2shSegwitAddress, val.FullPublicKey, nil, true)
		if err != nil {
			logger.Errorf("DB.UpdateMultisigAddrOnAddedPubkeyHistoryTable(%s) error: %s", accountType, err)
		}
	}

	return nil
}
