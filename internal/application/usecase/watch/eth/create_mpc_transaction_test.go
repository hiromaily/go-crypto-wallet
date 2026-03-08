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
	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	watchusecaseeth "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/eth"
	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
	ethapiamocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/mocks"
	storagemocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/transaction/mocks"
	pkgeth "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth"
	pkguuid "github.com/hiromaily/go-crypto-wallet/pkg/uuid"
	uuidmocks "github.com/hiromaily/go-crypto-wallet/pkg/uuid/mocks"
)

const (
	mpcFromAddr   = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	mpcToAddr     = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
	mpcUUIDStr    = "550e8400-e29b-41d4-a716-446655440001"
	mpcTestTxHash = "0x" + "aabbccdd" + "aabbccddaabbccddaabbccddaabbccddaabbccddaabbccddaabbccdd"
	// mpcTestTxHex is a pre-encoded EIP-1559 tx used as a fixture in tests.
	mpcTestTxHex = "0x02f856831e694080850218711a00850218711a00825208" +
		"940000000000000000000000000000000000000001880de0b6b3a764000080c0"
)

var mpcTestUUID = pkguuid.UUID{
	0x55, 0x0e, 0x84, 0x00, 0xe2, 0x9b, 0x41, 0xd4,
	0xa7, 0x16, 0x44, 0x66, 0x55, 0x44, 0x00, 0x01,
}

type createMPCTxDeps struct {
	txCreator   *ethapiamocks.MockTxCreator
	mpcFileRepo *storagemocks.MockMPCFileRepositorier
	uuidHandler *uuidmocks.MockUUIDHandler
}

func newCreateMPCTxDeps(t *testing.T) *createMPCTxDeps {
	t.Helper()
	return &createMPCTxDeps{
		txCreator:   ethapiamocks.NewMockTxCreator(t),
		mpcFileRepo: storagemocks.NewMockMPCFileRepositorier(t),
		uuidHandler: uuidmocks.NewMockUUIDHandler(t),
	}
}

func newCreateMPCTxUseCase(deps *createMPCTxDeps) watchusecase.CreateMPCTransactionUseCase {
	return watchusecaseeth.NewCreateMPCTransactionUseCase(
		deps.txCreator, deps.mpcFileRepo, deps.uuidHandler)
}

func validCreateMPCInput() watchusecase.CreateMPCTransactionInput {
	return watchusecase.CreateMPCTransactionInput{
		FromAddress: mpcFromAddr,
		ToAddress:   mpcToAddr,
		AmountEther: 0.1,
		ActionType:  "payment",
		Threshold:   2,
		PartyIDs:    []string{"node-1", "node-2", "node-3"},
	}
}

func mockRawTx() *domainETH.RawTx {
	return &domainETH.RawTx{
		UUID:  "tx-uuid",
		From:  mpcFromAddr,
		To:    mpcToAddr,
		Value: *new(big.Int).SetUint64(100000000000000000),
		Nonce: 5,
		TxHex: mpcTestTxHex,
		Hash:  mpcTestTxHash,
	}
}

func mockTxParams() *apieth.TxCreateParams {
	return &apieth.TxCreateParams{
		UUID:                 "tx-uuid",
		FromAddress:          mpcFromAddr,
		ToAddress:            mpcToAddr,
		Amount:               100000000000000000,
		Fee:                  21000 * 2000000000,
		GasLimit:             21000,
		Nonce:                5,
		EthTxType:            2,
		ChainID:              31337,
		MaxFeePerGas:         2000000000,
		MaxPriorityFeePerGas: 1000000000,
	}
}

func TestCreateMPCTransaction_Success(t *testing.T) {
	t.Parallel()

	const expectedPath = "/tmp/payment_mpc_" + mpcUUIDStr + ".json"
	deps := newCreateMPCTxDeps(t)

	deps.uuidHandler.EXPECT().GenerateV4().Return(mpcTestUUID, nil).Once()
	deps.txCreator.EXPECT().
		CreateRawTransactionEIP1559(mock.Anything, mpcFromAddr, mpcToAddr, mock.Anything, 0).
		Return(mockRawTx(), mockTxParams(), nil).Once()
	deps.mpcFileRepo.EXPECT().
		WriteETHMPCJSONFile(mock.AnythingOfType("*eth.ETHMPCTransactionFile"), false).
		Return(expectedPath, nil).Once()

	uc := newCreateMPCTxUseCase(deps)
	out, err := uc.Execute(context.Background(), validCreateMPCInput())

	require.NoError(t, err)
	assert.Equal(t, expectedPath, out.FilePath)
	assert.Equal(t, mpcUUIDStr, out.UUID)
	assert.Equal(t, mpcTestTxHash, out.TxHash)
}

func TestCreateMPCTransaction_FileContents(t *testing.T) {
	t.Parallel()

	deps := newCreateMPCTxDeps(t)

	deps.uuidHandler.EXPECT().GenerateV4().Return(mpcTestUUID, nil).Once()
	deps.txCreator.EXPECT().
		CreateRawTransactionEIP1559(mock.Anything, mpcFromAddr, mpcToAddr, mock.Anything, 0).
		Return(mockRawTx(), mockTxParams(), nil).Once()

	var capturedFile *dtoeth.ETHMPCTransactionFile
	deps.mpcFileRepo.EXPECT().
		WriteETHMPCJSONFile(mock.AnythingOfType("*eth.ETHMPCTransactionFile"), false).
		RunAndReturn(func(f *dtoeth.ETHMPCTransactionFile, _ bool) (string, error) {
			capturedFile = f
			return "/tmp/test.json", nil
		}).Once()

	uc := newCreateMPCTxUseCase(deps)
	_, err := uc.Execute(context.Background(), validCreateMPCInput())
	require.NoError(t, err)

	require.NotNil(t, capturedFile)
	assert.Equal(t, 1, capturedFile.Version)
	assert.Equal(t, "unsigned", capturedFile.TxType)
	assert.Equal(t, mpcUUIDStr, capturedFile.UUID)
	assert.Equal(t, "payment", capturedFile.ActionType)
	assert.Equal(t, mpcFromAddr, capturedFile.From)
	assert.Equal(t, mpcToAddr, capturedFile.To)
	assert.Equal(t, pkgeth.FromFloatEther(0.1).String(), capturedFile.Value)
	assert.Equal(t, uint64(5), capturedFile.Nonce)
	assert.Equal(t, uint64(21000), capturedFile.GasLimit)
	assert.Equal(t, uint64(31337), capturedFile.ChainID)
	assert.Equal(t, mpcTestTxHash, capturedFile.TxHash)
	assert.Equal(t, mpcTestTxHex, capturedFile.RawTxHex)
	assert.Equal(t, 2, capturedFile.Threshold)
	assert.Equal(t, []string{"node-1", "node-2", "node-3"}, capturedFile.PartyIDs)
	assert.Empty(t, capturedFile.SignedTxHex)
}

func TestCreateMPCTransaction_InvalidFromAddress(t *testing.T) {
	t.Parallel()

	deps := newCreateMPCTxDeps(t)
	uc := newCreateMPCTxUseCase(deps)

	input := validCreateMPCInput()
	input.FromAddress = "not-a-valid-address"

	_, err := uc.Execute(context.Background(), input)

	require.Error(t, err)
	assert.ErrorIs(t, err, dtoeth.ErrInvalidMPCAddress)
}

func TestCreateMPCTransaction_InvalidToAddress(t *testing.T) {
	t.Parallel()

	deps := newCreateMPCTxDeps(t)
	uc := newCreateMPCTxUseCase(deps)

	input := validCreateMPCInput()
	input.ToAddress = "0xNOT_HEX"

	_, err := uc.Execute(context.Background(), input)

	require.Error(t, err)
	assert.ErrorIs(t, err, dtoeth.ErrInvalidMPCAddress)
}

func TestCreateMPCTransaction_TxCreatorError(t *testing.T) {
	t.Parallel()

	deps := newCreateMPCTxDeps(t)
	txErr := errors.New("node connection refused")

	deps.uuidHandler.EXPECT().GenerateV4().Return(mpcTestUUID, nil).Once()
	deps.txCreator.EXPECT().
		CreateRawTransactionEIP1559(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil, txErr).Once()

	uc := newCreateMPCTxUseCase(deps)
	_, err := uc.Execute(context.Background(), validCreateMPCInput())

	require.Error(t, err)
	assert.ErrorIs(t, err, txErr)
}

func TestCreateMPCTransaction_WriteFileError(t *testing.T) {
	t.Parallel()

	deps := newCreateMPCTxDeps(t)
	writeErr := errors.New("disk full")

	deps.uuidHandler.EXPECT().GenerateV4().Return(mpcTestUUID, nil).Once()
	deps.txCreator.EXPECT().
		CreateRawTransactionEIP1559(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(mockRawTx(), mockTxParams(), nil).Once()
	deps.mpcFileRepo.EXPECT().
		WriteETHMPCJSONFile(mock.Anything, false).
		Return("", writeErr).Once()

	uc := newCreateMPCTxUseCase(deps)
	_, err := uc.Execute(context.Background(), validCreateMPCInput())

	require.Error(t, err)
	assert.ErrorIs(t, err, writeErr)
}
