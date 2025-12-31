package key

import (
	"strings"
	"testing"
)

func TestNewFingerprint(t *testing.T) {
	tests := []struct {
		name    string
		hex     string
		want    string
		wantErr bool
	}{
		{
			name:    "valid lowercase fingerprint",
			hex:     "a1b2c3d4",
			want:    "a1b2c3d4",
			wantErr: false,
		},
		{
			name:    "valid uppercase fingerprint",
			hex:     "A1B2C3D4",
			want:    "A1B2C3D4",
			wantErr: false,
		},
		{
			name:    "valid mixed case fingerprint",
			hex:     "a1B2c3D4",
			want:    "a1B2c3D4",
			wantErr: false,
		},
		{
			name:    "empty fingerprint",
			hex:     "",
			wantErr: true,
		},
		{
			name:    "too short",
			hex:     "a1b2c3",
			wantErr: true,
		},
		{
			name:    "too long",
			hex:     "a1b2c3d4e5",
			wantErr: true,
		},
		{
			name:    "invalid characters",
			hex:     "g1h2i3j4",
			wantErr: true,
		},
		{
			name:    "special characters",
			hex:     "a1b2-3d4",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFingerprint(tt.hex)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFingerprint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.String() != tt.want {
				t.Errorf("NewFingerprint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFingerprint_String(t *testing.T) {
	tests := []struct {
		name        string
		fingerprint Fingerprint
		want        string
	}{
		{
			name:        "lowercase fingerprint",
			fingerprint: Fingerprint("a1b2c3d4"),
			want:        "a1b2c3d4",
		},
		{
			name:        "uppercase fingerprint",
			fingerprint: Fingerprint("A1B2C3D4"),
			want:        "A1B2C3D4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fingerprint.String(); got != tt.want {
				t.Errorf("Fingerprint.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFingerprint_Bytes(t *testing.T) {
	tests := []struct {
		name        string
		fingerprint Fingerprint
		want        []byte
		wantErr     bool
	}{
		{
			name:        "valid lowercase fingerprint",
			fingerprint: Fingerprint("a1b2c3d4"),
			want:        []byte{0xa1, 0xb2, 0xc3, 0xd4},
			wantErr:     false,
		},
		{
			name:        "valid uppercase fingerprint",
			fingerprint: Fingerprint("A1B2C3D4"),
			want:        []byte{0xa1, 0xb2, 0xc3, 0xd4},
			wantErr:     false,
		},
		{
			name:        "valid mixed case fingerprint",
			fingerprint: Fingerprint("a1B2c3D4"),
			want:        []byte{0xa1, 0xb2, 0xc3, 0xd4},
			wantErr:     false,
		},
		{
			name:        "all zeros",
			fingerprint: Fingerprint("00000000"),
			want:        []byte{0x00, 0x00, 0x00, 0x00},
			wantErr:     false,
		},
		{
			name:        "all ones",
			fingerprint: Fingerprint("ffffffff"),
			want:        []byte{0xff, 0xff, 0xff, 0xff},
			wantErr:     false,
		},
		{
			name:        "empty fingerprint",
			fingerprint: Fingerprint(""),
			wantErr:     true,
		},
		{
			name:        "too short",
			fingerprint: Fingerprint("a1b2c3"),
			wantErr:     true,
		},
		{
			name:        "invalid hex characters",
			fingerprint: Fingerprint("g1h2i3j4"),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fingerprint.Bytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("Fingerprint.Bytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("Fingerprint.Bytes() length = %d, want %d", len(got), len(tt.want))
					return
				}
				for i := range got {
					if got[i] != tt.want[i] {
						t.Errorf("Fingerprint.Bytes()[%d] = %02x, want %02x", i, got[i], tt.want[i])
					}
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
		errMsg      string
	}{
		{
			name:        "valid fingerprint lowercase",
			fingerprint: "a1b2c3d4",
			wantErr:     false,
		},
		{
			name:        "valid fingerprint uppercase",
			fingerprint: "A1B2C3D4",
			wantErr:     false,
		},
		{
			name:        "valid fingerprint mixed case",
			fingerprint: "a1B2c3D4",
			wantErr:     false,
		},
		{
			name:        "empty fingerprint",
			fingerprint: "",
			wantErr:     true,
			errMsg:      "must be exactly 8 characters",
		},
		{
			name:        "too short",
			fingerprint: "a1b2c3",
			wantErr:     true,
			errMsg:      "must be exactly 8 characters",
		},
		{
			name:        "too long",
			fingerprint: "a1b2c3d4e5",
			wantErr:     true,
			errMsg:      "must be exactly 8 characters",
		},
		{
			name:        "invalid characters",
			fingerprint: "g1h2i3j4",
			wantErr:     true,
			errMsg:      "must be 8 hexadecimal characters",
		},
		{
			name:        "special characters",
			fingerprint: "a1b2-3d4",
			wantErr:     true,
			errMsg:      "must be 8 hexadecimal characters",
		},
		{
			name:        "spaces",
			fingerprint: "a1b2 3d4",
			wantErr:     true,
			errMsg:      "must be 8 hexadecimal characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFingerprint(tt.fingerprint)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFingerprint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateFingerprint() error = %v, want error containing %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}
