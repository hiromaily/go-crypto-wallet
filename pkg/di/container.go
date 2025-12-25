package di

import (
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	signusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	appdi "github.com/hiromaily/go-crypto-wallet/internal/di"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// Container is for DI container interface
// This interface is kept in pkg/di for backward compatibility
// but the implementation is in internal/di
type Container interface {
	NewWalleter() wallets.Watcher
	NewKeygener() wallets.Keygener
	NewSigner(authName string) wallets.Signer

	// Watch Use Cases
	NewWatchCreateTransactionUseCase() any
	NewWatchMonitorTransactionUseCase() any
	NewWatchSendTransactionUseCase() any
	NewWatchImportAddressUseCase() watchusecase.ImportAddressUseCase
	NewWatchCreatePaymentRequestUseCase() watchusecase.CreatePaymentRequestUseCase

	// Keygen Use Cases
	NewKeygenGenerateHDWalletUseCase() keygenusecase.GenerateHDWalletUseCase
	NewKeygenGenerateSeedUseCase() keygenusecase.GenerateSeedUseCase
	NewKeygenExportAddressUseCase() keygenusecase.ExportAddressUseCase
	NewKeygenImportPrivateKeyUseCase() keygenusecase.ImportPrivateKeyUseCase
	NewKeygenCreateMultisigAddressUseCase() keygenusecase.CreateMultisigAddressUseCase
	NewKeygenImportFullPubkeyUseCase() keygenusecase.ImportFullPubkeyUseCase
	NewKeygenGenerateKeyUseCase() keygenusecase.GenerateKeyUseCase
	NewKeygenSignTransactionUseCase() keygenusecase.SignTransactionUseCase

	// Sign Use Cases
	NewSignTransactionUseCase() signusecase.SignTransactionUseCase
	NewSignImportPrivateKeyUseCase(authType domainAccount.AuthType) signusecase.ImportPrivateKeyUseCase
	NewSignExportFullPubkeyUseCase(authType domainAccount.AuthType) signusecase.ExportFullPubkeyUseCase
	NewSignGenerateSeedUseCase() signusecase.GenerateSeedUseCase
	NewSignStoreSeedUseCase() signusecase.StoreSeedUseCase
	NewSignGenerateAuthKeyUseCase() signusecase.GenerateAuthKeyUseCase

	// Auth accessors
	AuthName() string
	AuthType() domainAccount.AuthType

	// Config accessors
	AddressType() address.AddrType
}

// NewContainer creates a new DI container
// This is a thin wrapper that delegates to internal/di
func NewContainer(
	conf *config.WalletRoot,
	accountConf *account.AccountRoot,
	walletType domainWallet.WalletType,
) Container {
	return appdi.NewContainer(conf, accountConf, walletType)
}
