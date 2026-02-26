package eth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	watchusecaseeth "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/eth"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	ethapiamocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/mocks"
	repomocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch/mocks"
	storagemocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/transaction/mocks"
)

// ethSendTxTestDependencies holds all mock dependencies for testing
type ethSendTxTestDependencies struct {
	ethClient    *ethapiamocks.MockTxSender
	txDetailRepo *repomocks.MockETHDetailTXRepositorier
	txFileRepo   *storagemocks.MockTransactionFileRepositorier
}

func newEthSendTxTestDeps(t *testing.T) *ethSendTxTestDependencies {
	t.Helper()
	return &ethSendTxTestDependencies{
		ethClient:    ethapiamocks.NewMockTxSender(t),
		txDetailRepo: repomocks.NewMockETHDetailTXRepositorier(t),
		txFileRepo:   storagemocks.NewMockTransactionFileRepositorier(t),
	}
}

func newEthSendTxUseCase(deps *ethSendTxTestDependencies) watchusecase.SendTransactionUseCase {
	return watchusecaseeth.NewSendTransactionUseCase(
		deps.ethClient,
		deps.txDetailRepo,
		deps.txFileRepo,
	)
}

// ---- Helper values ----

const (
	testETHSignedFilePath = "deposit_1_signed_1_1534744535097796209.json"
	testETHSignedTxHex    = "0x02f870820539808501" +
		"2a05f20085012a05f200825208940deadbeef1234567890deadbeef1234567890abcd88" +
		"0de0b6b3a764000080c001a07f7e5b4dabe9e3f7f8e0e3e5e2e9d8c7a6f5e4d3" +
		"c2b1a098765432101234567a0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0"
	testETHSentTxHash = "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	testETHUUID       = "01234567-89ab-cdef-0123-456789abcdef"
)

func testSignedETHFile() *dtoeth.ETHTransactionFile {
	return &dtoeth.ETHTransactionFile{
		Version:              1,
		TxType:               string(domainTx.TxTypeSigned),
		EthTxType:            dtoeth.EthTxTypeEIP1559,
		UUID:                 testETHUUID,
		ChainID:              31337,
		Nonce:                0,
		From:                 testETHFromAddr,
		To:                   testETHToAddr,
		Value:                "999000000000000000",
		Gas:                  21000,
		MaxFeePerGas:         "20000000000",
		MaxPriorityFeePerGas: "2000000000",
		RawTxHex:             testETHRawTx,
		SignedTxHex:          testETHSignedTxHex,
	}
}

// ---- Tests ----

func TestETHSendTransactionUseCase_Execute_InvalidFilePath(t *testing.T) {
	t.Parallel()
	deps := newEthSendTxTestDeps(t)
	useCase := newEthSendTxUseCase(deps)

	deps.txFileRepo.EXPECT().
		ValidateFilePath("invalid-path", domainTx.TxTypeSigned).
		Return(domainTx.ActionType(""), domainTx.TxType(""), int64(0), 0,
			errors.New("invalid file path format"))

	output, err := useCase.Execute(context.Background(), watchusecase.SendTransactionInput{
		FilePath: "invalid-path",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "fail to call txFileRepo.ValidateFilePath()")
	assert.Empty(t, output.TxID)
}

func TestETHSendTransactionUseCase_Execute_ReadJSONError(t *testing.T) {
	t.Parallel()
	deps := newEthSendTxTestDeps(t)
	useCase := newEthSendTxUseCase(deps)

	deps.txFileRepo.EXPECT().
		ValidateFilePath(testETHSignedFilePath, domainTx.TxTypeSigned).
		Return(domainTx.ActionTypeDeposit, domainTx.TxTypeSigned, int64(1), 1, nil)
	deps.txFileRepo.EXPECT().
		ReadETHJSONFile(testETHSignedFilePath).
		Return(nil, errors.New("file not found"))

	output, err := useCase.Execute(context.Background(), watchusecase.SendTransactionInput{
		FilePath: testETHSignedFilePath,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "fail to call txFileRepo.ReadETHJSONFile()")
	assert.Empty(t, output.TxID)
}

func TestETHSendTransactionUseCase_Execute_EmptySignedTxHex(t *testing.T) {
	t.Parallel()
	deps := newEthSendTxTestDeps(t)
	useCase := newEthSendTxUseCase(deps)

	txFile := testSignedETHFile()
	txFile.SignedTxHex = "" // missing signed hex

	deps.txFileRepo.EXPECT().
		ValidateFilePath(testETHSignedFilePath, domainTx.TxTypeSigned).
		Return(domainTx.ActionTypeDeposit, domainTx.TxTypeSigned, int64(1), 1, nil)
	deps.txFileRepo.EXPECT().
		ReadETHJSONFile(testETHSignedFilePath).
		Return(txFile, nil)

	output, err := useCase.Execute(context.Background(), watchusecase.SendTransactionInput{
		FilePath: testETHSignedFilePath,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "signed_tx_hex is empty")
	assert.Empty(t, output.TxID)
}

func TestETHSendTransactionUseCase_Execute_SendError(t *testing.T) {
	t.Parallel()
	deps := newEthSendTxTestDeps(t)
	useCase := newEthSendTxUseCase(deps)

	txFile := testSignedETHFile()

	deps.txFileRepo.EXPECT().
		ValidateFilePath(testETHSignedFilePath, domainTx.TxTypeSigned).
		Return(domainTx.ActionTypeDeposit, domainTx.TxTypeSigned, int64(1), 1, nil)
	deps.txFileRepo.EXPECT().
		ReadETHJSONFile(testETHSignedFilePath).
		Return(txFile, nil)
	deps.ethClient.EXPECT().
		SendSignedRawTransaction(context.Background(), testETHSignedTxHex).
		Return("", errors.New("nonce too low"))

	output, err := useCase.Execute(context.Background(), watchusecase.SendTransactionInput{
		FilePath: testETHSignedFilePath,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "fail to call eth.SendSignedRawTransaction()")
	assert.Empty(t, output.TxID)
}

func TestETHSendTransactionUseCase_Execute_DBUpdateError(t *testing.T) {
	t.Parallel()
	deps := newEthSendTxTestDeps(t)
	useCase := newEthSendTxUseCase(deps)

	txFile := testSignedETHFile()

	deps.txFileRepo.EXPECT().
		ValidateFilePath(testETHSignedFilePath, domainTx.TxTypeSigned).
		Return(domainTx.ActionTypeDeposit, domainTx.TxTypeSigned, int64(1), 1, nil)
	deps.txFileRepo.EXPECT().
		ReadETHJSONFile(testETHSignedFilePath).
		Return(txFile, nil)
	deps.ethClient.EXPECT().
		SendSignedRawTransaction(context.Background(), testETHSignedTxHex).
		Return(testETHSentTxHash, nil)
	deps.txDetailRepo.EXPECT().
		UpdateAfterTxSent(testETHUUID, domainTx.TxTypeSent, testETHSignedTxHex, testETHSentTxHash).
		Return(int64(0), errors.New("db connection error"))

	output, err := useCase.Execute(context.Background(), watchusecase.SendTransactionInput{
		FilePath: testETHSignedFilePath,
	})

	// DB error is logged as warning but tx is already sent; return error with TxID
	require.Error(t, err)
	assert.Contains(t, err.Error(), "fail to call txDetailRepo.UpdateAfterTxSent()")
	assert.Equal(t, testETHSentTxHash, output.TxID)
}

func TestETHSendTransactionUseCase_Execute_Success(t *testing.T) {
	t.Parallel()
	deps := newEthSendTxTestDeps(t)
	useCase := newEthSendTxUseCase(deps)

	txFile := testSignedETHFile()

	deps.txFileRepo.EXPECT().
		ValidateFilePath(testETHSignedFilePath, domainTx.TxTypeSigned).
		Return(domainTx.ActionTypeDeposit, domainTx.TxTypeSigned, int64(1), 1, nil)
	deps.txFileRepo.EXPECT().
		ReadETHJSONFile(testETHSignedFilePath).
		Return(txFile, nil)
	deps.ethClient.EXPECT().
		SendSignedRawTransaction(context.Background(), testETHSignedTxHex).
		Return(testETHSentTxHash, nil)
	deps.txDetailRepo.EXPECT().
		UpdateAfterTxSent(testETHUUID, domainTx.TxTypeSent, testETHSignedTxHex, testETHSentTxHash).
		Return(int64(1), nil)

	output, err := useCase.Execute(context.Background(), watchusecase.SendTransactionInput{
		FilePath: testETHSignedFilePath,
	})

	require.NoError(t, err)
	assert.Equal(t, testETHSentTxHash, output.TxID)
}
