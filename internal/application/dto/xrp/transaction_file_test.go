package xrp

import (
	"testing"
)

func TestXRPTransactionFile_Validate(t *testing.T) {
	tests := []struct {
		name    string
		file    *XRPTransactionFile
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid transaction file",
			file: &XRPTransactionFile{
				Version:   "1.0.0",
				Chain:     "XRP",
				Network:   "mainnet",
				CreatedAt: "2024-02-14T00:00:00Z",
				Transactions: []XRPTransactionEntry{
					{
						UUID: "01234567-89ab-cdef-0123-456789abcdef",
						UnsignedData: TxInput{
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
			wantErr: false,
		},
		{
			name: "invalid version format",
			file: &XRPTransactionFile{
				Version:      "1.0",
				Chain:        "XRP",
				Network:      "mainnet",
				CreatedAt:    "2024-02-14T00:00:00Z",
				Transactions: []XRPTransactionEntry{},
			},
			wantErr: true,
			errMsg:  "version must follow semantic versioning pattern",
		},
		{
			name: "invalid chain",
			file: &XRPTransactionFile{
				Version:      "1.0.0",
				Chain:        "BTC",
				Network:      "mainnet",
				CreatedAt:    "2024-02-14T00:00:00Z",
				Transactions: []XRPTransactionEntry{},
			},
			wantErr: true,
			errMsg:  "chain must be XRP",
		},
		{
			name: "invalid network",
			file: &XRPTransactionFile{
				Version:      "1.0.0",
				Chain:        "XRP",
				Network:      "devnet",
				CreatedAt:    "2024-02-14T00:00:00Z",
				Transactions: []XRPTransactionEntry{},
			},
			wantErr: true,
			errMsg:  "network must be mainnet or testnet",
		},
		{
			name: "empty transactions",
			file: &XRPTransactionFile{
				Version:      "1.0.0",
				Chain:        "XRP",
				Network:      "mainnet",
				CreatedAt:    "2024-02-14T00:00:00Z",
				Transactions: []XRPTransactionEntry{},
			},
			wantErr: true,
			errMsg:  "transactions array cannot be empty",
		},
		{
			name: "signature count exceeds required",
			file: &XRPTransactionFile{
				Version:   "1.0.0",
				Chain:     "XRP",
				Network:   "mainnet",
				CreatedAt: "2024-02-14T00:00:00Z",
				Transactions: []XRPTransactionEntry{
					{
						UUID:               "01234567-89ab-cdef-0123-456789abcdef",
						SenderAccount:      "rN7n7otQDd6FczFgLdlqtyMVrn3HMfgnj",
						SenderAccountType:  "client",
						SignatureCount:     3,
						RequiredSignatures: 2,
						IsComplete:         true,
					},
				},
			},
			wantErr: true,
			errMsg:  "signatureCount cannot exceed requiredSignatures",
		},
		{
			name: "negative signature count",
			file: &XRPTransactionFile{
				Version:   "1.0.0",
				Chain:     "XRP",
				Network:   "mainnet",
				CreatedAt: "2024-02-14T00:00:00Z",
				Transactions: []XRPTransactionEntry{
					{
						UUID:               "01234567-89ab-cdef-0123-456789abcdef",
						SenderAccount:      "rN7n7otQDd6FczFgLdlqtyMVrn3HMfgnj",
						SenderAccountType:  "client",
						SignatureCount:     -1,
						RequiredSignatures: 1,
						IsComplete:         false,
					},
				},
			},
			wantErr: true,
			errMsg:  "signatureCount must be >= 0",
		},
		{
			name: "signedBlob should be null when signatureCount is 0",
			file: &XRPTransactionFile{
				Version:   "1.0.0",
				Chain:     "XRP",
				Network:   "mainnet",
				CreatedAt: "2024-02-14T00:00:00Z",
				Transactions: []XRPTransactionEntry{
					{
						UUID:               "01234567-89ab-cdef-0123-456789abcdef",
						SenderAccount:      "rN7n7otQDd6FczFgLdlqtyMVrn3HMfgnj",
						SenderAccountType:  "client",
						SignatureCount:     0,
						RequiredSignatures: 1,
						SignedBlob:         stringPtr("some-blob"),
						IsComplete:         false,
					},
				},
			},
			wantErr: true,
			errMsg:  "signedBlob must be null when signatureCount is 0",
		},
		{
			name: "valid multi-signature transaction",
			file: &XRPTransactionFile{
				Version:   "1.0.0",
				Chain:     "XRP",
				Network:   "testnet",
				CreatedAt: "2024-02-14T00:00:00Z",
				Transactions: []XRPTransactionEntry{
					{
						UUID: "01234567-89ab-cdef-0123-456789abcdef",
						UnsignedData: TxInput{
							Account:     "rN7n7otQDd6FczFgLdlqtyMVrn3HMfgnj",
							Destination: "rLHzPsX6oXkzU9jNbPGJcEKvxbqBKzhqgN",
							Amount:      "1000000",
							Fee:         "12",
						},
						SenderAccount:      "rN7n7otQDd6FczFgLdlqtyMVrn3HMfgnj",
						SenderAccountType:  "client",
						SignatureCount:     2,
						RequiredSignatures: 3,
						SignedBlob: stringPtr(
							"1200002280000000240000000161400000000000000168400000000000000C732103AB40A0490F9B7ED8DF29" +
								"D246BF2D6269820A0EE7742ACDD457BEA7C7D0931EDB74473045022100D184EB4AE5956FF600E7536E" +
								"E459345C7BBCF097A84CC61A93B9AF7197EDB98702201E5C9F83C9B1F4A6E3F1E9F1E9F1E9F1E9F1E9" +
								"F1E9F1E9F1E9F1E9F1E9F181145E7B112523F68D2F5E879DB4EAC51C6698A69304",
						),
						IsComplete: false,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.file.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("XRPTransactionFile.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("XRPTransactionFile.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestXRPTransactionEntry_Validate(t *testing.T) {
	tests := []struct {
		name    string
		entry   *XRPTransactionEntry
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid unsigned transaction entry",
			entry: &XRPTransactionEntry{
				UUID:               "01234567-89ab-cdef-0123-456789abcdef",
				SenderAccount:      "rN7n7otQDd6FczFgLdlqtyMVrn3HMfgnj",
				SenderAccountType:  "client",
				SignatureCount:     0,
				RequiredSignatures: 1,
				SignedBlob:         nil,
				IsComplete:         false,
			},
			wantErr: false,
		},
		{
			name: "valid signed transaction entry",
			entry: &XRPTransactionEntry{
				UUID:               "01234567-89ab-cdef-0123-456789abcdef",
				SenderAccount:      "rN7n7otQDd6FczFgLdlqtyMVrn3HMfgnj",
				SenderAccountType:  "client",
				SignatureCount:     1,
				RequiredSignatures: 1,
				SignedBlob:         stringPtr("signed-blob-hex"),
				IsComplete:         true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.entry.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("XRPTransactionEntry.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("XRPTransactionEntry.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
