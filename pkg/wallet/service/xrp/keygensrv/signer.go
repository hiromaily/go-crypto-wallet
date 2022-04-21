package keygensrv

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/coldrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp/xrp"
)

// Sign type
type Sign struct {
	xrp               xrpgrp.Rippler
	logger            *zap.Logger
	xrpAccountKeyRepo coldrepo.XRPAccountKeyRepositorier
	txFileRepo        tx.FileRepositorier
	wtype             wallet.WalletType
}

// NewSign returns sign object
func NewSign(
	xrp xrpgrp.Rippler,
	logger *zap.Logger,
	xrpAccountKeyRepo coldrepo.XRPAccountKeyRepositorier,
	txFileRepo tx.FileRepositorier,
	wtype wallet.WalletType,
) *Sign {
	return &Sign{
		xrp:               xrp,
		logger:            logger,
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
		return "", false, "", errors.Wrap(err, "fail to call txFileRepo.ReadFileSlice()")
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
			return "", false, "", errors.Wrap(err, "fail to call json.Unmarshal(txJSON)")
		}
		// TODO: get secret from database by txInput.Account
		// master_seed from xrp_account_key table
		secret, err := s.xrpAccountKeyRepo.GetSecret(senderAccount, txInput.Account)
		if err != nil {
			return "", false, "", errors.Wrap(err, "fail to call xrpAccountKeyRepo.GetSecret()")
		}

		// sign
		signedTxID, txBlob, err := s.xrp.SignTransaction(&txInput, secret)
		if err != nil {
			return "", false, "", errors.Wrap(err, "fail to call xrp.SignTransaction()")
		}
		s.logger.Debug("signed_tx", zap.String("uuid", uuid), zap.String("signed_tx_id", signedTxID), zap.String("signed_tx_blob", txBlob))
		txHexs = append(txHexs, fmt.Sprintf("%s,%s,%s", uuid, signedTxID, txBlob))
	}

	// write file
	path := s.txFileRepo.CreateFilePath(actionType, tx.TxTypeSigned, txID, signedCount+1)
	generatedFileName, err := s.txFileRepo.WriteFileSlice(path, txHexs)
	if err != nil {
		return "", false, "", errors.Wrap(err, "fail to call txFileRepo.WriteFileSlice()")
	}

	// return hexTx, isSigned, generatedFileName, nil
	return "", true, generatedFileName, nil
}
