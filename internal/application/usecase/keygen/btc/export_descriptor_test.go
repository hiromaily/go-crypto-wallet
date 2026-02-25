package btc_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	keygenusecasebtc "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/btc"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainBitcoin "github.com/hiromaily/go-crypto-wallet/internal/domain/bitcoin"
	coldmocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold/mocks"
	filemocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/transaction/mocks"
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
	accountKeyRepo := coldmocks.NewMockBTCAccountKeyRepositorier(t)
	accountKeyRepo.EXPECT().GetOneMaxID(domainAccount.AccountTypeDeposit).Return(
		&domainBitcoin.BTCAccountKey{ID: 1, KeyType: "bip44"}, nil)

	var writtenPath string
	var writtenData []byte
	writer := filemocks.NewMockDescriptorFileWriter(t)
	writer.EXPECT().WriteFile(mock.Anything, mock.Anything).
		Run(func(path string, data []byte) {
			writtenPath = path
			writtenData = data
		}).Return(nil)

	useCase := keygenusecasebtc.NewExportDescriptorUseCase(generator, writer, accountKeyRepo)

	output, err := useCase.Export(context.Background(), keygenusecase.ExportDescriptorInput{
		AccountType:   domainAccount.AccountTypeDeposit,
		OutputPath:    "/tmp/descriptors.txt",
		Format:        keygenusecase.DescriptorFormatText,
		IncludeChange: true,
	})
	require.NoError(t, err)
	require.Equal(t, "/tmp/descriptors.txt", output.FilePath)
	require.Equal(t, "/tmp/descriptors.txt", writtenPath)

	// accountKeyRepo returns bip44 type, so only legacy (pkh) descriptors should be exported
	expectedDescriptors := []string{
		"pkh-desc",
		"pkh-desc-change",
	}
	require.Equal(t, strings.Join(expectedDescriptors, "\n"), string(writtenData))
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
	accountKeyRepo := coldmocks.NewMockBTCAccountKeyRepositorier(t)
	accountKeyRepo.EXPECT().GetOneMaxID(domainAccount.AccountTypePayment).Return(
		&domainBitcoin.BTCAccountKey{ID: 1, KeyType: "bip44"}, nil)

	var writtenData []byte
	writer := filemocks.NewMockDescriptorFileWriter(t)
	writer.EXPECT().WriteFile(mock.Anything, mock.Anything).
		Run(func(_ string, data []byte) {
			writtenData = data
		}).Return(nil)

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
	require.NoError(t, json.Unmarshal(writtenData, &items))
	// accountKeyRepo returns bip44 type, so only 1 legacy descriptor is exported (no change)
	require.Len(t, items, 1)

	require.Equal(t, "pkh-desc", items[0].Descriptor)
	require.Equal(t, "now", items[0].Timestamp)
	require.Equal(t, []int{0, 1000}, items[0].Range)
	require.True(t, items[0].WatchOnly)
}

func TestExportDescriptorUseCase_InvalidFormat(t *testing.T) {
	t.Parallel()

	accountKeyRepo := coldmocks.NewMockBTCAccountKeyRepositorier(t)
	accountKeyRepo.EXPECT().GetOneMaxID(domainAccount.AccountTypeDeposit).Return(
		&domainBitcoin.BTCAccountKey{ID: 1, KeyType: "bip44"}, nil)
	writer := filemocks.NewMockDescriptorFileWriter(t)

	useCase := keygenusecasebtc.NewExportDescriptorUseCase(
		&stubDescriptorGenerator{},
		writer,
		accountKeyRepo,
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
