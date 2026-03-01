//nolint:revive // Bitcoin descriptors and xpub keys are necessarily long test data
package btc

import (
	"testing"

	"github.com/stretchr/testify/require"

	btcdescriptor "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/descriptor"
)

func TestNewDescriptorParser(t *testing.T) {
	parser := NewDescriptorParser()
	require.NotNil(t, parser, "NewDescriptorParser() returned nil")
}

func TestDescriptorParser_Parse(t *testing.T) {
	parser := NewDescriptorParser()

	tests := []struct {
		name         string
		descriptor   string
		wantType     btcdescriptor.DescriptorType
		wantKeyCount int
		wantErr      bool
	}{
		{
			name:         "empty descriptor",
			descriptor:   "",
			wantType:     btcdescriptor.DescriptorTypeUnknown,
			wantKeyCount: 0,
			wantErr:      true,
		},
		{
			name:         "PKH descriptor with full metadata",
			descriptor:   "pkh([a1b2c3d4/44'/0'/0']xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL/0/*)",
			wantType:     btcdescriptor.DescriptorTypePKH,
			wantKeyCount: 1,
			wantErr:      false,
		},
		{
			name:         "WPKH descriptor (Bech32)",
			descriptor:   "wpkh([a1b2c3d4/84'/0'/0']xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL/0/*)",
			wantType:     btcdescriptor.DescriptorTypeWPKH,
			wantKeyCount: 1,
			wantErr:      false,
		},
		{
			name:         "TR descriptor (Taproot)",
			descriptor:   "tr([a1b2c3d4/86'/0'/0']xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL/0/*)",
			wantType:     btcdescriptor.DescriptorTypeTR,
			wantKeyCount: 1,
			wantErr:      false,
		},
		{
			name:         "SH(WPKH) descriptor (P2SH-SegWit)",
			descriptor:   "sh(wpkh([a1b2c3d4/49'/0'/0']xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL/0/*))",
			wantType:     btcdescriptor.DescriptorTypeSHWPKH,
			wantKeyCount: 1,
			wantErr:      false,
		},
		{
			name:         "descriptor with checksum",
			descriptor:   "wpkh([a1b2c3d4/84'/0'/0']xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL/0/*)#abcdef12",
			wantType:     btcdescriptor.DescriptorTypeWPKH,
			wantKeyCount: 1,
			wantErr:      false,
		},
		{
			name:         "descriptor with tpub (testnet)",
			descriptor:   "wpkh([a1b2c3d4/84'/1'/0']tpubD6NzVbkrYhZ4XgiXtGrdW5XDAPFCL9h7we1vwNCpn8tGbBcgfVYjXyhWo4E1xkh56hjod1RhGjxbaTLV3X4FyWuejifB9jusQ46QzG87VKp/0/*)",
			wantType:     btcdescriptor.DescriptorTypeWPKH,
			wantKeyCount: 1,
			wantErr:      false,
		},
		{
			name:         "unknown descriptor type",
			descriptor:   "unknown([a1b2c3d4/44'/0'/0']xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL/0/*)",
			wantType:     btcdescriptor.DescriptorTypeUnknown,
			wantKeyCount: 0,
			wantErr:      true,
		},
		{
			name:         "descriptor without key",
			descriptor:   "pkh()",
			wantType:     btcdescriptor.DescriptorTypePKH,
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

func TestDescriptorParser_FormatDescriptor(t *testing.T) {
	tests := []struct {
		name       string
		descriptor *btcdescriptor.Descriptor
		wantScript string
		wantErr    bool
	}{
		{
			name: "valid descriptor without checksum",
			descriptor: &btcdescriptor.Descriptor{
				Type:   btcdescriptor.DescriptorTypeWPKH,
				Script: "wpkh([a1b2c3d4/84'/0'/0']xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL/0/*)",
				Keys: []btcdescriptor.DescriptorKey{
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
			descriptor: &btcdescriptor.Descriptor{
				Type:   btcdescriptor.DescriptorTypeWPKH,
				Script: "wpkh([a1b2c3d4/84'/0'/0']xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL/0/*)",
				Keys: []btcdescriptor.DescriptorKey{
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
