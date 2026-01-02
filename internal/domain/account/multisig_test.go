package account

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestMultisigConfig_IsMultisigAccount tests the IsMultisigAccount method.
func TestMultisigConfig_IsMultisigAccount(t *testing.T) {
	tests := []struct {
		name        string
		config      *MultisigConfig
		accountType AccountType
		want        bool
	}{
		{
			name:        "nil config returns false",
			config:      nil,
			accountType: AccountTypeDeposit,
			want:        false,
		},
		{
			name: "empty config returns false",
			config: &MultisigConfig{
				AccountMap: make(map[AccountType]map[int][]AuthType),
			},
			accountType: AccountTypeDeposit,
			want:        false,
		},
		{
			name: "client account is not multisig",
			config: &MultisigConfig{
				AccountMap: map[AccountType]map[int][]AuthType{
					AccountTypeDeposit: {
						2: {AuthType1, AuthType2, AuthType3, AuthType4, AuthType5},
					},
				},
			},
			accountType: AccountTypeClient,
			want:        false,
		},
		{
			name: "deposit account is multisig",
			config: &MultisigConfig{
				AccountMap: map[AccountType]map[int][]AuthType{
					AccountTypeDeposit: {
						2: {AuthType1, AuthType2, AuthType3, AuthType4, AuthType5},
					},
				},
			},
			accountType: AccountTypeDeposit,
			want:        true,
		},
		{
			name: "payment account is multisig",
			config: &MultisigConfig{
				AccountMap: map[AccountType]map[int][]AuthType{
					AccountTypePayment: {
						3: {AuthType1, AuthType2, AuthType3, AuthType4, AuthType5},
					},
				},
			},
			accountType: AccountTypePayment,
			want:        true,
		},
		{
			name: "stored account is multisig",
			config: &MultisigConfig{
				AccountMap: map[AccountType]map[int][]AuthType{
					AccountTypeStored: {
						4: {AuthType1, AuthType2, AuthType3, AuthType4, AuthType5},
					},
				},
			},
			accountType: AccountTypeStored,
			want:        true,
		},
		{
			name: "empty account type returns false",
			config: &MultisigConfig{
				AccountMap: map[AccountType]map[int][]AuthType{
					AccountTypeDeposit: {
						2: {AuthType1, AuthType2, AuthType3, AuthType4, AuthType5},
					},
				},
			},
			accountType: "",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.config.IsMultisigAccount(tt.accountType)
			require.Equal(t, tt.want, res, "IsMultisigAccount() result mismatch")
		})
	}
}

// TestNewMultisigConfig tests the NewMultisigConfig factory function.
// This verifies that the function correctly creates MultisigConfig from domain types.
func TestNewMultisigConfig(t *testing.T) {
	tests := []struct {
		name       string
		accountMap map[AccountType]map[int][]AuthType
		wantCount  int
	}{
		{
			name:       "nil map returns empty map",
			accountMap: nil,
			wantCount:  0,
		},
		{
			name: "single multisig config",
			accountMap: map[AccountType]map[int][]AuthType{
				AccountTypeDeposit: {
					2: {AuthType1, AuthType2},
				},
			},
			wantCount: 1,
		},
		{
			name: "multiple multisig configs",
			accountMap: map[AccountType]map[int][]AuthType{
				AccountTypeDeposit: {
					2: {AuthType1, AuthType2},
				},
				AccountTypePayment: {
					3: {AuthType1, AuthType2, AuthType3},
				},
			},
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			multi := NewMultisigConfig(tt.accountMap)
			require.NotNil(t, multi, "MultisigConfig should not be nil")
			require.Equal(t, tt.wantCount, len(multi.AccountMap), "AccountMap count mismatch")
		})
	}
}

// TestMultisigConfig_MultiAccounts tests the MultiAccounts method.
func TestMultisigConfig_MultiAccounts(t *testing.T) {
	tests := []struct {
		name   string
		config *MultisigConfig
		want   map[AccountType]map[int][]AuthType
	}{
		{
			name:   "nil config returns nil",
			config: nil,
			want:   nil,
		},
		{
			name: "returns account map",
			config: &MultisigConfig{
				AccountMap: map[AccountType]map[int][]AuthType{
					AccountTypeDeposit: {
						2: {AuthType1, AuthType2, AuthType3},
					},
					AccountTypePayment: {
						3: {AuthType1, AuthType2, AuthType3, AuthType4},
					},
				},
			},
			want: map[AccountType]map[int][]AuthType{
				AccountTypeDeposit: {
					2: {AuthType1, AuthType2, AuthType3},
				},
				AccountTypePayment: {
					3: {AuthType1, AuthType2, AuthType3, AuthType4},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.config.MultiAccounts()
			require.Equal(t, tt.want, res, "MultiAccounts() result mismatch")
		})
	}
}
