package btc_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	portsWallet "github.com/hiromaily/go-crypto-wallet/internal/application/ports/wallet"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	keygenusecasebtc "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/btc"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainAuth "github.com/hiromaily/go-crypto-wallet/internal/domain/auth"
	domainBitcoin "github.com/hiromaily/go-crypto-wallet/internal/domain/bitcoin"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	apibtcimpl "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/btc/btc"
	coldmocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold/mocks"
	infraKey "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key"
)

// createTestCoinStrategy creates a BTC coin strategy for testing
func createTestCoinStrategy(
	t *testing.T,
	conf *chaincfg.Params,
) portsWallet.CoinKeyStrategy {
	t.Helper()
	strategy, err := infraKey.CreateCoinKeyStrategy(domainCoin.BTC, conf)
	require.NoError(t, err, "Failed to create test coin strategy")
	return strategy
}

const (
	// Test account-level xpub/xpriv pair derived from deterministic seed (bytes.Repeat([]byte{0x01}, 32))
	// Derivation path: m/84'/0'/0' (BIP84 account 0 on mainnet)
	testDescriptorMainnetXpub = "xpub6CFtfy4QXsEUW5CtgE7mZe1Lvs15Yw7ctjdyaDRy89JdhtyM1wFf8uY2BdyJ3JmAFfrHdw77h" +
		"Eit1ebVXxB2dytGAvq9mmQJ2c83G1q8P7A"
	testDescriptorMainnetXpriv = "xprv9yGYGTXWhVgBHb8RaCamCW4cNqAb9UPmXWiNmq2MZomeq6eCUPwQb7DYLQ2prX4aFBTpF8Cpq" +
		"abzQVCeuhBNJo9YvKHjfnJBCeM9Vk8mN5H"
)

func TestNewGenerateDescriptorUseCase(t *testing.T) {
	t.Parallel()

	coinStrategy := createTestCoinStrategy(t, &chaincfg.MainNetParams)
	useCase := keygenusecasebtc.NewGenerateDescriptorUseCase(
		apibtcimpl.NewDescriptorService(&chaincfg.MainNetParams),
		&chaincfg.MainNetParams,
		nil,
		nil,
		nil,
		domainCoin.BTC,
		nil,
		coinStrategy,
	)
	require.NotNil(t, useCase)
}

func TestGenerateDescriptorUseCase_SingleSig(t *testing.T) {
	t.Parallel()

	descriptorService := apibtcimpl.NewDescriptorService(&chaincfg.MainNetParams)
	xpriv := testDescriptorMainnetXpriv
	accountKey := &domainBitcoin.BTCAccountKey{
		FullPublicKey:          testDescriptorMainnetXpub,
		KeyType:                string(domainKey.KeyTypeBIP84),
		Account:                domainAccount.AccountTypeDeposit,
		AccountExtendedPrivkey: &xpriv,
	}
	accountRepo := coldmocks.NewMockBTCAccountKeyRepositorier(t)
	accountRepo.EXPECT().GetOneMaxID(domainAccount.AccountTypeDeposit).Return(accountKey, nil)
	authRepo := coldmocks.NewMockAuthFullPubkeyRepositorier(t)

	coinStrategy := createTestCoinStrategy(t, &chaincfg.MainNetParams)
	useCase := keygenusecasebtc.NewGenerateDescriptorUseCase(
		descriptorService,
		&chaincfg.MainNetParams,
		authRepo,
		accountRepo,
		nil,
		domainCoin.BTC,
		nil,
		coinStrategy,
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
	xpub2 := newTestXpubFromSeed(t, 0x02)

	multiConfig := domainAccount.NewMultisigConfig(map[domainAccount.AccountType]map[int][]domainAccount.AuthType{
		domainAccount.AccountTypeDeposit: {
			2: {domainAccount.AuthType1, domainAccount.AuthType2},
		},
	})

	authRepo := coldmocks.NewMockAuthFullPubkeyRepositorier(t)
	authRepo.EXPECT().GetOneByPurpose(domainAccount.AuthType1, domainAuth.PurposeBIP84).Return(
		&domainAuth.AuthFullPubkey{ExtendedPubKey: testDescriptorMainnetXpub}, nil)
	authRepo.EXPECT().GetOneByPurpose(domainAccount.AuthType2, domainAuth.PurposeBIP84).Return(
		&domainAuth.AuthFullPubkey{ExtendedPubKey: xpub2}, nil)
	accountRepo := coldmocks.NewMockBTCAccountKeyRepositorier(t)
	seedRepo := coldmocks.NewMockSeedRepositorier(t)
	seedRepo.EXPECT().GetOne(mock.Anything).Return(
		&domainKey.Seed{Seed: hex.EncodeToString(testSeed)}, nil)

	descriptorService := apibtcimpl.NewDescriptorService(&chaincfg.MainNetParams)
	coinStrategy := createTestCoinStrategy(t, &chaincfg.MainNetParams)
	useCase := keygenusecasebtc.NewGenerateDescriptorUseCase(
		descriptorService,
		&chaincfg.MainNetParams,
		authRepo,
		accountRepo,
		seedRepo,
		domainCoin.BTC,
		multiConfig,
		coinStrategy,
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
	xpub3 := newTestXpubFromSeed(t, 0x03)

	multiConfig := domainAccount.NewMultisigConfig(map[domainAccount.AccountType]map[int][]domainAccount.AuthType{
		domainAccount.AccountTypeDeposit: {
			2: {domainAccount.AuthType1, domainAccount.AuthType2},
		},
	})

	authRepo := coldmocks.NewMockAuthFullPubkeyRepositorier(t)
	authRepo.EXPECT().GetOneByPurpose(domainAccount.AuthType1, domainAuth.PurposeBIP86).Return(
		&domainAuth.AuthFullPubkey{ExtendedPubKey: testDescriptorMainnetXpub}, nil)
	authRepo.EXPECT().GetOneByPurpose(domainAccount.AuthType2, domainAuth.PurposeBIP86).Return(
		&domainAuth.AuthFullPubkey{ExtendedPubKey: xpub3}, nil)
	accountRepo := coldmocks.NewMockBTCAccountKeyRepositorier(t)
	seedRepo := coldmocks.NewMockSeedRepositorier(t)
	seedRepo.EXPECT().GetOne(mock.Anything).Return(
		&domainKey.Seed{Seed: hex.EncodeToString(testSeed)}, nil)

	coinStrategy := createTestCoinStrategy(t, &chaincfg.MainNetParams)
	useCase := keygenusecasebtc.NewGenerateDescriptorUseCase(
		apibtcimpl.NewDescriptorService(&chaincfg.MainNetParams),
		&chaincfg.MainNetParams,
		authRepo,
		accountRepo,
		seedRepo,
		domainCoin.BTC,
		multiConfig,
		coinStrategy,
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

	accountRepo := coldmocks.NewMockBTCAccountKeyRepositorier(t)
	accountRepo.EXPECT().GetOneMaxID(domainAccount.AccountTypeDeposit).Return((*domainBitcoin.BTCAccountKey)(nil), nil)
	authRepo := coldmocks.NewMockAuthFullPubkeyRepositorier(t)

	coinStrategy := createTestCoinStrategy(t, &chaincfg.MainNetParams)
	useCase := keygenusecasebtc.NewGenerateDescriptorUseCase(
		apibtcimpl.NewDescriptorService(&chaincfg.MainNetParams),
		&chaincfg.MainNetParams,
		authRepo,
		accountRepo,
		nil,
		domainCoin.BTC,
		nil,
		coinStrategy,
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
	xpub2 := newTestXpubFromSeed(t, 0x02)

	multiConfig := domainAccount.NewMultisigConfig(map[domainAccount.AccountType]map[int][]domainAccount.AuthType{
		domainAccount.AccountTypePayment: {
			2: {domainAccount.AuthType1, domainAccount.AuthType2},
		},
	})

	authRepo := coldmocks.NewMockAuthFullPubkeyRepositorier(t)
	authRepo.EXPECT().GetOneByPurpose(domainAccount.AuthType1, domainAuth.PurposeBIP49).Return(
		&domainAuth.AuthFullPubkey{ExtendedPubKey: testDescriptorMainnetXpub}, nil)
	authRepo.EXPECT().GetOneByPurpose(domainAccount.AuthType2, domainAuth.PurposeBIP49).Return(
		&domainAuth.AuthFullPubkey{ExtendedPubKey: xpub2}, nil)
	accountRepo := coldmocks.NewMockBTCAccountKeyRepositorier(t)
	seedRepo := coldmocks.NewMockSeedRepositorier(t)
	seedRepo.EXPECT().GetOne(mock.Anything).Return(
		&domainKey.Seed{Seed: hex.EncodeToString(testSeed)}, nil)

	descriptorService := apibtcimpl.NewDescriptorService(&chaincfg.MainNetParams)
	coinStrategy := createTestCoinStrategy(t, &chaincfg.MainNetParams)
	useCase := keygenusecasebtc.NewGenerateDescriptorUseCase(
		descriptorService,
		&chaincfg.MainNetParams,
		authRepo,
		accountRepo,
		seedRepo,
		domainCoin.BTC,
		multiConfig,
		coinStrategy,
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
	xpubCount := strings.Count(output.Descriptor, "xpub")
	require.Equal(t, 3, xpubCount, "descriptor should contain 3 xpubs (keygen + auth1 + auth2)")
}

// TestGenerateDescriptorUseCase_MultisigMissingSeed tests error handling when seed is missing.
func TestGenerateDescriptorUseCase_MultisigMissingSeed(t *testing.T) {
	t.Parallel()

	multiConfig := domainAccount.NewMultisigConfig(map[domainAccount.AccountType]map[int][]domainAccount.AuthType{
		domainAccount.AccountTypeDeposit: {
			2: {domainAccount.AuthType1},
		},
	})

	authRepo := coldmocks.NewMockAuthFullPubkeyRepositorier(t)
	accountRepo := coldmocks.NewMockBTCAccountKeyRepositorier(t)
	seedRepo := coldmocks.NewMockSeedRepositorier(t)
	seedRepo.EXPECT().GetOne(mock.Anything).Return((*domainKey.Seed)(nil), nil)

	coinStrategy := createTestCoinStrategy(t, &chaincfg.MainNetParams)
	useCase := keygenusecasebtc.NewGenerateDescriptorUseCase(
		apibtcimpl.NewDescriptorService(&chaincfg.MainNetParams),
		&chaincfg.MainNetParams,
		authRepo,
		accountRepo,
		seedRepo,
		domainCoin.BTC,
		multiConfig,
		coinStrategy,
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

func newTestXpubFromSeed(t *testing.T, seedByte byte) string {
	t.Helper()

	seed := bytes.Repeat([]byte{seedByte}, hdkeychain.RecommendedSeedLen)
	master, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	require.NoError(t, err)

	xpub, err := master.Neuter()
	require.NoError(t, err)

	return xpub.String()
}
