package btc_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	watchusecasebtc "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/btc"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin/btc"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/mocks"
	infraKey "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key"
)

const (
	testDescriptorMainnetXpub = "xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL"
)

func TestImportDescriptorUseCase_ImportsAddresses(t *testing.T) {
	t.Parallel()

	fp, err := infraKey.FingerprintFromExtendedKey(testDescriptorMainnetXpub)
	require.NoError(t, err)

	desc := fmt.Sprintf("wpkh([%s/84'/0'/0']%s/0/*)", fp.String(), testDescriptorMainnetXpub)

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "desc.txt")
	require.NoError(t, os.WriteFile(filePath, []byte(desc), 0o600))

	repo := mocks.NewMockAddressRepositorier(t)
	repo.EXPECT().
		InsertBulk(mock.Anything, mock.Anything).
		Return(nil).
		Once()

	parser := btc.NewDescriptorParser()
	useCase := watchusecasebtc.NewImportDescriptorUseCase(
		parser,
		&chaincfg.MainNetParams,
		repo,
		domainCoin.BTC,
	)

	output, err := useCase.Import(context.Background(), watchusecase.ImportDescriptorInput{
		FilePath:    filePath,
		AccountType: domainAccount.AccountTypeDeposit,
		StartIndex:  0,
		Count:       2,
	})
	require.NoError(t, err)
	require.Equal(t, 1, output.DescriptorsImported)
	require.Equal(t, 2, output.AddressesGenerated)
	require.Empty(t, output.Errors)
}

func TestImportDescriptorUseCase_ValidateOnly(t *testing.T) {
	t.Parallel()

	fp, err := infraKey.FingerprintFromExtendedKey(testDescriptorMainnetXpub)
	require.NoError(t, err)

	desc := fmt.Sprintf("wpkh([%s/84'/0'/0']%s/0/*)", fp.String(), testDescriptorMainnetXpub)

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "desc.txt")
	require.NoError(t, os.WriteFile(filePath, []byte(desc), 0o600))

	repo := mocks.NewMockAddressRepositorier(t)
	// No insert expected

	parser := btc.NewDescriptorParser()
	useCase := watchusecasebtc.NewImportDescriptorUseCase(
		parser,
		&chaincfg.MainNetParams,
		repo,
		domainCoin.BTC,
	)

	output, err := useCase.Import(context.Background(), watchusecase.ImportDescriptorInput{
		FilePath:     filePath,
		AccountType:  domainAccount.AccountTypeDeposit,
		StartIndex:   0,
		Count:        1,
		ValidateOnly: true,
	})
	require.NoError(t, err)
	require.Equal(t, 0, output.DescriptorsImported)
	require.Equal(t, 0, output.AddressesGenerated)
	require.Empty(t, output.Errors)
}
