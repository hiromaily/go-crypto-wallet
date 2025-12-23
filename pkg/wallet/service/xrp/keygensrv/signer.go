package keygensrv

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/coldrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp/xrp"
)

// Sign type
type Sign struct {
	xrp               xrpgrp.Rippler
	xrpAccountKeyRepo coldrepo.XRPAccountKeyRepositorier
	txFileRepo        tx.FileRepositorier
	wtype             wallet.WalletType
}

// NewSign returns sign object
func NewSign(
	xrpAPI xrpgrp.Rippler,
	xrpAccountKeyRepo coldrepo.XRPAccountKeyRepositorier,
	txFileRepo tx.FileRepositorier,
	wtype wallet.WalletType,
) *Sign {
	return &Sign{
		xrp:               xrpAPI,
		xrpAccountKeyRepo: xrpAccountKeyRepo,
		txFileRepo:        txFileRepo,
		wtype:             wtype,
	}
}

// SignTx sign on tx in csv file
// - multisig equivalent functionality is not implemented yet in ETH
func (s *Sign) SignTx(filePath string) (string, bool, string, error) {
	// get tx_deposit_id from tx file name
	actionType, _, txID, signedCount, err := s.txFileRepo.ValidateFilePath(filePath, tx.TxTypeUnsigned)
	if err != nil {
		return "", false, "", err
	}

	var senderAccount account.AccountType

	// get hex tx from file
	data, err := s.txFileRepo.ReadFileSlice(filePath)
	if err != nil {
		return "", false, "", fmt.Errorf("fail to call txFileRepo.ReadFileSlice(): %w", err)
	}
	if len(data) > 1 {
		senderAccount = account.AccountType(data[0])
	} else {
		return "", false, "", errors.New("file is invalid")
	}
	serializedTxs := data[1:]

	txHexs := make([]string, 0, len(serializedTxs))
	for _, serializedTx := range serializedTxs {
		// uid, txJSON
		tmp := strings.SplitAfterN(serializedTx, ",", 2)
		if len(tmp) != 2 {
			return "", false, "", errors.New("data format is invalid in file")
		}
		uuid := strings.TrimRight(tmp[0], ",")
		txJSON := tmp[1]

		var txInput xrp.TxInput
		if err = json.Unmarshal([]byte(txJSON), &txInput); err != nil {
			return "", false, "", fmt.Errorf("fail to call json.Unmarshal(txJSON): %w", err)
		}
		// TODO: get secret from database by txInput.Account
		// master_seed from xrp_account_key table
		var secret string
		secret, err = s.xrpAccountKeyRepo.GetSecret(senderAccount, txInput.Account)
		if err != nil {
			return "", false, "", fmt.Errorf("fail to call xrpAccountKeyRepo.GetSecret(): %w", err)
		}

		// sign
		var signedTxID string
		var txBlob string
		signedTxID, txBlob, err = s.xrp.SignTransaction(&txInput, secret)
		if err != nil {
			return "", false, "", fmt.Errorf("fail to call xrp.SignTransaction(): %w", err)
		}
		logger.Debug("signed_tx",
			"uuid", uuid, "signed_tx_id", signedTxID, "signed_tx_blob", txBlob)
		txHexs = append(txHexs, fmt.Sprintf("%s,%s,%s", uuid, signedTxID, txBlob))
	}

	// write file
	path := s.txFileRepo.CreateFilePath(actionType, tx.TxTypeSigned, txID, signedCount+1)
	generatedFileName, err := s.txFileRepo.WriteFileSlice(path, txHexs)
	if err != nil {
		return "", false, "", fmt.Errorf("fail to call txFileRepo.WriteFileSlice(): %w", err)
	}

	// return hexTx, isSigned, generatedFileName, nil
	return "", true, generatedFileName, nil
}
