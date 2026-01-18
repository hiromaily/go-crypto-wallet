//nolint:revive // Test file contains long test data
package key

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key/strategy"

	"github.com/btcsuite/btcd/chaincfg"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
)

// Test seed (BIP39 mnemonic: "abandon abandon ... abandon about")
const testSeedHex = "5eb00bbddcf069084889a8ab9155568165f5c453ccb85e70811aaed6f6da5fc19a5ac40b389cd370d086206dec8aa6c43daea6690f20ad3d8d48b2d2ce9e38e4" //nolint:lll

func getTestSeed(t *testing.T) []byte {
	t.Helper()
	seed, err := hex.DecodeString(testSeedHex)
	if err != nil {
		t.Fatalf("Failed to decode test seed: %v", err)
	}
	return seed
}

func TestDescriptorGenerator_GetMasterFingerprint(t *testing.T) {
	seed := getTestSeed(t)

	tests := []struct {
		name    string
		purpose PurposeType
		wantErr bool
	}{
		{
			name:    "BIP44 master fingerprint",
			purpose: PurposeTypeBIP44,
			wantErr: false,
		},
		{
			name:    "BIP49 master fingerprint",
			purpose: PurposeTypeBIP49,
			wantErr: false,
		},
		{
			name:    "BIP84 master fingerprint",
			purpose: PurposeTypeBIP84,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create coin strategy
			coinStrategy, stratErr := strategy.CreateCoinKeyStrategy(domainCoin.BTC, &chaincfg.RegressionNetParams)
			if stratErr != nil {
				t.Fatalf("Failed to create coin strategy: %v", stratErr)
			}
			hdKey := NewHDKey(tt.purpose, domainCoin.BTC, &chaincfg.RegressionNetParams, coinStrategy)
			generator := NewDescriptorGenerator(hdKey, &chaincfg.RegressionNetParams)

			fingerprint, err := generator.GetMasterFingerprint(seed)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMasterFingerprint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify fingerprint format (8 hex characters)
				if len(fingerprint) != 8 {
					t.Errorf("GetMasterFingerprint() fingerprint length = %d, want 8", len(fingerprint))
				}

				// Verify all characters are hex
				for _, c := range fingerprint {
					if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
						t.Errorf("GetMasterFingerprint() invalid hex character: %c", c)
					}
				}

				t.Logf("Master fingerprint: %s", fingerprint)
			}
		})
	}
}

func TestDescriptorGenerator_GetAccountXPub(t *testing.T) {
	seed := getTestSeed(t)

	tests := []struct {
		name        string
		purpose     PurposeType
		accountType domainAccount.AccountType
		wantPrefix  string
		wantErr     bool
	}{
		{
			name:        "BIP44 deposit account xpub",
			purpose:     PurposeTypeBIP44,
			accountType: domainAccount.AccountTypeDeposit,
			wantPrefix:  "tpub", // Testnet/Regtest prefix
			wantErr:     false,
		},
		{
			name:        "BIP49 payment account xpub",
			purpose:     PurposeTypeBIP49,
			accountType: domainAccount.AccountTypePayment,
			wantPrefix:  "tpub",
			wantErr:     false,
		},
		{
			name:        "BIP84 stored account xpub",
			purpose:     PurposeTypeBIP84,
			accountType: domainAccount.AccountTypeStored,
			wantPrefix:  "tpub",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create coin strategy
			coinStrategy, stratErr := strategy.CreateCoinKeyStrategy(domainCoin.BTC, &chaincfg.RegressionNetParams)
			if stratErr != nil {
				t.Fatalf("Failed to create coin strategy: %v", stratErr)
			}
			hdKey := NewHDKey(tt.purpose, domainCoin.BTC, &chaincfg.RegressionNetParams, coinStrategy)
			generator := NewDescriptorGenerator(hdKey, &chaincfg.RegressionNetParams)

			xpub, err := generator.GetAccountXPub(seed, tt.accountType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountXPub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify xpub format
				if !strings.HasPrefix(xpub, tt.wantPrefix) {
					t.Errorf("GetAccountXPub() xpub prefix = %s, want prefix %s", xpub[:4], tt.wantPrefix)
				}

				// Verify reasonable xpub length (should be around 111 characters)
				if len(xpub) < 100 || len(xpub) > 120 {
					t.Errorf("GetAccountXPub() xpub length = %d, want 100-120", len(xpub))
				}

				t.Logf("Account xpub: %s", xpub)
			}
		})
	}
}

func TestDescriptorGenerator_GenerateAccountDescriptor(t *testing.T) {
	seed := getTestSeed(t)

	tests := []struct {
		name         string
		purpose      PurposeType
		accountType  domainAccount.AccountType
		withChecksum bool
		wantErr      bool
	}{
		{
			name:         "BIP49 P2SH-SegWit descriptor without checksum",
			purpose:      PurposeTypeBIP49,
			accountType:  domainAccount.AccountTypeDeposit,
			withChecksum: false,
			wantErr:      false,
		},
		{
			name:         "BIP49 P2SH-SegWit descriptor with checksum",
			purpose:      PurposeTypeBIP49,
			accountType:  domainAccount.AccountTypePayment,
			withChecksum: true,
			wantErr:      false,
		},
		{
			name:         "BIP84 Native SegWit descriptor",
			purpose:      PurposeTypeBIP84,
			accountType:  domainAccount.AccountTypeStored,
			withChecksum: true,
			wantErr:      false,
		},
		{
			name:         "BIP44 P2PKH descriptor",
			purpose:      PurposeTypeBIP44,
			accountType:  domainAccount.AccountTypeClient,
			withChecksum: true,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create coin strategy
			coinStrategy, stratErr := strategy.CreateCoinKeyStrategy(domainCoin.BTC, &chaincfg.RegressionNetParams)
			if stratErr != nil {
				t.Fatalf("Failed to create coin strategy: %v", stratErr)
			}
			hdKey := NewHDKey(tt.purpose, domainCoin.BTC, &chaincfg.RegressionNetParams, coinStrategy)
			generator := NewDescriptorGenerator(hdKey, &chaincfg.RegressionNetParams)

			descriptor, err := generator.GenerateAccountDescriptor(seed, tt.accountType, tt.withChecksum)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateAccountDescriptor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify descriptor structure
				if descriptor == nil {
					t.Fatal("GenerateAccountDescriptor() returned nil descriptor")
				}

				// Verify script is not empty
				if descriptor.Script == "" {
					t.Error("GenerateAccountDescriptor() descriptor script is empty")
				}

				// Verify keys are present
				if len(descriptor.Keys) != 1 {
					t.Errorf("GenerateAccountDescriptor() keys count = %d, want 1", len(descriptor.Keys))
				}

				// Verify checksum presence
				if tt.withChecksum && descriptor.Checksum == "" {
					t.Error("GenerateAccountDescriptor() checksum is empty but was requested")
				}
				if !tt.withChecksum && descriptor.Checksum != "" {
					t.Error("GenerateAccountDescriptor() checksum is present but was not requested")
				}

				// Verify descriptor format based on purpose
				expectedPrefix := ""
				switch tt.purpose {
				case PurposeTypeBIP44:
					expectedPrefix = "pkh("
				case PurposeTypeBIP49:
					expectedPrefix = "sh(wpkh("
				case PurposeTypeBIP84:
					expectedPrefix = "wpkh("
				case PurposeTypeBIP86:
					expectedPrefix = "tr("
				}

				if !strings.HasPrefix(descriptor.Script, expectedPrefix) {
					t.Errorf("GenerateAccountDescriptor() descriptor prefix = %s, want %s", descriptor.Script[:10], expectedPrefix)
				}

				t.Logf("Generated descriptor: %s", descriptor.Script)
				if descriptor.Checksum != "" {
					t.Logf("Checksum: %s", descriptor.Checksum)
				}
			}
		})
	}
}

func TestDescriptorGenerator_GenerateAccountDescriptorWithRange(t *testing.T) {
	seed := getTestSeed(t)

	tests := []struct {
		name         string
		purpose      PurposeType
		accountType  domainAccount.AccountType
		change       uint32
		wildcard     bool
		withChecksum bool
		wantContains string
		wantErr      bool
	}{
		{
			name:         "P2SH-SegWit with wildcard range",
			purpose:      PurposeTypeBIP49,
			accountType:  domainAccount.AccountTypeDeposit,
			change:       0,
			wildcard:     true,
			withChecksum: true,
			wantContains: "/0/*",
			wantErr:      false,
		},
		{
			name:         "Native SegWit with internal chain",
			purpose:      PurposeTypeBIP84,
			accountType:  domainAccount.AccountTypePayment,
			change:       1,
			wildcard:     true,
			withChecksum: true,
			wantContains: "/1/*",
			wantErr:      false,
		},
		{
			name:         "P2SH-SegWit with specific index",
			purpose:      PurposeTypeBIP49,
			accountType:  domainAccount.AccountTypeStored,
			change:       0,
			wildcard:     false,
			withChecksum: false,
			wantContains: "/0",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create coin strategy
			coinStrategy, stratErr := strategy.CreateCoinKeyStrategy(domainCoin.BTC, &chaincfg.RegressionNetParams)
			if stratErr != nil {
				t.Fatalf("Failed to create coin strategy: %v", stratErr)
			}
			hdKey := NewHDKey(tt.purpose, domainCoin.BTC, &chaincfg.RegressionNetParams, coinStrategy)
			generator := NewDescriptorGenerator(hdKey, &chaincfg.RegressionNetParams)

			descriptorStr, err := generator.GenerateAccountDescriptorWithRange(
				seed,
				tt.accountType,
				tt.change,
				tt.wildcard,
				tt.withChecksum,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateAccountDescriptorWithRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify descriptor contains range
				if !strings.Contains(descriptorStr, tt.wantContains) {
					t.Errorf("GenerateAccountDescriptorWithRange() descriptor does not contain %s", tt.wantContains)
				}

				// Verify checksum format if requested
				if tt.withChecksum {
					if !strings.Contains(descriptorStr, "#") {
						t.Error("GenerateAccountDescriptorWithRange() missing checksum separator '#'")
					}
					parts := strings.Split(descriptorStr, "#")
					if len(parts) != 2 {
						t.Error("GenerateAccountDescriptorWithRange() invalid checksum format")
					}
					if len(parts[1]) != 8 {
						t.Errorf("GenerateAccountDescriptorWithRange() checksum length = %d, want 8", len(parts[1]))
					}
				}

				t.Logf("Generated descriptor with range: %s", descriptorStr)
			}
		})
	}
}

func TestDescriptorGenerator_Integration(t *testing.T) {
	// Integration test: generate descriptor, verify it has all required components
	seed := getTestSeed(t)
	// Create coin strategy
	coinStrategy, stratErr := strategy.CreateCoinKeyStrategy(domainCoin.BTC, &chaincfg.RegressionNetParams)
	if stratErr != nil {
		t.Fatalf("Failed to create coin strategy: %v", stratErr)
	}
	hdKey := NewHDKey(PurposeTypeBIP49, domainCoin.BTC, &chaincfg.RegressionNetParams, coinStrategy)
	generator := NewDescriptorGenerator(hdKey, &chaincfg.RegressionNetParams)

	// Generate complete descriptor with range and checksum
	descriptorStr, err := generator.GenerateAccountDescriptorWithRange(
		seed,
		domainAccount.AccountTypeDeposit,
		0,    // external chain
		true, // wildcard
		true, // with checksum
	)
	if err != nil {
		t.Fatalf("Integration test failed: %v", err)
	}

	// Verify complete descriptor format
	// Expected format: sh(wpkh([fingerprint/49'/1'/0']tpub.../0/*))#checksum
	t.Logf("Complete descriptor: %s", descriptorStr)

	// Verify structure
	if !strings.HasPrefix(descriptorStr, "sh(wpkh([") {
		t.Errorf("Descriptor does not start with 'sh(wpkh([': %s", descriptorStr[:20])
	}

	if !strings.Contains(descriptorStr, "/49'/") {
		t.Error("Descriptor missing purpose /49'/")
	}

	if !strings.Contains(descriptorStr, "tpub") {
		t.Error("Descriptor missing xpub (tpub)")
	}

	if !strings.Contains(descriptorStr, "/0/*") {
		t.Error("Descriptor missing wildcard range /0/*")
	}

	if !strings.Contains(descriptorStr, "#") {
		t.Error("Descriptor missing checksum separator #")
	}

	// Verify checksum is exactly 8 characters
	parts := strings.Split(descriptorStr, "#")
	if len(parts) == 2 && len(parts[1]) != 8 {
		t.Errorf("Checksum length = %d, want 8", len(parts[1]))
	}
}
