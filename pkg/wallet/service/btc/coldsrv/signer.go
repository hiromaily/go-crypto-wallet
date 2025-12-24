package coldsrv

import (
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/wire"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/bitcoin"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/bitcoin/btc"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/repository/cold"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/storage/file"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/serial"
)

// Sign type
type Sign struct {
	btc             bitcoin.Bitcoiner
	accountKeyRepo  cold.AccountKeyRepositorier
	authKeyRepo     cold.AuthAccountKeyRepositorier
	txFileRepo      file.TransactionFileRepositorier
	multisigAccount account.MultisigAccounter
	wtype           domainWallet.WalletType
}

// NewSign returns sign object
func NewSign(
	btcAPI bitcoin.Bitcoiner,
	accountKeyRepo cold.AccountKeyRepositorier,
	authKeyRepo cold.AuthAccountKeyRepositorier,
	txFileRepo file.TransactionFileRepositorier,
	multisigAccount account.MultisigAccounter,
	wtype domainWallet.WalletType,
) *Sign {
	return &Sign{
		btc:             btcAPI,
		accountKeyRepo:  accountKeyRepo,
		authKeyRepo:     authKeyRepo,
		txFileRepo:      txFileRepo,
		multisigAccount: multisigAccount,
		wtype:           wtype,
	}
}

// SignTx sign on tx in csv file
// - logic would vary among account, addressType like multisig
// - returns tx, isSigned, generatedFileName, error
func (s *Sign) SignTx(filePath string) (string, bool, string, error) {
	// get tx_deposit_id from tx file name
	//  if payment_5_unsigned_0_1534466246366489473, 5 is target
	actionType, _, txID, signedCount, err := s.txFileRepo.ValidateFilePath(filePath, domainTx.TxTypeUnsigned)
	if err != nil {
		return "", false, "", err
	}

	// get hex tx from file
	data, err := s.txFileRepo.ReadFile(filePath)
	if err != nil {
		return "", false, "", err
	}

	var hex, encodedPrevsAddrs string
	tmp := strings.Split(data, ",")
	// file: hex, prev_address
	hex = tmp[0]
	if len(tmp) > 1 {
		encodedPrevsAddrs = tmp[1]
	}
	if encodedPrevsAddrs == "" {
		// it's required data since Bitcoin core ver17
		return "", false, "", errors.New("encodedPrevsAddrs must be set in csv file")
	}

	// sing
	hexTx, isSigned, newEncodedPrevsAddrs, err := s.sign(hex, encodedPrevsAddrs)
	if err != nil {
		return "", isSigned, "", err
	}

	// hexTx for save data as file
	saveData := hexTx

	// if sign is not finished because of multisig, signedCount should be increment
	txType := domainTx.TxTypeSigned
	if !isSigned {
		txType = domainTx.TxTypeUnsigned
		signedCount++
		if newEncodedPrevsAddrs != "" {
			saveData = fmt.Sprintf("%s,%s", saveData, newEncodedPrevsAddrs)
		}
	}

	// write file
	path := s.txFileRepo.CreateFilePath(actionType, txType, txID, signedCount)
	generatedFileName, err := s.txFileRepo.WriteFile(path, saveData)
	if err != nil {
		return "", isSigned, "", err
	}

	return hexTx, isSigned, generatedFileName, nil
}

// sign
// - coin is sent from sender account to receiver account. Sender's privKey(sender account) is required
// - [actionType:deposit]  [from] client [to] deposit, (not multisig addr)
// - [actionType:payment]  [from] payment [to] unknown, (multisig addr)
// - [actionType:transfer] [from] from [to] to, (multisig addr)
func (s *Sign) sign(hex, encodedPrevsAddrs string) (string, bool, string, error) {
	// get tx from hex
	msgTx, err := s.btc.ToMsgTx(hex)
	if err != nil {
		return "", false, "", err
	}

	// decode encodedPrevsAddrs string to btc.AddrsPrevTxs struct
	var prevsAddrs btc.PreviousTxs
	if err = serial.DecodeFromString(encodedPrevsAddrs, &prevsAddrs); err != nil {
		return "", false, "", err
	}

	// single signature address
	var (
		signedTx             *wire.MsgTx
		isSigned             bool
		newEncodedPrevsAddrs string
	)

	// sign
	if !s.multisigAccount.IsMultisigAccount(prevsAddrs.SenderAccount) {
		signedTx, isSigned, err = s.btc.SignRawTransaction(msgTx, prevsAddrs.PrevTxs)
	} else {
		signedTx, isSigned, newEncodedPrevsAddrs, err = s.signMultisig(msgTx, &prevsAddrs)
	}

	if newEncodedPrevsAddrs == "" {
		newEncodedPrevsAddrs = encodedPrevsAddrs
	}

	// after sign
	if err != nil {
		return "", false, "", err
	}

	hexTx, err := s.btc.ToHex(signedTx)
	if err != nil {
		return "", false, "", fmt.Errorf("fail to call s.btc.ToHex(signedTx): %w", err)
	}
	logger.Debug(
		"call btc.SignRawTransaction()",
		"hexTx", hexTx,
		"isSigned", isSigned)

	return hexTx, isSigned, newEncodedPrevsAddrs, nil
}

func (s *Sign) signMultisig(msgTx *wire.MsgTx, prevsAddrs *btc.PreviousTxs) (*wire.MsgTx, bool, string, error) {
	var wips []string
	var newEncodedPrevsAddrs string

	// get WIPs, RedeedScript
	switch s.wtype {
	case domainWallet.WalletTypeKeyGen:
		accountKeys, err := s.accountKeyRepo.GetAllMultiAddr(prevsAddrs.SenderAccount, prevsAddrs.Addrs)
		if err != nil {
			return nil, false, "", fmt.Errorf("fail to call accountKeyRepo.GetAllMultiAddr(): %w", err)
		}

		// retrieve WIPs
		for _, val := range accountKeys {
			wips = append(wips, val.WalletImportFormat)
		}

		// mapping redeemScript to PrevTxs
		for idx, val := range prevsAddrs.Addrs {
			rs := cold.GetRedeemScriptByAddress(accountKeys, val)
			if rs == "" {
				logger.Error("redeemScript can not be found")
				continue
			}
			prevsAddrs.PrevTxs[idx].RedeemScript = rs
		}

		// serialize prevsAddrs with redeemScript
		newEncodedPrevsAddrs, err = serial.EncodeToString(prevsAddrs)
		if err != nil {
			return nil, false, "", fmt.Errorf("fail to call serial.EncodeToString(): %w", err)
		}

	case domainWallet.WalletTypeSign:
		authKey, err := s.authKeyRepo.GetOne("")
		if err != nil {
			return nil, false, "", fmt.Errorf("fail to call authKeyRepo.GetOne(): %w", err)
		}
		// wip
		wips = []string{authKey.WalletImportFormat}
	case domainWallet.WalletTypeWatchOnly:
		return nil, false, "", fmt.Errorf("WalletType is invalid: %s", s.wtype.String())
	default:
		return nil, false, "", fmt.Errorf("WalletType is invalid: %s", s.wtype.String())
	}

	// sign
	signedTx, isSigned, err := s.btc.SignRawTransactionWithKey(msgTx, wips, prevsAddrs.PrevTxs)
	return signedTx, isSigned, newEncodedPrevsAddrs, err
}
