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
func (w *Wallet) AddMultisigAddressByAuthorization(accountType enum.AccountType) error {
	//accountチェック
	if accountType != enum.AccountTypeReceipt && accountType != enum.AccountTypePayment {
		logger.Info("AccountType should be AccountTypeReceipt or AccountTypePayment")
		return nil
	}

	//account_key_authorizationテーブルからAuthorizationのwallet_addressを取得
	authKeyTable, err := w.DB.GetOneByMaxID(enum.AccountTypeAuthorization)
	if err != nil {
		return errors.Errorf("DB.GetOneByMaxID(enum.AccountTypeAuthorization) error: %s", err)
	}
	//grok.Value(authKeyTable)

	//added_pubkey_history_xxxテーブルからwallet_addressを取得
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
				val.WalletAddress,          // receipt or payment address
				authKeyTable.WalletAddress, // authorization address
			},
			fmt.Sprintf("multi_%s", accountType),
		)
		if err != nil {
			logger.Errorf("BTC.CreateMultiSig(2,%s,%s) error: %s", val.WalletAddress, authKeyTable.WalletAddress, err)
			continue
		}

		//レスポンスをadded_pubkey_history_xxxテーブルに保存
		//err := w.DB.UpdateAddedPubkeyHistoryTableByMultisigAddr(accountType, "aaa", "bbb", val.WalletAddress, nil, true)
		err = w.DB.UpdateAddedPubkeyHistoryTableByMultisigAddr(accountType, resAddr.Address, resAddr.RedeemScript, val.WalletAddress, nil, true)
		if err != nil {
			logger.Errorf("DB.UpdateAddedPubkeyHistoryTableByMultisigAddr(%s) error: %s", accountType, err)
		}
	}

	return nil
}
