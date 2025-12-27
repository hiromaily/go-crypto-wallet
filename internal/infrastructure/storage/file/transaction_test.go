package file

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
)

func TestCreateFilePath(t *testing.T) {
	tests := []struct {
		name        string
		basePath    string
		actionType  domainTx.ActionType
		txType      domainTx.TxType
		txID        int64
		signedCount int
		wantPrefix  string
	}{
		{
			name:        "deposit unsigned",
			basePath:    "./data/tx/btc/",
			actionType:  domainTx.ActionTypeDeposit,
			txType:      domainTx.TxTypeUnsigned,
			txID:        8,
			signedCount: 0,
			wantPrefix:  "./data/tx/btc/deposit_8_unsigned_0_",
		},
		{
			name:        "payment signed",
			basePath:    "./data/tx/btc/",
			actionType:  domainTx.ActionTypePayment,
			txType:      domainTx.TxTypeSigned,
			txID:        42,
			signedCount: 2,
			wantPrefix:  "./data/tx/btc/payment_42_signed_2_",
		},
		{
			name:        "transfer unsigned",
			basePath:    "./data/tx/eth/",
			actionType:  domainTx.ActionTypeTransfer,
			txType:      domainTx.TxTypeUnsigned,
			txID:        123,
			signedCount: 0,
			wantPrefix:  "./data/tx/eth/transfer_123_unsigned_0_",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewTransactionFileRepository(tt.basePath)
			got := repo.CreateFilePath(tt.actionType, tt.txType, tt.txID, tt.signedCount)

			if !strings.HasPrefix(got, tt.wantPrefix) {
				t.Errorf("CreateFilePath() = %v, want prefix %v", got, tt.wantPrefix)
			}

			// Verify format: {actionType}_{txID}_{txType}_{signedCount}_{timestamp}
			// Should NOT have .psbt extension yet (added by WritePSBTFile)
			parts := strings.Split(strings.TrimPrefix(got, tt.basePath), "_")
			if len(parts) != 5 {
				t.Errorf("CreateFilePath() returned path with %d parts, want 5", len(parts))
			}
		})
	}
}

func TestGetFileNameType(t *testing.T) {
	tests := []struct {
		name            string
		filePath        string
		wantActionType  domainTx.ActionType
		wantTxType      domainTx.TxType
		wantTxID        int64
		wantSignedCount int
		wantErr         bool
	}{
		{
			name:            "valid PSBT file with extension",
			filePath:        "./data/tx/btc/deposit_8_unsigned_0_1534744535097796209.psbt",
			wantActionType:  domainTx.ActionTypeDeposit,
			wantTxType:      domainTx.TxTypeUnsigned,
			wantTxID:        8,
			wantSignedCount: 0,
			wantErr:         false,
		},
		{
			name:            "valid file without extension (legacy)",
			filePath:        "./data/tx/btc/payment_42_signed_2_1534744535097796209",
			wantActionType:  domainTx.ActionTypePayment,
			wantTxType:      domainTx.TxTypeSigned,
			wantTxID:        42,
			wantSignedCount: 2,
			wantErr:         false,
		},
		{
			name:            "just filename with .psbt",
			filePath:        "transfer_123_unsigned_0_1534744535097796209.psbt",
			wantActionType:  domainTx.ActionTypeTransfer,
			wantTxType:      domainTx.TxTypeUnsigned,
			wantTxID:        123,
			wantSignedCount: 0,
			wantErr:         false,
		},
		{
			name:     "invalid format - too few parts",
			filePath: "./data/tx/btc/deposit_8_unsigned.psbt",
			wantErr:  true,
		},
		{
			name:     "invalid format - too many parts",
			filePath: "./data/tx/btc/deposit_8_unsigned_0_123_extra.psbt",
			wantErr:  true,
		},
		{
			name:     "invalid action type",
			filePath: "./data/tx/btc/invalid_8_unsigned_0_1534744535097796209.psbt",
			wantErr:  true,
		},
		{
			name:     "invalid tx type",
			filePath: "./data/tx/btc/deposit_8_invalid_0_1534744535097796209.psbt",
			wantErr:  true,
		},
		{
			name:     "invalid txID - not a number",
			filePath: "./data/tx/btc/deposit_abc_unsigned_0_1534744535097796209.psbt",
			wantErr:  true,
		},
		{
			name:     "invalid signedCount - not a number",
			filePath: "./data/tx/btc/deposit_8_unsigned_abc_1534744535097796209.psbt",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewTransactionFileRepository("")
			got, err := repo.GetFileNameType(tt.filePath)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileNameType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.ActionType != tt.wantActionType {
					t.Errorf("GetFileNameType().ActionType = %v, want %v", got.ActionType, tt.wantActionType)
				}
				if got.TxType != tt.wantTxType {
					t.Errorf("GetFileNameType().TxType = %v, want %v", got.TxType, tt.wantTxType)
				}
				if got.TxID != tt.wantTxID {
					t.Errorf("GetFileNameType().TxID = %v, want %v", got.TxID, tt.wantTxID)
				}
				if got.SignedCount != tt.wantSignedCount {
					t.Errorf("GetFileNameType().SignedCount = %v, want %v", got.SignedCount, tt.wantSignedCount)
				}
			}
		})
	}
}

func TestValidateFilePath(t *testing.T) {
	tests := []struct {
		name            string
		filePath        string
		expectedTxType  domainTx.TxType
		wantActionType  domainTx.ActionType
		wantTxType      domainTx.TxType
		wantTxID        int64
		wantSignedCount int
		wantErr         bool
	}{
		{
			name:            "valid unsigned PSBT",
			filePath:        "./data/tx/btc/deposit_8_unsigned_0_1534744535097796209.psbt",
			expectedTxType:  domainTx.TxTypeUnsigned,
			wantActionType:  domainTx.ActionTypeDeposit,
			wantTxType:      domainTx.TxTypeUnsigned,
			wantTxID:        8,
			wantSignedCount: 0,
			wantErr:         false,
		},
		{
			name:            "valid signed PSBT",
			filePath:        "./data/tx/btc/payment_42_signed_2_1534744535097796209.psbt",
			expectedTxType:  domainTx.TxTypeSigned,
			wantActionType:  domainTx.ActionTypePayment,
			wantTxType:      domainTx.TxTypeSigned,
			wantTxID:        42,
			wantSignedCount: 2,
			wantErr:         false,
		},
		{
			name:           "mismatched tx type",
			filePath:       "./data/tx/btc/deposit_8_unsigned_0_1534744535097796209.psbt",
			expectedTxType: domainTx.TxTypeSigned,
			wantErr:        true,
		},
		{
			name:           "invalid file path",
			filePath:       "./data/tx/btc/invalid.psbt",
			expectedTxType: domainTx.TxTypeUnsigned,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewTransactionFileRepository("")
			actionType, txType, txID, signedCount, err := repo.ValidateFilePath(tt.filePath, tt.expectedTxType)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFilePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if actionType != tt.wantActionType {
					t.Errorf("ValidateFilePath() actionType = %v, want %v", actionType, tt.wantActionType)
				}
				if txType != tt.wantTxType {
					t.Errorf("ValidateFilePath() txType = %v, want %v", txType, tt.wantTxType)
				}
				if txID != tt.wantTxID {
					t.Errorf("ValidateFilePath() txID = %v, want %v", txID, tt.wantTxID)
				}
				if signedCount != tt.wantSignedCount {
					t.Errorf("ValidateFilePath() signedCount = %v, want %v", signedCount, tt.wantSignedCount)
				}
			}
		})
	}
}

func TestWritePSBTFile(t *testing.T) {
	// Create temp directory for tests
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		path        string
		psbtBase64  string
		wantErr     bool
		checkSuffix string
	}{
		{
			name:        "write valid PSBT",
			path:        filepath.Join(tempDir, "deposit_8_unsigned_0_"),
			psbtBase64:  "cHNidP8BAHECAAAAAe3o6gAHAAc0+4ywFkHaE8nzN/example+base64+data==",
			wantErr:     false,
			checkSuffix: ".psbt",
		},
		{
			name:        "write empty PSBT",
			path:        filepath.Join(tempDir, "payment_42_signed_2_"),
			psbtBase64:  "",
			wantErr:     false,
			checkSuffix: ".psbt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewTransactionFileRepository("")
			fileName, err := repo.WritePSBTFile(tt.path, tt.psbtBase64)

			if (err != nil) != tt.wantErr {
				t.Errorf("WritePSBTFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Check file was created
				if _, err := os.Stat(fileName); os.IsNotExist(err) {
					t.Errorf("WritePSBTFile() did not create file: %v", fileName)
				}

				// Check file has .psbt extension
				if !strings.HasSuffix(fileName, tt.checkSuffix) {
					t.Errorf("WritePSBTFile() file name = %v, want suffix %v", fileName, tt.checkSuffix)
				}

				// Read and verify content
				content, err := os.ReadFile(fileName) //nolint:gosec
				if err != nil {
					t.Errorf("Failed to read written file: %v", err)
				}
				if string(content) != tt.psbtBase64 {
					t.Errorf("WritePSBTFile() content = %v, want %v", string(content), tt.psbtBase64)
				}

				// Cleanup
				_ = os.Remove(fileName)
			}
		})
	}
}

func TestReadPSBTFile(t *testing.T) {
	// Create temp directory for tests
	tempDir := t.TempDir()

	// Create test files
	validPSBTPath := filepath.Join(tempDir, "test_valid.psbt")
	validPSBTContent := "cHNidP8BAHECAAAAAe3o6gAHAAc0+4ywFkHaE8nzN/example+base64+data=="
	if err := os.WriteFile(validPSBTPath, []byte(validPSBTContent), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	invalidExtPath := filepath.Join(tempDir, "test_invalid.txt")
	if err := os.WriteFile(invalidExtPath, []byte("test"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{
			name:    "read valid PSBT file",
			path:    validPSBTPath,
			want:    validPSBTContent,
			wantErr: false,
		},
		{
			name:    "read non-existent file",
			path:    filepath.Join(tempDir, "nonexistent.psbt"),
			want:    "",
			wantErr: true,
		},
		{
			name:    "read file with invalid extension",
			path:    invalidExtPath,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewTransactionFileRepository("")
			got, err := repo.ReadPSBTFile(tt.path)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadPSBTFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("ReadPSBTFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPSBTRoundTrip(t *testing.T) {
	// Create temp directory for tests
	tempDir := t.TempDir()

	repo := NewTransactionFileRepository(tempDir + "/")
	psbtData := "cHNidP8BAHECAAAAAe3o6gAHAAc0+4ywFkHaE8nzN/example+base64+data=="

	// Create file path
	path := repo.CreateFilePath(domainTx.ActionTypeDeposit, domainTx.TxTypeUnsigned, 8, 0)

	// Write PSBT
	fileName, err := repo.WritePSBTFile(path, psbtData)
	if err != nil {
		t.Fatalf("WritePSBTFile() error = %v", err)
	}
	defer func() { _ = os.Remove(fileName) }()

	// Read PSBT back
	gotData, err := repo.ReadPSBTFile(fileName)
	if err != nil {
		t.Fatalf("ReadPSBTFile() error = %v", err)
	}

	if gotData != psbtData {
		t.Errorf("Round trip failed: got %v, want %v", gotData, psbtData)
	}

	// Validate file path
	actionType, txType, txID, signedCount, err := repo.ValidateFilePath(fileName, domainTx.TxTypeUnsigned)
	if err != nil {
		t.Fatalf("ValidateFilePath() error = %v", err)
	}

	if actionType != domainTx.ActionTypeDeposit {
		t.Errorf("ValidateFilePath() actionType = %v, want %v", actionType, domainTx.ActionTypeDeposit)
	}
	if txType != domainTx.TxTypeUnsigned {
		t.Errorf("ValidateFilePath() txType = %v, want %v", txType, domainTx.TxTypeUnsigned)
	}
	if txID != 8 {
		t.Errorf("ValidateFilePath() txID = %v, want %v", txID, 8)
	}
	if signedCount != 0 {
		t.Errorf("ValidateFilePath() signedCount = %v, want %v", signedCount, 0)
	}
}
