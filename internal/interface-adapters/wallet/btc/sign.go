package btc

import (
	"context"
	"database/sql"

	signusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/address"
)

// BTCSign is sign wallet object
type BTCSign struct {
	BTC                     bitcoin.Bitcoiner
	dbConn                  *sql.DB
	authAccount             domainAccount.AuthType
	addrType                address.AddrType
	wtype                   domainWallet.WalletType
	generateSeedUseCase     signusecase.GenerateSeedUseCase
	storeSeedUseCase        signusecase.StoreSeedUseCase
	generateAuthKeyUseCase  signusecase.GenerateAuthKeyUseCase
	importPrivKeyUseCase    signusecase.ImportPrivateKeyUseCase
	exportFullPubkeyUseCase signusecase.ExportFullPubkeyUseCase
	signTxUseCase           signusecase.SignTransactionUseCase
}

// NewBTCSign returns Sign object
func NewBTCSign(
	btc bitcoin.Bitcoiner,
	dbConn *sql.DB,
	authAccount domainAccount.AuthType,
	addrType address.AddrType,
	generateSeedUseCase signusecase.GenerateSeedUseCase,
	storeSeedUseCase signusecase.StoreSeedUseCase,
	generateAuthKeyUseCase signusecase.GenerateAuthKeyUseCase,
	importPrivKeyUseCase signusecase.ImportPrivateKeyUseCase,
	exportFullPubkeyUseCase signusecase.ExportFullPubkeyUseCase,
	signTxUseCase signusecase.SignTransactionUseCase,
	wtype domainWallet.WalletType,
) *BTCSign {
	return &BTCSign{
		BTC:                     btc,
		dbConn:                  dbConn,
		authAccount:             authAccount,
		addrType:                addrType,
		wtype:                   wtype,
		generateSeedUseCase:     generateSeedUseCase,
		storeSeedUseCase:        storeSeedUseCase,
		generateAuthKeyUseCase:  generateAuthKeyUseCase,
		importPrivKeyUseCase:    importPrivKeyUseCase,
		exportFullPubkeyUseCase: exportFullPubkeyUseCase,
		signTxUseCase:           signTxUseCase,
	}
}

// GenerateSeed generates seed
func (s *BTCSign) GenerateSeed() ([]byte, error) {
	output, err := s.generateSeedUseCase.Generate(context.Background())
	if err != nil {
		return nil, err
	}
	return output.Seed, nil
}

// StoreSeed stores seed
func (s *BTCSign) StoreSeed(strSeed string) ([]byte, error) {
	output, err := s.storeSeedUseCase.Store(context.Background(), signusecase.StoreSeedInput{
		Seed: strSeed,
	})
	if err != nil {
		return nil, err
	}
	return output.Seed, nil
}

// GenerateAuthKey generates account keys
func (s *BTCSign) GenerateAuthKey(seed []byte, count uint32) ([]domainKey.WalletKey, error) {
	_, err := s.generateAuthKeyUseCase.Generate(context.Background(), signusecase.GenerateAuthKeyInput{
		AuthType: s.authAccount,
		Seed:     seed,
		Count:    count,
	})
	if err != nil {
		return nil, err
	}
	// Note: Use case returns count, not keys. Keys are stored in database.
	return nil, nil
}

// ImportPrivKey imports privKey
func (s *BTCSign) ImportPrivKey() error {
	return s.importPrivKeyUseCase.Import(context.Background(), signusecase.ImportPrivateKeyInput{})
}

// ExportFullPubkey exports full-pubkey
func (s *BTCSign) ExportFullPubkey() (string, error) {
	output, err := s.exportFullPubkeyUseCase.Export(context.Background())
	if err != nil {
		return "", err
	}
	return output.FileName, nil
}

// SignTx signs on transaction
func (s *BTCSign) SignTx(filePath string) (string, bool, string, error) {
	output, err := s.signTxUseCase.Sign(context.Background(), signusecase.SignTransactionInput{
		FilePath: filePath,
	})
	if err != nil {
		return "", false, "", err
	}

	return output.SignedData, output.IsComplete, output.NextFilePath, nil
}

// Done should be called before exit
func (s *BTCSign) Done() {
	_ = s.dbConn.Close() // Best effort cleanup
	s.BTC.Close()
}
