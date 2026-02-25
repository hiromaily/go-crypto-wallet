package eth_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	keygenusecaseeth "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/eth"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainEth "github.com/hiromaily/go-crypto-wallet/internal/domain/ethereum"
)

// stubETHAccountKeyRepo is a test stub for ETHAccountKeyRepositorier.
type stubETHAccountKeyRepo struct {
	key *domainEth.ETHAccountKey
	err error
}

func (*stubETHAccountKeyRepo) GetMaxIndex(_ context.Context, _ domainAccount.AccountType) (int64, error) {
	return 0, nil
}

func (s *stubETHAccountKeyRepo) GetOneMaxID(_ domainAccount.AccountType) (*domainEth.ETHAccountKey, error) {
	return s.key, s.err
}

func (*stubETHAccountKeyRepo) GetAllAddrStatus(
	_ domainAccount.AccountType, _ domainAddress.AddrStatus,
) ([]*domainEth.ETHAccountKey, error) {
	return nil, nil
}

func (*stubETHAccountKeyRepo) GetByAddress(_ string) (*domainEth.ETHAccountKey, error) {
	return nil, nil
}

func (*stubETHAccountKeyRepo) InsertBulk(_ []*domainEth.ETHAccountKey) error {
	return nil
}

func (*stubETHAccountKeyRepo) UpdateAddrStatus(
	_ domainAccount.AccountType, _ domainAddress.AddrStatus, _ []string,
) (int64, error) {
	return 0, nil
}

// stubAddressFileRepo is a test stub for AddressFileRepositorier.
type stubAddressFileRepo struct {
	path string
}

func (s *stubAddressFileRepo) CreateFilePath(_ domainAccount.AccountType) string {
	return s.path
}

func (*stubAddressFileRepo) ValidateFilePath(_ string, _ domainAccount.AccountType) error {
	return nil
}

func (*stubAddressFileRepo) ImportAddress(_ string) ([]string, error) {
	return nil, nil
}

// newTestAccountXpriv generates a valid BIP-32 extended private key for testing.
func newTestAccountXpriv(t *testing.T) string {
	t.Helper()
	seed := make([]byte, 32)
	// Use a deterministic seed for reproducible tests
	for i := range seed {
		seed[i] = byte(i + 1)
	}
	master, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	require.NoError(t, err)

	// Derive m/44'/60'/0' (BIP-44 ETH account 0)
	purpose, err := master.Derive(hdkeychain.HardenedKeyStart + 44)
	require.NoError(t, err)
	coinType, err := purpose.Derive(hdkeychain.HardenedKeyStart + 60)
	require.NoError(t, err)
	account, err := coinType.Derive(hdkeychain.HardenedKeyStart + 0)
	require.NoError(t, err)

	return account.String()
}

func TestExportFullPubkeyUseCase_Export_Success(t *testing.T) {
	t.Parallel()

	xpriv := newTestAccountXpriv(t)

	ethKey, err := domainEth.NewETHAccountKey(
		domainAccount.AccountTypeDeposit,
		"0xabcdef1234567890abcdef1234567890abcdef12",
		"0x04...",
		"privatekey",
		0,
	)
	require.NoError(t, err)
	ethKey.SetAccountExtendedPrivkey(xpriv)

	tmpFile, err := os.CreateTemp(t.TempDir(), "xpub-*.csv")
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())

	accountKeyRepo := &stubETHAccountKeyRepo{key: ethKey}
	fileRepo := &stubAddressFileRepo{path: tmpFile.Name()}

	useCase := keygenusecaseeth.NewExportFullPubkeyUseCase(
		accountKeyRepo,
		fileRepo,
		domainCoin.ETH,
	)

	output, err := useCase.Export(context.Background(), keygenusecase.ExportFullPubkeyInput{
		AccountType: domainAccount.AccountTypeDeposit,
	})
	require.NoError(t, err)
	require.Equal(t, tmpFile.Name(), output.FileName)

	// Verify file contents
	data, err := os.ReadFile(tmpFile.Name())
	require.NoError(t, err)
	line := strings.TrimSpace(string(data))

	// Format: coinTypeCode,accountType,44,xpub,derivationPath
	parts := strings.Split(line, ",")
	require.Len(t, parts, 5, "expected 5 CSV fields, got: %s", line)
	require.Equal(t, "eth", parts[0], "coinTypeCode")
	require.Equal(t, "deposit", parts[1], "accountType")
	require.Equal(t, "44", parts[2], "purpose")
	require.True(t, strings.HasPrefix(parts[3], "xpub"), "xpub should start with 'xpub', got: %s", parts[3])
	require.Equal(t, "m/44'/60'/0'", parts[4], "derivationPath")
}

func TestExportFullPubkeyUseCase_Export_MissingXpriv(t *testing.T) {
	t.Parallel()

	ethKey, err := domainEth.NewETHAccountKey(
		domainAccount.AccountTypeDeposit,
		"0xabcdef1234567890abcdef1234567890abcdef12",
		"0x04...",
		"privatekey",
		0,
	)
	require.NoError(t, err)
	// Do NOT set accountXpriv

	accountKeyRepo := &stubETHAccountKeyRepo{key: ethKey}
	fileRepo := &stubAddressFileRepo{path: "/tmp/unused.csv"}

	useCase := keygenusecaseeth.NewExportFullPubkeyUseCase(
		accountKeyRepo,
		fileRepo,
		domainCoin.ETH,
	)

	_, err = useCase.Export(context.Background(), keygenusecase.ExportFullPubkeyInput{
		AccountType: domainAccount.AccountTypeDeposit,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "accountXpriv not found")
}

func TestExportFullPubkeyUseCase_Export_RepoError(t *testing.T) {
	t.Parallel()

	accountKeyRepo := &stubETHAccountKeyRepo{err: os.ErrNotExist}
	fileRepo := &stubAddressFileRepo{path: "/tmp/unused.csv"}

	useCase := keygenusecaseeth.NewExportFullPubkeyUseCase(
		accountKeyRepo,
		fileRepo,
		domainCoin.ETH,
	)

	_, err := useCase.Export(context.Background(), keygenusecase.ExportFullPubkeyInput{
		AccountType: domainAccount.AccountTypeDeposit,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get ETH account key")
}
