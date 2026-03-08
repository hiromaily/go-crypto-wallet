package transaction

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
)

// Anvil pre-funded account addresses (valid EIP-55 checksum).
const (
	anvilAddr0 = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	anvilAddr1 = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
	zeroAddr   = "0x0000000000000000000000000000000000000000"
)

// validMultisigUnsigned returns a minimal valid unsigned ETH multisig transaction file.
func validMultisigUnsigned() *dtoeth.ETHMultisigTransactionFile {
	return &dtoeth.ETHMultisigTransactionFile{
		Version:        1,
		TxType:         "unsigned",
		UUID:           "550e8400-e29b-41d4-a716-446655440000",
		SafeAddress:    anvilAddr0,
		To:             anvilAddr1,
		Value:          "100000000000000000",
		Data:           "0x",
		Operation:      0,
		SafeTxGas:      "0",
		BaseGas:        "0",
		GasPrice:       "0",
		GasToken:       zeroAddr,
		RefundReceiver: zeroAddr,
		Nonce:          "0",
		ChainID:        31337,
		SafeTxHash:     "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
		Threshold:      2,
		Signatures:     []dtoeth.SignEntry{},
	}
}

// validMultisigSigned returns a valid fully-signed ETH multisig transaction file.
func validMultisigSigned() *dtoeth.ETHMultisigTransactionFile {
	f := validMultisigUnsigned()
	f.TxType = "signed"
	f.Signatures = []dtoeth.SignEntry{
		{
			SignerAddress: anvilAddr0,
			SignatureHex:  "0x" + strings.Repeat("ab", 65),
		},
		{
			SignerAddress: anvilAddr1,
			SignatureHex:  "0x" + strings.Repeat("cd", 65),
		},
	}
	return f
}

func TestCreateMultisigFilePath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		basePath    string
		actionType  domainTx.ActionType
		uuid        string
		signedCount int
		wantSuffix  string
	}{
		{
			name:        "unsigned proposal (count 0)",
			basePath:    "/data/tx/",
			actionType:  domainTx.ActionTypeDeposit,
			uuid:        "550e8400-e29b-41d4-a716-446655440000",
			signedCount: 0,
			wantSuffix:  "deposit_multisig_550e8400-e29b-41d4-a716-446655440000_0.json",
		},
		{
			name:        "after first signing (count 1)",
			basePath:    "/data/tx/",
			actionType:  domainTx.ActionTypePayment,
			uuid:        "aaaabbbb-cccc-dddd-eeee-ffff00001111",
			signedCount: 1,
			wantSuffix:  "payment_multisig_aaaabbbb-cccc-dddd-eeee-ffff00001111_1.json",
		},
		{
			name:        "fully signed (count 2)",
			basePath:    "/data/tx/",
			actionType:  domainTx.ActionTypeTransfer,
			uuid:        "00000000-0000-0000-0000-000000000001",
			signedCount: 2,
			wantSuffix:  "transfer_multisig_00000000-0000-0000-0000-000000000001_2.json",
		},
		{
			name:        "empty base path",
			basePath:    "",
			actionType:  domainTx.ActionTypeDeposit,
			uuid:        "test-uuid",
			signedCount: 0,
			wantSuffix:  "deposit_multisig_test-uuid_0.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := NewTransactionFileRepository(tt.basePath)
			got := repo.CreateMultisigFilePath(tt.actionType, tt.uuid, tt.signedCount)
			assert.True(t, strings.HasSuffix(got, tt.wantSuffix),
				"path %q should end with %q", got, tt.wantSuffix)
			assert.True(t, strings.HasSuffix(got, ".json"),
				"path should have .json extension")
		})
	}
}

func TestWriteETHMultisigJSONFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    *dtoeth.ETHMultisigTransactionFile
		wantErr bool
	}{
		{
			name:    "write valid unsigned multisig file",
			data:    validMultisigUnsigned(),
			wantErr: false,
		},
		{
			name:    "write valid signed multisig file",
			data:    validMultisigSigned(),
			wantErr: false,
		},
		{
			name:    "reject nil file",
			data:    nil,
			wantErr: true,
		},
		{
			name: "reject invalid file (version 0)",
			data: &dtoeth.ETHMultisigTransactionFile{
				Version: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tempDir := t.TempDir()
			repo := NewTransactionFileRepository("")
			path := filepath.Join(tempDir, "deposit_multisig_test-uuid_0.json")

			writtenPath, err := repo.WriteETHMultisigJSONFile(path, tt.data)

			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, path, writtenPath)

			// File must exist and contain valid JSON.
			content, err := os.ReadFile(writtenPath)
			require.NoError(t, err)
			assert.NotEmpty(t, content)
		})
	}
}

func TestReadETHMultisigJSONFile(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	repo := NewTransactionFileRepository("")

	// Pre-write a valid file to read back.
	validPath := filepath.Join(tempDir, "deposit_multisig_uuid1_0.json")
	_, err := repo.WriteETHMultisigJSONFile(validPath, validMultisigUnsigned())
	require.NoError(t, err)

	// Invalid JSON file.
	invalidJSONPath := filepath.Join(tempDir, "invalid.json")
	err = os.WriteFile(invalidJSONPath, []byte(`{"invalid": json}`), 0o644)
	require.NoError(t, err)

	// File with invalid extension.
	invalidExtPath := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(invalidExtPath, []byte(`{}`), 0o644)
	require.NoError(t, err)

	tests := []struct {
		name    string
		path    string
		wantErr bool
		checkFn func(t *testing.T, got *dtoeth.ETHMultisigTransactionFile)
	}{
		{
			name:    "read valid unsigned file",
			path:    validPath,
			wantErr: false,
			checkFn: func(t *testing.T, got *dtoeth.ETHMultisigTransactionFile) {
				t.Helper()
				assert.Equal(t, 1, got.Version)
				assert.Equal(t, "unsigned", got.TxType)
				assert.Equal(t, uint64(31337), got.ChainID)
				assert.Equal(t, 2, got.Threshold)
			},
		},
		{
			name:    "error on non-existent file",
			path:    filepath.Join(tempDir, "nonexistent.json"),
			wantErr: true,
		},
		{
			name:    "error on invalid JSON",
			path:    invalidJSONPath,
			wantErr: true,
		},
		{
			name:    "error on invalid extension",
			path:    invalidExtPath,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := repo.ReadETHMultisigJSONFile(tt.path)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			if tt.checkFn != nil {
				tt.checkFn(t, got)
			}
		})
	}
}

func TestMultisigJSONRoundTrip(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	repo := NewTransactionFileRepository(tempDir + "/")

	original := validMultisigUnsigned()
	path := repo.CreateMultisigFilePath(domainTx.ActionTypeDeposit, original.UUID, 0)

	writtenPath, err := repo.WriteETHMultisigJSONFile(path, original)
	require.NoError(t, err)
	defer func() { _ = os.Remove(writtenPath) }()

	got, err := repo.ReadETHMultisigJSONFile(writtenPath)
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, original.Version, got.Version)
	assert.Equal(t, original.TxType, got.TxType)
	assert.Equal(t, original.UUID, got.UUID)
	assert.Equal(t, original.SafeAddress, got.SafeAddress)
	assert.Equal(t, original.To, got.To)
	assert.Equal(t, original.Value, got.Value)
	assert.Equal(t, original.ChainID, got.ChainID)
	assert.Equal(t, original.SafeTxHash, got.SafeTxHash)
	assert.Equal(t, original.Threshold, got.Threshold)
	assert.Empty(t, got.Signatures)
}
