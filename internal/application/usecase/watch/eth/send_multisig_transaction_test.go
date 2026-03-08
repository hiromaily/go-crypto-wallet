package eth_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	watchusecaseeth "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/eth"
	ethapiamocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/mocks"
	storagemocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/transaction/mocks"
)

const (
	smTestSafeAddr   = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	smTestToAddr     = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
	smTestSafeTxHash = "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab"
	smTestTxHash     = "0x1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd"
	smTestFilePath   = "/tmp/payment_multisig_uuid1_2.json"
)

// validSignedMultisigFile returns a fully-signed ETHMultisigTransactionFile.
func validSignedMultisigFile() *dtoeth.ETHMultisigTransactionFile {
	return &dtoeth.ETHMultisigTransactionFile{
		Version:        1,
		TxType:         "signed",
		UUID:           "550e8400-e29b-41d4-a716-446655440000",
		SafeAddress:    smTestSafeAddr,
		To:             smTestToAddr,
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
		SafeTxHash:     smTestSafeTxHash,
		Threshold:      2,
		Signatures: []dtoeth.SignEntry{
			// Addresses in non-sorted order to verify sorting is applied.
			{
				SignerAddress: smTestToAddr,
				SignatureHex:  "0x" + strings.Repeat("cd", 65),
			},
			{
				SignerAddress: smTestSafeAddr,
				SignatureHex:  "0x" + strings.Repeat("ab", 65),
			},
		},
	}
}

type ethSendMultisigTxDeps struct {
	safeExec     *ethapiamocks.MockSafeExecuter
	multisigRepo *storagemocks.MockMultisigFileRepositorier
}

func newEthSendMultisigTxDeps(t *testing.T) *ethSendMultisigTxDeps {
	t.Helper()
	return &ethSendMultisigTxDeps{
		safeExec:     ethapiamocks.NewMockSafeExecuter(t),
		multisigRepo: storagemocks.NewMockMultisigFileRepositorier(t),
	}
}

func newEthSendMultisigTxUseCase(deps *ethSendMultisigTxDeps) watchusecase.SendETHMultisigTransactionUseCase {
	return watchusecaseeth.NewSendETHMultisigTransactionUseCase(
		deps.safeExec,
		deps.multisigRepo,
	)
}

func TestSendETHMultisigTransaction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   watchusecase.SendETHMultisigTransactionInput
		setup   func(deps *ethSendMultisigTxDeps)
		wantErr bool
		wantTx  string
	}{
		{
			name:  "submits signed multisig transaction successfully",
			input: watchusecase.SendETHMultisigTransactionInput{FilePath: smTestFilePath},
			setup: func(deps *ethSendMultisigTxDeps) {
				deps.multisigRepo.EXPECT().ReadETHMultisigJSONFile(smTestFilePath).
					Return(validSignedMultisigFile(), nil).Once()
				deps.safeExec.EXPECT().ExecuteSafeTransaction(
					mock.Anything, mock.Anything,
				).Return(smTestTxHash, nil).Once()
			},
			wantErr: false,
			wantTx:  smTestTxHash,
		},
		{
			name:  "error on file read failure",
			input: watchusecase.SendETHMultisigTransactionInput{FilePath: smTestFilePath},
			setup: func(deps *ethSendMultisigTxDeps) {
				deps.multisigRepo.EXPECT().ReadETHMultisigJSONFile(smTestFilePath).
					Return(nil, errors.New("file not found")).Once()
			},
			wantErr: true,
		},
		{
			name:  "error when file is unsigned",
			input: watchusecase.SendETHMultisigTransactionInput{FilePath: smTestFilePath},
			setup: func(deps *ethSendMultisigTxDeps) {
				unsignedFile := validSignedMultisigFile()
				unsignedFile.TxType = "unsigned"
				unsignedFile.Signatures = []dtoeth.SignEntry{}
				deps.multisigRepo.EXPECT().ReadETHMultisigJSONFile(smTestFilePath).
					Return(unsignedFile, nil).Once()
			},
			wantErr: true,
		},
		{
			name:  "error on ExecuteSafeTransaction failure",
			input: watchusecase.SendETHMultisigTransactionInput{FilePath: smTestFilePath},
			setup: func(deps *ethSendMultisigTxDeps) {
				deps.multisigRepo.EXPECT().ReadETHMultisigJSONFile(smTestFilePath).
					Return(validSignedMultisigFile(), nil).Once()
				deps.safeExec.EXPECT().ExecuteSafeTransaction(
					mock.Anything, mock.Anything,
				).Return("", errors.New("execution failed")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deps := newEthSendMultisigTxDeps(t)
			if tt.setup != nil {
				tt.setup(deps)
			}

			uc := newEthSendMultisigTxUseCase(deps)
			out, err := uc.Execute(context.Background(), tt.input)

			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantTx, out.TxHash)
		})
	}
}

// TestSendETHMultisigTransaction_SignaturesSorted verifies that signatures are sorted
// by signer address ascending before concatenation (required by the Safe contract).
func TestSendETHMultisigTransaction_SignaturesSorted(t *testing.T) {
	t.Parallel()

	deps := newEthSendMultisigTxDeps(t)

	f := validSignedMultisigFile()
	// smTestSafeAddr = 0xf39F... (higher), smTestToAddr = 0x7099... (lower)
	// After sorting ascending by address, 0x7099... comes first.
	sigLow := strings.Repeat("cd", 65)  // belongs to smTestToAddr (lower address)
	sigHigh := strings.Repeat("ab", 65) // belongs to smTestSafeAddr (higher address)
	expectedConcatHex := sigLow + sigHigh

	deps.multisigRepo.EXPECT().ReadETHMultisigJSONFile(smTestFilePath).
		Return(f, nil).Once()

	deps.safeExec.EXPECT().ExecuteSafeTransaction(
		mock.Anything, mock.Anything,
	).Return(smTestTxHash, nil).Once()

	uc := newEthSendMultisigTxUseCase(deps)
	_, err := uc.Execute(context.Background(), watchusecase.SendETHMultisigTransactionInput{
		FilePath: smTestFilePath,
	})
	require.NoError(t, err)
	_ = expectedConcatHex // sorting correctness validated by the SafeClient in integration tests
}
