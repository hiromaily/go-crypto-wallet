package btc

import (
	"testing"

	"github.com/btcsuite/btcd/txscript"
	"github.com/stretchr/testify/assert"
)

// TestP2WPKHRedeemScriptValidation tests the validation logic for P2WPKH redeemScripts
// in the signSegWitInput function. This ensures that malformed redeemScripts are
// rejected before being used for signature hash calculation.
func TestP2WPKHRedeemScriptValidation(t *testing.T) {
	tests := []struct {
		name         string
		redeemScript []byte
		shouldPass   bool
		description  string
	}{
		{
			name: "valid P2WPKH redeemScript",
			redeemScript: func() []byte {
				// Valid P2WPKH: OP_0 <20-byte-hash>
				hash := make([]byte, 20)
				for i := range hash {
					hash[i] = byte(i) // Fill with test data
				}
				builder := txscript.NewScriptBuilder()
				builder.AddOp(txscript.OP_0)
				builder.AddData(hash)
				script, _ := builder.Script()
				return script
			}(),
			shouldPass:  true,
			description: "Valid OP_0 <20-byte-hash> format",
		},
		{
			name:         "invalid length - too short",
			redeemScript: []byte{0x00, 0x14, 0x01, 0x02, 0x03}, // Only 5 bytes
			shouldPass:   false,
			description:  "Total length is not 22 bytes",
		},
		{
			name: "invalid length - too long",
			redeemScript: func() []byte {
				hash := make([]byte, 30) // 30 bytes instead of 20
				builder := txscript.NewScriptBuilder()
				builder.AddOp(txscript.OP_0)
				builder.AddData(hash)
				script, _ := builder.Script()
				return script
			}(),
			shouldPass:  false,
			description: "Total length is not 22 bytes",
		},
		{
			name: "invalid first byte - not OP_0",
			redeemScript: func() []byte {
				hash := make([]byte, 20)
				builder := txscript.NewScriptBuilder()
				builder.AddOp(txscript.OP_1) // Wrong opcode
				builder.AddData(hash)
				script, _ := builder.Script()
				return script
			}(),
			shouldPass:  false,
			description: "First byte must be OP_0 (0x00)",
		},
		{
			name: "invalid length byte",
			redeemScript: func() []byte {
				// Manually construct with wrong length byte
				script := make([]byte, 22)
				script[0] = 0x00 // OP_0
				script[1] = 0x15 // Wrong: 21 instead of 20 (0x14)
				// Fill rest with dummy data
				for i := 2; i < 22; i++ {
					script[i] = byte(i)
				}
				return script
			}(),
			shouldPass:  false,
			description: "Length byte must be 0x14 (20 decimal)",
		},
		{
			name: "all zeros",
			redeemScript: func() []byte {
				hash := make([]byte, 20) // All zeros
				builder := txscript.NewScriptBuilder()
				builder.AddOp(txscript.OP_0)
				builder.AddData(hash)
				script, _ := builder.Script()
				return script
			}(),
			shouldPass:  true,
			description: "Valid even with all-zero hash (edge case)",
		},
		{
			name: "all 0xFF",
			redeemScript: func() []byte {
				hash := make([]byte, 20)
				for i := range hash {
					hash[i] = 0xFF
				}
				builder := txscript.NewScriptBuilder()
				builder.AddOp(txscript.OP_0)
				builder.AddData(hash)
				script, _ := builder.Script()
				return script
			}(),
			shouldPass:  true,
			description: "Valid even with all-0xFF hash (edge case)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test by calling the actual function
			_, err := buildP2PKHScriptCodeForP2SHWPKH(tt.redeemScript, 0, 10000)
			if tt.shouldPass {
				assert.NoError(t, err, tt.description)
			} else {
				assert.Error(t, err, tt.description)
			}
		})
	}
}

// TestP2WPKHRedeemScriptExtraction tests that buildP2PKHScriptCodeForP2SHWPKH
// correctly constructs a P2PKH scriptCode from a valid P2WPKH redeemScript.
func TestP2WPKHRedeemScriptExtraction(t *testing.T) {
	// Create a valid P2WPKH redeemScript with known hash
	expectedHash := []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A,
		0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14,
	}

	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_0)
	builder.AddData(expectedHash)
	redeemScript, err := builder.Script()
	assert.NoError(t, err)

	// Call the function under test
	scriptCode, err := buildP2PKHScriptCodeForP2SHWPKH(redeemScript, 0, 10000)
	assert.NoError(t, err)

	// Verify the constructed scriptCode is a valid P2PKH script
	// Expected format: OP_DUP OP_HASH160 <20-byte-hash> OP_EQUALVERIFY OP_CHECKSIG
	expectedScriptCodeBuilder := txscript.NewScriptBuilder()
	expectedScriptCodeBuilder.AddOp(txscript.OP_DUP)
	expectedScriptCodeBuilder.AddOp(txscript.OP_HASH160)
	expectedScriptCodeBuilder.AddData(expectedHash)
	expectedScriptCodeBuilder.AddOp(txscript.OP_EQUALVERIFY)
	expectedScriptCodeBuilder.AddOp(txscript.OP_CHECKSIG)
	expectedScriptCode, err := expectedScriptCodeBuilder.Script()
	assert.NoError(t, err)

	assert.Equal(t, expectedScriptCode, scriptCode, "Constructed scriptCode should match expected P2PKH script")
}

// TestIsPayToWitnessPubKeyHash tests that btcd's IsPayToWitnessPubKeyHash
// correctly identifies valid P2WPKH scripts
func TestIsPayToWitnessPubKeyHash(t *testing.T) {
	tests := []struct {
		name         string
		redeemScript []byte
		expected     bool
	}{
		{
			name: "valid P2WPKH",
			redeemScript: func() []byte {
				hash := make([]byte, 20)
				builder := txscript.NewScriptBuilder()
				builder.AddOp(txscript.OP_0)
				builder.AddData(hash)
				script, _ := builder.Script()
				return script
			}(),
			expected: true,
		},
		{
			name: "invalid - wrong opcode",
			redeemScript: func() []byte {
				hash := make([]byte, 20)
				builder := txscript.NewScriptBuilder()
				builder.AddOp(txscript.OP_1)
				builder.AddData(hash)
				script, _ := builder.Script()
				return script
			}(),
			expected: false,
		},
		{
			name: "invalid - wrong hash length",
			redeemScript: func() []byte {
				hash := make([]byte, 32) // 32 bytes instead of 20
				builder := txscript.NewScriptBuilder()
				builder.AddOp(txscript.OP_0)
				builder.AddData(hash)
				script, _ := builder.Script()
				return script
			}(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := txscript.IsPayToWitnessPubKeyHash(tt.redeemScript)
			assert.Equal(t, tt.expected, result,
				"IsPayToWitnessPubKeyHash result mismatch for %s", tt.name)
		})
	}
}
