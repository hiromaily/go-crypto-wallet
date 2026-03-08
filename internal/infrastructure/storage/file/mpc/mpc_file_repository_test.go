package mpc_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
	mpcstore "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/mpc"
)

// validMPCFile returns a minimal valid unsigned ETHMPCTransactionFile for tests.
// Addresses are valid EIP-55 checksummed.
func validMPCFile() *dtoeth.ETHMPCTransactionFile {
	return &dtoeth.ETHMPCTransactionFile{
		Version:              1,
		TxType:               "unsigned",
		UUID:                 "550e8400-e29b-41d4-a716-446655440000",
		ActionType:           "payment",
		From:                 "0x0742D35CC6634c0532925A3b844bc9E7595f0Beb",
		To:                   "0x1234567890AbcdEF1234567890aBcdef12345678",
		Value:                "1000000000000000000",
		Nonce:                0,
		GasLimit:             21000,
		ChainID:              1,
		MaxFeePerGas:         "30000000000",
		MaxPriorityFeePerGas: "2000000000",
		TxHash:               "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
		RawTxHex:             "0x02f86c",
		Threshold:            2,
		PartyIDs:             []string{"node-1", "node-2", "node-3"},
	}
}

func TestCreateMPCFilePath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		basePath   string
		actionType string
		uuid       string
		isSigned   bool
		wantSuffix string
	}{
		{
			name:       "unsigned file",
			basePath:   "/data/tx/",
			actionType: "payment",
			uuid:       "550e8400-e29b-41d4-a716-446655440000",
			isSigned:   false,
			wantSuffix: "payment_mpc_550e8400-e29b-41d4-a716-446655440000.json",
		},
		{
			name:       "signed file",
			basePath:   "/data/tx/",
			actionType: "payment",
			uuid:       "550e8400-e29b-41d4-a716-446655440000",
			isSigned:   true,
			wantSuffix: "payment_mpc_550e8400-e29b-41d4-a716-446655440000_signed.json",
		},
		{
			name:       "deposit unsigned",
			basePath:   "",
			actionType: "deposit",
			uuid:       "test-uuid",
			isSigned:   false,
			wantSuffix: "deposit_mpc_test-uuid.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := mpcstore.NewMPCFileRepository(tt.basePath)
			got := repo.CreateMPCFilePath(tt.actionType, tt.uuid, tt.isSigned)
			assert.True(t, len(got) >= len(tt.wantSuffix), "path %q should contain %q", got, tt.wantSuffix)
			assert.Contains(t, got, tt.wantSuffix, "path %q should end with %q", got, tt.wantSuffix)
		})
	}
}

func TestWriteETHMPCJSONFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		file     *dtoeth.ETHMPCTransactionFile
		isSigned bool
		wantErr  bool
	}{
		{
			name:     "write valid unsigned file",
			file:     validMPCFile(),
			isSigned: false,
			wantErr:  false,
		},
		{
			name: "write valid signed file",
			file: func() *dtoeth.ETHMPCTransactionFile {
				f := validMPCFile()
				f.TxType = "signed"
				f.SignedTxHex = "0x02f86c8204d28459682f008459682f14825208"
				return f
			}(),
			isSigned: true,
			wantErr:  false,
		},
		{
			name:     "reject nil file",
			file:     nil,
			isSigned: false,
			wantErr:  true,
		},
		{
			name: "reject invalid file (version 0)",
			file: &dtoeth.ETHMPCTransactionFile{
				Version: 0,
			},
			isSigned: false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tempDir := t.TempDir()
			repo := mpcstore.NewMPCFileRepository(tempDir + "/")

			writtenPath, err := repo.WriteETHMPCJSONFile(tt.file, tt.isSigned)

			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, writtenPath)
			assert.True(t, filepath.IsAbs(writtenPath) || len(writtenPath) > 0)

			// File must exist and contain content.
			content, err := os.ReadFile(writtenPath)
			require.NoError(t, err)
			assert.NotEmpty(t, content)
		})
	}
}

func TestReadETHMPCJSONFile(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	repo := mpcstore.NewMPCFileRepository("")

	// Write a valid file to read back.
	f := validMPCFile()
	writtenPath, err := repo.WriteETHMPCJSONFile(f, false)
	require.NoError(t, err)

	// Write an invalid JSON file.
	invalidJSONPath := filepath.Join(tempDir, "invalid.json")
	err = os.WriteFile(invalidJSONPath, []byte(`{"invalid": json}`), 0o644)
	require.NoError(t, err)

	// Write a file with invalid extension.
	invalidExtPath := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(invalidExtPath, []byte(`{}`), 0o644)
	require.NoError(t, err)

	tests := []struct {
		name    string
		path    string
		wantErr bool
		checkFn func(t *testing.T, got *dtoeth.ETHMPCTransactionFile)
	}{
		{
			name:    "read valid file back",
			path:    writtenPath,
			wantErr: false,
			checkFn: func(t *testing.T, got *dtoeth.ETHMPCTransactionFile) {
				t.Helper()
				assert.Equal(t, f.Version, got.Version)
				assert.Equal(t, f.TxType, got.TxType)
				assert.Equal(t, f.UUID, got.UUID)
				assert.Equal(t, f.From, got.From)
				assert.Equal(t, f.To, got.To)
				assert.Equal(t, f.ChainID, got.ChainID)
				assert.Equal(t, f.Threshold, got.Threshold)
				assert.Equal(t, f.PartyIDs, got.PartyIDs)
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
			got, err := repo.ReadETHMPCJSONFile(tt.path)
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

func TestMPCFileRepositoryRoundTrip(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	repo := mpcstore.NewMPCFileRepository(tempDir + "/")

	original := validMPCFile()
	writtenPath, err := repo.WriteETHMPCJSONFile(original, false)
	require.NoError(t, err)

	got, err := repo.ReadETHMPCJSONFile(writtenPath)
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, original.Version, got.Version)
	assert.Equal(t, original.TxType, got.TxType)
	assert.Equal(t, original.UUID, got.UUID)
	assert.Equal(t, original.ActionType, got.ActionType)
	assert.Equal(t, original.From, got.From)
	assert.Equal(t, original.To, got.To)
	assert.Equal(t, original.Value, got.Value)
	assert.Equal(t, original.ChainID, got.ChainID)
	assert.Equal(t, original.Threshold, got.Threshold)
	assert.Equal(t, original.PartyIDs, got.PartyIDs)
	assert.Equal(t, original.TxHash, got.TxHash)
	assert.Equal(t, original.RawTxHex, got.RawTxHex)
	assert.Empty(t, got.SignedTxHex)
}
