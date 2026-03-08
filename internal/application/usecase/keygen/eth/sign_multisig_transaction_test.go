package eth_test

import (
	"context"
	"encoding/hex"
	"errors"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	keygeneth "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/eth"
	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
	coldmocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold/mocks"
	storagemocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/transaction/mocks"
)

const (
	signTestSafeAddr = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	signTestToAddr   = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
	signTestFilePath = "/tmp/payment_multisig_550e8400-e29b-41d4-a716-446655440000_0.json"
)

// buildSignTestFile returns an unsigned multisig file with a correctly computed safeTxHash.
func buildSignTestFile(t *testing.T) *dtoeth.ETHMultisigTransactionFile {
	t.Helper()
	f := &dtoeth.ETHMultisigTransactionFile{
		Version:        1,
		TxType:         "unsigned",
		UUID:           "550e8400-e29b-41d4-a716-446655440000",
		SafeAddress:    signTestSafeAddr,
		To:             signTestToAddr,
		Value:          "100000000000000000",
		Data:           "0x",
		Operation:      0,
		SafeTxGas:      "0",
		BaseGas:        "0",
		GasPrice:       "0",
		GasToken:       "0x0000000000000000000000000000000000000000",
		RefundReceiver: "0x0000000000000000000000000000000000000000",
		Nonce:          "3",
		ChainID:        31337,
		Threshold:      1, // 1-of-1 so signing makes it complete
		Signatures:     []dtoeth.SignEntry{},
	}
	// Compute the real safeTxHash so the use case doesn't abort with ErrSafeTxHashMismatch.
	hash, err := keygeneth.ComputeSafeTxHash(f)
	require.NoError(t, err)
	f.SafeTxHash = "0x" + hex.EncodeToString(hash[:])
	return f
}

// buildSignTestFile2of2 returns an unsigned multisig file with threshold=2.
func buildSignTestFile2of2(t *testing.T) *dtoeth.ETHMultisigTransactionFile {
	t.Helper()
	f := buildSignTestFile(t)
	f.Threshold = 2
	return f
}

type signMultisigTxDeps struct {
	multisigRepo   *storagemocks.MockMultisigFileRepositorier
	accountKeyRepo *coldmocks.MockETHAccountKeyRepositorier
}

func newSignMultisigTxDeps(t *testing.T) *signMultisigTxDeps {
	t.Helper()
	return &signMultisigTxDeps{
		multisigRepo:   storagemocks.NewMockMultisigFileRepositorier(t),
		accountKeyRepo: coldmocks.NewMockETHAccountKeyRepositorier(t),
	}
}

func newSignMultisigTxUseCase(deps *signMultisigTxDeps) keygenusecase.SignMultisigTransactionUseCase {
	return keygeneth.NewSignMultisigTransactionUseCase(
		deps.multisigRepo,
		deps.accountKeyRepo,
	)
}

// newTestAccountKey returns a valid ETHAccountKey with a real BIP-44 xpriv.
// The signer address in the key is set to addr derived from the Anvil account 0 private key.
// The use case uses input.SignerAddress (not the derived key's address) for the SignEntry.
func newTestAccountKey(t *testing.T) *domainETH.ETHAccountKey {
	t.Helper()
	privKeyHex := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	privKeyBytes, err := hex.DecodeString(privKeyHex)
	require.NoError(t, err)
	privKey, err := crypto.ToECDSA(privKeyBytes)
	require.NoError(t, err)
	addr := crypto.PubkeyToAddress(privKey.PublicKey).Hex()

	return &domainETH.ETHAccountKey{
		Address:                addr,
		AccountExtendedPrivkey: newTestAccountXprivForSign(t),
		Idx:                    0,
	}
}

func TestSignMultisigTransaction_Success(t *testing.T) {
	t.Parallel()

	deps := newSignMultisigTxDeps(t)
	f := buildSignTestFile(t)

	// Use the Anvil account 0 private key.
	privKeyHex := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	privKeyBytes, err := hex.DecodeString(privKeyHex)
	require.NoError(t, err)
	privKey, err := crypto.ToECDSA(privKeyBytes)
	require.NoError(t, err)
	signerAddr := crypto.PubkeyToAddress(privKey.PublicKey).Hex()

	// We need the signing address to match the derived key from xpriv at index 0.
	// Rather than using HD derivation in the test, we use a known xpriv that derives
	// the Anvil account 0 key at m/44'/60'/0'/0/0.
	// This is a simplification: in tests we inject an xpriv and accept that the
	// derived key may not match signerAddr — we just verify the overall flow.
	// The important unit being tested is: read → hash verify → sign → write.
	key := newTestAccountKey(t)
	// Override the address to match what GetByAddress will be called with.
	key.Address = signerAddr

	newPath := "/tmp/payment_multisig_550e8400-e29b-41d4-a716-446655440000_1.json"

	deps.multisigRepo.EXPECT().ReadETHMultisigJSONFile(signTestFilePath).Return(f, nil).Once()
	deps.accountKeyRepo.EXPECT().GetByAddress(signerAddr).Return(key, nil).Once()
	deps.multisigRepo.EXPECT().CreateMultisigFilePath(mock.Anything, f.UUID, 1).Return(newPath).Once()
	deps.multisigRepo.EXPECT().WriteETHMultisigJSONFile(newPath, mock.Anything).Return(newPath, nil).Once()

	uc := newSignMultisigTxUseCase(deps)
	out, err := uc.Sign(context.Background(), keygenusecase.SignMultisigTransactionInput{
		FilePath:      signTestFilePath,
		SignerAddress: signerAddr,
	})

	require.NoError(t, err)
	assert.Equal(t, newPath, out.FilePath)
	assert.True(t, out.IsComplete, "1-of-1 threshold: signing completes the proposal")
	assert.Equal(t, 1, out.SignCount)
	assert.Equal(t, 1, out.Threshold)
}

func TestSignMultisigTransaction_NotComplete(t *testing.T) {
	t.Parallel()

	deps := newSignMultisigTxDeps(t)
	f := buildSignTestFile2of2(t)

	privKeyHex := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	privKeyBytes, err := hex.DecodeString(privKeyHex)
	require.NoError(t, err)
	privKey, err := crypto.ToECDSA(privKeyBytes)
	require.NoError(t, err)
	signerAddr := crypto.PubkeyToAddress(privKey.PublicKey).Hex()

	key := newTestAccountKey(t)
	key.Address = signerAddr

	newPath := "/tmp/payment_multisig_550e8400-e29b-41d4-a716-446655440000_1.json"

	deps.multisigRepo.EXPECT().ReadETHMultisigJSONFile(signTestFilePath).Return(f, nil).Once()
	deps.accountKeyRepo.EXPECT().GetByAddress(signerAddr).Return(key, nil).Once()
	deps.multisigRepo.EXPECT().CreateMultisigFilePath(mock.Anything, f.UUID, 1).Return(newPath).Once()
	deps.multisigRepo.EXPECT().WriteETHMultisigJSONFile(newPath, mock.Anything).Return(newPath, nil).Once()

	uc := newSignMultisigTxUseCase(deps)
	out, err := uc.Sign(context.Background(), keygenusecase.SignMultisigTransactionInput{
		FilePath:      signTestFilePath,
		SignerAddress: signerAddr,
	})

	require.NoError(t, err)
	assert.False(t, out.IsComplete, "2-of-2 with only 1 sig: not complete")
	assert.Equal(t, 1, out.SignCount)
	assert.Equal(t, 2, out.Threshold)
}

func TestSignMultisigTransaction_FileReadError(t *testing.T) {
	t.Parallel()

	deps := newSignMultisigTxDeps(t)
	deps.multisigRepo.EXPECT().ReadETHMultisigJSONFile(signTestFilePath).
		Return(nil, errors.New("file not found")).Once()

	uc := newSignMultisigTxUseCase(deps)
	_, err := uc.Sign(context.Background(), keygenusecase.SignMultisigTransactionInput{
		FilePath:      signTestFilePath,
		SignerAddress: signTestSafeAddr,
	})

	require.Error(t, err)
}

func TestSignMultisigTransaction_HashMismatch(t *testing.T) {
	t.Parallel()

	deps := newSignMultisigTxDeps(t)

	f := buildSignTestFile(t)
	f.SafeTxHash = "0x" + strings.Repeat("00", 32) // tampered hash

	deps.multisigRepo.EXPECT().ReadETHMultisigJSONFile(signTestFilePath).Return(f, nil).Once()

	uc := newSignMultisigTxUseCase(deps)
	_, err := uc.Sign(context.Background(), keygenusecase.SignMultisigTransactionInput{
		FilePath:      signTestFilePath,
		SignerAddress: signTestSafeAddr,
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, dtoeth.ErrSafeTxHashMismatch)
}

func TestSignMultisigTransaction_DuplicateSigner(t *testing.T) {
	t.Parallel()

	deps := newSignMultisigTxDeps(t)

	f := buildSignTestFile(t)
	// Signer already present in Signatures.
	f.Signatures = []dtoeth.SignEntry{
		{SignerAddress: signTestSafeAddr, SignatureHex: "0x" + strings.Repeat("ab", 65)},
	}
	// Add the signature so TxType can stay "unsigned" (threshold=1 but 1 sig present; adjust)
	f.Threshold = 2 // keep unsigned with 1 sig

	deps.multisigRepo.EXPECT().ReadETHMultisigJSONFile(signTestFilePath).Return(f, nil).Once()

	uc := newSignMultisigTxUseCase(deps)
	_, err := uc.Sign(context.Background(), keygenusecase.SignMultisigTransactionInput{
		FilePath:      signTestFilePath,
		SignerAddress: signTestSafeAddr,
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, dtoeth.ErrDuplicateSigner)
}

func TestSignMultisigTransaction_KeyNotFound(t *testing.T) {
	t.Parallel()

	deps := newSignMultisigTxDeps(t)
	f := buildSignTestFile(t)

	deps.multisigRepo.EXPECT().ReadETHMultisigJSONFile(signTestFilePath).Return(f, nil).Once()
	deps.accountKeyRepo.EXPECT().GetByAddress(signTestSafeAddr).
		Return(nil, errors.New("key not found")).Once()

	uc := newSignMultisigTxUseCase(deps)
	_, err := uc.Sign(context.Background(), keygenusecase.SignMultisigTransactionInput{
		FilePath:      signTestFilePath,
		SignerAddress: signTestSafeAddr,
	})

	require.Error(t, err)
}

func TestSignMultisigTransaction_WriteError(t *testing.T) {
	t.Parallel()

	deps := newSignMultisigTxDeps(t)
	f := buildSignTestFile(t)

	privKeyHex := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	privKeyBytes, err := hex.DecodeString(privKeyHex)
	require.NoError(t, err)
	privKey, err := crypto.ToECDSA(privKeyBytes)
	require.NoError(t, err)
	signerAddr := crypto.PubkeyToAddress(privKey.PublicKey).Hex()

	key := newTestAccountKey(t)
	key.Address = signerAddr

	newPath := "/tmp/payment_multisig_550e8400-e29b-41d4-a716-446655440000_1.json"

	deps.multisigRepo.EXPECT().ReadETHMultisigJSONFile(signTestFilePath).Return(f, nil).Once()
	deps.accountKeyRepo.EXPECT().GetByAddress(signerAddr).Return(key, nil).Once()
	deps.multisigRepo.EXPECT().CreateMultisigFilePath(mock.Anything, f.UUID, 1).Return(newPath).Once()
	deps.multisigRepo.EXPECT().WriteETHMultisigJSONFile(newPath, mock.Anything).
		Return("", errors.New("disk full")).Once()

	uc := newSignMultisigTxUseCase(deps)
	_, err = uc.Sign(context.Background(), keygenusecase.SignMultisigTransactionInput{
		FilePath:      signTestFilePath,
		SignerAddress: signerAddr,
	})

	require.Error(t, err)
}

// TestSignMultisigTransaction_SignatureAppended verifies that the generated signature
// is appended with the correct signer address and a v byte adjusted by +27.
func TestSignMultisigTransaction_SignatureAppended(t *testing.T) {
	t.Parallel()

	deps := newSignMultisigTxDeps(t)
	f := buildSignTestFile(t)

	privKeyHex := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	privKeyBytes, err := hex.DecodeString(privKeyHex)
	require.NoError(t, err)
	privKey, err := crypto.ToECDSA(privKeyBytes)
	require.NoError(t, err)
	signerAddr := crypto.PubkeyToAddress(privKey.PublicKey).Hex()

	key := newTestAccountKey(t)
	key.Address = signerAddr

	newPath := "/tmp/payment_multisig_550e8400-e29b-41d4-a716-446655440000_1.json"

	var capturedFile *dtoeth.ETHMultisigTransactionFile
	deps.multisigRepo.EXPECT().ReadETHMultisigJSONFile(signTestFilePath).Return(f, nil).Once()
	deps.accountKeyRepo.EXPECT().GetByAddress(signerAddr).Return(key, nil).Once()
	deps.multisigRepo.EXPECT().CreateMultisigFilePath(mock.Anything, f.UUID, 1).Return(newPath).Once()
	captureFn := mock.MatchedBy(func(ff *dtoeth.ETHMultisigTransactionFile) bool {
		capturedFile = ff
		return true
	})
	deps.multisigRepo.EXPECT().WriteETHMultisigJSONFile(newPath, captureFn).Return(newPath, nil).Once()

	uc := newSignMultisigTxUseCase(deps)
	_, err = uc.Sign(context.Background(), keygenusecase.SignMultisigTransactionInput{
		FilePath:      signTestFilePath,
		SignerAddress: signerAddr,
	})
	require.NoError(t, err)
	require.NotNil(t, capturedFile)

	require.Len(t, capturedFile.Signatures, 1)
	entry := capturedFile.Signatures[0]

	// Signer address must match.
	assert.Equal(t, signerAddr, entry.SignerAddress)

	// Signature must be 0x-prefixed 65-byte hex (130 hex chars after 0x).
	assert.True(t, strings.HasPrefix(entry.SignatureHex, "0x"))
	sigBytes, err := hex.DecodeString(strings.TrimPrefix(entry.SignatureHex, "0x"))
	require.NoError(t, err)
	assert.Len(t, sigBytes, 65, "signature must be 65 bytes")

	// v byte (last byte) must be 27 or 28 (Safe-compatible), not 0 or 1.
	v := sigBytes[64]
	assert.True(t, v == 27 || v == 28, "v byte must be 27 or 28; got %d", v)
}
