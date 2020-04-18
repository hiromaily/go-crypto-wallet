package coldwallet

import (
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb/coldrepo"
	"github.com/hiromaily/go-bitcoin/pkg/serial"
	"github.com/hiromaily/go-bitcoin/pkg/tx"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/btc"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/types"
)

// sing on unsigned transaction

// SignTx sign on tx in csv file
// - logic would vary among account, addressType like multisig
// - returns tx, isSigned, generatedFileName, error
func (w *ColdWallet) SignTx(filePath string) (string, bool, string, error) {

	// get tx_receipt_id from tx file name
	//  if payment_5_unsigned_0_1534466246366489473, 5 is target
	actionType, _, txReceiptID, signedCount, err := w.txFileRepo.ValidateFilePath(filePath, tx.TxTypeUnsigned)
	if err != nil {
		return "", false, "", err
	}

	// get hex tx from file
	data, err := w.txFileRepo.ReadFile(filePath)
	if err != nil {
		return "", false, "", err
	}

	var hex, encodedPrevsAddrs string
	tmp := strings.Split(data, ",")
	//file: hex, prev_address
	hex = tmp[0]
	if len(tmp) > 1 {
		encodedPrevsAddrs = tmp[1]
	}
	if encodedPrevsAddrs == "" {
		//it's required data since Bitcoin core ver17
		return "", false, "", errors.New("encodedPrevsAddrs must be set in csv file")
	}

	// sing
	hexTx, isSigned, newEncodedPrevsAddrs, err := w.sign(hex, encodedPrevsAddrs)
	if err != nil {
		return "", isSigned, "", err
	}

	// hexTx for save data as file
	saveData := hexTx

	// if sign is not finished because of multisig, signedCount should be increment
	txType := tx.TxTypeSigned
	if !isSigned {
		txType = tx.TxTypeUnsigned
		signedCount++
		if newEncodedPrevsAddrs != "" {
			saveData = fmt.Sprintf("%s,%s", saveData, newEncodedPrevsAddrs)
		}
	}

	// write file
	path := w.txFileRepo.CreateFilePath(actionType, txType, txReceiptID, signedCount)
	generatedFileName, err := w.txFileRepo.WriteFile(path, saveData)
	if err != nil {
		return "", isSigned, "", err
	}

	return hexTx, isSigned, generatedFileName, nil
}

// sign
// - coin is sent [from] account to [to] account then sernder's privKey(from account) is required
// - [actionType:receipt]  [from] client [to] receipt, (not multisig addr)
// - [actionType:payment]  [from] payment [to] unknown, (multisig addr)
// - [actionType:transfer] [from] from [to] to, (multisig addr)
// TODO:transfer action is not implemented yet ??
func (w *ColdWallet) sign(hex, encodedPrevsAddrs string) (string, bool, string, error) {
	// get tx from hex
	msgTx, err := w.btc.ToMsgTx(hex)
	if err != nil {
		return "", false, "", err
	}

	var (
		signedTx             *wire.MsgTx
		isSigned             bool
		prevsAddrs           btc.AddrsPrevTxs
		accountKeys          []coldrepo.AccountKeyTable
		wips                 []string
		newEncodedPrevsAddrs string
	)

	// decode encodedPrevsAddrs to prevsAddrs
	serial.DecodeFromString(encodedPrevsAddrs, &prevsAddrs)

	// get WIPs, RedeedScript
	// - logic vary between keygen wallet and sign wallet (actually multisig or not)
	// - sign wallet requires AccountTypeAuthorization
	// - TODO: wips is not used if action is receipt because client account is not multisig address
	switch w.wtype {
	case types.WalletTypeKeyGen:
		//TODO: if ActionType==`transfer`, address for from account is required
		// => logic is changed. addrsPrevs.SenderAccount is used for getting sender information
		// address must be multisig address
		// get data from account_key_table
		accountKeys, err = w.storager.GetAllAccountKeyByMultiAddrs(prevsAddrs.SenderAccount, prevsAddrs.Addrs)
		if err != nil {
			return "", false, "", errors.Errorf("DB.GetWIPByMultiAddrs() error: %s", err)
		}
	case types.WalletTypeSignature:
		// sign wallet is used from 2nd signature, only multisig address
		// get data from account_key_authorization table
		// TODO: client account doesn't have multisig address, so this code could be skipped
		accountKey, err := w.storager.GetOneByMaxIDOnAccountKeyTable(account.AccountTypeAuthorization)
		if err != nil {
			return "", false, "", errors.Wrap(err, "fail to call storager.GetOneByMaxIDOnAccountKeyTable()")
		}
		accountKeys = append(accountKeys, *accountKey)
	default:
		return "", false, "", errors.Errorf("WalletType is invalid: %s", w.wtype.String())
	}

	// retrieve WIPs
	for _, val := range accountKeys {
		wips = append(wips, val.WalletImportFormat)
	}

	//if sender account is multisig account
	if account.AccountTypeMultisig[prevsAddrs.SenderAccount] {
		switch w.wtype {
		case types.WalletTypeKeyGen:
			// mapping redeemScript to PrevTxs
			for idx, val := range prevsAddrs.Addrs {
				rs := coldrepo.GetRedeedScriptByAddress(accountKeys, val)
				if rs == "" {
					w.logger.Error("redeemScript can not be found")
					continue
				}
				prevsAddrs.PrevTxs[idx].RedeemScript = rs
			}
			//grok.Value(prevsAddrs)

			// serialize prevsAddrs with redeemScript
			newEncodedPrevsAddrs, err = serial.EncodeToString(prevsAddrs)
			if err != nil {
				return "", false, "", errors.Errorf("serial.EncodeToString(): error: %s", err)
			}
		case types.WalletTypeSignature:
			newEncodedPrevsAddrs = encodedPrevsAddrs
		default:
			return "", false, "", errors.Errorf("WalletType is invalid: %s", w.wtype.String())
		}
	}

	//sign
	if account.AccountTypeMultisig[prevsAddrs.SenderAccount] {
		//wips is required
		signedTx, isSigned, err = w.btc.SignRawTransactionWithKey(msgTx, wips, prevsAddrs.PrevTxs)
	} else {
		signedTx, isSigned, err = w.btc.SignRawTransaction(msgTx, prevsAddrs.PrevTxs)
	}
	if err != nil {
		return "", false, "", err
	}
	hexTx, err := w.btc.ToHex(signedTx)
	if err != nil {
		return "", false, "", errors.Errorf("w.BTC.ToHex(msgTx): error: %s", err)
	}
	w.logger.Debug(
		"call btc.SignRawTransaction()",
		zap.String("hexTx", hexTx),
		zap.Bool("isSigned", isSigned))

	return hexTx, isSigned, newEncodedPrevsAddrs, nil
}
