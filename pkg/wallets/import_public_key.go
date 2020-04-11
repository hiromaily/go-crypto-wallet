package wallets

import (
	"strings"
	"time"

	"github.com/bookerzzz/grok"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/key"
)

// ImportPublicKeyForColdWallet2 csvファイルからpublicアドレスをimportする for Cold Wallet2
func (w *Wallet) ImportPublicKeyForColdWallet2(fileName string, accountType account.AccountType) error {
	if w.wtype != WalletTypeSignature {
		return errors.New("it's available on Coldwallet2")
	}

	//accountチェック
	//multisigであればこのチェックはOK
	//if accountType != enum.AccountTypeReceipt && accountType != enum.AccountTypePayment {
	//	logger.Info("AccountType should be AccountTypeReceipt or AccountTypePayment")
	//	return nil
	//}
	if !account.AccountTypeMultisig[accountType] {
		w.logger.Info("This func is for only account witch uses multiaddress")
		return nil
	}

	//TODO:ImportするファイルのaccountTypeもチェックしたほうがBetter
	//e.g. ./data/pubkey/receipt
	tmp := strings.Split(strings.Split(fileName, "_")[0], "/")
	if tmp[len(tmp)-1] != string(accountType) {
		return errors.Errorf("mismatching between accountType(%s) and file prefix [%s]", accountType, tmp[0])
	}

	//ファイル読み込み(full public key)
	pubKeys, err := key.ImportPubKey(fileName)
	if err != nil {
		return errors.Errorf("key.ImportPubKey() error: %s", err)
	}

	//added_pubkey_history_receiptテーブルにInsert
	addedPubkeyHistorys := make([]AddedPubkeyHistoryTable, len(pubKeys))
	for i, key := range pubKeys {
		inner := strings.Split(key, ",")

		//ここでは、FullPublicKeyをセットする必要がある
		addedPubkeyHistorys[i] = AddedPubkeyHistoryTable{
			FullPublicKey:         inner[2],
			AuthAddress1:          "",
			AuthAddress2:          "",
			WalletMultisigAddress: "",
			RedeemScript:          "",
		}
	}
	//TODO:Upsertに変えたほうがいいか？Insert済の場合、エラーが出る
	err = w.storager.InsertAddedPubkeyHistoryTable(accountType, addedPubkeyHistorys, nil, true)
	if err != nil {
		return errors.Errorf("DB.InsertAddedPubkeyHistoryTable() error: %s", err)
	}

	return nil
}

// ImportMultisigAddrForColdWallet1 coldwallet2でexportされたmultisigアドレス情報をimportする for Cold Wallet1
func (w *Wallet) ImportMultisigAddrForColdWallet1(fileName string, accountType account.AccountType) error {
	if w.wtype != WalletTypeKeyGen {
		return errors.New("it's available on Coldwallet1")
	}

	//accountチェック
	//multisigであればこのチェックはOK
	//if accountType != enum.AccountTypeReceipt && accountType != enum.AccountTypePayment {
	//	logger.Info("AccountType should be AccountTypeReceipt or AccountTypePayment")
	//	return nil
	//}
	if !account.AccountTypeMultisig[accountType] {
		w.logger.Info("This func is for only account witch uses multiaddress")
		return nil
	}

	//TODO:ImportするファイルのaccountTypeもチェックしたほうがBetter
	//e.g. ./data/pubkey/receipt
	tmp := strings.Split(strings.Split(fileName, "_")[0], "/")
	if tmp[len(tmp)-1] != string(accountType) {
		return errors.Errorf("mismatching between accountType(%s) and file prefix [%s]", accountType, tmp[0])
	}

	//ファイル読み込み(full public key)
	pubKeys, err := key.ImportPubKey(fileName)
	if err != nil {
		return errors.Errorf("key.ImportPubKey() error: %s", err)
	}

	//added_pubkey_history_receiptテーブルにInsert
	accountKeyTable := make([]AccountKeyTable, len(pubKeys))

	tm := time.Now()
	for i, key := range pubKeys {
		//TODO:とりあえず、1カラムのデータを前提でコーディングしておく
		inner := strings.Split(key, ",")
		//csvファイル
		//tmpData := []string{
		//	record.FullPublicKey,
		//	record.AuthAddress1,
		//	record.AuthAddress2,
		//	record.WalletMultisigAddress,
		//	record.RedeemScript,
		//}

		//Upsertをかけるには情報が不足しているので、とりあえず1行ずつUpdateする
		accountKeyTable[i] = AccountKeyTable{
			FullPublicKey:         inner[0],
			WalletMultisigAddress: inner[3],
			RedeemScript:          inner[4],
			KeyStatus:             enum.KeyStatusValue[enum.KeyStatusMultiAddressImported],
			UpdatedAt:             &tm,
		}
	}
	//Update
	err = w.storager.UpdateMultisigAddrOnAccountKeyTableByFullPubKey(accountType, accountKeyTable, nil, true)
	if err != nil {
		return errors.Errorf("DB.UpdateMultisigAddrOnAccountKeyTableByFullPubKey() error: %s", err)
	}

	return nil
}

