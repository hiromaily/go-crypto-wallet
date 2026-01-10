//nolint:revive // Bitcoin descriptors and xpub keys are necessarily long test data
package btc

import (
	"testing"

	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
)

func TestNewDescriptorParser(t *testing.T) {
	parser := NewDescriptorParser()
	if parser == nil {
		t.Fatal("NewDescriptorParser() returned nil")
	}
	if parser.keyRegex == nil {
		t.Fatal("DescriptorParser.keyRegex is nil")
	}
}

func TestDescriptorParser_Parse(t *testing.T) {
	parser := NewDescriptorParser()

	tests := []struct {
		name         string
		descriptor   string
		wantType     domainWallet.DescriptorType
		wantKeyCount int
		wantErr      bool
	}{
		{
			name:         "empty descriptor",
			descriptor:   "",
			wantType:     domainWallet.DescriptorTypeUnknown,
			wantKeyCount: 0,
			wantErr:      true,
		},
		{
			name:         "PKH descriptor with full metadata",
			descriptor:   "pkh([a1b2c3d4/44'/0'/0']xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL/0/*)",
			wantType:     domainWallet.DescriptorTypePKH,
			wantKeyCount: 1,
			wantErr:      false,
		},
		{
			name:         "WPKH descriptor (Bech32)",
			descriptor:   "wpkh([a1b2c3d4/84'/0'/0']xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL/0/*)",
			wantType:     domainWallet.DescriptorTypeWPKH,
			wantKeyCount: 1,
			wantErr:      false,
		},
		{
			name:         "TR descriptor (Taproot)",
			descriptor:   "tr([a1b2c3d4/86'/0'/0']xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL/0/*)",
			wantType:     domainWallet.DescriptorTypeTR,
			wantKeyCount: 1,
			wantErr:      false,
		},
		{
			name:         "SH(WPKH) descriptor (P2SH-SegWit)",
			descriptor:   "sh(wpkh([a1b2c3d4/49'/0'/0']xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL/0/*))",
			wantType:     domainWallet.DescriptorTypeSHWPKH,
			wantKeyCount: 1,
			wantErr:      false,
		},
		{
			name:         "descriptor with checksum",
			descriptor:   "wpkh([a1b2c3d4/84'/0'/0']xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL/0/*)#abcdef12",
			wantType:     domainWallet.DescriptorTypeWPKH,
			wantKeyCount: 1,
			wantErr:      false,
		},
		{
			name:         "descriptor with tpub (testnet)",
			descriptor:   "wpkh([a1b2c3d4/84'/1'/0']tpubD6NzVbkrYhZ4XgiXtGrdW5XDAPFCL9h7we1vwNCpn8tGbBcgfVYjXyhWo4E1xkh56hjod1RhGjxbaTLV3X4FyWuejifB9jusQ46QzG87VKp/0/*)",
			wantType:     domainWallet.DescriptorTypeWPKH,
			wantKeyCount: 1,
			wantErr:      false,
		},
		{
			name:         "unknown descriptor type",
			descriptor:   "unknown([a1b2c3d4/44'/0'/0']xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL/0/*)",
			wantType:     domainWallet.DescriptorTypeUnknown,
			wantKeyCount: 0,
			wantErr:      true,
		},
		{
			name:         "descriptor without key",
			descriptor:   "pkh()",
			wantType:     domainWallet.DescriptorTypePKH,
			wantKeyCount: 0,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc, err := parser.Parse(tt.descriptor)
			if (err != nil) != tt.wantErr {
				t.Errorf("DescriptorParser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			if desc.Type != tt.wantType {
				t.Errorf("DescriptorParser.Parse() descriptor type = %v, want %v", desc.Type, tt.wantType)
			}

			if len(desc.Keys) != tt.wantKeyCount {
				t.Errorf("DescriptorParser.Parse() key count = %d, want %d", len(desc.Keys), tt.wantKeyCount)
			}

			// Verify key extraction for successful parses
			if !tt.wantErr && tt.wantKeyCount > 0 {
				key := desc.Keys[0]
				if key.ExtendedPubKey == "" {
					t.Error("DescriptorParser.Parse() extracted key has empty extended public key")
				}
			}
		})
	}
}

func TestDescriptorParser_DetermineType(t *testing.T) {
	tests := []struct {
		name       string
		descriptor string
		wantType   domainWallet.DescriptorType
		wantErr    bool
	}{
		{"PKH", "pkh(key)", domainWallet.DescriptorTypePKH, false},
		{"SHWPKH", "sh(wpkh(key))", domainWallet.DescriptorTypeSHWPKH, false},
		{"WPKH", "wpkh(key)", domainWallet.DescriptorTypeWPKH, false},
		{"TR", "tr(key)", domainWallet.DescriptorTypeTR, false},
		{"WSH", "wsh(sortedmulti(2,key1,key2))", domainWallet.DescriptorTypeWSH, false},
		{"Unknown", "unknown(key)", domainWallet.DescriptorTypeUnknown, true},
		{"Empty", "", domainWallet.DescriptorTypeUnknown, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, err := determineType(tt.descriptor)
			if (err != nil) != tt.wantErr {
				t.Errorf("determineType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotType != tt.wantType {
				t.Errorf("determineType() = %v, want %v", gotType, tt.wantType)
			}
		})
	}
}

func TestDescriptorParser_ExtractKeys(t *testing.T) {
	parser := NewDescriptorParser()

	tests := []struct {
		name            string
		descriptor      string
		wantKeyCount    int
		wantFingerprint string
		wantPath        string
		wantErr         bool
	}{
		{
			name:            "key with full metadata",
			descriptor:      "pkh([a1b2c3d4/44'/0'/0']xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL/0/*)",
			wantKeyCount:    1,
			wantFingerprint: "a1b2c3d4",
			wantPath:        "/0/*", // Only derivation path from xpub, not origin path
			wantErr:         false,
		},
		{
			name:            "key without path after xpub",
			descriptor:      "pkh([a1b2c3d4/44'/0'/0']xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL)",
			wantKeyCount:    1,
			wantFingerprint: "a1b2c3d4",
			wantPath:        "", // No derivation path after xpub
			wantErr:         false,
		},
		{
			name:            "multiple keys (multisig)",
			descriptor:      "wsh(sortedmulti(2,[a1b2c3d4/48'/0'/0'/2']xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL,[b2c3d4e5/48'/0'/0'/2']xpub6D4BDPcP2GT577Vvch3R8wDkScZWzQzMMUm3PWbmWvVJrZwQY4VUNgqFJPMM3No2dFDFGTsxxpG5uJh7n7epu4trkrX7x7DogT5Uv6fcLW5/0/*))",
			wantKeyCount:    2,
			wantFingerprint: "a1b2c3d4",
			wantPath:        "", // First key has no derivation path after xpub (only second key has /0/*)
			wantErr:         false,
		},
		{
			name:         "no keys",
			descriptor:   "pkh()",
			wantKeyCount: 0,
			wantErr:      true,
		},
		{
			name:            "issue #297 - multisig with origin path and derivation path",
			descriptor:      "wsh(sortedmulti(2,[4a91842f/84'/1'/1]tpubDDed8zT1hv9vC5YpYDdufsNU5rGVr36q1ZrEM8Sx4y8NqEfCrqEfYfZjXXJP9f3ybvjkemCgvM6N7iJGrQhPHbKvDZZQWRpPJ7cj8Hio6Ty/0/*,[5bf99bc2/84'/1'/1]tpubDDHe7z2hv9vC5YpYDdufsNU5rGVr36q1ZrEM8Sx4y8NqEfCrqEfYfZjXXJP9f3ybvjkemCgvM6N7iJGrQhPHbKvDZZQWRpPJ7cj8Hio6Ty/0/*))",
			wantKeyCount:    2,
			wantFingerprint: "4a91842f",
			wantPath:        "/0/*", // Should only contain derivation path FROM xpub, not origin path
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys, err := parser.extractKeys(tt.descriptor)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			if len(keys) != tt.wantKeyCount {
				t.Errorf("extractKeys() key count = %d, want %d", len(keys), tt.wantKeyCount)
				return
			}

			if tt.wantKeyCount > 0 {
				firstKey := keys[0]
				if tt.wantFingerprint != "" && firstKey.Fingerprint != tt.wantFingerprint {
					t.Errorf("extractKeys() fingerprint = %s, want %s", firstKey.Fingerprint, tt.wantFingerprint)
				}
				if tt.wantPath != "" && firstKey.DerivationPath != tt.wantPath {
					t.Errorf("extractKeys() path = %s, want %s", firstKey.DerivationPath, tt.wantPath)
				}
				if firstKey.ExtendedPubKey == "" {
					t.Error("extractKeys() extended public key is empty")
				}

				// For issue #297 test case, validate second key as well to ensure
				// both keys in the multisig descriptor are parsed correctly
				if tt.name == "issue #297 - multisig with origin path and derivation path" && len(keys) > 1 {
					secondKey := keys[1]
					expectedFingerprint := "5bf99bc2"
					expectedPath := "/0/*"
					if secondKey.Fingerprint != expectedFingerprint {
						t.Errorf("extractKeys() second key fingerprint = %s, want %s", secondKey.Fingerprint, expectedFingerprint)
					}
					if secondKey.DerivationPath != expectedPath {
						t.Errorf("extractKeys() second key path = %s, want %s", secondKey.DerivationPath, expectedPath)
					}
					if secondKey.ExtendedPubKey == "" {
						t.Error("extractKeys() second key extended public key is empty")
					}
				}
			}
		})
	}
}

func TestDescriptorParser_FormatDescriptor(t *testing.T) {
	tests := []struct {
		name       string
		descriptor *domainWallet.Descriptor
		wantScript string
		wantErr    bool
	}{
		{
			name: "valid descriptor without checksum",
			descriptor: &domainWallet.Descriptor{
				Type:   domainWallet.DescriptorTypeWPKH,
				Script: "wpkh([a1b2c3d4/84'/0'/0']xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL/0/*)",
				Keys: []domainWallet.DescriptorKey{
					{
						Fingerprint:    "a1b2c3d4",
						DerivationPath: "/0/*", // Only derivation path from xpub
						ExtendedPubKey: "xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL",
					},
				},
				Checksum: "",
			},
			wantScript: "wpkh([a1b2c3d4/84'/0'/0']xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL/0/*)",
			wantErr:    false,
		},
		{
			name: "valid descriptor with checksum",
			descriptor: &domainWallet.Descriptor{
				Type:   domainWallet.DescriptorTypeWPKH,
				Script: "wpkh([a1b2c3d4/84'/0'/0']xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL/0/*)",
				Keys: []domainWallet.DescriptorKey{
					{
						Fingerprint:    "a1b2c3d4",
						DerivationPath: "/0/*", // Only derivation path from xpub
						ExtendedPubKey: "xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL",
					},
				},
				Checksum: "abcd1234",
			},
			wantScript: "wpkh([a1b2c3d4/84'/0'/0']xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL/0/*)#abcd1234",
			wantErr:    false,
		},
		{
			name:       "invalid descriptor",
			descriptor: nil,
			wantScript: "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotScript, err := FormatDescriptor(tt.descriptor)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatDescriptor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotScript != tt.wantScript {
				t.Errorf("FormatDescriptor() = %v, want %v", gotScript, tt.wantScript)
			}
		})
	}
}
