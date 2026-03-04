package xrp_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	"github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/xrp"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	xrpapiamocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/xrp/mocks"
	repomocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch/mocks"
	storagemocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/transaction/mocks"
	"github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/xrplgo"
)

const (
	// testSignedBlob1 is a sample signed XRP transaction blob for testing
	testSignedBlob1 = "1200002280000000240000000361D4838D7EA4C68000000000000000000000000055534400000000004B4E9C06" +
		"F24296074F7BC48F92A97916C6DC5EA968400000000000000A732103AB40A0490F9B7ED8DF29D246BF2D6269820A0EE774" +
		"2ACDD457BEA7C7D0931EDB7446304402200E5C2DD81FDF0BE9AB2A8D797885ED49E804DBF28E806604D878756410CA98B" +
		"102203349581946B0DDD06B36B35DED53E2449A5D656EDFE4071297A5043D0C0A842081146AA8D9E8FD102C92C6FB32D7" +
		"A5F00C1C26DA0BFA83140000000000000000000000000455553440000000000000000000000000000000000000000001"

	// testSignedBlob2 is another sample signed XRP transaction blob for testing (different sequence)
	testSignedBlob2 = "1200002280000000240000000461D4838D7EA4C68000000000000000000000000055534400000000004B4E9C06" +
		"F24296074F7BC48F92A97916C6DC5EA968400000000000000A732103AB40A0490F9B7ED8DF29D246BF2D6269820A0EE774" +
		"2ACDD457BEA7C7D0931EDB7446304402200E5C2DD81FDF0BE9AB2A8D797885ED49E804DBF28E806604D878756410CA98B" +
		"102203349581946B0DDD06B36B35DED53E2449A5D656EDFE4071297A5043D0C0A842081146AA8D9E8FD102C92C6FB32D7" +
		"A5F00C1C26DA0BFA83140000000000000000000000000455553440000000000000000000000000000000000000000002"
)

// sendTestDependencies holds all mock dependencies for SendTransactionUseCase testing
type sendTestDependencies struct {
	submitter    *xrpapiamocks.MockTransactionSubmitter
	txDetailRepo *repomocks.MockXRPDetailTXRepositorier
	txFileRepo   *storagemocks.MockTransactionFileRepositorier
}

// newSendTestDependencies creates all mock dependencies
func newSendTestDependencies(t *testing.T) *sendTestDependencies {
	t.Helper()
	return &sendTestDependencies{
		submitter:    xrpapiamocks.NewMockTransactionSubmitter(t),
		txDetailRepo: repomocks.NewMockXRPDetailTXRepositorier(t),
		txFileRepo:   storagemocks.NewMockTransactionFileRepositorier(t),
	}
}

// createSendUseCase creates a new SendTransactionUseCase with the test dependencies
func createSendUseCase(deps *sendTestDependencies) watchusecase.SendTransactionUseCase {
	return xrp.NewSendTransactionUseCase(
		deps.submitter,
		deps.txDetailRepo,
		deps.txFileRepo,
	)
}

func TestSendTransactionUseCase_Execute_InvalidFilePath(t *testing.T) {
	t.Parallel()
	deps := newSendTestDependencies(t)
	useCase := createSendUseCase(deps)

	// Setup mocks
	deps.txFileRepo.EXPECT().ValidateFilePath("invalid-path", domainTx.TxTypeSigned).
		Return(domainTx.ActionType(""), domainTx.TxType(""), int64(0), 0, errors.New("invalid file path"))

	input := watchusecase.SendTransactionInput{
		FilePath: "invalid-path",
	}

	output, err := useCase.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "fail to call txFileRepo.ValidateFilePath()")
	assert.Empty(t, output.TxID)
}

func TestSendTransactionUseCase_Execute_ReadError(t *testing.T) {
	t.Parallel()
	deps := newSendTestDependencies(t)
	useCase := createSendUseCase(deps)

	// Setup mocks
	deps.txFileRepo.EXPECT().ValidateFilePath("valid-path.csv", domainTx.TxTypeSigned).
		Return(domainTx.ActionTypeDeposit, domainTx.TxTypeSigned, int64(42), 1, nil)
	deps.txFileRepo.EXPECT().ReadFileSlice("valid-path.csv").
		Return(nil, errors.New("read error"))

	input := watchusecase.SendTransactionInput{
		FilePath: "valid-path.csv",
	}

	output, err := useCase.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "fail to call txFileRepo.ReadFileSlice()")
	assert.Empty(t, output.TxID)
}

func TestSendTransactionUseCase_Execute_InvalidLineFormat(t *testing.T) {
	t.Parallel()
	deps := newSendTestDependencies(t)
	useCase := createSendUseCase(deps)

	// Setup mocks: file contains a line with wrong format (missing txBlob field)
	deps.txFileRepo.EXPECT().ValidateFilePath("signed.csv", domainTx.TxTypeSigned).
		Return(domainTx.ActionTypeDeposit, domainTx.TxTypeSigned, int64(42), 1, nil)
	deps.txFileRepo.EXPECT().ReadFileSlice("signed.csv").
		Return([]string{"uuid-only-no-blob"}, nil)

	input := watchusecase.SendTransactionInput{
		FilePath: "signed.csv",
	}

	output, err := useCase.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid signed tx line format")
	assert.Empty(t, output.TxID)
}

func TestSendTransactionUseCase_Execute_EmptySignedBlob(t *testing.T) {
	t.Parallel()
	deps := newSendTestDependencies(t)
	useCase := createSendUseCase(deps)

	// Setup mocks: file contains a line with empty blob (third field empty)
	deps.txFileRepo.EXPECT().ValidateFilePath("signed.csv", domainTx.TxTypeSigned).
		Return(domainTx.ActionTypeDeposit, domainTx.TxTypeSigned, int64(42), 1, nil)
	deps.txFileRepo.EXPECT().ReadFileSlice("signed.csv").
		Return([]string{"uuid-1,txhash-1,"}, nil)

	input := watchusecase.SendTransactionInput{
		FilePath: "signed.csv",
	}

	output, err := useCase.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid signed tx line format")
	assert.Empty(t, output.TxID)
}

func TestSendTransactionUseCase_Execute_EmptyFile(t *testing.T) {
	t.Parallel()
	deps := newSendTestDependencies(t)
	useCase := createSendUseCase(deps)

	// Setup mocks: ReadFileSlice returns an empty slice
	deps.txFileRepo.EXPECT().ValidateFilePath("signed.csv", domainTx.TxTypeSigned).
		Return(domainTx.ActionTypeDeposit, domainTx.TxTypeSigned, int64(42), 0, nil)
	deps.txFileRepo.EXPECT().ReadFileSlice("signed.csv").
		Return([]string{}, nil)

	input := watchusecase.SendTransactionInput{
		FilePath: "signed.csv",
	}

	output, err := useCase.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no signed transactions found in file")
	assert.Empty(t, output.TxID)
}

func TestSendTransactionUseCase_Execute_SubmissionError(t *testing.T) {
	t.Parallel()
	deps := newSendTestDependencies(t)
	useCase := createSendUseCase(deps)

	// Setup mocks
	deps.txFileRepo.EXPECT().ValidateFilePath("signed.csv", domainTx.TxTypeSigned).
		Return(domainTx.ActionTypeDeposit, domainTx.TxTypeSigned, int64(42), 1, nil)
	deps.txFileRepo.EXPECT().ReadFileSlice("signed.csv").
		Return([]string{"01234567-89ab-cdef-0123-456789abcdef,txhash," + testSignedBlob1}, nil)
	deps.submitter.EXPECT().SubmitTransaction(mock.Anything, testSignedBlob1).
		Return((*xrplgo.SentTx)(nil), uint64(0), errors.New("tefPAST_SEQ: sequence number already used"))

	input := watchusecase.SendTransactionInput{
		FilePath: "signed.csv",
	}

	output, err := useCase.Execute(context.Background(), input)

	// Should not return error for individual transaction failures (logged as warnings)
	require.NoError(t, err)
	assert.Empty(t, output.TxID)
}

func TestSendTransactionUseCase_Execute_Success(t *testing.T) {
	t.Parallel()
	deps := newSendTestDependencies(t)
	useCase := createSendUseCase(deps)

	uuid := "01234567-89ab-cdef-0123-456789abcdef"
	txHash := "1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF"
	csvLine := uuid + ",txhash," + testSignedBlob1

	sentTx := &xrplgo.SentTx{
		ResultCode:    "tesSUCCESS",
		ResultMessage: "The transaction was applied.",
		TxBlob:        testSignedBlob1,
		TxJSON: xrplgo.TxInput{
			Hash:               txHash,
			LastLedgerSequence: 12345,
		},
	}

	txInfo := &xrplgo.TxInfo{
		Outcome: xrplgo.TxOutcome{
			Result: "tesSUCCESS",
		},
	}

	// Setup mocks
	deps.txFileRepo.EXPECT().ValidateFilePath("signed.csv", domainTx.TxTypeSigned).
		Return(domainTx.ActionTypeDeposit, domainTx.TxTypeSigned, int64(42), 1, nil)
	deps.txFileRepo.EXPECT().ReadFileSlice("signed.csv").
		Return([]string{csvLine}, nil)
	deps.submitter.EXPECT().SubmitTransaction(mock.Anything, testSignedBlob1).
		Return(sentTx, uint64(12340), nil)
	deps.submitter.EXPECT().WaitValidation(mock.Anything, uint64(12345)).
		Return(uint64(12345), nil)
	deps.submitter.EXPECT().GetTransaction(mock.Anything, txHash, uint64(12340)).
		Return(txInfo, nil)
	deps.txDetailRepo.EXPECT().UpdateAfterTxSent(
		uuid,
		domainTx.TxTypeSent,
		txHash,
		testSignedBlob1,
		uint64(12340),
	).Return(int64(1), nil)

	input := watchusecase.SendTransactionInput{
		FilePath: "signed.csv",
	}

	output, err := useCase.Execute(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, txHash, output.TxID, "should return transaction hash")
}

func TestSendTransactionUseCase_Execute_MultipleTransactions(t *testing.T) {
	t.Parallel()
	deps := newSendTestDependencies(t)
	useCase := createSendUseCase(deps)

	txHash1 := "HASH1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890AB"
	txHash2 := "HASH2234567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890AB"
	csvLine1 := "uuid-1,txhash1," + testSignedBlob1
	csvLine2 := "uuid-2,txhash2," + testSignedBlob2

	sentTx1 := &xrplgo.SentTx{
		ResultCode:    "tesSUCCESS",
		ResultMessage: "The transaction was applied.",
		TxBlob:        testSignedBlob1,
		TxJSON: xrplgo.TxInput{
			Hash:               txHash1,
			LastLedgerSequence: 12345,
		},
	}

	sentTx2 := &xrplgo.SentTx{
		ResultCode:    "tesSUCCESS",
		ResultMessage: "The transaction was applied.",
		TxBlob:        testSignedBlob2,
		TxJSON: xrplgo.TxInput{
			Hash:               txHash2,
			LastLedgerSequence: 12346,
		},
	}

	txInfo1 := &xrplgo.TxInfo{
		Outcome: xrplgo.TxOutcome{Result: "tesSUCCESS"},
	}
	txInfo2 := &xrplgo.TxInfo{
		Outcome: xrplgo.TxOutcome{Result: "tesSUCCESS"},
	}

	// Setup mocks
	deps.txFileRepo.EXPECT().ValidateFilePath("signed.csv", domainTx.TxTypeSigned).
		Return(domainTx.ActionTypeDeposit, domainTx.TxTypeSigned, int64(42), 2, nil)
	deps.txFileRepo.EXPECT().ReadFileSlice("signed.csv").
		Return([]string{csvLine1, csvLine2}, nil)

	// Transaction 1
	deps.submitter.EXPECT().SubmitTransaction(mock.Anything, testSignedBlob1).
		Return(sentTx1, uint64(12340), nil)
	deps.submitter.EXPECT().WaitValidation(mock.Anything, uint64(12345)).
		Return(uint64(12345), nil)
	deps.submitter.EXPECT().GetTransaction(mock.Anything, txHash1, uint64(12340)).
		Return(txInfo1, nil)
	deps.txDetailRepo.EXPECT().
		UpdateAfterTxSent("uuid-1", domainTx.TxTypeSent, txHash1, testSignedBlob1, uint64(12340)).
		Return(int64(1), nil)

	// Transaction 2
	deps.submitter.EXPECT().SubmitTransaction(mock.Anything, testSignedBlob2).
		Return(sentTx2, uint64(12341), nil)
	deps.submitter.EXPECT().WaitValidation(mock.Anything, uint64(12346)).
		Return(uint64(12346), nil)
	deps.submitter.EXPECT().GetTransaction(mock.Anything, txHash2, uint64(12341)).
		Return(txInfo2, nil)
	deps.txDetailRepo.EXPECT().
		UpdateAfterTxSent("uuid-2", domainTx.TxTypeSent, txHash2, testSignedBlob2, uint64(12341)).
		Return(int64(1), nil)

	input := watchusecase.SendTransactionInput{
		FilePath: "signed.csv",
	}

	output, err := useCase.Execute(context.Background(), input)

	require.NoError(t, err)
	// For multiple transactions, return the first successful transaction hash
	assert.NotEmpty(t, output.TxID, "should return a transaction hash")
	// The hash should be one of the submitted transactions (order may vary due to concurrency)
	assert.True(t,
		output.TxID == txHash1 || output.TxID == txHash2,
		"should return one of the submitted transaction hashes",
	)
}
