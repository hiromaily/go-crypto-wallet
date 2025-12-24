package wallets

import (
	domainKey "github.com/hiromaily/go-crypto-wallet/pkg/domain/key"
)

// Signer is for signature wallet service interface
type Signer interface {
	GenerateSeed() ([]byte, error)
	StoreSeed(strSeed string) ([]byte, error)
	GenerateAuthKey(seed []byte, count uint32) ([]domainKey.WalletKey, error)
	ImportPrivKey() error
	ExportFullPubkey() (string, error)
	SignTx(filePath string) (string, bool, string, error)
	Done()
}
