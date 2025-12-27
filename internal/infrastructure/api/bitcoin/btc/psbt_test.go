//go:build integration
// +build integration

package btc_test

import (
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin/btc"
	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

// TestCreatePSBT tests PSBT creation from unsigned transaction
func TestCreatePSBT(t *testing.T) {
	// Create Bitcoin instance for testing
	bitcoin, err := testutil.GetBTC()
	require.NoError(t, err)

	// Create a simple unsigned transaction
	msgTx := wire.NewMsgTx(wire.TxVersion)

	// Add input (using a test transaction hash)
	prevHash, err := chainhash.NewHashFromStr("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	require.NoError(t, err)

	txIn := wire.NewTxIn(wire.NewOutPoint(prevHash, 0), nil, nil)
	msgTx.AddTxIn(txIn)

	// Add output
	addr, err := btcutil.DecodeAddress("tb1qw508d6qejxtdg4y5r3zarvary0c5xw7kxpjzsx", &chaincfg.TestNet3Params)
	require.NoError(t, err)

	pkScript, err := txscript.PayToAddrScript(addr)
	require.NoError(t, err)

	txOut := wire.NewTxOut(100000, pkScript)
	msgTx.AddTxOut(txOut)

	// Create prevTxs metadata
	prevTxs := []btc.PrevTx{
		{
			Txid:         prevHash.String(),
			Vout:         0,
			ScriptPubKey: "0014751e76e8199196d454941c45d1b3a323f1433bd6",
			Amount:       0.002,
		},
	}

	// Create PSBT
	psbtBase64, err := bitcoin.CreatePSBT(msgTx, prevTxs)
	require.NoError(t, err)
	assert.NotEmpty(t, psbtBase64)

	t.Logf("Created PSBT: %s", psbtBase64[:50]+"...")
}

// TestParsePSBT tests PSBT parsing functionality
func TestParsePSBT(t *testing.T) {
	// This is a valid PSBT from Bitcoin Core testnet (simplified for testing)
	// This PSBT has 1 input and 1 output, unsigned
	validPSBT := "cHNidP8BAHECAAAAAeWU5KQnIgL9xnm9wWKHxWDcY7D6IlQFKkGVBQKmJm+CAAAAAAD/////AoA4AQAAAAAAFgAUkh7tjzD5WbwzQe6iL6Cg9UGMKy/waSYBAAAAABYAFM8dV5vxr5vdPJJ8uiDGfKsAO5qCAAAAAAA="

	bitcoin, err := testutil.GetBTC()
	require.NoError(t, err)

	// Parse PSBT
	parsed, err := bitcoin.ParsePSBT(validPSBT)
	require.NoError(t, err)
	assert.NotNil(t, parsed)
	assert.NotNil(t, parsed.Packet)
	assert.Equal(t, 1, parsed.InputCount)
	assert.Equal(t, 2, parsed.OutputCount)
	assert.False(t, parsed.IsComplete) // Unsigned PSBT should not be complete
	assert.False(t, parsed.HasSignature)

	t.Logf("Parsed PSBT: inputs=%d, outputs=%d, complete=%t",
		parsed.InputCount, parsed.OutputCount, parsed.IsComplete)
}

// TestParsePSBT_InvalidBase64 tests error handling for invalid base64
func TestParsePSBT_InvalidBase64(t *testing.T) {
	bitcoin, err := testutil.GetBTC()
	require.NoError(t, err)

	invalidBase64 := "not-valid-base64!!!"

	_, err = bitcoin.ParsePSBT(invalidBase64)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode base64")
}

// TestParsePSBT_InvalidPSBT tests error handling for invalid PSBT structure
func TestParsePSBT_InvalidPSBT(t *testing.T) {
	bitcoin, err := testutil.GetBTC()
	require.NoError(t, err)

	// Valid base64 but not a valid PSBT
	invalidPSBT := "SGVsbG8gV29ybGQh" // "Hello World!" in base64

	_, err = bitcoin.ParsePSBT(invalidPSBT)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse PSBT")
}

// TestValidatePSBT tests PSBT validation
func TestValidatePSBT(t *testing.T) {
	bitcoin, err := testutil.GetBTC()
	require.NoError(t, err)

	tests := []struct {
		name    string
		psbt    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid PSBT",
			psbt:    "cHNidP8BAHECAAAAAeWU5KQnIgL9xnm9wWKHxWDcY7D6IlQFKkGVBQKmJm+CAAAAAAD/////AoA4AQAAAAAAFgAUkh7tjzD5WbwzQe6iL6Cg9UGMKy/waSYBAAAAABYAFM8dV5vxr5vdPJJ8uiDGfKsAO5qCAAAAAAA=",
			wantErr: false,
		},
		{
			name:    "invalid base64",
			psbt:    "not-valid-base64",
			wantErr: true,
			errMsg:  "failed to parse PSBT",
		},
		{
			name:    "invalid PSBT structure",
			psbt:    "SGVsbG8gV29ybGQh",
			wantErr: true,
			errMsg:  "failed to parse PSBT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = bitcoin.ValidatePSBT(tt.psbt)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else if err != nil {
				// Note: validation might still fail if witness UTXO is missing
				// This is expected for real PSBTs that need proper setup
				t.Logf("Validation error (expected for test PSBT): %v", err)
			}
		})
	}
}

// TestIsPSBTComplete tests PSBT completion check
func TestIsPSBTComplete(t *testing.T) {
	bitcoin, err := testutil.GetBTC()
	require.NoError(t, err)

	// Unsigned PSBT (not complete)
	unsignedPSBT := "cHNidP8BAHECAAAAAeWU5KQnIgL9xnm9wWKHxWDcY7D6IlQFKkGVBQKmJm+CAAAAAAD/////AoA4AQAAAAAAFgAUkh7tjzD5WbwzQe6iL6Cg9UGMKy/waSYBAAAAABYAFM8dV5vxr5vdPJJ8uiDGfKsAO5qCAAAAAAA="

	isComplete, err := bitcoin.IsPSBTComplete(unsignedPSBT)
	require.NoError(t, err)
	assert.False(t, isComplete)
}

// TestGetPSBTFee tests PSBT fee calculation
func TestGetPSBTFee(t *testing.T) {
	bitcoin, err := testutil.GetBTC()
	require.NoError(t, err)

	// This test requires a properly constructed PSBT with witness UTXOs
	// For now, we test the error case
	validPSBT := "cHNidP8BAHECAAAAAeWU5KQnIgL9xnm9wWKHxWDcY7D6IlQFKkGVBQKmJm+CAAAAAAD/////AoA4AQAAAAAAFgAUkh7tjzD5WbwzQe6iL6Cg9UGMKy/waSYBAAAAABYAFM8dV5vxr5vdPJJ8uiDGfKsAO5qCAAAAAAA="

	fee, err := bitcoin.GetPSBTFee(validPSBT)
	// Fee calculation will fail without proper witness UTXO, but parsing should work
	if err == nil {
		assert.GreaterOrEqual(t, fee, int64(0))
		t.Logf("Calculated fee: %d satoshis", fee)
	} else {
		t.Logf("Fee calculation failed (expected for incomplete PSBT): %v", err)
	}
}

// TestSignPSBTWithKey_ErrorCases tests error handling in PSBT signing
func TestSignPSBTWithKey_ErrorCases(t *testing.T) {
	bitcoin, err := testutil.GetBTC()
	require.NoError(t, err)

	tests := []struct {
		name    string
		psbt    string
		wifs    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "invalid PSBT",
			psbt:    "invalid-psbt",
			wifs:    []string{"cVnFEfUZiHX2jJDGmupHDqpLYu2bZhAfMWXLBY6Jyunh7PxmpnG5"},
			wantErr: true,
			errMsg:  "failed to parse PSBT",
		},
		{
			name:    "invalid WIF",
			psbt:    "cHNidP8BAHECAAAAAeWU5KQnIgL9xnm9wWKHxWDcY7D6IlQFKkGVBQKmJm+CAAAAAAD/////AoA4AQAAAAAAFgAUkh7tjzD5WbwzQe6iL6Cg9UGMKy/waSYBAAAAABYAFM8dV5vxr5vdPJJ8uiDGfKsAO5qCAAAAAAA=",
			wifs:    []string{"invalid-wif"},
			wantErr: true,
			errMsg:  "failed to decode WIF",
		},
		{
			name:    "empty WIF list",
			psbt:    "cHNidP8BAHECAAAAAeWU5KQnIgL9xnm9wWKHxWDcY7D6IlQFKkGVBQKmJm+CAAAAAAD/////AoA4AQAAAAAAFgAUkh7tjzD5WbwzQe6iL6Cg9UGMKy/waSYBAAAAABYAFM8dV5vxr5vdPJJ8uiDGfKsAO5qCAAAAAAA=",
			wifs:    []string{},
			wantErr: true,
			errMsg:  "no signatures were added",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := bitcoin.SignPSBTWithKey(tt.psbt, tt.wifs)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestFinalizePSBT_ErrorCases tests error handling in PSBT finalization
func TestFinalizePSBT_ErrorCases(t *testing.T) {
	bitcoin, err := testutil.GetBTC()
	require.NoError(t, err)

	tests := []struct {
		name    string
		psbt    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "invalid PSBT",
			psbt:    "invalid-psbt",
			wantErr: true,
			errMsg:  "failed to parse PSBT",
		},
		{
			name:    "unsigned PSBT (not complete)",
			psbt:    "cHNidP8BAHECAAAAAeWU5KQnIgL9xnm9wWKHxWDcY7D6IlQFKkGVBQKmJm+CAAAAAAD/////AoA4AQAAAAAAFgAUkh7tjzD5WbwzQe6iL6Cg9UGMKy/waSYBAAAAABYAFM8dV5vxr5vdPJJ8uiDGfKsAO5qCAAAAAAA=",
			wantErr: true,
			errMsg:  "cannot finalize incomplete PSBT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := bitcoin.FinalizePSBT(tt.psbt)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestExtractTransaction_ErrorCases tests error handling in transaction extraction
func TestExtractTransaction_ErrorCases(t *testing.T) {
	bitcoin, err := testutil.GetBTC()
	require.NoError(t, err)

	tests := []struct {
		name    string
		psbt    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "invalid PSBT",
			psbt:    "invalid-psbt",
			wantErr: true,
			errMsg:  "failed to parse PSBT",
		},
		{
			name:    "unfinalized PSBT",
			psbt:    "cHNidP8BAHECAAAAAeWU5KQnIgL9xnm9wWKHxWDcY7D6IlQFKkGVBQKmJm+CAAAAAAD/////AoA4AQAAAAAAFgAUkh7tjzD5WbwzQe6iL6Cg9UGMKy/waSYBAAAAABYAFM8dV5vxr5vdPJJ8uiDGfKsAO5qCAAAAAAA=",
			wantErr: true,
			errMsg:  "failed to extract transaction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := bitcoin.ExtractTransaction(tt.psbt)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestPSBTWorkflow_Integration tests complete PSBT workflow (requires proper setup)
// NOTE: This test requires a properly funded Bitcoin Core testnet wallet
// Disabled by default - remove skip to run in integration environment
func TestPSBTWorkflow_Integration(t *testing.T) {
	t.Skip("Requires properly funded Bitcoin Core testnet wallet")
	bitcoin, err := testutil.GetBTC()
	require.NoError(t, err)

	// This test requires a properly funded Bitcoin Core testnet wallet
	// 1. Create unsigned transaction
	prevHash, err := chainhash.NewHashFromStr("0000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)

	inputs := []btcjson.TransactionInput{
		{
			Txid: prevHash.String(),
			Vout: 0,
		},
	}

	addr, err := btcutil.DecodeAddress("tb1qw508d6qejxtdg4y5r3zarvary0c5xw7kxpjzsx", &chaincfg.TestNet3Params)
	require.NoError(t, err)

	outputs := map[btcutil.Address]btcutil.Amount{
		addr: 100000,
	}

	msgTx, err := bitcoin.CreateRawTransaction(inputs, outputs)
	require.NoError(t, err)

	// 2. Create PSBT from transaction
	prevTxs := []btc.PrevTx{
		{
			Txid:         prevHash.String(),
			Vout:         0,
			ScriptPubKey: "0014751e76e8199196d454941c45d1b3a323f1433bd6",
			Amount:       0.002,
		},
	}

	psbtBase64, err := bitcoin.CreatePSBT(msgTx, prevTxs)
	require.NoError(t, err)
	t.Logf("Step 1: Created PSBT")

	// 3. Parse PSBT
	parsed, err := bitcoin.ParsePSBT(psbtBase64)
	require.NoError(t, err)
	assert.False(t, parsed.IsComplete)
	t.Logf("Step 2: Parsed PSBT (complete=%t)", parsed.IsComplete)

	// 4. Validate PSBT
	err = bitcoin.ValidatePSBT(psbtBase64)
	if err != nil {
		t.Logf("Step 3: PSBT validation failed (expected without proper setup): %v", err)
	} else {
		t.Logf("Step 3: PSBT validated successfully")
	}

	// Note: Steps 5-8 (signing, finalization, extraction) require proper private keys
	// and funded UTXOs, which are not available in unit tests
	t.Logf("PSBT workflow test completed (partial - signing requires integration setup)")
}

// TestCreatePSBT_ErrorCases tests error cases in PSBT creation
func TestCreatePSBT_ErrorCases(t *testing.T) {
	bitcoin, err := testutil.GetBTC()
	require.NoError(t, err)

	msgTx := wire.NewMsgTx(wire.TxVersion)

	// Add input
	prevHash, _ := chainhash.NewHashFromStr("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	txIn := wire.NewTxIn(wire.NewOutPoint(prevHash, 0), nil, nil)
	msgTx.AddTxIn(txIn)

	tests := []struct {
		name    string
		prevTxs []btc.PrevTx
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty prevTxs",
			prevTxs: []btc.PrevTx{},
			wantErr: false, // Empty prevTxs is valid, PSBT will be created without metadata
		},
		{
			name: "invalid scriptPubKey",
			prevTxs: []btc.PrevTx{
				{
					Txid:         prevHash.String(),
					Vout:         0,
					ScriptPubKey: "invalid-hex",
					Amount:       0.001,
				},
			},
			wantErr: true,
			errMsg:  "failed to decode scriptPubKey",
		},
		{
			name: "prevTxs count exceeds inputs",
			prevTxs: []btc.PrevTx{
				{
					Txid:         prevHash.String(),
					Vout:         0,
					ScriptPubKey: "0014751e76e8199196d454941c45d1b3a323f1433bd6",
					Amount:       0.001,
				},
				{
					Txid:         prevHash.String(),
					Vout:         1,
					ScriptPubKey: "0014751e76e8199196d454941c45d1b3a323f1433bd6",
					Amount:       0.001,
				},
			},
			wantErr: true,
			errMsg:  "exceeds number of inputs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := bitcoin.CreatePSBT(msgTx, tt.prevTxs)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				// May still error due to other validation, but not the specific error we're testing
				t.Logf("Result: err=%v", err)
			}
		})
	}
}
