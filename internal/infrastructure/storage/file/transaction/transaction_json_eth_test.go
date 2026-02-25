package transaction

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
)

// validETHUnsignedLegacy returns a valid unsigned legacy ETH transaction file.
func validETHUnsignedLegacy() *dtoeth.ETHTransactionFile {
	return &dtoeth.ETHTransactionFile{
		Version:   1,
		TxType:    "unsigned",
		EthTxType: 0,
		ChainID:   1,
		Nonce:     42,
		From:      "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
		To:        "0x1234567890abcdef1234567890abcdef12345678",
		Value:     "1000000000000000000",
		Gas:       21000,
		GasPrice:  "20000000000",
		RawTxHex:  "0xf86c2a8504a817c80082520894123456789...",
	}
}

// validETHUnsignedEIP1559 returns a valid unsigned EIP-1559 ETH transaction file.
func validETHUnsignedEIP1559() *dtoeth.ETHTransactionFile {
	return &dtoeth.ETHTransactionFile{
		Version:              1,
		TxType:               "unsigned",
		EthTxType:            2,
		ChainID:              1,
		Nonce:                42,
		From:                 "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
		To:                   "0x1234567890abcdef1234567890abcdef12345678",
		Value:                "1000000000000000000",
		Gas:                  21000,
		MaxFeePerGas:         "30000000000",
		MaxPriorityFeePerGas: "2000000000",
		RawTxHex:             "0x02f86c01...",
	}
}

// validETHSignedEIP1559 returns a valid signed EIP-1559 ETH transaction file.
func validETHSignedEIP1559() *dtoeth.ETHTransactionFile {
	return &dtoeth.ETHTransactionFile{
		Version:              1,
		TxType:               "signed",
		EthTxType:            2,
		ChainID:              1,
		Nonce:                42,
		From:                 "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
		To:                   "0x1234567890abcdef1234567890abcdef12345678",
		Value:                "1000000000000000000",
		Gas:                  21000,
		MaxFeePerGas:         "30000000000",
		MaxPriorityFeePerGas: "2000000000",
		RawTxHex:             "0x02f86c01...",
		SignedTxHex:          "0x02f86c018504a817c800...",
	}
}

func TestWriteETHJSONFile(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		path        string
		data        *dtoeth.ETHTransactionFile
		wantErr     bool
		checkSuffix string
	}{
		{
			name:        "write valid unsigned legacy ETH transaction",
			path:        filepath.Join(tempDir, "deposit_8_unsigned_0_"),
			data:        validETHUnsignedLegacy(),
			wantErr:     false,
			checkSuffix: ".json",
		},
		{
			name:        "write valid unsigned EIP-1559 ETH transaction",
			path:        filepath.Join(tempDir, "deposit_9_unsigned_0_"),
			data:        validETHUnsignedEIP1559(),
			wantErr:     false,
			checkSuffix: ".json",
		},
		{
			name:        "write valid signed EIP-1559 ETH transaction",
			path:        filepath.Join(tempDir, "deposit_10_signed_1_"),
			data:        validETHSignedEIP1559(),
			wantErr:     false,
			checkSuffix: ".json",
		},
		{
			name: "reject invalid transaction (zero version)",
			path: filepath.Join(tempDir, "deposit_11_unsigned_0_"),
			data: &dtoeth.ETHTransactionFile{
				Version: 0,
			},
			wantErr: true,
		},
		{
			name: "reject invalid transaction (missing RawTxHex)",
			path: filepath.Join(tempDir, "deposit_12_unsigned_0_"),
			data: &dtoeth.ETHTransactionFile{
				Version:   1,
				TxType:    "unsigned",
				EthTxType: 0,
				ChainID:   1,
				From:      "0xSender",
				To:        "0xRecipient",
				Value:     "1000",
				Gas:       21000,
				GasPrice:  "20000000000",
				RawTxHex:  "", // missing
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewTransactionFileRepository("")
			fileName, err := repo.WriteETHJSONFile(tt.path, tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteETHJSONFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// File must exist
				if _, err := os.Stat(fileName); os.IsNotExist(err) {
					t.Errorf("WriteETHJSONFile() did not create file: %v", fileName)
				}

				// File must have .json extension
				if !strings.HasSuffix(fileName, tt.checkSuffix) {
					t.Errorf("WriteETHJSONFile() file name = %v, want suffix %v", fileName, tt.checkSuffix)
				}

				// File must contain valid JSON
				content, err := os.ReadFile(fileName)
				if err != nil {
					t.Errorf("Failed to read written file: %v", err)
				}
				if len(content) == 0 {
					t.Error("WriteETHJSONFile() wrote empty file")
				}

				var readBack dtoeth.ETHTransactionFile
				if err := json.Unmarshal(content, &readBack); err != nil {
					t.Errorf("WriteETHJSONFile() wrote invalid JSON: %v", err)
				}
				if err := readBack.Validate(); err != nil {
					t.Errorf("WriteETHJSONFile() wrote invalid transaction data: %v", err)
				}

				_ = os.Remove(fileName)
			}
		})
	}
}

func TestReadETHJSONFile(t *testing.T) {
	tempDir := t.TempDir()

	// Create valid legacy JSON file
	validLegacyPath := filepath.Join(tempDir, "test_legacy_valid.json")
	validLegacyContent := `{
		"version": 1,
		"tx_type": "unsigned",
		"eth_tx_type": 0,
		"chain_id": 1,
		"nonce": 42,
		"from": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
		"to": "0x1234567890abcdef1234567890abcdef12345678",
		"value": "1000000000000000000",
		"gas": 21000,
		"gas_price": "20000000000",
		"raw_tx_hex": "0xf86c2a..."
	}`
	if err := os.WriteFile(validLegacyPath, []byte(validLegacyContent), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create valid EIP-1559 JSON file
	validEIP1559Path := filepath.Join(tempDir, "test_eip1559_valid.json")
	validEIP1559Content := `{
		"version": 1,
		"tx_type": "unsigned",
		"eth_tx_type": 2,
		"chain_id": 31337,
		"nonce": 0,
		"from": "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		"to": "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
		"value": "500000000000000000",
		"gas": 21000,
		"max_fee_per_gas": "30000000000",
		"max_priority_fee_per_gas": "2000000000",
		"raw_tx_hex": "0x02f86c..."
	}`
	if err := os.WriteFile(validEIP1559Path, []byte(validEIP1559Content), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create invalid JSON file
	invalidJSONPath := filepath.Join(tempDir, "test_invalid.json")
	if err := os.WriteFile(invalidJSONPath, []byte(`{"invalid": json}`), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create file with invalid data (version 0)
	invalidDataPath := filepath.Join(tempDir, "test_invalid_data.json")
	invalidDataContent := `{"version": 0, "tx_type": "unsigned"}`
	if err := os.WriteFile(invalidDataPath, []byte(invalidDataContent), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create file with wrong extension
	invalidExtPath := filepath.Join(tempDir, "test_invalid.txt")
	if err := os.WriteFile(invalidExtPath, []byte(validLegacyContent), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name        string
		path        string
		wantErr     bool
		checkChain  uint64
		checkTxType string
	}{
		{
			name:        "read valid legacy ETH transaction file",
			path:        validLegacyPath,
			wantErr:     false,
			checkChain:  1,
			checkTxType: "unsigned",
		},
		{
			name:        "read valid EIP-1559 ETH transaction file",
			path:        validEIP1559Path,
			wantErr:     false,
			checkChain:  31337,
			checkTxType: "unsigned",
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
			name:    "read file with invalid data (version 0)",
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
			got, err := repo.ReadETHJSONFile(tt.path)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadETHJSONFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if err := got.Validate(); err != nil {
					t.Errorf("ReadETHJSONFile() returned invalid data: %v", err)
				}
				if got.ChainID != tt.checkChain {
					t.Errorf("ReadETHJSONFile() chain_id = %v, want %v", got.ChainID, tt.checkChain)
				}
				if got.TxType != tt.checkTxType {
					t.Errorf("ReadETHJSONFile() tx_type = %v, want %v", got.TxType, tt.checkTxType)
				}
			}
		})
	}
}

func TestETHJSONRoundTrip_Legacy(t *testing.T) {
	tempDir := t.TempDir()
	repo := NewTransactionFileRepository(tempDir + "/")

	txData := validETHUnsignedLegacy()
	path := repo.CreateFilePath(domainTx.ActionTypeDeposit, domainTx.TxTypeUnsigned, 8, 0)

	fileName, err := repo.WriteETHJSONFile(path, txData)
	if err != nil {
		t.Fatalf("WriteETHJSONFile() error = %v", err)
	}
	defer func() { _ = os.Remove(fileName) }()

	got, err := repo.ReadETHJSONFile(fileName)
	if err != nil {
		t.Fatalf("ReadETHJSONFile() error = %v", err)
	}

	if got.Version != txData.Version {
		t.Errorf("Round trip Version = %v, want %v", got.Version, txData.Version)
	}
	if got.EthTxType != txData.EthTxType {
		t.Errorf("Round trip EthTxType = %v, want %v", got.EthTxType, txData.EthTxType)
	}
	if got.ChainID != txData.ChainID {
		t.Errorf("Round trip ChainID = %v, want %v", got.ChainID, txData.ChainID)
	}
	if got.From != txData.From {
		t.Errorf("Round trip From = %v, want %v", got.From, txData.From)
	}
	if got.GasPrice != txData.GasPrice {
		t.Errorf("Round trip GasPrice = %v, want %v", got.GasPrice, txData.GasPrice)
	}
	if got.RawTxHex != txData.RawTxHex {
		t.Errorf("Round trip RawTxHex = %v, want %v", got.RawTxHex, txData.RawTxHex)
	}
}

func TestETHJSONRoundTrip_EIP1559(t *testing.T) {
	tempDir := t.TempDir()
	repo := NewTransactionFileRepository(tempDir + "/")

	txData := validETHUnsignedEIP1559()
	path := repo.CreateFilePath(domainTx.ActionTypeDeposit, domainTx.TxTypeUnsigned, 9, 0)

	fileName, err := repo.WriteETHJSONFile(path, txData)
	if err != nil {
		t.Fatalf("WriteETHJSONFile() error = %v", err)
	}
	defer func() { _ = os.Remove(fileName) }()

	got, err := repo.ReadETHJSONFile(fileName)
	if err != nil {
		t.Fatalf("ReadETHJSONFile() error = %v", err)
	}

	if got.EthTxType != 2 {
		t.Errorf("Round trip EthTxType = %v, want 2", got.EthTxType)
	}
	if got.MaxFeePerGas != txData.MaxFeePerGas {
		t.Errorf("Round trip MaxFeePerGas = %v, want %v", got.MaxFeePerGas, txData.MaxFeePerGas)
	}
	if got.MaxPriorityFeePerGas != txData.MaxPriorityFeePerGas {
		t.Errorf("Round trip MaxPriorityFeePerGas = %v, want %v", got.MaxPriorityFeePerGas, txData.MaxPriorityFeePerGas)
	}
	if got.GasPrice != "" {
		t.Errorf("Round trip GasPrice should be empty for EIP-1559, got %v", got.GasPrice)
	}
}

func TestETHJSONRoundTrip_Signed(t *testing.T) {
	tempDir := t.TempDir()
	repo := NewTransactionFileRepository(tempDir + "/")

	txData := validETHSignedEIP1559()
	path := repo.CreateFilePath(domainTx.ActionTypeDeposit, domainTx.TxTypeSigned, 10, 1)

	fileName, err := repo.WriteETHJSONFile(path, txData)
	if err != nil {
		t.Fatalf("WriteETHJSONFile() error = %v", err)
	}
	defer func() { _ = os.Remove(fileName) }()

	got, err := repo.ReadETHJSONFile(fileName)
	if err != nil {
		t.Fatalf("ReadETHJSONFile() error = %v", err)
	}

	if got.TxType != "signed" {
		t.Errorf("Round trip TxType = %v, want signed", got.TxType)
	}
	if got.SignedTxHex != txData.SignedTxHex {
		t.Errorf("Round trip SignedTxHex = %v, want %v", got.SignedTxHex, txData.SignedTxHex)
	}
}

func TestReadETHJSONFilePathTraversal(t *testing.T) {
	tempDir := t.TempDir()
	jsonDir := filepath.Join(tempDir, "json")
	if err := os.MkdirAll(jsonDir, 0o700); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create a valid ETH JSON file
	validPath := filepath.Join(jsonDir, "test_valid.json")
	validContent := `{
		"version": 1,
		"tx_type": "unsigned",
		"eth_tx_type": 2,
		"chain_id": 31337,
		"nonce": 0,
		"from": "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		"to": "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
		"value": "1000",
		"gas": 21000,
		"max_fee_per_gas": "30000000000",
		"max_priority_fee_per_gas": "2000000000",
		"raw_tx_hex": "0x02f86c..."
	}`
	if err := os.WriteFile(validPath, []byte(validContent), 0o644); err != nil {
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
			readPath: validPath,
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
			readPath: validPath,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewTransactionFileRepository(tt.basePath)
			got, err := repo.ReadETHJSONFile(tt.readPath)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadETHJSONFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errContain != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errContain) {
					t.Errorf("ReadETHJSONFile() error = %v, want error containing %q", err, tt.errContain)
				}
			}
			if !tt.wantErr && got == nil {
				t.Error("ReadETHJSONFile() returned nil data")
			}
		})
	}
}

func TestWriteETHJSONFilePathTraversal(t *testing.T) {
	tempDir := t.TempDir()
	jsonDir := filepath.Join(tempDir, "json")
	if err := os.MkdirAll(jsonDir, 0o700); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	validData := validETHUnsignedEIP1559()

	tests := []struct {
		name       string
		basePath   string
		writePath  string
		wantErr    bool
		errContain string
	}{
		{
			name:      "write file within allowed directory",
			basePath:  jsonDir + "/",
			writePath: filepath.Join(jsonDir, "test_"),
			wantErr:   false,
		},
		{
			name:       "reject path traversal with ../",
			basePath:   jsonDir + "/",
			writePath:  filepath.Join(jsonDir, "../../../etc/passwd_"),
			wantErr:    true,
			errContain: "path traversal attempt detected",
		},
		{
			name:       "reject absolute path outside base",
			basePath:   jsonDir + "/",
			writePath:  "/tmp/malicious_",
			wantErr:    true,
			errContain: "path traversal attempt detected",
		},
		{
			name:      "allow write when basePath is empty (test mode)",
			basePath:  "",
			writePath: filepath.Join(tempDir, "test_mode_"),
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewTransactionFileRepository(tt.basePath)
			fileName, err := repo.WriteETHJSONFile(tt.writePath, validData)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteETHJSONFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errContain != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errContain) {
					t.Errorf("WriteETHJSONFile() error = %v, want error containing %q", err, tt.errContain)
				}
			}
			if !tt.wantErr {
				if !strings.HasSuffix(fileName, ".json") {
					t.Errorf("WriteETHJSONFile() file name = %v, want suffix .json", fileName)
				}
				_ = os.Remove(fileName)
			}
		})
	}
}
