package coldwallet

import (
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb/coldrepo"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/types"
)

// ImportPubKey import pubKey from csv file for sign wallet
//  only multisig account is available
func (w *ColdWallet) ImportPubKey(fileName string, accountType account.AccountType) error {
	if w.wtype != types.WalletTypeSignature {
		return errors.New("it's available on sign wallet")
	}

	//validate account, only multisig account is ok
	if !account.AccountTypeMultisig[accountType] {
		w.logger.Info("multisig address can be imported, but this account is not")
		return nil
	}

	//validate file name
	if err := w.addrFileStorager.ValidateFilePath(fileName, accountType); err != nil {
		return err
	}

	// read file for full public key
	pubKeys, err := w.addrFileStorager.ImportPubKey(fileName)
	if err != nil {
		return errors.Wrapf(err, "fail to call fileStorager.ImportPubKey() fileName: %s", fileName)
	}

	// insert full pubKey into added_pubkey_history_table
	addedPubkeyHistorys := make([]coldrepo.AddedPubkeyHistoryTable, len(pubKeys))
	for i, key := range pubKeys {
		inner := strings.Split(key, ",")
		//FullPublicKey is required
		addedPubkeyHistorys[i] = coldrepo.AddedPubkeyHistoryTable{
			FullPublicKey:         inner[2],
			AuthAddress1:          "",
			AuthAddress2:          "",
			WalletMultisigAddress: "",
			RedeemScript:          "",
		}
	}
	//TODO:Upsert would be better to prevent error which occur when data is already inserted
	err = w.storager.InsertAddedPubkeyHistoryTable(accountType, addedPubkeyHistorys, nil, true)
	if err != nil {
		return errors.Wrap(err, "fail to call storager.InsertAddedPubkeyHistoryTable()")
	}

	return nil
}

// ImportMultisigAddress coldwallet2でexportされたmultisigアドレス情報をimportする for Cold Wallet1
func (w *ColdWallet) ImportMultisigAddress(fileName string, accountType account.AccountType) error {
	if w.wtype != types.WalletTypeKeyGen {
		return errors.New("it's available on Coldwallet1")
	}

	//accountチェック
	//multisigであればこのチェックはOK
	//if accountType != ctype.AccountTypeReceipt && accountType != ctype.AccountTypePayment {
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
	pubKeys, err := w.addrFileStorager.ImportPubKey(fileName)
	if err != nil {
		return errors.Errorf("key.ImportPubKey() error: %s", err)
	}

	//added_pubkey_history_receiptテーブルにInsert
	accountKeyTable := make([]coldrepo.AccountKeyTable, len(pubKeys))

	tm := time.Now()
	for i, pubkey := range pubKeys {
		//TODO:とりあえず、1カラムのデータを前提でコーディングしておく
		inner := strings.Split(pubkey, ",")
		//csvファイル
		//tmpData := []string{
		//	record.FullPublicKey,
		//	record.AuthAddress1,
		//	record.AuthAddress2,
		//	record.WalletMultisigAddress,
		//	record.RedeemScript,
		//}

		//Upsertをかけるには情報が不足しているので、とりあえず1行ずつUpdateする
		accountKeyTable[i] = coldrepo.AccountKeyTable{
			FullPublicKey:         inner[0],
			WalletMultisigAddress: inner[3],
			RedeemScript:          inner[4],
			AddressStatus:         address.AddressStatusValue[address.AddressStatusMultiAddressImported],
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
