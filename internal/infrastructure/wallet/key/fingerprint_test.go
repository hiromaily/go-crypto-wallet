package key

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
)

// BIP32 Test Vector 1 from https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki#test-vectors
const (
	testVector1Seed = "000102030405060708090a0b0c0d0e0f"
	// Master key for test vector 1
	//nolint:revive // Bitcoin xpub keys are necessarily long (111 chars)
	testVector1MasterXpub = "xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet8"
	// Expected fingerprint: 3442193e (first 4 bytes of HASH160 of master pubkey)
	testVector1Fingerprint = "3442193e"
)

// BIP32 Test Vector 2
const (
	//nolint:revive // BIP32 test vector seed is necessarily long (128 chars)
	testVector2Seed = "fffcf9f6f3f0edeae7e4e1dedbd8d5d2cfccc9c6c3c0bdbab7b4b1aeaba8a5a29f9c999693908d8a8784817e7b7875726f6c696663605d5a5754514e4b484542"
	//nolint:revive // Bitcoin xpub keys are necessarily long (111 chars)
	testVector2MasterXpub = "xpub661MyMwAqRbcFW31YEwpkMuc5THy2PSt5bDMsktWQcFF8syAmRUapSCGu8ED9W6oDMSgv6Zz8idoc4a6mr8BDzTJY47LJhkJ8UB7WEGuduB"
	// Expected fingerprint: bd16bee5
	testVector2Fingerprint = "bd16bee5"
)

func TestCalculateMasterFingerprint(t *testing.T) {
	tests := []struct {
		name            string
		setupPubKey     func() *btcec.PublicKey
		wantFingerprint string
		wantErr         bool
		errMsg          string
	}{
		{
			name: "BIP32 test vector 1",
			setupPubKey: func() *btcec.PublicKey {
				// Parse master xpub from test vector 1
				key, err := hdkeychain.NewKeyFromString(testVector1MasterXpub)
				if err != nil {
					t.Fatalf("failed to parse test vector 1 xpub: %v", err)
				}
				pubKey, err := key.ECPubKey()
				if err != nil {
					t.Fatalf("failed to get public key: %v", err)
				}
				return pubKey
			},
			wantFingerprint: testVector1Fingerprint,
			wantErr:         false,
		},
		{
			name: "BIP32 test vector 2",
			setupPubKey: func() *btcec.PublicKey {
				// Parse master xpub from test vector 2
				key, err := hdkeychain.NewKeyFromString(testVector2MasterXpub)
				if err != nil {
					t.Fatalf("failed to parse test vector 2 xpub: %v", err)
				}
				pubKey, err := key.ECPubKey()
				if err != nil {
					t.Fatalf("failed to get public key: %v", err)
				}
				return pubKey
			},
			wantFingerprint: testVector2Fingerprint,
			wantErr:         false,
		},
		{
			name: "nil public key",
			setupPubKey: func() *btcec.PublicKey {
				return nil
			},
			wantErr: true,
			errMsg:  "master public key is nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pubKey := tt.setupPubKey()
			got, err := CalculateMasterFingerprint(pubKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateMasterFingerprint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf(
						"CalculateMasterFingerprint() error = %v, want error containing %v",
						err.Error(),
						tt.errMsg,
					)
				}
				return
			}
			if !tt.wantErr {
				if got.String() != tt.wantFingerprint {
					t.Errorf("CalculateMasterFingerprint() = %v, want %v", got, tt.wantFingerprint)
				}
				// Verify fingerprint length
				if len(got.String()) != 8 {
					t.Errorf("CalculateMasterFingerprint() fingerprint length = %d, want 8", len(got.String()))
				}
			}
		})
	}
}

func TestFingerprintFromExtendedKey(t *testing.T) {
	tests := []struct {
		name            string
		extendedKey     string
		wantFingerprint string
		wantErr         bool
		errMsg          string
	}{
		{
			name:            "BIP32 test vector 1 master xpub",
			extendedKey:     testVector1MasterXpub,
			wantFingerprint: testVector1Fingerprint,
			wantErr:         false,
		},
		{
			name:            "BIP32 test vector 2 master xpub",
			extendedKey:     testVector2MasterXpub,
			wantFingerprint: testVector2Fingerprint,
			wantErr:         false,
		},
		{
			name:        "empty extended key",
			extendedKey: "",
			wantErr:     true,
			errMsg:      "extended key is empty",
		},
		{
			name:        "invalid extended key",
			extendedKey: "invalid-xpub",
			wantErr:     true,
			errMsg:      "failed to parse extended key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FingerprintFromExtendedKey(tt.extendedKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("FingerprintFromExtendedKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf(
						"FingerprintFromExtendedKey() error = %v, want error containing %v",
						err.Error(),
						tt.errMsg,
					)
				}
				return
			}
			if !tt.wantErr {
				if got.String() != tt.wantFingerprint {
					t.Errorf("FingerprintFromExtendedKey() = %v, want %v", got, tt.wantFingerprint)
				}
			}
		})
	}
}

func TestFingerprintFromSeed(t *testing.T) {
	tests := []struct {
		name            string
		seedHex         string
		params          *chaincfg.Params
		wantFingerprint string
		wantErr         bool
		errMsg          string
	}{
		{
			name:            "BIP32 test vector 1",
			seedHex:         testVector1Seed,
			params:          &chaincfg.MainNetParams,
			wantFingerprint: testVector1Fingerprint,
			wantErr:         false,
		},
		{
			name:            "BIP32 test vector 2",
			seedHex:         testVector2Seed,
			params:          &chaincfg.MainNetParams,
			wantFingerprint: testVector2Fingerprint,
			wantErr:         false,
		},
		{
			name:    "empty seed",
			seedHex: "",
			params:  &chaincfg.MainNetParams,
			wantErr: true,
			errMsg:  "seed is empty",
		},
		{
			name:    "nil params",
			seedHex: testVector1Seed,
			params:  nil,
			wantErr: true,
			errMsg:  "network params is nil",
		},
		{
			name:    "invalid seed hex",
			seedHex: "invalid-hex",
			params:  &chaincfg.MainNetParams,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var seed []byte
			var err error
			if tt.seedHex != "" {
				seed, err = hex.DecodeString(tt.seedHex)
				if err != nil && !tt.wantErr {
					t.Fatalf("failed to decode seed hex: %v", err)
				}
			}

			got, err := FingerprintFromSeed(seed, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FingerprintFromSeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("FingerprintFromSeed() error = %v, want error containing %v", err.Error(), tt.errMsg)
				}
				return
			}
			if !tt.wantErr {
				if got.String() != tt.wantFingerprint {
					t.Errorf("FingerprintFromSeed() = %v, want %v", got, tt.wantFingerprint)
				}
			}
		})
	}
}

// TestFingerprintConsistency verifies that all three fingerprint calculation methods
// produce the same result for the same master key.
func TestFingerprintConsistency(t *testing.T) {
	seedBytes, err := hex.DecodeString(testVector1Seed)
	if err != nil {
		t.Fatalf("failed to decode test vector seed: %v", err)
	}

	// Method 1: From seed
	fp1, err := FingerprintFromSeed(seedBytes, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatalf("FingerprintFromSeed() error = %v", err)
	}

	// Method 2: From extended key
	fp2, err := FingerprintFromExtendedKey(testVector1MasterXpub)
	if err != nil {
		t.Fatalf("FingerprintFromExtendedKey() error = %v", err)
	}

	// Method 3: From public key directly
	key, err := hdkeychain.NewKeyFromString(testVector1MasterXpub)
	if err != nil {
		t.Fatalf("failed to parse xpub: %v", err)
	}
	pubKey, err := key.ECPubKey()
	if err != nil {
		t.Fatalf("failed to get public key: %v", err)
	}
	fp3, err := CalculateMasterFingerprint(pubKey)
	if err != nil {
		t.Fatalf("CalculateMasterFingerprint() error = %v", err)
	}

	// All three methods should produce the same fingerprint
	if fp1.String() != fp2.String() {
		t.Errorf("FingerprintFromSeed() = %v, FingerprintFromExtendedKey() = %v, want equal", fp1, fp2)
	}
	if fp1.String() != fp3.String() {
		t.Errorf("FingerprintFromSeed() = %v, CalculateMasterFingerprint() = %v, want equal", fp1, fp3)
	}
	if fp2.String() != fp3.String() {
		t.Errorf("FingerprintFromExtendedKey() = %v, CalculateMasterFingerprint() = %v, want equal", fp2, fp3)
	}

	// Verify against expected test vector
	if fp1.String() != testVector1Fingerprint {
		t.Errorf("fingerprint = %v, want %v (from BIP32 test vector)", fp1, testVector1Fingerprint)
	}
}
