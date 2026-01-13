package btc_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	keygenusecasebtc "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/btc"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainBitcoin "github.com/hiromaily/go-crypto-wallet/internal/domain/bitcoin"
)

func TestExportDescriptorUseCase_TextFormat(t *testing.T) {
	t.Parallel()

	generator := &stubDescriptorGenerator{
		descriptors: map[domainAddress.AddrType]string{
			domainAddress.AddrTypeTaproot:    "tr-desc",
			domainAddress.AddrTypeBech32:     "wpkh-desc",
			domainAddress.AddrTypeP2shSegwit: "shwpkh-desc",
			domainAddress.AddrTypeLegacy:     "pkh-desc",
		},
	}
	writer := &stubDescriptorFileWriter{}
	accountKeyRepo := &stubAccountKeyRepo{}

	useCase := keygenusecasebtc.NewExportDescriptorUseCase(generator, writer, accountKeyRepo)

	output, err := useCase.Export(context.Background(), keygenusecase.ExportDescriptorInput{
		AccountType:   domainAccount.AccountTypeDeposit,
		OutputPath:    "/tmp/descriptors.txt",
		Format:        keygenusecase.DescriptorFormatText,
		IncludeChange: true,
	})
	require.NoError(t, err)
	require.Equal(t, "/tmp/descriptors.txt", output.FilePath)
	require.Equal(t, output.FilePath, writer.path)

	expectedDescriptors := []string{
		"tr-desc",
		"tr-desc-change",
		"wpkh-desc",
		"wpkh-desc-change",
		"shwpkh-desc",
		"shwpkh-desc-change",
		"pkh-desc",
		"pkh-desc-change",
	}
	require.Equal(t, strings.Join(expectedDescriptors, "\n"), string(writer.data))
}

func TestExportDescriptorUseCase_BitcoinCoreFormat(t *testing.T) {
	t.Parallel()

	generator := &stubDescriptorGenerator{
		descriptors: map[domainAddress.AddrType]string{
			domainAddress.AddrTypeTaproot:    "tr-desc",
			domainAddress.AddrTypeBech32:     "wpkh-desc",
			domainAddress.AddrTypeP2shSegwit: "shwpkh-desc",
			domainAddress.AddrTypeLegacy:     "pkh-desc",
		},
	}
	writer := &stubDescriptorFileWriter{}
	accountKeyRepo := &stubAccountKeyRepo{}

	useCase := keygenusecasebtc.NewExportDescriptorUseCase(generator, writer, accountKeyRepo)

	output, err := useCase.Export(context.Background(), keygenusecase.ExportDescriptorInput{
		AccountType: domainAccount.AccountTypePayment,
		OutputPath:  "/tmp/descriptors.json",
		Format:      keygenusecase.DescriptorFormatBitcoinCore,
	})
	require.NoError(t, err)
	require.Equal(t, "/tmp/descriptors.json", output.FilePath)

	var items []struct {
		Descriptor string `json:"desc"`
		Timestamp  string `json:"timestamp"`
		Range      []int  `json:"range"`
		WatchOnly  bool   `json:"watchonly"`
	}
	require.NoError(t, json.Unmarshal(writer.data, &items))
	require.Len(t, items, 4)

	require.Equal(t, "tr-desc", items[0].Descriptor)
	require.Equal(t, "now", items[0].Timestamp)
	require.Equal(t, []int{0, 1000}, items[0].Range)
	require.True(t, items[0].WatchOnly)

	require.Equal(t, "wpkh-desc", items[1].Descriptor)
	require.Equal(t, "shwpkh-desc", items[2].Descriptor)
	require.Equal(t, "pkh-desc", items[3].Descriptor)
}

func TestExportDescriptorUseCase_InvalidFormat(t *testing.T) {
	t.Parallel()

	useCase := keygenusecasebtc.NewExportDescriptorUseCase(
		&stubDescriptorGenerator{},
		&stubDescriptorFileWriter{},
		&stubAccountKeyRepo{},
	)

	_, err := useCase.Export(context.Background(), keygenusecase.ExportDescriptorInput{
		AccountType: domainAccount.AccountTypeDeposit,
		OutputPath:  "/tmp/descriptors.invalid",
		Format:      keygenusecase.DescriptorFormat("unsupported"),
	})
	require.Error(t, err)
}

type stubDescriptorGenerator struct {
	descriptors map[domainAddress.AddrType]string
}

func (s *stubDescriptorGenerator) Generate(
	_ context.Context,
	input keygenusecase.GenerateDescriptorInput,
) (keygenusecase.GenerateDescriptorOutput, error) {
	base, ok := s.descriptors[input.AddressType]
	if !ok {
		return keygenusecase.GenerateDescriptorOutput{}, fmt.Errorf("descriptor not found for %s", input.AddressType)
	}

	if input.IsChange {
		base += "-change"
	}

	return keygenusecase.GenerateDescriptorOutput{
		Descriptor:  base,
		AccountType: input.AccountType,
		AddressType: input.AddressType,
	}, nil
}

type stubDescriptorFileWriter struct {
	path string
	data []byte
}

func (s *stubDescriptorFileWriter) WriteFile(path string, data []byte) error {
	s.path = path
	s.data = data
	return nil
}

type stubAccountKeyRepo struct{}

func (*stubAccountKeyRepo) GetMaxIndex(_ domainAccount.AccountType) (int64, error) {
	return 0, nil
}

func (*stubAccountKeyRepo) GetOneMaxID(_ domainAccount.AccountType) (*domainBitcoin.BtcAccountKey, error) {
	// Return a mock account key with legacy/bip44 type for testing
	// This simulates having keys available for all address types
	return &domainBitcoin.BtcAccountKey{
		ID:      1,
		KeyType: "bip44",
	}, nil
}

func (*stubAccountKeyRepo) GetAllAddrStatus(
	_ domainAccount.AccountType, _ domainAddress.AddrStatus,
) ([]*domainBitcoin.BtcAccountKey, error) {
	return nil, nil
}

func (*stubAccountKeyRepo) GetAllMultiAddr(
	_ domainAccount.AccountType, _ []string,
) ([]*domainBitcoin.BtcAccountKey, error) {
	return nil, nil
}

func (*stubAccountKeyRepo) InsertBulk(_ []*domainBitcoin.BtcAccountKey) error {
	return nil
}

func (*stubAccountKeyRepo) UpdateAddr(
	_ domainAccount.AccountType, _, _ string,
) (int64, error) {
	return 0, nil
}

func (*stubAccountKeyRepo) UpdateAddrStatus(
	_ domainAccount.AccountType, _ domainAddress.AddrStatus, _ []string,
) (int64, error) {
	return 0, nil
}

func (*stubAccountKeyRepo) UpdateMultisigAddr(
	_ domainAccount.AccountType, _ *domainBitcoin.BtcAccountKey,
) (int64, error) {
	return 0, nil
}

func (*stubAccountKeyRepo) UpdateMultisigAddrs(
	_ domainAccount.AccountType, _ []*domainBitcoin.BtcAccountKey,
) (int64, error) {
	return 0, nil
}
