package eth_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
)

// validMPCFile returns a minimal valid unsigned ETHMPCTransactionFile for use as a test base.
// Addresses are EIP-55 checksummed; all required fields are populated.
func validMPCFile() dtoeth.ETHMPCTransactionFile {
	return dtoeth.ETHMPCTransactionFile{
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

func TestETHMPCTransactionFile_Validate_ValidUnsigned(t *testing.T) {
	f := validMPCFile()
	require.NoError(t, f.Validate())
}

func TestETHMPCTransactionFile_Validate_ValidSigned(t *testing.T) {
	f := validMPCFile()
	f.TxType = "signed"
	f.SignedTxHex = "0x02f86c8204d28459682f008459682f14825208"
	require.NoError(t, f.Validate())
}

func TestETHMPCTransactionFile_Validate_HeaderErrors(t *testing.T) {
	tests := []struct {
		name    string
		modify  func(*dtoeth.ETHMPCTransactionFile)
		wantErr error
	}{
		{
			name:    "version zero",
			modify:  func(f *dtoeth.ETHMPCTransactionFile) { f.Version = 0 },
			wantErr: dtoeth.ErrInvalidMPCVersion,
		},
		{
			name:    "invalid tx_type",
			modify:  func(f *dtoeth.ETHMPCTransactionFile) { f.TxType = "pending" },
			wantErr: dtoeth.ErrInvalidMPCTxType,
		},
		{
			name:    "zero chain_id",
			modify:  func(f *dtoeth.ETHMPCTransactionFile) { f.ChainID = 0 },
			wantErr: dtoeth.ErrInvalidMPCChainID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := validMPCFile()
			tt.modify(&f)
			err := f.Validate()
			require.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestETHMPCTransactionFile_Validate_RequiredFields(t *testing.T) {
	tests := []struct {
		name    string
		modify  func(*dtoeth.ETHMPCTransactionFile)
		wantErr error
	}{
		{
			name:    "empty UUID",
			modify:  func(f *dtoeth.ETHMPCTransactionFile) { f.UUID = "" },
			wantErr: dtoeth.ErrEmptyMPCUUID,
		},
		{
			name:    "empty From",
			modify:  func(f *dtoeth.ETHMPCTransactionFile) { f.From = "" },
			wantErr: dtoeth.ErrEmptyMPCFrom,
		},
		{
			name:    "invalid From (non-checksum)",
			modify:  func(f *dtoeth.ETHMPCTransactionFile) { f.From = "0xabcd" },
			wantErr: dtoeth.ErrInvalidMPCAddress,
		},
		{
			name:    "empty To",
			modify:  func(f *dtoeth.ETHMPCTransactionFile) { f.To = "" },
			wantErr: dtoeth.ErrEmptyMPCTo,
		},
		{
			name:    "invalid To (non-checksum)",
			modify:  func(f *dtoeth.ETHMPCTransactionFile) { f.To = "0xabcd" },
			wantErr: dtoeth.ErrInvalidMPCAddress,
		},
		{
			name:    "empty Value",
			modify:  func(f *dtoeth.ETHMPCTransactionFile) { f.Value = "" },
			wantErr: dtoeth.ErrEmptyMPCValue,
		},
		{
			name:    "zero GasLimit",
			modify:  func(f *dtoeth.ETHMPCTransactionFile) { f.GasLimit = 0 },
			wantErr: dtoeth.ErrZeroMPCGasLimit,
		},
		{
			name:    "empty TxHash",
			modify:  func(f *dtoeth.ETHMPCTransactionFile) { f.TxHash = "" },
			wantErr: dtoeth.ErrEmptyMPCTxHash,
		},
		{
			name:    "empty RawTxHex",
			modify:  func(f *dtoeth.ETHMPCTransactionFile) { f.RawTxHex = "" },
			wantErr: dtoeth.ErrEmptyMPCRawTxHex,
		},
		{
			name:    "threshold zero",
			modify:  func(f *dtoeth.ETHMPCTransactionFile) { f.Threshold = 0 },
			wantErr: dtoeth.ErrInvalidMPCThreshold,
		},
		{
			name:    "PartyIDs fewer than Threshold",
			modify:  func(f *dtoeth.ETHMPCTransactionFile) { f.PartyIDs = []string{"node-1"} },
			wantErr: dtoeth.ErrMPCPartyIDsBelowThreshold,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := validMPCFile()
			tt.modify(&f)
			err := f.Validate()
			require.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestETHMPCTransactionFile_Validate_SignedConsistency(t *testing.T) {
	t.Run("signed type requires SignedTxHex", func(t *testing.T) {
		f := validMPCFile()
		f.TxType = "signed"
		// SignedTxHex intentionally left empty
		err := f.Validate()
		require.Error(t, err)
		assert.ErrorIs(t, err, dtoeth.ErrMissingMPCSignedTxHex)
	})
}
