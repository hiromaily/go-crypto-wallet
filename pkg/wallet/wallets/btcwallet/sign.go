package btcwallet

import (
	"database/sql"

	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainKey "github.com/hiromaily/go-crypto-wallet/pkg/domain/key"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/bitcoin"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service"
	btcsignsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/sign/btc"
)

// BTCSign is sign wallet object
type BTCSign struct {
	BTC         bitcoin.Bitcoiner
	dbConn      *sql.DB
	authAccount domainAccount.AuthType
	addrType    address.AddrType
	wtype       domainWallet.WalletType
	service.Seeder
	service.HDWalleter
	btcsignsrv.PrivKeyer
	service.FullPubkeyExporter
	service.Signer
}

// NewBTCSign returns Sign object
func NewBTCSign(
	btc bitcoin.Bitcoiner,
	dbConn *sql.DB,
	authAccount domainAccount.AuthType,
	addrType address.AddrType,
	seeder service.Seeder,
	hdWallter service.HDWalleter,
	privKeyer btcsignsrv.PrivKeyer,
	fullPubkeyExporter service.FullPubkeyExporter,
	signer service.Signer,
	wtype domainWallet.WalletType,
) *BTCSign {
	return &BTCSign{
		BTC:                btc,
		dbConn:             dbConn,
		authAccount:        authAccount,
		addrType:           addrType,
		wtype:              wtype,
		Seeder:             seeder,
		HDWalleter:         hdWallter,
		PrivKeyer:          privKeyer,
		FullPubkeyExporter: fullPubkeyExporter,
		Signer:             signer,
	}
}

// GenerateSeed generates seed
func (s *BTCSign) GenerateSeed() ([]byte, error) {
	return s.Seeder.Generate()
}

// StoreSeed stores seed
func (s *BTCSign) StoreSeed(strSeed string) ([]byte, error) {
	return s.Seeder.Store(strSeed)
}

// GenerateAuthKey generates account keys
func (s *BTCSign) GenerateAuthKey(seed []byte, count uint32) ([]domainKey.WalletKey, error) {
	return s.HDWalleter.Generate(s.authAccount.AccountType(), seed, count)
}

// ImportPrivKey imports privKey
func (s *BTCSign) ImportPrivKey() error {
	return s.PrivKeyer.Import()
}

// ExportFullPubkey exports full-pubkey
func (s *BTCSign) ExportFullPubkey() (string, error) {
	return s.FullPubkeyExporter.ExportFullPubkey()
}

// SignTx signs on transaction
func (s *BTCSign) SignTx(filePath string) (string, bool, string, error) {
	return s.Signer.SignTx(filePath)
}

// Done should be called before exit
func (s *BTCSign) Done() {
	_ = s.dbConn.Close() // Best effort cleanup
	s.BTC.Close()
}
