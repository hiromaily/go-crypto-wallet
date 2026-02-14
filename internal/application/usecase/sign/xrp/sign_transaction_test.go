package xrp

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	file "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	signusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	domainXrp "github.com/hiromaily/go-crypto-wallet/internal/domain/xrp"
)

// MockXRPSignClient is a mock implementation of xrpSignClient for testing
type MockXRPSignClient struct {
	mock.Mock
}

func (m *MockXRPSignClient) SignTransaction(
	ctx context.Context,
	txInput *dtoxrp.TxInput,
	secret string,
) (string, string, error) {
	args := m.Called(ctx, txInput, secret)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockXRPSignClient) SignTransactionNative(
	ctx context.Context,
	txInput *dtoxrp.TxInput,
	secret string,
	isMultiSig bool,
	existingSignedBlob *string,
) (string, string, error) {
	args := m.Called(ctx, txInput, secret, isMultiSig, existingSignedBlob)
	return args.String(0), args.String(1), args.Error(2)
}

// MockXRPAccountKeyRepo is a mock implementation of XRPAccountKeyRepositorier
type MockXRPAccountKeyRepo struct {
	mock.Mock
}

func (m *MockXRPAccountKeyRepo) GetAllAddrStatus(
	ctx context.Context,
	accountType domainAccount.AccountType,
	addrStatus domainAddress.AddrStatus,
) ([]*domainXrp.XRPAccountKey, error) {
	args := m.Called(ctx, accountType, addrStatus)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domainXrp.XRPAccountKey), args.Error(1)
}

func (m *MockXRPAccountKeyRepo) GetSecret(
	ctx context.Context,
	accountType domainAccount.AccountType,
	address string,
) (string, error) {
	args := m.Called(ctx, accountType, address)
	return args.String(0), args.Error(1)
}

func (m *MockXRPAccountKeyRepo) InsertBulk(ctx context.Context, items []*domainXrp.XRPAccountKey) error {
	args := m.Called(ctx, items)
	return args.Error(0)
}

func (m *MockXRPAccountKeyRepo) UpdateAddrStatus(
	ctx context.Context,
	accountType domainAccount.AccountType,
	addrStatus domainAddress.AddrStatus,
	strWIFs []string,
) (int64, error) {
	args := m.Called(ctx, accountType, addrStatus, strWIFs)
	return args.Get(0).(int64), args.Error(1)
}

// MockTransactionFileRepo is a mock implementation of TransactionFileRepositorier
type MockTransactionFileRepo struct {
	mock.Mock
}

func (m *MockTransactionFileRepo) CreateFilePath(
	actionType domainTx.ActionType,
	txType domainTx.TxType,
	txID int64,
	signedCount int,
) string {
	args := m.Called(actionType, txType, txID, signedCount)
	return args.String(0)
}

func (m *MockTransactionFileRepo) GetFileNameType(filePath string) (*file.FileName, error) {
	args := m.Called(filePath)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*file.FileName), args.Error(1)
}

func (m *MockTransactionFileRepo) ValidateFilePath(
	filePath string,
	expectedTxType domainTx.TxType,
) (domainTx.ActionType, domainTx.TxType, int64, int, error) {
	args := m.Called(filePath, expectedTxType)
	return args.Get(0).(domainTx.ActionType),
		args.Get(1).(domainTx.TxType),
		args.Get(2).(int64),
		args.Int(3),
		args.Error(4)
}

func (m *MockTransactionFileRepo) ReadFile(path string) (string, error) {
	args := m.Called(path)
	return args.String(0), args.Error(1)
}

func (m *MockTransactionFileRepo) ReadFileSlice(path string) ([]string, error) {
	args := m.Called(path)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockTransactionFileRepo) WriteFile(path, hexTx string) (string, error) {
	args := m.Called(path, hexTx)
	return args.String(0), args.Error(1)
}

func (m *MockTransactionFileRepo) WriteFileSlice(path string, data []string) (string, error) {
	args := m.Called(path, data)
	return args.String(0), args.Error(1)
}

func (m *MockTransactionFileRepo) ReadPSBTFile(path string) (string, error) {
	args := m.Called(path)
	return args.String(0), args.Error(1)
}

func (m *MockTransactionFileRepo) WritePSBTFile(path, psbtBase64 string) (string, error) {
	args := m.Called(path, psbtBase64)
	return args.String(0), args.Error(1)
}

func (m *MockTransactionFileRepo) WriteHexFile(path, hexTx string) (string, error) {
	args := m.Called(path, hexTx)
	return args.String(0), args.Error(1)
}

func (m *MockTransactionFileRepo) ReadXRPJSONFile(path string) (*dtoxrp.XRPTransactionFile, error) {
	args := m.Called(path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dtoxrp.XRPTransactionFile), args.Error(1)
}

func (m *MockTransactionFileRepo) WriteXRPJSONFile(path string, data *dtoxrp.XRPTransactionFile) (string, error) {
	args := m.Called(path, data)
	return args.String(0), args.Error(1)
}

func TestSignTransactionUseCase_Sign_SingleSignature(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockXRP := new(MockXRPSignClient)
	mockKeyRepo := new(MockXRPAccountKeyRepo)
	mockFileRepo := new(MockTransactionFileRepo)

	usecase := NewSignTransactionUseCase(
		mockXRP,
		mockKeyRepo,
		mockFileRepo,
		domainWallet.WalletTypeSign,
	)

	input := signusecase.SignTransactionInput{
		FilePath: "/path/to/unsigned.json",
	}

	// Mock GetFileNameType
	mockFileRepo.On("GetFileNameType", input.FilePath).
		Return(&file.FileName{
			ActionType:  domainTx.ActionTypeDeposit,
			TxType:      domainTx.TxTypeUnsigned,
			TxID:        1,
			SignedCount: 0,
		}, nil)

	// Mock ReadXRPJSONFile - unsigned transaction
	txInput := dtoxrp.TxInput{
		Account:            "rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH",
		Destination:        "rLHzPsX6oXkzU9J1kSwq7X4X8J1Xa8AkDC",
		Amount:             "1000000",
		Fee:                "12",
		Sequence:           42,
		LastLedgerSequence: 1000000,
		TransactionType:    "Payment",
	}
	signedBlobNil := (*string)(nil)
	unsignedTxFile := &dtoxrp.XRPTransactionFile{
		Version:   "1.0.0",
		Chain:     "XRP",
		Network:   "testnet",
		CreatedAt: "2024-02-14T00:00:00Z",
		Transactions: []dtoxrp.XRPTransactionEntry{
			{
				UUID:               "018d5e3e-1234-7890-abcd-ef1234567890",
				UnsignedData:       txInput,
				SenderAccount:      "rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH",
				SenderAccountType:  "client",
				SignatureCount:     0,
				RequiredSignatures: 1,
				SignedBlob:         signedBlobNil,
				IsComplete:         false,
			},
		},
	}
	mockFileRepo.On("ReadXRPJSONFile", input.FilePath).Return(unsignedTxFile, nil)

	// Mock GetSecret
	mockKeyRepo.On("GetSecret", ctx, domainAccount.AccountType("client"), "rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH").
		Return("sEdTM1uX8pu2do5XvTnutH6HsouMaM2", nil)

	// Mock SignTransactionNative (no existing blob for first/only signature)
	var nilBlob *string
	mockXRP.On("SignTransactionNative", ctx, &txInput, "sEdTM1uX8pu2do5XvTnutH6HsouMaM2", false, nilBlob).
		Return(
			"3E3C6E8F1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF12345678",
			"120000228000000024000000002A61400000000000F424068400000000000000C732103ABCDEF...",
			nil,
		)

	// Mock CreateFilePath
	mockFileRepo.On("CreateFilePath", domainTx.ActionTypeDeposit, domainTx.TxTypeSigned, int64(1), 1).
		Return("/path/to/signed_1.json")

	// Mock WriteXRPJSONFile
	mockFileRepo.On(
		"WriteXRPJSONFile",
		"/path/to/signed_1.json",
		mock.MatchedBy(func(data *dtoxrp.XRPTransactionFile) bool {
			// Verify the signed file has correct signature count and completion status
			return data.Transactions[0].SignatureCount == 1 &&
				data.Transactions[0].IsComplete == true &&
				data.Transactions[0].SignedBlob != nil
		}),
	).Return("/path/to/signed_1.json", nil)

	// Act
	output, err := usecase.Sign(ctx, input)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "/path/to/signed_1.json", output.NextFilePath)
	assert.True(t, output.IsComplete)
	mockFileRepo.AssertExpectations(t)
	mockKeyRepo.AssertExpectations(t)
	mockXRP.AssertExpectations(t)
}

func TestSignTransactionUseCase_Sign_MultiSignature_FirstSigner(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockXRP := new(MockXRPSignClient)
	mockKeyRepo := new(MockXRPAccountKeyRepo)
	mockFileRepo := new(MockTransactionFileRepo)

	usecase := NewSignTransactionUseCase(
		mockXRP,
		mockKeyRepo,
		mockFileRepo,
		domainWallet.WalletTypeSign,
	)

	input := signusecase.SignTransactionInput{
		FilePath: "/path/to/unsigned.json",
	}

	// Define nil blob for mock expectations (no existing signatures)
	var nilBlob *string

	// Mock GetFileNameType
	mockFileRepo.On("GetFileNameType", input.FilePath).
		Return(&file.FileName{
			ActionType:  domainTx.ActionTypePayment,
			TxType:      domainTx.TxTypeUnsigned,
			TxID:        2,
			SignedCount: 0,
		}, nil)

	// Mock ReadXRPJSONFile - unsigned multi-sig transaction
	txInput := dtoxrp.TxInput{
		Account:            "rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH",
		Destination:        "rLHzPsX6oXkzU9J1kSwq7X4X8J1Xa8AkDC",
		Amount:             "1000000",
		Fee:                "36", // Multi-sig fee: (2+1) * 12
		Sequence:           42,
		LastLedgerSequence: 1000000,
		TransactionType:    "Payment",
	}
	signedBlobNil := (*string)(nil)
	unsignedTxFile := &dtoxrp.XRPTransactionFile{
		Version:   "1.0.0",
		Chain:     "XRP",
		Network:   "testnet",
		CreatedAt: "2024-02-14T00:00:00Z",
		Transactions: []dtoxrp.XRPTransactionEntry{
			{
				UUID:               "018d5e3e-5678-9012-abcd-ef1234567890",
				UnsignedData:       txInput,
				SenderAccount:      "rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH",
				SenderAccountType:  "payment",
				SignatureCount:     0,
				RequiredSignatures: 2, // 2-of-3 multi-sig
				SignedBlob:         signedBlobNil,
				IsComplete:         false,
			},
		},
	}
	mockFileRepo.On("ReadXRPJSONFile", input.FilePath).Return(unsignedTxFile, nil)

	// Mock GetSecret
	mockKeyRepo.On("GetSecret", ctx, domainAccount.AccountType("payment"), "rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH").
		Return("sEdSigner1Secret", nil)

	// Mock SignTransactionNative (multi-sig = true for RequiredSignatures > 1, no existing blob for first signature)
	mockXRP.On("SignTransactionNative", ctx, &txInput, "sEdSigner1Secret", true, nilBlob).
		Return(
			"4F4D7F9A1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF12345678",
			"120000228000000024000000002A61400000000000F424068400000000000000C732103MULTISIG...",
			nil,
		)

	// Mock CreateFilePath
	mockFileRepo.On("CreateFilePath", domainTx.ActionTypePayment, domainTx.TxTypeSigned, int64(2), 1).
		Return("/path/to/signed_1.json")

	// Mock WriteXRPJSONFile
	mockFileRepo.On(
		"WriteXRPJSONFile",
		"/path/to/signed_1.json",
		mock.MatchedBy(func(data *dtoxrp.XRPTransactionFile) bool {
			// First signature: count=1, not complete (needs 2)
			return data.Transactions[0].SignatureCount == 1 &&
				data.Transactions[0].IsComplete == false &&
				data.Transactions[0].SignedBlob != nil
		}),
	).Return("/path/to/signed_1.json", nil)

	// Act
	output, err := usecase.Sign(ctx, input)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "/path/to/signed_1.json", output.NextFilePath)
	assert.False(t, output.IsComplete) // Not complete yet, needs second signature
	mockFileRepo.AssertExpectations(t)
	mockKeyRepo.AssertExpectations(t)
	mockXRP.AssertExpectations(t)
}

func TestSignTransactionUseCase_Sign_SkipAlreadyComplete(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockXRP := new(MockXRPSignClient)
	mockKeyRepo := new(MockXRPAccountKeyRepo)
	mockFileRepo := new(MockTransactionFileRepo)

	usecase := NewSignTransactionUseCase(
		mockXRP,
		mockKeyRepo,
		mockFileRepo,
		domainWallet.WalletTypeSign,
	)

	input := signusecase.SignTransactionInput{
		FilePath: "/path/to/signed_1.json",
	}

	// Mock GetFileNameType
	mockFileRepo.On("GetFileNameType", input.FilePath).
		Return(&file.FileName{
			ActionType:  domainTx.ActionTypeDeposit,
			TxType:      domainTx.TxTypeSigned,
			TxID:        1,
			SignedCount: 1,
		}, nil)

	// Mock ReadXRPJSONFile - already complete transaction
	signedBlob := "120000228000000024000000002A61400000000000F424068400000000000000C732103COMPLETE..."
	completeTxFile := &dtoxrp.XRPTransactionFile{
		Version:   "1.0.0",
		Chain:     "XRP",
		Network:   "testnet",
		CreatedAt: "2024-02-14T00:00:00Z",
		Transactions: []dtoxrp.XRPTransactionEntry{
			{
				UUID:               "018d5e3e-9012-3456-abcd-ef1234567890",
				UnsignedData:       dtoxrp.TxInput{},
				SenderAccount:      "rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH",
				SenderAccountType:  "client",
				SignatureCount:     1,
				RequiredSignatures: 1,
				SignedBlob:         &signedBlob,
				IsComplete:         true, // Already complete
			},
		},
	}
	mockFileRepo.On("ReadXRPJSONFile", input.FilePath).Return(completeTxFile, nil)

	// No signing or secret retrieval should happen

	// Mock CreateFilePath
	mockFileRepo.On("CreateFilePath", domainTx.ActionTypeDeposit, domainTx.TxTypeSigned, int64(1), 1).
		Return("/path/to/signed_1.json")

	// Mock WriteXRPJSONFile (file unchanged)
	mockFileRepo.On("WriteXRPJSONFile", "/path/to/signed_1.json", completeTxFile).
		Return("/path/to/signed_1.json", nil)

	// Act
	output, err := usecase.Sign(ctx, input)

	// Assert
	require.NoError(t, err)
	assert.True(t, output.IsComplete)
	assert.Equal(t, "/path/to/signed_1.json", output.NextFilePath)

	// Verify no signing operations were attempted
	mockXRP.AssertNotCalled(t, "SignTransactionNative")
	mockKeyRepo.AssertNotCalled(t, "GetSecret")
	mockFileRepo.AssertExpectations(t)
}

func TestSignTransactionUseCase_Sign_ErrorSecretRetrieval(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockXRP := new(MockXRPSignClient)
	mockKeyRepo := new(MockXRPAccountKeyRepo)
	mockFileRepo := new(MockTransactionFileRepo)

	usecase := NewSignTransactionUseCase(
		mockXRP,
		mockKeyRepo,
		mockFileRepo,
		domainWallet.WalletTypeSign,
	)

	input := signusecase.SignTransactionInput{
		FilePath: "/path/to/unsigned.json",
	}

	// Mock GetFileNameType
	mockFileRepo.On("GetFileNameType", input.FilePath).
		Return(&file.FileName{
			ActionType:  domainTx.ActionTypeDeposit,
			TxType:      domainTx.TxTypeUnsigned,
			TxID:        1,
			SignedCount: 0,
		}, nil)

	// Mock ReadXRPJSONFile
	txInput := dtoxrp.TxInput{
		Account:         "rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH",
		TransactionType: "Payment",
	}
	signedBlobNil := (*string)(nil)
	unsignedTxFile := &dtoxrp.XRPTransactionFile{
		Version:   "1.0.0",
		Chain:     "XRP",
		Network:   "testnet",
		CreatedAt: "2024-02-14T00:00:00Z",
		Transactions: []dtoxrp.XRPTransactionEntry{
			{
				UUID:               "018d5e3e-1234-7890-abcd-ef1234567890",
				UnsignedData:       txInput,
				SenderAccount:      "rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH",
				SenderAccountType:  "client",
				SignatureCount:     0,
				RequiredSignatures: 1,
				SignedBlob:         signedBlobNil,
				IsComplete:         false,
			},
		},
	}
	mockFileRepo.On("ReadXRPJSONFile", input.FilePath).Return(unsignedTxFile, nil)

	// Mock GetSecret returns error
	mockKeyRepo.On("GetSecret", ctx, domainAccount.AccountType("client"), "rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH").
		Return("", errors.New("secret not found"))

	// Act
	output, err := usecase.Sign(ctx, input)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get secret for transaction 018d5e3e-1234-7890-abcd-ef1234567890")
	assert.Equal(t, "", output.NextFilePath)
	mockFileRepo.AssertExpectations(t)
	mockKeyRepo.AssertExpectations(t)
}

func TestSignTransactionUseCase_Sign_ErrorSigning(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockXRP := new(MockXRPSignClient)
	mockKeyRepo := new(MockXRPAccountKeyRepo)
	mockFileRepo := new(MockTransactionFileRepo)

	usecase := NewSignTransactionUseCase(
		mockXRP,
		mockKeyRepo,
		mockFileRepo,
		domainWallet.WalletTypeSign,
	)

	input := signusecase.SignTransactionInput{
		FilePath: "/path/to/unsigned.json",
	}

	// Define nil blob for mock expectations (no existing signatures)
	var nilBlob *string

	// Mock GetFileNameType
	mockFileRepo.On("GetFileNameType", input.FilePath).
		Return(&file.FileName{
			ActionType:  domainTx.ActionTypeDeposit,
			TxType:      domainTx.TxTypeUnsigned,
			TxID:        1,
			SignedCount: 0,
		}, nil)

	// Mock ReadXRPJSONFile
	txInput := dtoxrp.TxInput{
		Account:         "rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH",
		TransactionType: "Payment",
		Sequence:        42,
	}
	signedBlobNil := (*string)(nil)
	unsignedTxFile := &dtoxrp.XRPTransactionFile{
		Version:   "1.0.0",
		Chain:     "XRP",
		Network:   "testnet",
		CreatedAt: "2024-02-14T00:00:00Z",
		Transactions: []dtoxrp.XRPTransactionEntry{
			{
				UUID:               "018d5e3e-5678-9012-abcd-ef1234567890",
				UnsignedData:       txInput,
				SenderAccount:      "rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH",
				SenderAccountType:  "client",
				SignatureCount:     0,
				RequiredSignatures: 1,
				SignedBlob:         signedBlobNil,
				IsComplete:         false,
			},
		},
	}
	mockFileRepo.On("ReadXRPJSONFile", input.FilePath).Return(unsignedTxFile, nil)

	// Mock GetSecret
	mockKeyRepo.On("GetSecret", ctx, domainAccount.AccountType("client"), "rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH").
		Return("sEdTM1uX8pu2do5XvTnutH6HsouMaM2", nil)

	// Mock SignTransactionNative returns error (no existing blob)
	mockXRP.On("SignTransactionNative", ctx, &txInput, "sEdTM1uX8pu2do5XvTnutH6HsouMaM2", false, nilBlob).
		Return("", "", errors.New("invalid transaction"))

	// Act
	output, err := usecase.Sign(ctx, input)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to sign transaction 018d5e3e-5678-9012-abcd-ef1234567890")
	assert.Equal(t, "", output.NextFilePath)
	mockFileRepo.AssertExpectations(t)
	mockKeyRepo.AssertExpectations(t)
	mockXRP.AssertExpectations(t)
}

func TestSignTransactionUseCase_Sign_EmptyTransactions(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockXRP := new(MockXRPSignClient)
	mockKeyRepo := new(MockXRPAccountKeyRepo)
	mockFileRepo := new(MockTransactionFileRepo)

	usecase := NewSignTransactionUseCase(
		mockXRP,
		mockKeyRepo,
		mockFileRepo,
		domainWallet.WalletTypeSign,
	)

	input := signusecase.SignTransactionInput{
		FilePath: "/path/to/empty.json",
	}

	// Mock GetFileNameType
	mockFileRepo.On("GetFileNameType", input.FilePath).
		Return(&file.FileName{
			ActionType:  domainTx.ActionTypeDeposit,
			TxType:      domainTx.TxTypeUnsigned,
			TxID:        1,
			SignedCount: 0,
		}, nil)

	// Mock ReadXRPJSONFile - file with no transactions
	emptyTxFile := &dtoxrp.XRPTransactionFile{
		Version:      "1.0.0",
		Chain:        "XRP",
		Network:      "testnet",
		CreatedAt:    "2024-02-14T00:00:00Z",
		Transactions: []dtoxrp.XRPTransactionEntry{}, // Empty array
	}
	mockFileRepo.On("ReadXRPJSONFile", input.FilePath).Return(emptyTxFile, nil)

	// Act
	output, err := usecase.Sign(ctx, input)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid transaction file: transactions array cannot be empty")
	assert.Equal(t, "", output.NextFilePath)
	mockFileRepo.AssertExpectations(t)
	// No signing or secret retrieval should happen
	mockXRP.AssertNotCalled(t, "SignTransactionNative")
	mockKeyRepo.AssertNotCalled(t, "GetSecret")
}

func TestSignTransactionUseCase_Sign_MultiSignature_SecondSigner(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockXRP := new(MockXRPSignClient)
	mockKeyRepo := new(MockXRPAccountKeyRepo)
	mockFileRepo := new(MockTransactionFileRepo)

	usecase := NewSignTransactionUseCase(
		mockXRP,
		mockKeyRepo,
		mockFileRepo,
		domainWallet.WalletTypeSign,
	)

	input := signusecase.SignTransactionInput{
		FilePath: "/path/to/signed_1.json",
	}

	// Mock GetFileNameType - reading a file that already has 1 signature
	mockFileRepo.On("GetFileNameType", input.FilePath).
		Return(&file.FileName{
			ActionType:  domainTx.ActionTypePayment,
			TxType:      domainTx.TxTypeSigned, // Already signed once
			TxID:        2,
			SignedCount: 1, // Has 1 signature already
		}, nil)

	// Mock ReadXRPJSONFile - partially signed transaction (1 of 2 signatures)
	txInput := dtoxrp.TxInput{
		Account:            "rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH",
		Destination:        "rLHzPsX6oXkzU9J1kSwq7X4X8J1Xa8AkDC",
		Amount:             "1000000",
		Fee:                "36", // Multi-sig fee: (2+1) * 12
		Sequence:           42,
		LastLedgerSequence: 1000000,
		TransactionType:    "Payment",
	}
	existingSignedBlob := "120000228000000024000000002A61400000000000F424068400000000000024C7" +
		"3210FIRST_SIGNATURE_BLOB"
	partiallySigned := &dtoxrp.XRPTransactionFile{
		Version:   "1.0.0",
		Chain:     "XRP",
		Network:   "testnet",
		CreatedAt: "2024-02-14T00:00:00Z",
		Transactions: []dtoxrp.XRPTransactionEntry{
			{
				UUID:               "018d5e3e-abcd-ef12-3456-789012345678",
				UnsignedData:       txInput,
				SenderAccount:      "rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH",
				SenderAccountType:  "payment",
				SignatureCount:     1,                   // Already has 1 signature
				RequiredSignatures: 2,                   // Needs 2 total
				SignedBlob:         &existingSignedBlob, // Existing signed blob
				IsComplete:         false,               // Not complete yet
			},
		},
	}
	mockFileRepo.On("ReadXRPJSONFile", input.FilePath).Return(partiallySigned, nil)

	// Mock GetSecret for second signer
	mockKeyRepo.On("GetSecret", ctx, domainAccount.AccountType("payment"), "rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH").
		Return("sEdSigner2Secret", nil)

	// Mock SignTransactionNative - should receive existing signed blob to add signature
	// The key here is that we pass the existing blob so signatures can be accumulated
	mockXRP.On(
		"SignTransactionNative",
		ctx,
		&txInput,
		"sEdSigner2Secret",
		true,                // multi-sig
		&existingSignedBlob, // Pass existing blob for signature accumulation
	).Return(
		"5F5E8FAB1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF12345678",
		"120000228000000024000000002A61400000000000F424068400000000000024C7"+
			"3210FIRST_SIGNATURE_BLOB7321SECOND_SIGNATURE_BLOB", // Combined signatures
		nil,
	)

	// Mock CreateFilePath - should increment to signed_2
	mockFileRepo.On("CreateFilePath", domainTx.ActionTypePayment, domainTx.TxTypeSigned, int64(2), 2).
		Return("/path/to/signed_2.json")

	// Mock WriteXRPJSONFile - verify second signature added and complete
	mockFileRepo.On(
		"WriteXRPJSONFile",
		"/path/to/signed_2.json",
		mock.MatchedBy(func(data *dtoxrp.XRPTransactionFile) bool {
			// Second signature: count=2, complete (has all required signatures)
			return data.Transactions[0].SignatureCount == 2 &&
				data.Transactions[0].IsComplete == true &&
				data.Transactions[0].SignedBlob != nil &&
				// Verify the blob contains both signatures
				strings.Contains(*data.Transactions[0].SignedBlob, "FIRST_SIGNATURE_BLOB") &&
				strings.Contains(*data.Transactions[0].SignedBlob, "SECOND_SIGNATURE_BLOB")
		}),
	).Return("/path/to/signed_2.json", nil)

	// Act
	output, err := usecase.Sign(ctx, input)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "/path/to/signed_2.json", output.NextFilePath)
	assert.True(t, output.IsComplete) // Now complete with 2 signatures
	mockFileRepo.AssertExpectations(t)
	mockKeyRepo.AssertExpectations(t)
	mockXRP.AssertExpectations(t)
}
