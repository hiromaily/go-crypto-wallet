package btcwallet

import (
	"context"
	"database/sql"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin"
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
)

// BTCKeygen is keygen wallet object
type BTCKeygen struct {
	BTC                       bitcoin.Bitcoiner
	dbConn                    *sql.DB
	addrType                  address.AddrType
	wtype                     domainWallet.WalletType
	generateSeedUseCase       keygenusecase.GenerateSeedUseCase
	generateHDWalletUseCase   keygenusecase.GenerateHDWalletUseCase
	importPrivKeyUseCase      keygenusecase.ImportPrivateKeyUseCase
	importFullPubkeyUseCase   keygenusecase.ImportFullPubkeyUseCase
	createMultisigAddrUseCase keygenusecase.CreateMultisigAddressUseCase
	exportAddressUseCase      keygenusecase.ExportAddressUseCase
	signTxUseCase             keygenusecase.SignTransactionUseCase
}

// NewBTCKeygen returns Keygen object
func NewBTCKeygen(
	btc bitcoin.Bitcoiner,
	dbConn *sql.DB,
	addrType address.AddrType,
	generateSeedUseCase keygenusecase.GenerateSeedUseCase,
	generateHDWalletUseCase keygenusecase.GenerateHDWalletUseCase,
	importPrivKeyUseCase keygenusecase.ImportPrivateKeyUseCase,
	importFullPubkeyUseCase keygenusecase.ImportFullPubkeyUseCase,
	createMultisigAddrUseCase keygenusecase.CreateMultisigAddressUseCase,
	exportAddressUseCase keygenusecase.ExportAddressUseCase,
	signTxUseCase keygenusecase.SignTransactionUseCase,
	wtype domainWallet.WalletType,
) *BTCKeygen {
	return &BTCKeygen{
		BTC:                       btc,
		dbConn:                    dbConn,
		addrType:                  addrType,
		wtype:                     wtype,
		generateSeedUseCase:       generateSeedUseCase,
		generateHDWalletUseCase:   generateHDWalletUseCase,
		importPrivKeyUseCase:      importPrivKeyUseCase,
		importFullPubkeyUseCase:   importFullPubkeyUseCase,
		createMultisigAddrUseCase: createMultisigAddrUseCase,
		exportAddressUseCase:      exportAddressUseCase,
		signTxUseCase:             signTxUseCase,
	}
}

// GenerateSeed generates seed
func (k *BTCKeygen) GenerateSeed() ([]byte, error) {
	output, err := k.generateSeedUseCase.Generate(context.Background())
	if err != nil {
		return nil, err
	}
	return output.Seed, nil
}

// StoreSeed stores seed
func (k *BTCKeygen) StoreSeed(strSeed string) ([]byte, error) {
	output, err := k.generateSeedUseCase.Store(context.Background(), keygenusecase.StoreSeedInput{
		Seed: strSeed,
	})
	if err != nil {
		return nil, err
	}
	return output.Seed, nil
}

// GenerateAccountKey generates account keys
func (k *BTCKeygen) GenerateAccountKey(
	accountType domainAccount.AccountType, seed []byte, count uint32, _ bool,
) ([]domainKey.WalletKey, error) {
	_, err := k.generateHDWalletUseCase.Generate(context.Background(), keygenusecase.GenerateHDWalletInput{
		AccountType: accountType,
		Seed:        seed,
		Count:       count,
	})
	if err != nil {
		return nil, err
	}
	// Note: The original implementation returns keys but use case returns count
	// For now, return nil as the keys are stored in the database
	return nil, nil
}

// ImportPrivKey imports privKey
func (k *BTCKeygen) ImportPrivKey(accountType domainAccount.AccountType) error {
	return k.importPrivKeyUseCase.Import(context.Background(), keygenusecase.ImportPrivateKeyInput{
		AccountType: accountType,
	})
}

// ImportFullPubKey imports full-pubkey
func (k *BTCKeygen) ImportFullPubKey(fileName string) error {
	return k.importFullPubkeyUseCase.Import(context.Background(), keygenusecase.ImportFullPubkeyInput{
		FileName: fileName,
	})
}

// CreateMultisigAddress creates multi sig address returns Multisiger interface
func (k *BTCKeygen) CreateMultisigAddress(accountType domainAccount.AccountType) error {
	return k.createMultisigAddrUseCase.Create(context.Background(), keygenusecase.CreateMultisigAddressInput{
		AccountType: accountType,
		AddressType: k.addrType,
	})
}

// ExportAddress exports address
func (k *BTCKeygen) ExportAddress(accountType domainAccount.AccountType) (string, error) {
	output, err := k.exportAddressUseCase.Export(context.Background(), keygenusecase.ExportAddressInput{
		AccountType: accountType,
	})
	if err != nil {
		return "", err
	}
	return output.FileName, nil
}

// SignTx signs on transaction
func (k *BTCKeygen) SignTx(filePath string) (string, bool, string, error) {
	output, err := k.signTxUseCase.Sign(context.Background(), keygenusecase.SignTransactionInput{
		FilePath: filePath,
	})
	if err != nil {
		return "", false, "", err
	}

	// Determine if signing is done based on signed vs unsigned count
	isDone := output.UnsignedCount == 0

	return output.FilePath, isDone, "", nil
}

// Done should be called before exit
func (k *BTCKeygen) Done() {
	_ = k.dbConn.Close() // Best effort cleanup
	k.BTC.Close()
}
