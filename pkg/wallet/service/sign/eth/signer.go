package eth

import (
	"errors"
	"fmt"

	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum/eth"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum/ethtx"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/storage/file"
	"github.com/hiromaily/go-crypto-wallet/pkg/serial"
)

// Sign type
type Sign struct {
	eth        ethereum.Ethereumer
	txFileRepo file.TransactionFileRepositorier
	wtype      domainWallet.WalletType
}

// NewSign returns sign object
func NewSign(
	ethAPI ethereum.Ethereumer,
	txFileRepo file.TransactionFileRepositorier,
	wtype domainWallet.WalletType,
) *Sign {
	return &Sign{
		eth:        ethAPI,
		txFileRepo: txFileRepo,
		wtype:      wtype,
	}
}

// SignTx sign on tx in csv file
// - multisig equivalent functionality is not implemented yet in ETH
func (s *Sign) SignTx(filePath string) (string, bool, string, error) {
	// get tx_deposit_id from tx file name
	actionType, _, txID, signedCount, err := s.txFileRepo.ValidateFilePath(filePath, domainTx.TxTypeUnsigned)
	if err != nil {
		return "", false, "", err
	}

	// get hex tx from file
	data, err := s.txFileRepo.ReadFileSlice(filePath)
	if err != nil {
		return "", false, "", fmt.Errorf("fail to call txFileRepo.ReadFileSlice(): %w", err)
	}
	if len(data) <= 1 {
		return "", false, "", errors.New("file is invalid")
	}
	serializedTxs := data[1:]

	txHexs := make([]string, 0, len(serializedTxs))
	for _, serializedTx := range serializedTxs {
		var rawTx ethtx.RawTx
		if err = serial.DecodeFromString(serializedTx, &rawTx); err != nil {
			return "", false, "", fmt.Errorf("fail to call serial.DecodeFromString(): %w", err)
		}
		// sign
		var signedRawTx *ethtx.RawTx
		signedRawTx, err = s.eth.SignOnRawTransaction(&rawTx, eth.Password)
		if err != nil {
			return "", false, "", fmt.Errorf("fail to call eth.SignOnRawTransaction(): %w", err)
		}
		txHexs = append(txHexs, fmt.Sprintf("%s,%s", rawTx.UUID, signedRawTx.TxHex))
	}

	// write file
	path := s.txFileRepo.CreateFilePath(actionType, domainTx.TxTypeSigned, txID, signedCount+1)
	generatedFileName, err := s.txFileRepo.WriteFileSlice(path, txHexs)
	if err != nil {
		return "", false, "", fmt.Errorf("fail to call txFileRepo.WriteFileSlice(): %w", err)
	}

	// return hexTx, isSigned, generatedFileName, nil
	return "", true, generatedFileName, nil
}
