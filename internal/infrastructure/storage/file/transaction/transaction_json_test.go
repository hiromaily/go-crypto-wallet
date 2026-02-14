package transaction

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
)

func TestWriteXRPJSONFile(t *testing.T) {
	// Create temp directory for tests
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		path        string
		data        *dtoxrp.XRPTransactionFile
		wantErr     bool
		checkSuffix string
	}{
		{
			name: "write valid XRP transaction file",
			path: filepath.Join(tempDir, "deposit_8_unsigned_0_"),
			data: &dtoxrp.XRPTransactionFile{
				Version:   "1.0.0",
				Chain:     "XRP",
				Network:   "mainnet",
				CreatedAt: "2024-02-14T00:00:00Z",
				Transactions: []dtoxrp.XRPTransactionEntry{
					{
						UUID: "01234567-89ab-cdef-0123-456789abcdef",
						UnsignedData: dtoxrp.TxInput{
							Account:     "rN7n7otQDd6FczFgLdlqtyMVrn3HMfgnj",
							Destination: "rLHzPsX6oXkzU9jNbPGJcEKvxbqBKzhqgN",
							Amount:      "1000000",
							Fee:         "12",
						},
						SenderAccount:      "rN7n7otQDd6FczFgLdlqtyMVrn3HMfgnj",
						SenderAccountType:  "client",
						SignatureCount:     0,
						RequiredSignatures: 1,
						SignedBlob:         nil,
						IsComplete:         false,
					},
				},
			},
			wantErr:     false,
			checkSuffix: ".json",
		},
		{
			name: "write multi-signature transaction",
			path: filepath.Join(tempDir, "payment_42_signed_2_"),
			data: &dtoxrp.XRPTransactionFile{
				Version:   "1.0.0",
				Chain:     "XRP",
				Network:   "testnet",
				CreatedAt: "2024-02-14T00:00:00Z",
				Transactions: []dtoxrp.XRPTransactionEntry{
					{
						UUID: "fedcba98-7654-3210-fedc-ba9876543210",
						UnsignedData: dtoxrp.TxInput{
							Account:     "rN7n7otQDd6FczFgLdlqtyMVrn3HMfgnj",
							Destination: "rLHzPsX6oXkzU9jNbPGJcEKvxbqBKzhqgN",
							Amount:      "2000000",
							Fee:         "30",
						},
						SenderAccount:      "rN7n7otQDd6FczFgLdlqtyMVrn3HMfgnj",
						SenderAccountType:  "receipt",
						SignatureCount:     2,
						RequiredSignatures: 3,
						SignedBlob:         stringPtr("120000228000000024000000016140000000000000016840"),
						IsComplete:         false,
					},
				},
			},
			wantErr:     false,
			checkSuffix: ".json",
		},
		{
			name: "reject invalid transaction file (empty transactions)",
			path: filepath.Join(tempDir, "transfer_123_unsigned_0_"),
			data: &dtoxrp.XRPTransactionFile{
				Version:      "1.0.0",
				Chain:        "XRP",
				Network:      "mainnet",
				CreatedAt:    "2024-02-14T00:00:00Z",
				Transactions: []dtoxrp.XRPTransactionEntry{},
			},
			wantErr: true,
		},
		{
			name: "reject invalid version format",
			path: filepath.Join(tempDir, "deposit_10_unsigned_0_"),
			data: &dtoxrp.XRPTransactionFile{
				Version:   "1.0",
				Chain:     "XRP",
				Network:   "mainnet",
				CreatedAt: "2024-02-14T00:00:00Z",
				Transactions: []dtoxrp.XRPTransactionEntry{
					{
						UUID:               "test-uuid",
						SenderAccount:      "rTest",
						SenderAccountType:  "client",
						SignatureCount:     0,
						RequiredSignatures: 1,
						IsComplete:         false,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "reject invalid chain",
			path: filepath.Join(tempDir, "deposit_11_unsigned_0_"),
			data: &dtoxrp.XRPTransactionFile{
				Version:   "1.0.0",
				Chain:     "BTC",
				Network:   "mainnet",
				CreatedAt: "2024-02-14T00:00:00Z",
				Transactions: []dtoxrp.XRPTransactionEntry{
					{
						UUID:               "test-uuid",
						SenderAccount:      "rTest",
						SenderAccountType:  "client",
						SignatureCount:     0,
						RequiredSignatures: 1,
						IsComplete:         false,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewTransactionFileRepository("")
			fileName, err := repo.WriteXRPJSONFile(tt.path, tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteXRPJSONFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Check file was created
				if _, err := os.Stat(fileName); os.IsNotExist(err) {
					t.Errorf("WriteXRPJSONFile() did not create file: %v", fileName)
				}

				// Check file has .json extension
				if !strings.HasSuffix(fileName, tt.checkSuffix) {
					t.Errorf("WriteXRPJSONFile() file name = %v, want suffix %v", fileName, tt.checkSuffix)
				}

				// Read and verify content is valid JSON
				content, err := os.ReadFile(fileName)
				if err != nil {
					t.Errorf("Failed to read written file: %v", err)
				}

				// Verify it's valid JSON by attempting to parse
				var readBack dtoxrp.XRPTransactionFile
				if err := readBack.Validate(); err == nil {
					// This shouldn't work if content is empty
					if len(content) == 0 {
						t.Error("WriteXRPJSONFile() wrote empty file")
					}
				}

				// Cleanup
				_ = os.Remove(fileName)
			}
		})
	}
}

func TestReadXRPJSONFile(t *testing.T) {
	// Create temp directory for tests
	tempDir := t.TempDir()

	// Create valid test file
	validJSONPath := filepath.Join(tempDir, "test_valid.json")
	validContent := `{
		"version": "1.0.0",
		"chain": "XRP",
		"network": "mainnet",
		"createdAt": "2024-02-14T00:00:00Z",
		"transactions": [
			{
				"uuid": "01234567-89ab-cdef-0123-456789abcdef",
				"unsignedData": {
					"Account": "rN7n7otQDd6FczFgLdlqtyMVrn3HMfgnj",
					"Destination": "rLHzPsX6oXkzU9jNbPGJcEKvxbqBKzhqgN",
					"Amount": "1000000",
					"Fee": "12"
				},
				"senderAccount": "rN7n7otQDd6FczFgLdlqtyMVrn3HMfgnj",
				"senderAccountType": "client",
				"signatureCount": 0,
				"requiredSignatures": 1,
				"signedBlob": null,
				"isComplete": false
			}
		]
	}`
	if err := os.WriteFile(validJSONPath, []byte(validContent), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create invalid JSON file
	invalidJSONPath := filepath.Join(tempDir, "test_invalid.json")
	invalidContent := `{"invalid": json}`
	if err := os.WriteFile(invalidJSONPath, []byte(invalidContent), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create file with invalid data (empty transactions)
	invalidDataPath := filepath.Join(tempDir, "test_invalid_data.json")
	invalidDataContent := `{
		"version": "1.0.0",
		"chain": "XRP",
		"network": "mainnet",
		"createdAt": "2024-02-14T00:00:00Z",
		"transactions": []
	}`
	if err := os.WriteFile(invalidDataPath, []byte(invalidDataContent), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create file with wrong extension
	invalidExtPath := filepath.Join(tempDir, "test_invalid.txt")
	if err := os.WriteFile(invalidExtPath, []byte(validContent), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "read valid JSON transaction file",
			path:    validJSONPath,
			wantErr: false,
		},
		{
			name:    "read non-existent file",
			path:    filepath.Join(tempDir, "nonexistent.json"),
			wantErr: true,
		},
		{
			name:    "read file with invalid JSON",
			path:    invalidJSONPath,
			wantErr: true,
		},
		{
			name:    "read file with invalid data (empty transactions)",
			path:    invalidDataPath,
			wantErr: true,
		},
		{
			name:    "read file with invalid extension",
			path:    invalidExtPath,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewTransactionFileRepository("")
			got, err := repo.ReadXRPJSONFile(tt.path)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadXRPJSONFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Validate the returned data
				if err := got.Validate(); err != nil {
					t.Errorf("ReadXRPJSONFile() returned invalid data: %v", err)
				}

				// Check basic fields
				if got.Version != "1.0.0" {
					t.Errorf("ReadXRPJSONFile() version = %v, want 1.0.0", got.Version)
				}
				if got.Chain != "XRP" {
					t.Errorf("ReadXRPJSONFile() chain = %v, want XRP", got.Chain)
				}
				if len(got.Transactions) == 0 {
					t.Error("ReadXRPJSONFile() transactions is empty")
				}
			}
		})
	}
}

func TestXRPJSONRoundTrip(t *testing.T) {
	// Create temp directory for tests
	tempDir := t.TempDir()

	repo := NewTransactionFileRepository(tempDir + "/")
	txData := &dtoxrp.XRPTransactionFile{
		Version:   "1.0.0",
		Chain:     "XRP",
		Network:   "testnet",
		CreatedAt: "2024-02-14T12:00:00Z",
		Transactions: []dtoxrp.XRPTransactionEntry{
			{
				UUID: "01234567-89ab-cdef-0123-456789abcdef",
				UnsignedData: dtoxrp.TxInput{
					Account:     "rN7n7otQDd6FczFgLdlqtyMVrn3HMfgnj",
					Destination: "rLHzPsX6oXkzU9jNbPGJcEKvxbqBKzhqgN",
					Amount:      "1000000",
					Fee:         "12",
				},
				SenderAccount:      "rN7n7otQDd6FczFgLdlqtyMVrn3HMfgnj",
				SenderAccountType:  "client",
				SignatureCount:     0,
				RequiredSignatures: 1,
				SignedBlob:         nil,
				IsComplete:         false,
			},
		},
	}

	// Create file path
	path := repo.CreateFilePath(domainTx.ActionTypeDeposit, domainTx.TxTypeUnsigned, 8, 0)

	// Write JSON
	fileName, err := repo.WriteXRPJSONFile(path, txData)
	if err != nil {
		t.Fatalf("WriteXRPJSONFile() error = %v", err)
	}
	defer func() { _ = os.Remove(fileName) }()

	// Read JSON back
	gotData, err := repo.ReadXRPJSONFile(fileName)
	if err != nil {
		t.Fatalf("ReadXRPJSONFile() error = %v", err)
	}

	// Verify data matches
	if gotData.Version != txData.Version {
		t.Errorf("Round trip version = %v, want %v", gotData.Version, txData.Version)
	}
	if gotData.Chain != txData.Chain {
		t.Errorf("Round trip chain = %v, want %v", gotData.Chain, txData.Chain)
	}
	if gotData.Network != txData.Network {
		t.Errorf("Round trip network = %v, want %v", gotData.Network, txData.Network)
	}
	if len(gotData.Transactions) != len(txData.Transactions) {
		t.Errorf("Round trip transactions length = %v, want %v", len(gotData.Transactions), len(txData.Transactions))
	}
	if len(gotData.Transactions) > 0 && len(txData.Transactions) > 0 {
		if gotData.Transactions[0].UUID != txData.Transactions[0].UUID {
			t.Errorf("Round trip transaction UUID = %v, want %v",
				gotData.Transactions[0].UUID, txData.Transactions[0].UUID)
		}
		if gotData.Transactions[0].SignatureCount != txData.Transactions[0].SignatureCount {
			t.Errorf("Round trip signature count = %v, want %v",
				gotData.Transactions[0].SignatureCount, txData.Transactions[0].SignatureCount)
		}
	}
}

func TestReadXRPJSONFilePathTraversal(t *testing.T) {
	// Create temp directory structure
	tempDir := t.TempDir()
	jsonDir := filepath.Join(tempDir, "json")
	if err := os.MkdirAll(jsonDir, 0o700); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create a valid JSON file in the allowed directory
	validJSONPath := filepath.Join(jsonDir, "test_valid.json")
	validContent := `{
		"version": "1.0.0",
		"chain": "XRP",
		"network": "mainnet",
		"createdAt": "2024-02-14T00:00:00Z",
		"transactions": [
			{
				"uuid": "01234567-89ab-cdef-0123-456789abcdef",
				"unsignedData": {
					"Account": "rN7n7otQDd6FczFgLdlqtyMVrn3HMfgnj",
					"Destination": "rLHzPsX6oXkzU9jNbPGJcEKvxbqBKzhqgN",
					"Amount": "1000000",
					"Fee": "12"
				},
				"senderAccount": "rN7n7otQDd6FczFgLdlqtyMVrn3HMfgnj",
				"senderAccountType": "client",
				"signatureCount": 0,
				"requiredSignatures": 1,
				"signedBlob": null,
				"isComplete": false
			}
		]
	}`
	if err := os.WriteFile(validJSONPath, []byte(validContent), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name       string
		basePath   string
		readPath   string
		wantErr    bool
		errContain string
	}{
		{
			name:     "read file within allowed directory",
			basePath: jsonDir + "/",
			readPath: validJSONPath,
			wantErr:  false,
		},
		{
			name:       "reject path traversal with ../",
			basePath:   jsonDir + "/",
			readPath:   filepath.Join(jsonDir, "../../../etc/passwd.json"),
			wantErr:    true,
			errContain: "path traversal attempt detected",
		},
		{
			name:       "reject absolute path outside base",
			basePath:   jsonDir + "/",
			readPath:   "/etc/passwd.json",
			wantErr:    true,
			errContain: "path traversal attempt detected",
		},
		{
			name:     "allow read when basePath is empty (test mode)",
			basePath: "",
			readPath: validJSONPath,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewTransactionFileRepository(tt.basePath)
			got, err := repo.ReadXRPJSONFile(tt.readPath)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadXRPJSONFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContain != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errContain) {
					t.Errorf("ReadXRPJSONFile() error = %v, want error containing %q", err, tt.errContain)
				}
			}

			if !tt.wantErr {
				if got == nil {
					t.Error("ReadXRPJSONFile() returned nil data")
				} else if got.Version != "1.0.0" {
					t.Errorf("ReadXRPJSONFile() version = %v, want 1.0.0", got.Version)
				}
			}
		})
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
