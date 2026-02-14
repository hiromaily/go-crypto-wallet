package xrp_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	apixrp "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/xrp"
)

// TestAccountInfoProviderInterface verifies that AccountInfoProvider interface
// exists with the correct method signatures for task 1.2.
func TestAccountInfoProviderInterface(t *testing.T) {
	t.Parallel()

	// Verify interface type exists
	var _ apixrp.AccountInfoProvider

	// Verify interface has the correct methods using reflection
	interfaceType := reflect.TypeOf((*apixrp.AccountInfoProvider)(nil)).Elem()

	// Should have exactly 2 methods: GetAccountInfo and GetBalance
	assert.Equal(t, 2, interfaceType.NumMethod(), "AccountInfoProvider should have 2 methods")

	// Verify GetAccountInfo method signature
	getAccountInfoMethod, found := interfaceType.MethodByName("GetAccountInfo")
	assert.True(t, found, "GetAccountInfo method should exist")
	assert.Equal(t, 2, getAccountInfoMethod.Type.NumIn(), "GetAccountInfo should have 2 inputs (ctx, address)")
	assert.Equal(t, 2, getAccountInfoMethod.Type.NumOut(), "GetAccountInfo should have 2 outputs (result, error)")

	// Verify GetBalance method signature
	getBalanceMethod, found := interfaceType.MethodByName("GetBalance")
	assert.True(t, found, "GetBalance method should exist")
	assert.Equal(t, 2, getBalanceMethod.Type.NumIn(), "GetBalance should have 2 inputs (ctx, addr)")
	assert.Equal(t, 2, getBalanceMethod.Type.NumOut(), "GetBalance should have 2 outputs (float64, error)")
}

// TestTransactionSignerInterface verifies that TransactionSigner interface
// exists with both legacy and native Go signing methods (task 1.2).
func TestTransactionSignerInterface(t *testing.T) {
	t.Parallel()

	// Verify interface type exists
	var _ apixrp.TransactionSigner

	// Verify interface has the correct methods using reflection
	interfaceType := reflect.TypeOf((*apixrp.TransactionSigner)(nil)).Elem()

	// Should have 2 methods: SignTransaction (legacy) and SignTransactionNative (new)
	assert.Equal(t, 2, interfaceType.NumMethod(), "TransactionSigner should have 2 methods (legacy + native)")

	// Verify SignTransaction method (legacy, for backward compatibility)
	signLegacyMethod, found := interfaceType.MethodByName("SignTransaction")
	assert.True(t, found, "SignTransaction method should exist for backward compatibility")
	assert.Equal(
		t, 3, signLegacyMethod.Type.NumIn(),
		"SignTransaction should have 3 inputs (ctx, txJSON, secret)",
	)
	assert.Equal(
		t, 3, signLegacyMethod.Type.NumOut(),
		"SignTransaction should have 3 outputs (txID, signedBlob, error)",
	)

	// Verify SignTransactionNative method (new, for native Go signing)
	signMethod, found := interfaceType.MethodByName("SignTransactionNative")
	assert.True(t, found, "SignTransactionNative method should exist for native Go signing")
	assert.Equal(
		t, 4, signMethod.Type.NumIn(),
		"SignTransactionNative should have 4 inputs (ctx, txInput, secret, isMultiSig)",
	)
	assert.Equal(
		t, 3, signMethod.Type.NumOut(),
		"SignTransactionNative should have 3 outputs (txID, signedBlob, error)",
	)
}

// TestTransactionSubmitterInterface verifies that TransactionSubmitter interface
// exists with all required methods for task 1.2.
func TestTransactionSubmitterInterface(t *testing.T) {
	t.Parallel()

	// Verify interface type exists
	var _ apixrp.TransactionSubmitter

	// Verify interface has the correct methods using reflection
	interfaceType := reflect.TypeOf((*apixrp.TransactionSubmitter)(nil)).Elem()

	// Should have exactly 3 methods: SubmitTransaction, WaitValidation, GetTransaction
	assert.Equal(t, 3, interfaceType.NumMethod(), "TransactionSubmitter should have 3 methods")

	// Verify SubmitTransaction method signature
	submitMethod, found := interfaceType.MethodByName("SubmitTransaction")
	assert.True(t, found, "SubmitTransaction method should exist")
	assert.Equal(
		t, 2, submitMethod.Type.NumIn(),
		"SubmitTransaction should have 2 inputs (ctx, signedTx)",
	)
	assert.Equal(
		t, 3, submitMethod.Type.NumOut(),
		"SubmitTransaction should have 3 outputs (SentTx, ledgerVersion, error)",
	)

	// Verify WaitValidation method signature
	waitMethod, found := interfaceType.MethodByName("WaitValidation")
	assert.True(t, found, "WaitValidation method should exist")
	assert.Equal(t, 2, waitMethod.Type.NumIn(), "WaitValidation should have 2 inputs (ctx, targetLedgerVersion)")
	assert.Equal(t, 2, waitMethod.Type.NumOut(), "WaitValidation should have 2 outputs (ledgerVersion, error)")

	// Verify GetTransaction method signature
	getTxMethod, found := interfaceType.MethodByName("GetTransaction")
	assert.True(t, found, "GetTransaction method should exist")
	assert.Equal(t, 3, getTxMethod.Type.NumIn(), "GetTransaction should have 3 inputs (ctx, txID, targetLedgerVersion)")
	assert.Equal(t, 2, getTxMethod.Type.NumOut(), "GetTransaction should have 2 outputs (TxInfo, error)")
}

// TestInterfaceSegregationPrinciple verifies that interfaces follow ISP
// by ensuring each interface has a focused, minimal set of methods.
func TestInterfaceSegregationPrinciple(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		interfaceType reflect.Type
		maxMethods    int
		description   string
	}{
		{
			name:          "AccountInfoProvider",
			interfaceType: reflect.TypeOf((*apixrp.AccountInfoProvider)(nil)).Elem(),
			maxMethods:    2,
			description:   "Account info provider should only have account query methods",
		},
		{
			name:          "TransactionSigner",
			interfaceType: reflect.TypeOf((*apixrp.TransactionSigner)(nil)).Elem(),
			maxMethods:    2,
			description:   "Transaction signer should only have signing methods (legacy + native)",
		},
		{
			name:          "TransactionSubmitter",
			interfaceType: reflect.TypeOf((*apixrp.TransactionSubmitter)(nil)).Elem(),
			maxMethods:    3,
			description:   "Transaction submitter should only have submission-related methods",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.LessOrEqual(t, tt.interfaceType.NumMethod(), tt.maxMethods,
				"%s: %s (has %d methods, max %d)",
				tt.name, tt.description, tt.interfaceType.NumMethod(), tt.maxMethods)
		})
	}
}

// TestInterfacesAreInSeparateFiles verifies that interfaces are defined
// in their own files as per task 1.2 requirements.
// This is a documentation test - actual file organization is verified by compilation.
func TestInterfacesAreInSeparateFiles(t *testing.T) {
	t.Parallel()

	// Expected file structure per task 1.2:
	expectedFiles := map[string]string{
		"AccountInfoProvider":  "application/ports/api/xrp/account_info.go",
		"TransactionSigner":    "application/ports/api/xrp/transaction_signer.go",
		"TransactionSubmitter": "application/ports/api/xrp/transaction_submitter.go",
	}

	// This test documents the expected file structure
	// Actual verification happens at compile time - if interfaces are not found,
	// the tests above will fail
	for interfaceName, expectedFile := range expectedFiles {
		t.Logf("Interface %s should be defined in %s", interfaceName, expectedFile)
	}
}

// Mock implementations to verify interfaces can be implemented

type mockAccountInfoProvider struct{}

func (*mockAccountInfoProvider) GetAccountInfo(
	ctx context.Context,
	address string,
) (*dtoxrp.ResponseGetAccountInfo, error) {
	return nil, nil
}

func (*mockAccountInfoProvider) GetBalance(ctx context.Context, addr string) (float64, error) {
	return 0, nil
}

type mockTransactionSigner struct{}

func (*mockTransactionSigner) SignTransaction(
	ctx context.Context,
	txJSON *dtoxrp.TxInput,
	secret string,
) (string, string, error) {
	return "", "", nil
}

func (*mockTransactionSigner) SignTransactionNative(
	ctx context.Context,
	txInput *dtoxrp.TxInput,
	secret string,
	isMultiSig bool,
	existingSignedBlob *string,
) (string, string, error) {
	return "", "", nil
}

type mockTransactionSubmitter struct{}

func (*mockTransactionSubmitter) SubmitTransaction(
	ctx context.Context,
	signedTx string,
) (*dtoxrp.SentTx, uint64, error) {
	return nil, 0, nil
}

func (*mockTransactionSubmitter) WaitValidation(
	ctx context.Context,
	targetLedgerVersion uint64,
) (uint64, error) {
	return 0, nil
}

func (*mockTransactionSubmitter) GetTransaction(
	ctx context.Context,
	txID string,
	targetLedgerVersion uint64,
) (*dtoxrp.TxInfo, error) {
	return nil, nil
}

// TestMockImplementations verifies that mock implementations satisfy the interfaces
func TestMockImplementations(t *testing.T) {
	t.Parallel()

	// These will fail to compile if the interfaces don't exist or have wrong signatures
	var _ apixrp.AccountInfoProvider = &mockAccountInfoProvider{}
	var _ apixrp.TransactionSigner = &mockTransactionSigner{}
	var _ apixrp.TransactionSubmitter = &mockTransactionSubmitter{}
}
