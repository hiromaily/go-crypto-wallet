package xrp_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	"github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/xrp"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	xrpapiamocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/xrp/mocks"
	repomocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch/mocks"
	storagemocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/transaction/mocks"
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

func TestSendTransactionUseCase_Execute_JSONParsingError(t *testing.T) {
	deps := newSendTestDependencies(t)
	useCase := createSendUseCase(deps)

	// Setup mocks
	deps.txFileRepo.EXPECT().ValidateFilePath("valid-path.json", domainTx.TxTypeSigned).
		Return(domainTx.ActionTypeDeposit, domainTx.TxTypeSigned, int64(42), 1, nil)
	deps.txFileRepo.EXPECT().ReadXRPJSONFile("valid-path.json").
		Return((*dtoxrp.XRPTransactionFile)(nil), errors.New("invalid JSON format"))

	input := watchusecase.SendTransactionInput{
		FilePath: "valid-path.json",
	}

	output, err := useCase.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "fail to call txFileRepo.ReadXRPJSONFile()")
	assert.Contains(t, err.Error(), "invalid JSON format")
	assert.Empty(t, output.TxID)
}

func TestSendTransactionUseCase_Execute_IncompleteTransaction(t *testing.T) {
	deps := newSendTestDependencies(t)
	useCase := createSendUseCase(deps)

	signedBlob := testSignedBlob1

	txFile := &dtoxrp.XRPTransactionFile{
		Version:   "1.0.0",
		Chain:     "XRP",
		Network:   "testnet",
		CreatedAt: "2024-02-14T00:00:00Z",
		Transactions: []dtoxrp.XRPTransactionEntry{
			{
				UUID:               "01234567-89ab-cdef-0123-456789abcdef",
				SenderAccount:      "rSenderAddress123",
				SenderAccountType:  "client",
				SignatureCount:     1,
				RequiredSignatures: 2, // Incomplete: 1/2 signatures
				SignedBlob:         &signedBlob,
				IsComplete:         false, // Not complete
			},
		},
	}

	// Setup mocks
	deps.txFileRepo.EXPECT().ValidateFilePath("signed.json", domainTx.TxTypeSigned).
		Return(domainTx.ActionTypeDeposit, domainTx.TxTypeSigned, int64(42), 1, nil)
	deps.txFileRepo.EXPECT().ReadXRPJSONFile("signed.json").
		Return(txFile, nil)

	input := watchusecase.SendTransactionInput{
		FilePath: "signed.json",
	}

	output, err := useCase.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "transaction is incomplete")
	assert.Contains(t, err.Error(), "requires 2 signatures but has 1")
	assert.Empty(t, output.TxID)
}

func TestSendTransactionUseCase_Execute_NullSignedBlob(t *testing.T) {
	deps := newSendTestDependencies(t)
	useCase := createSendUseCase(deps)

	txFile := &dtoxrp.XRPTransactionFile{
		Version:   "1.0.0",
		Chain:     "XRP",
		Network:   "testnet",
		CreatedAt: "2024-02-14T00:00:00Z",
		Transactions: []dtoxrp.XRPTransactionEntry{
			{
				UUID:               "01234567-89ab-cdef-0123-456789abcdef",
				SenderAccount:      "rSenderAddress123",
				SenderAccountType:  "client",
				SignatureCount:     1,
				RequiredSignatures: 1,
				SignedBlob:         nil, // Null signed blob
				IsComplete:         true,
			},
		},
	}

	// Setup mocks
	deps.txFileRepo.EXPECT().ValidateFilePath("signed.json", domainTx.TxTypeSigned).
		Return(domainTx.ActionTypeDeposit, domainTx.TxTypeSigned, int64(42), 1, nil)
	deps.txFileRepo.EXPECT().ReadXRPJSONFile("signed.json").
		Return(txFile, nil)

	input := watchusecase.SendTransactionInput{
		FilePath: "signed.json",
	}

	output, err := useCase.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "signedBlob is null for transaction")
	assert.Empty(t, output.TxID)
}

func TestSendTransactionUseCase_Execute_EmptySignedBlob(t *testing.T) {
	deps := newSendTestDependencies(t)
	useCase := createSendUseCase(deps)

	emptyBlob := ""
	txFile := &dtoxrp.XRPTransactionFile{
		Version:   "1.0.0",
		Chain:     "XRP",
		Network:   "testnet",
		CreatedAt: "2024-02-14T00:00:00Z",
		Transactions: []dtoxrp.XRPTransactionEntry{
			{
				UUID:               "01234567-89ab-cdef-0123-456789abcdef",
				SenderAccount:      "rSenderAddress123",
				SenderAccountType:  "client",
				SignatureCount:     1,
				RequiredSignatures: 1,
				SignedBlob:         &emptyBlob, // Empty signed blob
				IsComplete:         true,
			},
		},
	}

	// Setup mocks
	deps.txFileRepo.EXPECT().ValidateFilePath("signed.json", domainTx.TxTypeSigned).
		Return(domainTx.ActionTypeDeposit, domainTx.TxTypeSigned, int64(42), 1, nil)
	deps.txFileRepo.EXPECT().ReadXRPJSONFile("signed.json").
		Return(txFile, nil)

	input := watchusecase.SendTransactionInput{
		FilePath: "signed.json",
	}

	output, err := useCase.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "signedBlob is empty for transaction")
	assert.Empty(t, output.TxID)
}

func TestSendTransactionUseCase_Execute_SubmissionError(t *testing.T) {
	deps := newSendTestDependencies(t)
	useCase := createSendUseCase(deps)

	signedBlob := testSignedBlob1

	txFile := &dtoxrp.XRPTransactionFile{
		Version:   "1.0.0",
		Chain:     "XRP",
		Network:   "testnet",
		CreatedAt: "2024-02-14T00:00:00Z",
		Transactions: []dtoxrp.XRPTransactionEntry{
			{
				UUID:               "01234567-89ab-cdef-0123-456789abcdef",
				SenderAccount:      "rSenderAddress123",
				SenderAccountType:  "client",
				SignatureCount:     1,
				RequiredSignatures: 1,
				SignedBlob:         &signedBlob,
				IsComplete:         true,
			},
		},
	}

	// Setup mocks
	deps.txFileRepo.EXPECT().ValidateFilePath("signed.json", domainTx.TxTypeSigned).
		Return(domainTx.ActionTypeDeposit, domainTx.TxTypeSigned, int64(42), 1, nil)
	deps.txFileRepo.EXPECT().ReadXRPJSONFile("signed.json").
		Return(txFile, nil)
	deps.submitter.EXPECT().SubmitTransaction(mock.Anything, signedBlob).
		Return((*dtoxrp.SentTx)(nil), uint64(0), errors.New("tefPAST_SEQ: sequence number already used"))

	input := watchusecase.SendTransactionInput{
		FilePath: "signed.json",
	}

	output, err := useCase.Execute(context.Background(), input)

	// Should not return error for individual transaction failures (logged as warnings)
	require.NoError(t, err)
	assert.Empty(t, output.TxID)
}

func TestSendTransactionUseCase_Execute_Success(t *testing.T) {
	deps := newSendTestDependencies(t)
	useCase := createSendUseCase(deps)

	signedBlob := testSignedBlob1
	txHash := "1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF"

	txFile := &dtoxrp.XRPTransactionFile{
		Version:   "1.0.0",
		Chain:     "XRP",
		Network:   "testnet",
		CreatedAt: "2024-02-14T00:00:00Z",
		Transactions: []dtoxrp.XRPTransactionEntry{
			{
				UUID:               "01234567-89ab-cdef-0123-456789abcdef",
				SenderAccount:      "rSenderAddress123",
				SenderAccountType:  "client",
				SignatureCount:     1,
				RequiredSignatures: 1,
				SignedBlob:         &signedBlob,
				IsComplete:         true,
			},
		},
	}

	sentTx := &dtoxrp.SentTx{
		ResultCode:    "tesSUCCESS",
		ResultMessage: "The transaction was applied.",
		TxBlob:        signedBlob,
		TxJSON: dtoxrp.TxInput{
			Hash:               txHash,
			LastLedgerSequence: 12345,
		},
	}

	txInfo := &dtoxrp.TxInfo{
		Outcome: dtoxrp.TxOutcome{
			Result: "tesSUCCESS",
		},
	}

	// Setup mocks
	deps.txFileRepo.EXPECT().ValidateFilePath("signed.json", domainTx.TxTypeSigned).
		Return(domainTx.ActionTypeDeposit, domainTx.TxTypeSigned, int64(42), 1, nil)
	deps.txFileRepo.EXPECT().ReadXRPJSONFile("signed.json").
		Return(txFile, nil)
	deps.submitter.EXPECT().SubmitTransaction(mock.Anything, signedBlob).
		Return(sentTx, uint64(12340), nil)
	deps.submitter.EXPECT().WaitValidation(mock.Anything, uint64(12345)).
		Return(uint64(12345), nil)
	deps.submitter.EXPECT().GetTransaction(mock.Anything, txHash, uint64(12340)).
		Return(txInfo, nil)
	deps.txDetailRepo.EXPECT().UpdateAfterTxSent(
		"01234567-89ab-cdef-0123-456789abcdef",
		domainTx.TxTypeSent,
		txHash,
		signedBlob,
		uint64(12340),
	).Return(int64(1), nil)

	input := watchusecase.SendTransactionInput{
		FilePath: "signed.json",
	}

	output, err := useCase.Execute(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, txHash, output.TxID, "should return transaction hash")
}

func TestSendTransactionUseCase_Execute_MultipleTransactions(t *testing.T) {
	deps := newSendTestDependencies(t)
	useCase := createSendUseCase(deps)

	signedBlob1 := testSignedBlob1
	signedBlob2 := testSignedBlob2
	txHash1 := "HASH1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890AB"
	txHash2 := "HASH2234567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890AB"

	txFile := &dtoxrp.XRPTransactionFile{
		Version:   "1.0.0",
		Chain:     "XRP",
		Network:   "testnet",
		CreatedAt: "2024-02-14T00:00:00Z",
		Transactions: []dtoxrp.XRPTransactionEntry{
			{
				UUID:               "uuid-1",
				SenderAccount:      "rSenderAddress1",
				SenderAccountType:  "client",
				SignatureCount:     1,
				RequiredSignatures: 1,
				SignedBlob:         &signedBlob1,
				IsComplete:         true,
			},
			{
				UUID:               "uuid-2",
				SenderAccount:      "rSenderAddress2",
				SenderAccountType:  "client",
				SignatureCount:     1,
				RequiredSignatures: 1,
				SignedBlob:         &signedBlob2,
				IsComplete:         true,
			},
		},
	}

	sentTx1 := &dtoxrp.SentTx{
		ResultCode:    "tesSUCCESS",
		ResultMessage: "The transaction was applied.",
		TxBlob:        signedBlob1,
		TxJSON: dtoxrp.TxInput{
			Hash:               txHash1,
			LastLedgerSequence: 12345,
		},
	}

	sentTx2 := &dtoxrp.SentTx{
		ResultCode:    "tesSUCCESS",
		ResultMessage: "The transaction was applied.",
		TxBlob:        signedBlob2,
		TxJSON: dtoxrp.TxInput{
			Hash:               txHash2,
			LastLedgerSequence: 12346,
		},
	}

	txInfo1 := &dtoxrp.TxInfo{
		Outcome: dtoxrp.TxOutcome{Result: "tesSUCCESS"},
	}
	txInfo2 := &dtoxrp.TxInfo{
		Outcome: dtoxrp.TxOutcome{Result: "tesSUCCESS"},
	}

	// Setup mocks
	deps.txFileRepo.EXPECT().ValidateFilePath("signed.json", domainTx.TxTypeSigned).
		Return(domainTx.ActionTypeDeposit, domainTx.TxTypeSigned, int64(42), 2, nil)
	deps.txFileRepo.EXPECT().ReadXRPJSONFile("signed.json").
		Return(txFile, nil)

	// Transaction 1
	deps.submitter.EXPECT().SubmitTransaction(mock.Anything, signedBlob1).
		Return(sentTx1, uint64(12340), nil)
	deps.submitter.EXPECT().WaitValidation(mock.Anything, uint64(12345)).
		Return(uint64(12345), nil)
	deps.submitter.EXPECT().GetTransaction(mock.Anything, txHash1, uint64(12340)).
		Return(txInfo1, nil)
	deps.txDetailRepo.EXPECT().UpdateAfterTxSent("uuid-1", domainTx.TxTypeSent, txHash1, signedBlob1, uint64(12340)).
		Return(int64(1), nil)

	// Transaction 2
	deps.submitter.EXPECT().SubmitTransaction(mock.Anything, signedBlob2).
		Return(sentTx2, uint64(12341), nil)
	deps.submitter.EXPECT().WaitValidation(mock.Anything, uint64(12346)).
		Return(uint64(12346), nil)
	deps.submitter.EXPECT().GetTransaction(mock.Anything, txHash2, uint64(12341)).
		Return(txInfo2, nil)
	deps.txDetailRepo.EXPECT().UpdateAfterTxSent("uuid-2", domainTx.TxTypeSent, txHash2, signedBlob2, uint64(12341)).
		Return(int64(1), nil)

	input := watchusecase.SendTransactionInput{
		FilePath: "signed.json",
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
