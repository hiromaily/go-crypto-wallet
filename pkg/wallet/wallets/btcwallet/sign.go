package btcwallet

import (
	"database/sql"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/key"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/btc/coldsrv/signsrv"
)

// BTCSign is sign wallet object
type BTCSign struct {
	BTC         btcgrp.Bitcoiner
	dbConn      *sql.DB
	authAccount account.AuthType
	addrType    address.AddrType
	wtype       wallet.WalletType
	service.Seeder
	service.HDWalleter
	signsrv.PrivKeyer
	service.FullPubkeyExporter
	service.Signer
}

// NewBTCSign returns Sign object
func NewBTCSign(
	btc btcgrp.Bitcoiner,
	dbConn *sql.DB,
	authAccount account.AuthType,
	addrType address.AddrType,
	seeder service.Seeder,
	hdWallter service.HDWalleter,
	privKeyer signsrv.PrivKeyer,
	fullPubkeyExporter service.FullPubkeyExporter,
	signer service.Signer,
	wtype wallet.WalletType,
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
func (s *BTCSign) GenerateAuthKey(seed []byte, count uint32) ([]key.WalletKey, error) {
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
	// nolint:errcheck
	s.dbConn.Close()
	s.BTC.Close()
}
