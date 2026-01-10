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

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	watchusecasebtc "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/btc"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin/btc"
	btcmocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin/mocks"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/mocks"
	infraKey "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key"
)

const (
	testDescriptorMainnetXpub = "xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4" +
		"ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL"
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

	btcClient := btcmocks.NewMockBitcoiner(t)
	btcClient.EXPECT().
		ImportDescriptors(mock.MatchedBy(func(reqs []dtobtc.ImportDescriptorsRequest) bool {
			return len(reqs) == 1 && reqs[0].Active && reqs[0].Watchonly
		})).
		Return([]dtobtc.ImportDescriptorsResponse{{Success: true}}, nil).
		Once()

	parser := btc.NewDescriptorParser()
	useCase := watchusecasebtc.NewImportDescriptorUseCase(
		btcClient,
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

	btcClient := btcmocks.NewMockBitcoiner(t)
	// No ImportDescriptors call expected in ValidateOnly mode

	parser := btc.NewDescriptorParser()
	useCase := watchusecasebtc.NewImportDescriptorUseCase(
		btcClient,
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

func TestImportDescriptorUseCase_ImportsMultisigAddresses(t *testing.T) {
	t.Parallel()

	xpub2 := "xpub6D4BDPcP2GT577Vvch3R8wDkScZWzQzMMUm3PWbmWvVJrZwQY4VUNgqFJPMM3No2dFDFGTsxxpG5uJh7n7epu4trkrX7" +
		"x7DogT5Uv6fcLW5"
	desc := fmt.Sprintf(
		"wsh(sortedmulti(2,[a1b2c3d4/48'/0'/0'/2']%s/0/*,[b2c3d4e5/48'/0'/0'/2']%s/0/*))",
		testDescriptorMainnetXpub,
		xpub2,
	)

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "desc.txt")
	require.NoError(t, os.WriteFile(filePath, []byte(desc), 0o600))

	repo := mocks.NewMockAddressRepositorier(t)
	repo.EXPECT().
		InsertBulk(mock.Anything, mock.Anything).
		Return(nil).
		Once()

	btcClient := btcmocks.NewMockBitcoiner(t)
	btcClient.EXPECT().
		ImportDescriptors(mock.MatchedBy(func(reqs []dtobtc.ImportDescriptorsRequest) bool {
			return len(reqs) == 1 && reqs[0].Active && reqs[0].Watchonly
		})).
		Return([]dtobtc.ImportDescriptorsResponse{{Success: true}}, nil).
		Once()

	parser := btc.NewDescriptorParser()
	useCase := watchusecasebtc.NewImportDescriptorUseCase(
		btcClient,
		parser,
		&chaincfg.MainNetParams,
		repo,
		domainCoin.BTC,
	)

	output, err := useCase.Import(context.Background(), watchusecase.ImportDescriptorInput{
		FilePath:    filePath,
		AccountType: domainAccount.AccountTypeDeposit,
		StartIndex:  0,
		Count:       1,
	})
	require.NoError(t, err)
	require.Equal(t, 1, output.DescriptorsImported)
	require.Equal(t, 1, output.AddressesGenerated)
	require.Empty(t, output.Errors)
}
