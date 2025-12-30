package wallet

import (
	"testing"
)

func TestDescriptorType_String(t *testing.T) {
	tests := []struct {
		name     string
		descType DescriptorType
		want     string
	}{
		{"PKH", DescriptorTypePKH, "pkh"},
		{"SHWPKH", DescriptorTypeSHWPKH, "sh(wpkh)"},
		{"WPKH", DescriptorTypeWPKH, "wpkh"},
		{"TR", DescriptorTypeTR, "tr"},
		{"WSH", DescriptorTypeWSH, "wsh"},
		{"Unknown", DescriptorTypeUnknown, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.descType.String(); got != tt.want {
				t.Errorf("DescriptorType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateDescriptor(t *testing.T) {
	tests := []struct {
		name    string
		desc    *Descriptor
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil descriptor",
			desc:    nil,
			wantErr: true,
			errMsg:  "descriptor is nil",
		},
		{
			name: "unknown descriptor type",
			desc: &Descriptor{
				Type:   DescriptorTypeUnknown,
				Script: "unknown(xpub...)",
				Keys:   []DescriptorKey{{ExtendedPubKey: "xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL"}},
			},
			wantErr: true,
			errMsg:  "unknown descriptor type",
		},
		{
			name: "empty script",
			desc: &Descriptor{
				Type:   DescriptorTypePKH,
				Script: "",
				Keys:   []DescriptorKey{{ExtendedPubKey: "xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL"}},
			},
			wantErr: true,
			errMsg:  "descriptor script is empty",
		},
		{
			name: "no keys",
			desc: &Descriptor{
				Type:   DescriptorTypePKH,
				Script: "pkh(xpub...)",
				Keys:   []DescriptorKey{},
			},
			wantErr: true,
			errMsg:  "descriptor must have at least one key",
		},
		{
			name: "invalid key",
			desc: &Descriptor{
				Type:   DescriptorTypePKH,
				Script: "pkh([a1b2c3d4/44'/0'/0']xpub.../0/*)",
				Keys:   []DescriptorKey{{Fingerprint: "a1b2c3d4", DerivationPath: "/44'/0'/0'", ExtendedPubKey: ""}},
			},
			wantErr: true,
			errMsg:  "invalid key at index 0",
		},
		{
			name: "valid descriptor",
			desc: &Descriptor{
				Type:   DescriptorTypePKH,
				Script: "pkh([a1b2c3d4/44'/0'/0']xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL/0/*)",
				Keys: []DescriptorKey{
					{
						Fingerprint:    "a1b2c3d4",
						DerivationPath: "/44'/0'/0'",
						ExtendedPubKey: "xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDescriptor(tt.desc)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDescriptor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if err.Error() != tt.errMsg {
					if !contains(err.Error(), tt.errMsg) {
						t.Errorf("ValidateDescriptor() error = %v, want error containing %v", err.Error(), tt.errMsg)
					}
				}
			}
		})
	}
}

func TestValidateDescriptorKey(t *testing.T) {
	tests := []struct {
		name    string
		key     *DescriptorKey
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil key",
			key:     nil,
			wantErr: true,
			errMsg:  "descriptor key is nil",
		},
		{
			name: "empty extended public key",
			key: &DescriptorKey{
				Fingerprint:    "a1b2c3d4",
				DerivationPath: "/44'/0'/0'",
				ExtendedPubKey: "",
			},
			wantErr: true,
			errMsg:  "extended public key is empty",
		},
		{
			name: "invalid fingerprint",
			key: &DescriptorKey{
				Fingerprint:    "invalid",
				DerivationPath: "/44'/0'/0'",
				ExtendedPubKey: "xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL",
			},
			wantErr: true,
			errMsg:  "invalid fingerprint",
		},
		{
			name: "invalid derivation path",
			key: &DescriptorKey{
				Fingerprint:    "a1b2c3d4",
				DerivationPath: "44'/0'/0'", // Missing leading /
				ExtendedPubKey: "xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL",
			},
			wantErr: true,
			errMsg:  "invalid derivation path",
		},
		{
			name: "invalid extended public key",
			key: &DescriptorKey{
				Fingerprint:    "a1b2c3d4",
				DerivationPath: "/44'/0'/0'",
				ExtendedPubKey: "invalid-xpub",
			},
			wantErr: true,
			errMsg:  "invalid extended public key",
		},
		{
			name: "valid key with all fields",
			key: &DescriptorKey{
				Fingerprint:    "a1b2c3d4",
				DerivationPath: "/44'/0'/0'",
				ExtendedPubKey: "xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL",
			},
			wantErr: false,
		},
		{
			name: "valid key without optional fields",
			key: &DescriptorKey{
				Fingerprint:    "",
				DerivationPath: "",
				ExtendedPubKey: "xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDescriptorKey(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDescriptorKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateDescriptorKey() error = %v, want error containing %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestValidateFingerprint(t *testing.T) {
	tests := []struct {
		name        string
		fingerprint string
		wantErr     bool
	}{
		{"valid fingerprint lowercase", "a1b2c3d4", false},
		{"valid fingerprint uppercase", "A1B2C3D4", false},
		{"valid fingerprint mixed case", "a1B2c3D4", false},
		{"empty fingerprint", "", true},
		{"too short", "a1b2c3", true},
		{"too long", "a1b2c3d4e5", true},
		{"invalid characters", "g1h2i3j4", true},
		{"special characters", "a1b2-3d4", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFingerprint(tt.fingerprint)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFingerprint() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateDerivationPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"valid BIP44 path", "/44'/0'/0'", false},
		{"valid BIP84 path", "/84'/0'/0'", false},
		{"valid BIP86 path", "/86'/0'/0'", false},
		{"valid with wildcard", "/44'/0'/0'/0/*", false},
		{"valid with m prefix", "m/44'/0'/0'", false},
		{"valid hardened with h", "/44h/0h/0h", false},
		{"empty path", "", true},
		{"missing leading slash", "44'/0'/0'", true},
		{"wildcard in middle", "/44'/*/0'", true},
		{"invalid segment", "/44'/0'/abc'", true},
		{"double slash", "/44'//0'/0'", true},
		{"trailing slash", "/44'/0'/0'/", true}, // Trailing slash creates empty segment
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDerivationPath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDerivationPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateExtendedPubKey(t *testing.T) {
	tests := []struct {
		name    string
		xpub    string
		wantErr bool
	}{
		{
			"valid xpub",
			"xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL",
			false,
		},
		{
			"valid tpub",
			"tpubD6NzVbkrYhZ4XgiXtGrdW5XDAPFCL9h7we1vwNCpn8tGbBcgfVYjXyhWo4E1xkh56hjod1RhGjxbaTLV3X4FyWuejifB9jusQ46QzG87VKp",
			false,
		},
		{"empty xpub", "", true},
		{"invalid prefix", "apub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL", true},
		{"too short", "xpub123", true},
		{"invalid base58 chars", "xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcE0", true}, // Contains 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateExtendedPubKey(tt.xpub)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateExtendedPubKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
