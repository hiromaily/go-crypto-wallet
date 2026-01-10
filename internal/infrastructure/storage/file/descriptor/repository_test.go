//nolint:revive // Test file contains long test data
package descriptor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
)

func getTestDescriptors() []*domainWallet.Descriptor {
	return []*domainWallet.Descriptor{
		{
			Type:   domainWallet.DescriptorTypeSHWPKH,
			Script: "sh(wpkh([a1b2c3d4/49'/0'/0']xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL/0/*))",
			Keys: []domainWallet.DescriptorKey{
				{
					Fingerprint:    "a1b2c3d4",
					DerivationPath: "/49'/0'/0'",
					ExtendedPubKey: "xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL",
				},
			},
			Checksum: "abcd1234",
		},
		{
			Type:   domainWallet.DescriptorTypeWPKH,
			Script: "wpkh([a1b2c3d4/84'/0'/0']xpub6CUGRUonZSQ4TWtTMmzXdrXDtypWKiKrhko4egpiMZbpiaQL2jkwSB1icqYh2cfDfVxdx4df189oLKnC5fSwqPfgyP3hooxujYzAu3fDVmz/0/*)",
			Keys: []domainWallet.DescriptorKey{
				{
					Fingerprint:    "a1b2c3d4",
					DerivationPath: "/84'/0'/0'",
					ExtendedPubKey: "xpub6CUGRUonZSQ4TWtTMmzXdrXDtypWKiKrhko4egpiMZbpiaQL2jkwSB1icqYh2cfDfVxdx4df189oLKnC5fSwqPfgyP3hooxujYzAu3fDVmz",
				},
			},
			Checksum: "efgh5678",
		},
	}
}

func TestDescriptorFileRepository_WriteDescriptorsPlainText(t *testing.T) {
	// Create temporary directory for test files
	tmpDir := t.TempDir()
	repo := NewDescriptorFileRepository(tmpDir)

	descriptors := getTestDescriptors()
	filePath := filepath.Join(tmpDir, "test_descriptors.txt")

	tests := []struct {
		name         string
		withChecksum bool
		wantErr      bool
	}{
		{
			name:         "Write descriptors with checksums",
			withChecksum: true,
			wantErr:      false,
		},
		{
			name:         "Write descriptors without checksums",
			withChecksum: false,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.WriteDescriptorsPlainText(filePath, descriptors, tt.withChecksum)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteDescriptorsPlainText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify file was created
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Error("WriteDescriptorsPlainText() file was not created")
				}

				// Read back and verify content
				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Fatalf("Failed to read written file: %v", err)
				}

				lines := strings.Split(string(content), "\n")
				// Verify we have at least 2 lines (2 descriptors + trailing newline)
				if len(lines) < 2 {
					t.Errorf("WriteDescriptorsPlainText() lines count = %d, want >= 2", len(lines))
				}

				// Verify checksum presence
				if tt.withChecksum {
					if !strings.Contains(string(content), "#") {
						t.Error("WriteDescriptorsPlainText() checksum missing when requested")
					}
				}

				t.Logf("Written content:\n%s", content)

				// Clean up for next test
				_ = os.Remove(filePath)
			}
		})
	}
}

func TestDescriptorFileRepository_WriteDescriptorsBitcoinCoreJSON(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewDescriptorFileRepository(tmpDir)

	descriptors := getTestDescriptors()
	filePath := filepath.Join(tmpDir, "test_descriptors.json")

	err := repo.WriteDescriptorsBitcoinCoreJSON(
		filePath,
		descriptors,
		"now",
		"test_account",
		true,
		true,
	)
	if err != nil {
		t.Fatalf("WriteDescriptorsBitcoinCoreJSON() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("WriteDescriptorsBitcoinCoreJSON() file was not created")
	}

	// Read back and verify content is valid JSON
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read written file: %v", err)
	}

	// Verify JSON structure
	if !strings.HasPrefix(string(content), "[") {
		t.Error("WriteDescriptorsBitcoinCoreJSON() output is not a JSON array")
	}

	// Verify required fields are present
	requiredFields := []string{"desc", "timestamp", "active", "watchonly", "label"}
	for _, field := range requiredFields {
		if !strings.Contains(string(content), field) {
			t.Errorf("WriteDescriptorsBitcoinCoreJSON() missing required field: %s", field)
		}
	}

	t.Logf("Written JSON:\n%s", content)
}

func TestDescriptorFileRepository_ReadDescriptorsPlainText(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewDescriptorFileRepository(tmpDir)

	// Create test file
	descriptors := getTestDescriptors()
	filePath := filepath.Join(tmpDir, "test_read.txt")
	err := repo.WriteDescriptorsPlainText(filePath, descriptors, true)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	// Test reading
	read, err := repo.ReadDescriptorsPlainText(filePath)
	if err != nil {
		t.Fatalf("ReadDescriptorsPlainText() error = %v", err)
	}

	// Verify count
	if len(read) != len(descriptors) {
		t.Errorf("ReadDescriptorsPlainText() count = %d, want %d", len(read), len(descriptors))
	}

	// Verify checksums are preserved
	for _, desc := range read {
		if !strings.Contains(desc, "#") {
			t.Error("ReadDescriptorsPlainText() checksum not preserved")
		}
	}

	t.Logf("Read descriptors: %v", read)
}

func TestDescriptorFileRepository_ReadDescriptorsBitcoinCoreJSON(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewDescriptorFileRepository(tmpDir)

	// Create test file
	descriptors := getTestDescriptors()
	filePath := filepath.Join(tmpDir, "test_read.json")
	err := repo.WriteDescriptorsBitcoinCoreJSON(
		filePath,
		descriptors,
		"now",
		"test_account",
		true,
		true,
	)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	// Test reading
	read, err := repo.ReadDescriptorsBitcoinCoreJSON(filePath)
	if err != nil {
		t.Fatalf("ReadDescriptorsBitcoinCoreJSON() error = %v", err)
	}

	// Verify count
	if len(read) != len(descriptors) {
		t.Errorf("ReadDescriptorsBitcoinCoreJSON() count = %d, want %d", len(read), len(descriptors))
	}

	// Verify structure
	for i, item := range read {
		if item.Descriptor == "" {
			t.Errorf("ReadDescriptorsBitcoinCoreJSON() item %d missing descriptor", i)
		}
		if item.Timestamp == nil {
			t.Errorf("ReadDescriptorsBitcoinCoreJSON() item %d missing timestamp", i)
		}
		if !item.Active {
			t.Errorf("ReadDescriptorsBitcoinCoreJSON() item %d active = false, want true", i)
		}
		if !item.Watchonly {
			t.Errorf("ReadDescriptorsBitcoinCoreJSON() item %d watchonly = false, want true", i)
		}
	}

	t.Logf("Read descriptors: %+v", read)
}

func TestDescriptorFileRepository_AutoDetectAndRead(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewDescriptorFileRepository(tmpDir)

	descriptors := getTestDescriptors()

	tests := []struct {
		name    string
		format  string // "txt" or "json"
		wantErr bool
	}{
		{
			name:    "Auto-detect plain text format",
			format:  "txt",
			wantErr: false,
		},
		{
			name:    "Auto-detect JSON format",
			format:  "json",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filePath string
			var err error

			// Create file in specified format
			if tt.format == "txt" {
				filePath = filepath.Join(tmpDir, "test_auto_"+tt.format+".txt")
				err = repo.WriteDescriptorsPlainText(filePath, descriptors, true)
			} else {
				filePath = filepath.Join(tmpDir, "test_auto_"+tt.format+".json")
				err = repo.WriteDescriptorsBitcoinCoreJSON(
					filePath, descriptors, "now", "test", true, true,
				)
			}
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			// Test auto-detection
			read, err := repo.AutoDetectAndRead(filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("AutoDetectAndRead() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(read) != len(descriptors) {
					t.Errorf("AutoDetectAndRead() count = %d, want %d", len(read), len(descriptors))
				}
				t.Logf("Auto-detected and read %d descriptors from %s format", len(read), tt.format)
			}
		})
	}
}

func TestDescriptorFileRepository_CreateFilePath(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewDescriptorFileRepository(tmpDir)

	tests := []struct {
		name        string
		accountName string
		extension   string
		wantContain string
	}{
		{
			name:        "Create plain text file path",
			accountName: "deposit",
			extension:   ".txt",
			wantContain: "deposit_descriptors.txt",
		},
		{
			name:        "Create JSON file path",
			accountName: "payment",
			extension:   ".json",
			wantContain: "payment_descriptors.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := repo.CreateFilePath(tt.accountName, tt.extension)
			if !strings.Contains(got, tt.wantContain) {
				t.Errorf("CreateFilePath() = %v, want to contain %v", got, tt.wantContain)
			}
			t.Logf("Created file path: %s", got)
		})
	}
}

func TestDescriptorFileRepository_RoundTrip(t *testing.T) {
	// Integration test: write and read back descriptors in both formats
	tmpDir := t.TempDir()
	repo := NewDescriptorFileRepository(tmpDir)
	descriptors := getTestDescriptors()

	tests := []struct {
		name   string
		format string
	}{
		{
			name:   "Round trip plain text",
			format: "txt",
		},
		{
			name:   "Round trip JSON",
			format: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := filepath.Join(tmpDir, "roundtrip_"+tt.format+"."+tt.format)

			// Write
			var writeErr error
			if tt.format == "txt" {
				writeErr = repo.WriteDescriptorsPlainText(filePath, descriptors, true)
			} else {
				writeErr = repo.WriteDescriptorsBitcoinCoreJSON(
					filePath, descriptors, "now", "test", true, true,
				)
			}
			if writeErr != nil {
				t.Fatalf("Write failed: %v", writeErr)
			}

			// Read back using auto-detect
			read, err := repo.AutoDetectAndRead(filePath)
			if err != nil {
				t.Fatalf("Read failed: %v", err)
			}

			// Verify count matches
			if len(read) != len(descriptors) {
				t.Errorf("Round trip count mismatch: got %d, want %d", len(read), len(descriptors))
			}

			// Verify each descriptor
			for i, readDesc := range read {
				// Extract descriptor without checksum for comparison
				originalDesc := descriptors[i].Script
				readDescWithoutChecksum := strings.Split(readDesc, "#")[0]

				if !strings.Contains(readDescWithoutChecksum, originalDesc) {
					t.Errorf("Round trip descriptor %d mismatch", i)
				}
			}

			t.Logf("Round trip successful for %s format: %d descriptors", tt.format, len(read))
		})
	}
}
