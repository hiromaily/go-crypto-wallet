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
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
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
		domainCoin.BTC,
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
		domainCoin.BTC,
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

	testSeed := bytes.Repeat([]byte{0x01}, hdkeychain.RecommendedSeedLen)
	signers := map[domainAccount.AuthType]*domainAuth.AuthFullPubkey{
		domainAccount.AuthType1: {ExtendedPubKey: testDescriptorMainnetXpub},
		domainAccount.AuthType2: {ExtendedPubKey: newTestXpubFromSeed(t, 0x02)},
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
		&stubSeedRepo{seed: testSeed},
		domainCoin.BTC,
		multiConfig,
	)

	output, err := useCase.Generate(context.Background(), keygenusecase.GenerateDescriptorInput{
		AccountType:  domainAccount.AccountTypeDeposit,
		AddressType:  domainAddress.AddrTypeBech32,
		IsChange:     false,
		RequiredSigs: 2,
	})
	require.NoError(t, err)
	require.True(t, output.IsMultisig)

	// Verify descriptor format (multisig descriptor with derived account xpubs)
	require.Contains(t, output.Descriptor, "wsh(sortedmulti(2,")
	require.Contains(t, output.Descriptor, "/84'/0'/0]") // BIP84 path for Native SegWit
}

func TestGenerateDescriptorUseCase_TaprootScriptPath(t *testing.T) {
	t.Parallel()

	testSeed := bytes.Repeat([]byte{0x01}, hdkeychain.RecommendedSeedLen)
	signers := map[domainAccount.AuthType]*domainAuth.AuthFullPubkey{
		domainAccount.AuthType1: {ExtendedPubKey: testDescriptorMainnetXpub},
		domainAccount.AuthType2: {ExtendedPubKey: newTestXpubFromSeed(t, 0x03)},
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
		&stubSeedRepo{seed: testSeed},
		domainCoin.BTC,
		multiConfig,
	)

	output, err := useCase.Generate(context.Background(), keygenusecase.GenerateDescriptorInput{
		AccountType:  domainAccount.AccountTypeDeposit,
		AddressType:  domainAddress.AddrTypeTaproot,
		IsChange:     false,
		RequiredSigs: 2,
	})
	require.NoError(t, err)
	require.True(t, output.IsMultisig)

	// Verify descriptor format (taproot script-path descriptor with derived account xpubs)
	require.Contains(t, output.Descriptor, "tr(")
	require.Contains(t, output.Descriptor, "/86'/0'/0']") // BIP86 path for Taproot
}

func TestGenerateDescriptorUseCase_MissingAccountKey(t *testing.T) {
	t.Parallel()

	useCase := keygenusecasebtc.NewGenerateDescriptorUseCase(
		btc.NewDescriptorService(&chaincfg.MainNetParams),
		&chaincfg.MainNetParams,
		&stubAuthRepo{},
		&stubAccountRepo{key: nil},
		nil,
		domainCoin.BTC,
		nil,
	)

	_, err := useCase.Generate(context.Background(), keygenusecase.GenerateDescriptorInput{
		AccountType: domainAccount.AccountTypeDeposit,
		AddressType: domainAddress.AddrTypeLegacy,
		IsChange:    false,
	})
	require.Error(t, err)
}

// TestGenerateDescriptorUseCase_MultisigWithKeygenKey tests that multisig descriptors
// include the keygen wallet's own key along with auth keys.
// This is the critical fix for issue #320 - ensuring 2-of-3 multisig instead of 2-of-2.
func TestGenerateDescriptorUseCase_MultisigWithKeygenKey(t *testing.T) {
	t.Parallel()

	testSeed := bytes.Repeat([]byte{0x01}, hdkeychain.RecommendedSeedLen)
	signers := map[domainAccount.AuthType]*domainAuth.AuthFullPubkey{
		domainAccount.AuthType1: {ExtendedPubKey: testDescriptorMainnetXpub},
		domainAccount.AuthType2: {ExtendedPubKey: newTestXpubFromSeed(t, 0x02)},
	}

	multiConfig := domainAccount.NewMultisigConfig(map[domainAccount.AccountType]map[int][]domainAccount.AuthType{
		domainAccount.AccountTypePayment: {
			2: {domainAccount.AuthType1, domainAccount.AuthType2},
		},
	})

	descriptorService := btc.NewDescriptorService(&chaincfg.MainNetParams)
	useCase := keygenusecasebtc.NewGenerateDescriptorUseCase(
		descriptorService,
		&chaincfg.MainNetParams,
		&stubAuthRepo{items: signers},
		&stubAccountRepo{},
		&stubSeedRepo{seed: testSeed},
		domainCoin.BTC,
		multiConfig,
	)

	// Test P2SH-SegWit (BIP49) descriptor generation
	output, err := useCase.Generate(context.Background(), keygenusecase.GenerateDescriptorInput{
		AccountType:  domainAccount.AccountTypePayment,
		AddressType:  domainAddress.AddrTypeP2shSegwit,
		IsChange:     false,
		RequiredSigs: 2,
	})
	require.NoError(t, err)
	require.True(t, output.IsMultisig)

	// Verify descriptor contains sh(wsh(sortedmulti(2, ...)))
	require.Contains(t, output.Descriptor, "sh(wsh(sortedmulti(2,")

	// Verify BIP49 path for P2SH-SegWit (payment account = 1)
	require.Contains(t, output.Descriptor, "/49'/0'/1]") // account=1 (payment)

	// Count occurrences of extended public keys in descriptor
	// Should have 3 xpubs: keygen + auth1 + auth2
	xpubCount := bytes.Count([]byte(output.Descriptor), []byte("xpub"))
	require.Equal(t, 3, xpubCount, "descriptor should contain 3 xpubs (keygen + auth1 + auth2)")
}

// TestGenerateDescriptorUseCase_MultisigMissingSeed tests error handling when seed is missing.
func TestGenerateDescriptorUseCase_MultisigMissingSeed(t *testing.T) {
	t.Parallel()

	signers := map[domainAccount.AuthType]*domainAuth.AuthFullPubkey{
		domainAccount.AuthType1: {ExtendedPubKey: testDescriptorMainnetXpub},
	}

	multiConfig := domainAccount.NewMultisigConfig(map[domainAccount.AccountType]map[int][]domainAccount.AuthType{
		domainAccount.AccountTypeDeposit: {
			2: {domainAccount.AuthType1},
		},
	})

	useCase := keygenusecasebtc.NewGenerateDescriptorUseCase(
		btc.NewDescriptorService(&chaincfg.MainNetParams),
		&chaincfg.MainNetParams,
		&stubAuthRepo{items: signers},
		&stubAccountRepo{},
		&stubSeedRepo{seed: nil}, // No seed available
		domainCoin.BTC,
		multiConfig,
	)

	_, err := useCase.Generate(context.Background(), keygenusecase.GenerateDescriptorInput{
		AccountType:  domainAccount.AccountTypeDeposit,
		AddressType:  domainAddress.AddrTypeBech32,
		IsChange:     false,
		RequiredSigs: 2,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "seed not found")
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

func (s *stubAuthRepo) GetOneByPurpose(
	authType domainAccount.AuthType,
	_ domainAuth.Purpose,
) (*domainAuth.AuthFullPubkey, error) {
	// For tests, ignore purpose and return the same result as GetOne
	return s.GetOne(authType)
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

type stubSeedRepo struct {
	seed []byte
	err  error
}

func (s *stubSeedRepo) GetOne(context.Context) (*domainKey.Seed, error) {
	if s.err != nil {
		return nil, s.err
	}
	if s.seed == nil {
		return nil, nil
	}
	return &domainKey.Seed{
		Seed: fmt.Sprintf("%x", s.seed),
	}, nil
}

func (*stubSeedRepo) Insert(context.Context, string) error { return nil }
