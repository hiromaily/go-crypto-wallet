package btc_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	keygenusecasebtc "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/btc"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainAuth "github.com/hiromaily/go-crypto-wallet/internal/domain/auth"
	domainBitcoin "github.com/hiromaily/go-crypto-wallet/internal/domain/bitcoin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin/btc"
	infraKey "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key"
)

const (
	testDescriptorMainnetXpub = "xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4AL" +
		"HY2grBGRjaDMzQLcgJvLJuZZvRcEL"
)

func TestNewGenerateDescriptorUseCase(t *testing.T) {
	t.Parallel()

	useCase := keygenusecasebtc.NewGenerateDescriptorUseCase(
		btc.NewDescriptorService(&chaincfg.MainNetParams),
		&chaincfg.MainNetParams,
		nil,
		nil,
		nil,
	)
	require.NotNil(t, useCase)
}

func TestGenerateDescriptorUseCase_SingleSig(t *testing.T) {
	t.Parallel()

	descriptorService := btc.NewDescriptorService(&chaincfg.MainNetParams)
	accountRepo := &stubAccountRepo{
		key: &domainBitcoin.BtcAccountKey{
			FullPublicKey: testDescriptorMainnetXpub,
			Account:       domainAccount.AccountTypeDeposit,
		},
	}

	useCase := keygenusecasebtc.NewGenerateDescriptorUseCase(
		descriptorService,
		&chaincfg.MainNetParams,
		&stubAuthRepo{},
		accountRepo,
		nil,
	)

	fp, err := infraKey.FingerprintFromExtendedKey(testDescriptorMainnetXpub)
	require.NoError(t, err)

	output, err := useCase.Generate(context.Background(), keygenusecase.GenerateDescriptorInput{
		AccountType: domainAccount.AccountTypeDeposit,
		AddressType: domainAddress.AddrTypeBech32,
		IsChange:    false,
	})
	require.NoError(t, err)
	require.Equal(t, "wpkh(["+fp.String()+"/84'/0'/0']"+testDescriptorMainnetXpub+"/0/*)", output.Descriptor)
	require.False(t, output.IsMultisig)
}

func TestGenerateDescriptorUseCase_MultisigWsh(t *testing.T) {
	t.Parallel()

	signers := map[domainAccount.AuthType]*domainAuth.AuthFullPubkey{
		domainAccount.AuthType1: {FullPublicKey: testDescriptorMainnetXpub},
		domainAccount.AuthType2: {FullPublicKey: newTestXpubFromSeed(t, 0x02)},
	}

	multiConfig := domainAccount.NewMultisigConfig(map[domainAccount.AccountType]map[int][]domainAccount.AuthType{
		domainAccount.AccountTypeDeposit: {
			2: {domainAccount.AuthType1, domainAccount.AuthType2},
		},
	})

	descriptorService := btc.NewDescriptorService(&chaincfg.MainNetParams)
	useCase := keygenusecasebtc.NewGenerateDescriptorUseCase(
		descriptorService,
		&chaincfg.MainNetParams,
		&stubAuthRepo{items: signers},
		&stubAccountRepo{},
		multiConfig,
	)

	fp1, err := infraKey.FingerprintFromExtendedKey(signers[domainAccount.AuthType1].FullPublicKey)
	require.NoError(t, err)
	fp2, err := infraKey.FingerprintFromExtendedKey(signers[domainAccount.AuthType2].FullPublicKey)
	require.NoError(t, err)

	// Build expected with deterministic ordering (descriptor service sorts)
	xpub1, _ := hdkeychain.NewKeyFromString(signers[domainAccount.AuthType1].FullPublicKey)
	xpub2, _ := hdkeychain.NewKeyFromString(signers[domainAccount.AuthType2].FullPublicKey)
	expected, err := descriptorService.GenerateMultisigDescriptor(
		2,
		[]btc.MultisigSigner{
			{Fingerprint: fp1.String(), DerivationPath: "/48'/0'/0'/2'", ExtendedKey: xpub1},
			{Fingerprint: fp2.String(), DerivationPath: "/48'/0'/0'/2'", ExtendedKey: xpub2},
		},
		false,
	)
	require.NoError(t, err)

	output, err := useCase.Generate(context.Background(), keygenusecase.GenerateDescriptorInput{
		AccountType:  domainAccount.AccountTypeDeposit,
		AddressType:  domainAddress.AddrTypeBech32,
		IsChange:     false,
		RequiredSigs: 2,
	})
	require.NoError(t, err)
	require.True(t, output.IsMultisig)
	require.Equal(t, expected, output.Descriptor)
}

func TestGenerateDescriptorUseCase_TaprootScriptPath(t *testing.T) {
	t.Parallel()

	signers := map[domainAccount.AuthType]*domainAuth.AuthFullPubkey{
		domainAccount.AuthType1: {FullPublicKey: testDescriptorMainnetXpub},
		domainAccount.AuthType2: {FullPublicKey: newTestXpubFromSeed(t, 0x03)},
	}

	multiConfig := domainAccount.NewMultisigConfig(map[domainAccount.AccountType]map[int][]domainAccount.AuthType{
		domainAccount.AccountTypeDeposit: {
			2: {domainAccount.AuthType1, domainAccount.AuthType2},
		},
	})

	useCase := keygenusecasebtc.NewGenerateDescriptorUseCase(
		btc.NewDescriptorService(&chaincfg.MainNetParams),
		&chaincfg.MainNetParams,
		&stubAuthRepo{items: signers},
		&stubAccountRepo{},
		multiConfig,
	)

	fp1, err := infraKey.FingerprintFromExtendedKey(signers[domainAccount.AuthType1].FullPublicKey)
	require.NoError(t, err)
	fp2, err := infraKey.FingerprintFromExtendedKey(signers[domainAccount.AuthType2].FullPublicKey)
	require.NoError(t, err)

	output, err := useCase.Generate(context.Background(), keygenusecase.GenerateDescriptorInput{
		AccountType:  domainAccount.AccountTypeDeposit,
		AddressType:  domainAddress.AddrTypeTaproot,
		IsChange:     false,
		RequiredSigs: 2,
	})
	require.NoError(t, err)
	require.True(t, output.IsMultisig)
	xpub1, _ := hdkeychain.NewKeyFromString(signers[domainAccount.AuthType1].FullPublicKey)
	xpub2, _ := hdkeychain.NewKeyFromString(signers[domainAccount.AuthType2].FullPublicKey)
	expected, err := btc.NewDescriptorService(&chaincfg.MainNetParams).GenerateTaprootScriptPathDescriptor(
		[]btc.MultisigSigner{
			{Fingerprint: fp1.String(), DerivationPath: "/86'/0'/0'", ExtendedKey: xpub1},
			{Fingerprint: fp2.String(), DerivationPath: "/86'/0'/0'", ExtendedKey: xpub2},
		},
		false,
	)
	require.NoError(t, err)
	require.Equal(t, expected, output.Descriptor)
}

func TestGenerateDescriptorUseCase_MissingAccountKey(t *testing.T) {
	t.Parallel()

	useCase := keygenusecasebtc.NewGenerateDescriptorUseCase(
		btc.NewDescriptorService(&chaincfg.MainNetParams),
		&chaincfg.MainNetParams,
		&stubAuthRepo{},
		&stubAccountRepo{key: nil},
		nil,
	)

	_, err := useCase.Generate(context.Background(), keygenusecase.GenerateDescriptorInput{
		AccountType: domainAccount.AccountTypeDeposit,
		AddressType: domainAddress.AddrTypeLegacy,
		IsChange:    false,
	})
	require.Error(t, err)
}

type stubAccountRepo struct {
	key *domainBitcoin.BtcAccountKey
	err error
}

func (*stubAccountRepo) GetMaxIndex(domainAccount.AccountType) (int64, error) { return 0, nil }
func (s *stubAccountRepo) GetOneMaxID(domainAccount.AccountType) (*domainBitcoin.BtcAccountKey, error) {
	return s.key, s.err
}

func (*stubAccountRepo) GetAllAddrStatus(
	domainAccount.AccountType,
	domainAddress.AddrStatus,
) ([]*domainBitcoin.BtcAccountKey, error) {
	return nil, nil
}

func (*stubAccountRepo) GetAllMultiAddr(domainAccount.AccountType, []string) ([]*domainBitcoin.BtcAccountKey, error) {
	return nil, nil
}
func (*stubAccountRepo) InsertBulk([]*domainBitcoin.BtcAccountKey) error { return nil }
func (*stubAccountRepo) UpdateAddr(domainAccount.AccountType, string, string) (int64, error) {
	return 0, nil
}

func (*stubAccountRepo) UpdateAddrStatus(
	domainAccount.AccountType,
	domainAddress.AddrStatus,
	[]string,
) (int64, error) {
	return 0, nil
}

func (*stubAccountRepo) UpdateMultisigAddr(domainAccount.AccountType, *domainBitcoin.BtcAccountKey) (int64, error) {
	return 0, nil
}

func (*stubAccountRepo) UpdateMultisigAddrs(
	domainAccount.AccountType,
	[]*domainBitcoin.BtcAccountKey,
) (int64, error) {
	return 0, nil
}

type stubAuthRepo struct {
	items map[domainAccount.AuthType]*domainAuth.AuthFullPubkey
}

func (s *stubAuthRepo) GetOne(authType domainAccount.AuthType) (*domainAuth.AuthFullPubkey, error) {
	if s.items == nil {
		return nil, fmt.Errorf("auth not found: %s", authType.String())
	}
	val, ok := s.items[authType]
	if !ok {
		return nil, fmt.Errorf("auth not found: %s", authType.String())
	}
	return val, nil
}
func (*stubAuthRepo) Insert(domainAccount.AuthType, string) error   { return nil }
func (*stubAuthRepo) InsertBulk([]*domainAuth.AuthFullPubkey) error { return nil }

func newTestXpubFromSeed(t *testing.T, seedByte byte) string {
	t.Helper()

	seed := bytes.Repeat([]byte{seedByte}, hdkeychain.RecommendedSeedLen)
	master, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	require.NoError(t, err)

	xpub, err := master.Neuter()
	require.NoError(t, err)

	return xpub.String()
}
