package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	keygenusecasebtc "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/btc"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	watchusecasebtc "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/btc"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainAuth "github.com/hiromaily/go-crypto-wallet/internal/domain/auth"
	domainBitcoin "github.com/hiromaily/go-crypto-wallet/internal/domain/bitcoin"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin/btc"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key"
)

const descriptorWorkflowXpub = "xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL"

func TestDescriptorWorkflow_EndToEnd(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	fp, err := key.FingerprintFromExtendedKey(descriptorWorkflowXpub)
	require.NoError(t, err)

	accountRepo := &stubAccountRepo{
		key: &domainBitcoin.BtcAccountKey{
			FullPublicKey: descriptorWorkflowXpub,
			Account:       domainAccount.AccountTypeDeposit,
		},
	}

	generator := keygenusecasebtc.NewGenerateDescriptorUseCase(
		btc.NewDescriptorService(&chaincfg.MainNetParams),
		&chaincfg.MainNetParams,
		&stubAuthFullPubkeyRepo{fingerprint: fp},
		accountRepo,
		domainAccount.NewMultisigConfig(nil),
	)

	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "descriptors.json")
	writer := &recordingFileWriter{}

	exporter := keygenusecasebtc.NewExportDescriptorUseCase(generator, writer)
	exported, err := exporter.Export(ctx, keygenusecase.ExportDescriptorInput{
		AccountType:   domainAccount.AccountTypeDeposit,
		OutputPath:    outputPath,
		Format:        keygenusecase.DescriptorFormatBitcoinCore,
		IncludeChange: true,
	})
	require.NoError(t, err)
	require.NotEmpty(t, exported.FilePath)

	// Ensure the descriptors file was written for watch wallet import.
	_, statErr := os.Stat(exported.FilePath)
	require.NoError(t, statErr)

	descriptorCount, filterErr := filterNonTaprootDescriptors(exported.FilePath)
	require.NoError(t, filterErr)

	addrRepo := &recordingAddressRepo{}
	importer := watchusecasebtc.NewImportDescriptorUseCase(
		btc.NewDescriptorParser(),
		&chaincfg.MainNetParams,
		addrRepo,
		domainCoin.BTC,
	)

	imported, err := importer.Import(ctx, watchusecase.ImportDescriptorInput{
		FilePath:    exported.FilePath,
		AccountType: domainAccount.AccountTypeDeposit,
		StartIndex:  0,
		Count:       2,
	})
	require.NoError(t, err)
	require.Empty(t, imported.Errors)
	require.Equal(t, descriptorCount, imported.DescriptorsImported, "expect receive and change descriptors for supported address types")
	require.Equal(t, descriptorCount*2, imported.AddressesGenerated, "two addresses per descriptor")
	require.Len(t, addrRepo.inserted, imported.AddressesGenerated)
}

func filterNonTaprootDescriptors(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}

	var coreFormat []map[string]any
	if jsonErr := json.Unmarshal(data, &coreFormat); jsonErr != nil {
		return 0, fmt.Errorf("parse exported descriptors: %w", jsonErr)
	}
	if len(coreFormat) == 0 {
		return 0, fmt.Errorf("no descriptors found in %s", path)
	}

	var filtered []map[string]any
	for _, rec := range coreFormat {
		desc, _ := rec["desc"].(string)
		if strings.HasPrefix(desc, "tr(") {
			continue
		}
		filtered = append(filtered, rec)
	}

	if len(filtered) != len(coreFormat) {
		encoded, marshalErr := json.MarshalIndent(filtered, "", "  ")
		if marshalErr != nil {
			return 0, marshalErr
		}
		if writeErr := os.WriteFile(path, encoded, 0o600); writeErr != nil {
			return 0, writeErr
		}
	}

	return len(filtered), nil
}

type stubAuthFullPubkeyRepo struct {
	fingerprint domainKey.Fingerprint
}

func (s *stubAuthFullPubkeyRepo) GetOne(domainAccount.AuthType) (*domainAuth.AuthFullPubkey, error) {
	return &domainAuth.AuthFullPubkey{
		FullPublicKey: descriptorWorkflowXpub,
		Fingerprint:   &s.fingerprint,
	}, nil
}

func (s *stubAuthFullPubkeyRepo) Insert(domainAccount.AuthType, string) error {
	return nil
}

func (s *stubAuthFullPubkeyRepo) InsertBulk([]*domainAuth.AuthFullPubkey) error {
	return nil
}

type stubAccountRepo struct {
	key *domainBitcoin.BtcAccountKey
}

func (s *stubAccountRepo) GetMaxIndex(domainAccount.AccountType) (int64, error) {
	return 0, nil
}

func (s *stubAccountRepo) GetOneMaxID(domainAccount.AccountType) (*domainBitcoin.BtcAccountKey, error) {
	return s.key, nil
}

func (s *stubAccountRepo) GetAllAddrStatus(domainAccount.AccountType, domainAddress.AddrStatus) ([]*domainBitcoin.BtcAccountKey, error) {
	return nil, nil
}

func (s *stubAccountRepo) GetAllMultiAddr(domainAccount.AccountType, []string) ([]*domainBitcoin.BtcAccountKey, error) {
	return nil, nil
}

func (s *stubAccountRepo) InsertBulk([]*domainBitcoin.BtcAccountKey) error {
	return nil
}

func (s *stubAccountRepo) UpdateAddr(domainAccount.AccountType, string, string) (int64, error) {
	return 0, nil
}

func (s *stubAccountRepo) UpdateAddrStatus(domainAccount.AccountType, domainAddress.AddrStatus, []string) (int64, error) {
	return 0, nil
}

func (s *stubAccountRepo) UpdateMultisigAddr(domainAccount.AccountType, *domainBitcoin.BtcAccountKey) (int64, error) {
	return 0, nil
}

func (s *stubAccountRepo) UpdateMultisigAddrs(domainAccount.AccountType, []*domainBitcoin.BtcAccountKey) (int64, error) {
	return 0, nil
}

type recordingFileWriter struct{}

func (recordingFileWriter) WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0o600)
}

type recordingAddressRepo struct {
	inserted []*domainAddress.Address
}

func (r *recordingAddressRepo) GetAll(domainAccount.AccountType) ([]*domainAddress.Address, error) {
	return nil, nil
}

func (r *recordingAddressRepo) GetAllAddress(domainAccount.AccountType) ([]string, error) {
	return nil, nil
}

func (r *recordingAddressRepo) GetOneUnAllocated(domainAccount.AccountType) (*domainAddress.Address, error) {
	return nil, nil
}

func (r *recordingAddressRepo) InsertBulk(_ context.Context, items []*domainAddress.Address) error {
	r.inserted = append(r.inserted, items...)
	return nil
}

func (r *recordingAddressRepo) UpdateIsAllocated(bool, string) (int64, error) {
	return 0, nil
}
