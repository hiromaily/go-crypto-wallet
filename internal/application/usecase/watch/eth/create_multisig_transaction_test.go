package eth_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	watchusecaseeth "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/eth"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	ethapiamocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/mocks"
	storagemocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/transaction/mocks"
	pkgeth "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth"
	pkguuid "github.com/hiromaily/go-crypto-wallet/pkg/uuid"
	uuidmocks "github.com/hiromaily/go-crypto-wallet/pkg/uuid/mocks"
)

const (
	msTestSafeAddr   = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	msTestToAddr     = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
	msTestSafeTxHash = "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab"
	msTestChainID    = uint64(31337)
	msTestUUIDStr    = "550e8400-e29b-41d4-a716-446655440000"
)

// msTestUUID is the parsed UUID for test expectations.
var msTestUUID = pkguuid.UUID{
	0x55, 0x0e, 0x84, 0x00, 0xe2, 0x9b, 0x41, 0xd4,
	0xa7, 0x16, 0x44, 0x66, 0x55, 0x44, 0x00, 0x00,
}

type ethCreateMultisigTxDeps struct {
	nonceReader  *ethapiamocks.MockSafeNonceReader
	hashComputer *ethapiamocks.MockSafeTxHashComputer
	multisigRepo *storagemocks.MockMultisigFileRepositorier
	uuidHandler  *uuidmocks.MockUUIDHandler
}

func newEthCreateMultisigTxDeps(t *testing.T) *ethCreateMultisigTxDeps {
	t.Helper()
	return &ethCreateMultisigTxDeps{
		nonceReader:  ethapiamocks.NewMockSafeNonceReader(t),
		hashComputer: ethapiamocks.NewMockSafeTxHashComputer(t),
		multisigRepo: storagemocks.NewMockMultisigFileRepositorier(t),
		uuidHandler:  uuidmocks.NewMockUUIDHandler(t),
	}
}

func newEthCreateMultisigTxUseCase(deps *ethCreateMultisigTxDeps) watchusecase.CreateETHMultisigTransactionUseCase {
	return watchusecaseeth.NewCreateETHMultisigTransactionUseCase(
		deps.nonceReader,
		deps.hashComputer,
		deps.multisigRepo,
		deps.uuidHandler,
		msTestChainID,
	)
}

func TestCreateETHMultisigTransaction(t *testing.T) {
	t.Parallel()

	const expectedFilePath = "/tmp/payment_multisig_" + msTestUUIDStr + "_0.json"

	tests := []struct {
		name    string
		input   watchusecase.CreateETHMultisigTransactionInput
		setup   func(deps *ethCreateMultisigTxDeps)
		wantErr bool
		checkFn func(t *testing.T, out watchusecase.CreateETHMultisigTransactionOutput)
	}{
		{
			name: "creates unsigned multisig proposal for payment",
			input: watchusecase.CreateETHMultisigTransactionInput{
				SafeAddress: msTestSafeAddr,
				To:          msTestToAddr,
				AmountEther: 0.1,
				Threshold:   2,
				ActionType:  "payment",
			},
			setup: func(deps *ethCreateMultisigTxDeps) {
				deps.uuidHandler.EXPECT().GenerateV4().Return(msTestUUID, nil).Once()
				deps.nonceReader.EXPECT().GetSafeNonce(
					mock.Anything, msTestSafeAddr,
				).Return(big.NewInt(0), nil).Once()
				deps.hashComputer.EXPECT().GetSafeTxHash(
					mock.Anything, msTestSafeAddr, msTestToAddr,
					mock.Anything, mock.Anything, mock.Anything, mock.Anything,
				).Return(msTestSafeTxHash, nil).Once()
				deps.multisigRepo.EXPECT().CreateMultisigFilePath(
					domainTx.ActionTypePayment, msTestUUIDStr, 0,
				).Return(expectedFilePath).Once()
				deps.multisigRepo.EXPECT().WriteETHMultisigJSONFile(
					expectedFilePath, mock.AnythingOfType("*eth.ETHMultisigTransactionFile"),
				).Return(expectedFilePath, nil).Once()
			},
			wantErr: false,
			checkFn: func(t *testing.T, out watchusecase.CreateETHMultisigTransactionOutput) {
				t.Helper()
				assert.Equal(t, expectedFilePath, out.FilePath)
				assert.Equal(t, msTestUUIDStr, out.UUID)
			},
		},
		{
			name: "creates unsigned multisig proposal for deposit",
			input: watchusecase.CreateETHMultisigTransactionInput{
				SafeAddress: msTestSafeAddr,
				To:          msTestToAddr,
				AmountEther: 1.5,
				Threshold:   1,
				ActionType:  "deposit",
			},
			setup: func(deps *ethCreateMultisigTxDeps) {
				const depositPath = "/tmp/deposit_multisig_" + msTestUUIDStr + "_0.json"
				deps.uuidHandler.EXPECT().GenerateV4().Return(msTestUUID, nil).Once()
				deps.nonceReader.EXPECT().GetSafeNonce(
					mock.Anything, msTestSafeAddr,
				).Return(big.NewInt(5), nil).Once()
				deps.hashComputer.EXPECT().GetSafeTxHash(
					mock.Anything, msTestSafeAddr, msTestToAddr,
					mock.Anything, mock.Anything, mock.Anything, mock.Anything,
				).Return(msTestSafeTxHash, nil).Once()
				deps.multisigRepo.EXPECT().CreateMultisigFilePath(
					domainTx.ActionTypeDeposit, msTestUUIDStr, 0,
				).Return(depositPath).Once()
				deps.multisigRepo.EXPECT().WriteETHMultisigJSONFile(
					depositPath, mock.AnythingOfType("*eth.ETHMultisigTransactionFile"),
				).Return(depositPath, nil).Once()
			},
			wantErr: false,
			checkFn: func(t *testing.T, out watchusecase.CreateETHMultisigTransactionOutput) {
				t.Helper()
				assert.NotEmpty(t, out.FilePath)
				assert.Equal(t, msTestUUIDStr, out.UUID)
			},
		},
		{
			name: "error on UUID generation failure",
			input: watchusecase.CreateETHMultisigTransactionInput{
				SafeAddress: msTestSafeAddr,
				To:          msTestToAddr,
				AmountEther: 0.1,
				Threshold:   2,
				ActionType:  "deposit",
			},
			setup: func(deps *ethCreateMultisigTxDeps) {
				deps.uuidHandler.EXPECT().GenerateV4().Return(pkguuid.UUID{}, errors.New("uuid error")).Once()
			},
			wantErr: true,
		},
		{
			name: "error on GetSafeNonce failure",
			input: watchusecase.CreateETHMultisigTransactionInput{
				SafeAddress: msTestSafeAddr,
				To:          msTestToAddr,
				AmountEther: 0.5,
				Threshold:   2,
				ActionType:  "transfer",
			},
			setup: func(deps *ethCreateMultisigTxDeps) {
				deps.uuidHandler.EXPECT().GenerateV4().Return(msTestUUID, nil).Once()
				deps.nonceReader.EXPECT().GetSafeNonce(
					mock.Anything, msTestSafeAddr,
				).Return(nil, errors.New("nonce fetch error")).Once()
			},
			wantErr: true,
		},
		{
			name: "error on GetSafeTxHash failure",
			input: watchusecase.CreateETHMultisigTransactionInput{
				SafeAddress: msTestSafeAddr,
				To:          msTestToAddr,
				AmountEther: 1.0,
				Threshold:   2,
				ActionType:  "payment",
			},
			setup: func(deps *ethCreateMultisigTxDeps) {
				deps.uuidHandler.EXPECT().GenerateV4().Return(msTestUUID, nil).Once()
				deps.nonceReader.EXPECT().GetSafeNonce(
					mock.Anything, msTestSafeAddr,
				).Return(big.NewInt(3), nil).Once()
				deps.hashComputer.EXPECT().GetSafeTxHash(
					mock.Anything, mock.Anything, mock.Anything,
					mock.Anything, mock.Anything, mock.Anything, mock.Anything,
				).Return("", errors.New("hash computation error")).Once()
			},
			wantErr: true,
		},
		{
			name: "error on WriteETHMultisigJSONFile failure",
			input: watchusecase.CreateETHMultisigTransactionInput{
				SafeAddress: msTestSafeAddr,
				To:          msTestToAddr,
				AmountEther: 0.1,
				Threshold:   2,
				ActionType:  "payment",
			},
			setup: func(deps *ethCreateMultisigTxDeps) {
				deps.uuidHandler.EXPECT().GenerateV4().Return(msTestUUID, nil).Once()
				deps.nonceReader.EXPECT().GetSafeNonce(
					mock.Anything, msTestSafeAddr,
				).Return(big.NewInt(0), nil).Once()
				deps.hashComputer.EXPECT().GetSafeTxHash(
					mock.Anything, mock.Anything, mock.Anything,
					mock.Anything, mock.Anything, mock.Anything, mock.Anything,
				).Return(msTestSafeTxHash, nil).Once()
				deps.multisigRepo.EXPECT().CreateMultisigFilePath(
					mock.Anything, mock.Anything, mock.Anything,
				).Return(expectedFilePath).Once()
				deps.multisigRepo.EXPECT().WriteETHMultisigJSONFile(
					mock.Anything, mock.AnythingOfType("*eth.ETHMultisigTransactionFile"),
				).Return("", errors.New("write error")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deps := newEthCreateMultisigTxDeps(t)
			if tt.setup != nil {
				tt.setup(deps)
			}

			uc := newEthCreateMultisigTxUseCase(deps)
			out, err := uc.Execute(context.Background(), tt.input)

			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.checkFn != nil {
				tt.checkFn(t, out)
			}
		})
	}
}

// TestCreateETHMultisigTransaction_FileContents verifies the written file has correct fields.
func TestCreateETHMultisigTransaction_FileContents(t *testing.T) {
	t.Parallel()

	const expectedPath = "/tmp/payment_multisig_" + msTestUUIDStr + "_0.json"

	deps := newEthCreateMultisigTxDeps(t)
	deps.uuidHandler.EXPECT().GenerateV4().Return(msTestUUID, nil).Once()
	deps.nonceReader.EXPECT().GetSafeNonce(mock.Anything, msTestSafeAddr).Return(big.NewInt(7), nil).Once()
	deps.hashComputer.EXPECT().GetSafeTxHash(
		mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
	).Return(msTestSafeTxHash, nil).Once()
	deps.multisigRepo.EXPECT().CreateMultisigFilePath(
		domainTx.ActionTypePayment, msTestUUIDStr, 0,
	).Return(expectedPath).Once()

	// Capture the file written so we can inspect its fields.
	var capturedFile *dtoeth.ETHMultisigTransactionFile
	deps.multisigRepo.EXPECT().WriteETHMultisigJSONFile(
		expectedPath, mock.AnythingOfType("*eth.ETHMultisigTransactionFile"),
	).RunAndReturn(func(path string, f *dtoeth.ETHMultisigTransactionFile) (string, error) {
		capturedFile = f
		return path, nil
	}).Once()

	uc := newEthCreateMultisigTxUseCase(deps)
	out, err := uc.Execute(context.Background(), watchusecase.CreateETHMultisigTransactionInput{
		SafeAddress: msTestSafeAddr,
		To:          msTestToAddr,
		AmountEther: 0.1,
		Threshold:   2,
		ActionType:  "payment",
	})
	require.NoError(t, err)
	assert.Equal(t, expectedPath, out.FilePath)
	assert.Equal(t, msTestUUIDStr, out.UUID)

	require.NotNil(t, capturedFile)
	assert.Equal(t, 1, capturedFile.Version)
	assert.Equal(t, "unsigned", capturedFile.TxType)
	assert.Equal(t, msTestUUIDStr, capturedFile.UUID)
	assert.Equal(t, msTestSafeAddr, capturedFile.SafeAddress)
	assert.Equal(t, msTestToAddr, capturedFile.To)
	// 0.1 ETH in Wei via float64 conversion (matches pkgeth.FromFloatEther behaviour)
	assert.Equal(t, pkgeth.FromFloatEther(0.1).String(), capturedFile.Value)
	assert.Equal(t, "0x", capturedFile.Data)
	assert.Equal(t, uint8(0), capturedFile.Operation)
	assert.Equal(t, "7", capturedFile.Nonce)
	assert.Equal(t, msTestChainID, capturedFile.ChainID)
	assert.Equal(t, msTestSafeTxHash, capturedFile.SafeTxHash)
	assert.Equal(t, 2, capturedFile.Threshold)
	assert.Empty(t, capturedFile.Signatures)
}
