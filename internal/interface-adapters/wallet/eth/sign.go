package eth

import (
	"context"
	"database/sql"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// ETHSign sign wallet object
type ETHSign struct {
	ETH                   apieth.ETHLifecycle
	dbConn                *sql.DB
	wtype                 domainWallet.WalletType
	signMultisigTxUseCase keygenusecase.SignMultisigTransactionUseCase
	signerAddress         string
}

// NewETHSign returns ETHSign object
func NewETHSign(
	eth apieth.ETHLifecycle,
	dbConn *sql.DB,
	walletType domainWallet.WalletType,
	signMultisigTxUseCase keygenusecase.SignMultisigTransactionUseCase,
) *ETHSign {
	return &ETHSign{
		ETH:                   eth,
		dbConn:                dbConn,
		wtype:                 walletType,
		signMultisigTxUseCase: signMultisigTxUseCase,
	}
}

// SetSignerAddress sets the signer address used for ETH Safe multisig signing.
func (s *ETHSign) SetSignerAddress(addr string) {
	s.signerAddress = addr
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

// SignTx signs on transaction.
// For ETH Safe multisig files (detected by the "safe_address" JSON field), routes to the
// multisig signing use case when a signer address has been set via SetSignerAddress.
// Single-sig ETH signing is Keygen-only; this wallet only handles multisig.
func (s *ETHSign) SignTx(filePath string) (string, bool, string, error) {
	if s.signMultisigTxUseCase != nil && s.signerAddress != "" && isETHMultisigFile(filePath) {
		output, err := s.signMultisigTxUseCase.Sign(context.Background(),
			keygenusecase.SignMultisigTransactionInput{
				FilePath:      filePath,
				SignerAddress: s.signerAddress,
			})
		if err != nil {
			return "", false, "", err
		}
		return output.FilePath, output.IsComplete, "", nil
	}
	logger.Info("no functionality for SignTx() in ETH single-sig Sign wallet")
	return "", false, "", nil
}

// Done should be called before exit
func (s *ETHSign) Done() {
	_ = s.dbConn.Close() // Best effort cleanup
	s.ETH.Close()
}
