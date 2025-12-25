package ethwallet

import (
	"database/sql"

	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/ethereum"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// ETHSign keygen wallet object
type ETHSign struct {
	ETH    ethereum.Ethereumer
	dbConn *sql.DB
	wtype  domainWallet.WalletType
}

// NewETHSign returns ETHSign object
func NewETHSign(
	eth ethereum.Ethereumer,
	dbConn *sql.DB,
	walletType domainWallet.WalletType,
) *ETHSign {
	return &ETHSign{
		ETH:    eth,
		dbConn: dbConn,
		wtype:  walletType,
	}
}

// GenerateSeed generates seed
func (*ETHSign) GenerateSeed() ([]byte, error) {
	logger.Info("no functionality for CreateMultisigAddress() in ETH")
	return nil, nil
}

// StoreSeed stores seed
func (*ETHSign) StoreSeed(_ string) ([]byte, error) {
	logger.Info("no functionality for CreateMultisigAddress() in ETH")
	return nil, nil
}

// GenerateAuthKey generates account keys
func (*ETHSign) GenerateAuthKey(_ []byte, _ uint32) ([]domainKey.WalletKey, error) {
	logger.Info("no functionality for CreateMultisigAddress() in ETH")
	return nil, nil
}

// ImportPrivKey imports privKey
func (*ETHSign) ImportPrivKey() error {
	logger.Info("no functionality for CreateMultisigAddress() in ETH")
	return nil
}

// ExportFullPubkey exports full-pubkey
func (*ETHSign) ExportFullPubkey() (string, error) {
	logger.Info("no functionality for CreateMultisigAddress() in ETH")
	return "", nil
}

// SignTx signs on transaction
func (*ETHSign) SignTx(_ string) (string, bool, string, error) {
	logger.Info("no functionality for CreateMultisigAddress() in ETH")
	return "", false, "", nil
}

// Done should be called before exit
func (s *ETHSign) Done() {
	_ = s.dbConn.Close() // Best effort cleanup
	s.ETH.Close()
}
