package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	watchusecasebtc "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/btc"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin/btc"
)

func TestDescriptorImportWorkflow_SingleKey(t *testing.T) {
	t.Parallel()

	desc := "wpkh([a1b2c3d4/84'/0'/0']xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb" +
		"5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL/0/*)"
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "single.txt")
	require.NoError(t, os.WriteFile(filePath, []byte(desc), 0o600))

	addrRepo := &recordingAddressRepo{}
	importer := watchusecasebtc.NewImportDescriptorUseCase(
		btc.NewDescriptorParser(),
		&chaincfg.MainNetParams,
		addrRepo,
		domainCoin.BTC,
	)

	output, err := importer.Import(context.Background(), watchusecase.ImportDescriptorInput{
		FilePath:    filePath,
		AccountType: domainAccount.AccountTypeDeposit,
		StartIndex:  0,
		Count:       2,
	})
	require.NoError(t, err)
	require.Empty(t, output.Errors)
	require.Equal(t, 1, output.DescriptorsImported)
	require.Equal(t, 2, output.AddressesGenerated)
	require.Len(t, addrRepo.inserted, output.AddressesGenerated)
}

func TestDescriptorImportWorkflow_Multisig(t *testing.T) {
	t.Parallel()

	desc := "wsh(sortedmulti(2,[a1b2c3d4/48'/0'/0'/2']" +
		"xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJ" +
		"vLJuZZvRcEL/0/*,[b2c3d4e5/48'/0'/0'/2']" +
		"xpub6D4BDPcP2GT577Vvch3R8wDkScZWzQzMMUm3PWbmWvVJrZwQY4VUNgqFJPMM3No2dFDFGTsxxpG5uJh7n7epu4trkrX7x7Do" +
		"gT5Uv6fcLW5/0/*))"
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "multisig.txt")
	require.NoError(t, os.WriteFile(filePath, []byte(desc), 0o600))

	addrRepo := &recordingAddressRepo{}
	importer := watchusecasebtc.NewImportDescriptorUseCase(
		btc.NewDescriptorParser(),
		&chaincfg.MainNetParams,
		addrRepo,
		domainCoin.BTC,
	)

	output, err := importer.Import(context.Background(), watchusecase.ImportDescriptorInput{
		FilePath:    filePath,
		AccountType: domainAccount.AccountTypeDeposit,
		StartIndex:  0,
		Count:       1,
	})
	require.NoError(t, err)
	require.Empty(t, output.Errors)
	require.Equal(t, 1, output.DescriptorsImported)
	require.Equal(t, 1, output.AddressesGenerated)
	require.Len(t, addrRepo.inserted, output.AddressesGenerated)
}

type recordingAddressRepo struct {
	inserted []*domainAddress.Address
}

func (*recordingAddressRepo) GetAll(domainAccount.AccountType) ([]*domainAddress.Address, error) {
	return nil, nil
}

func (*recordingAddressRepo) GetAllAddress(domainAccount.AccountType) ([]string, error) {
	return nil, nil
}

func (*recordingAddressRepo) GetOneUnAllocated(domainAccount.AccountType) (*domainAddress.Address, error) {
	return nil, nil
}

func (r *recordingAddressRepo) InsertBulk(_ context.Context, items []*domainAddress.Address) error {
	r.inserted = append(r.inserted, items...)
	return nil
}

func (*recordingAddressRepo) UpdateIsAllocated(bool, string) (int64, error) {
	return 0, nil
}
