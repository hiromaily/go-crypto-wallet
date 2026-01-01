package account

import (
	"testing"

	"github.com/stretchr/testify/require"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	configutil "github.com/hiromaily/go-crypto-wallet/pkg/config/testutil"
)

// TestNewMultisigConfig tests the conversion from AccountMultisig config to domain MultisigConfig.
// This is an integration test that verifies the conversion logic works correctly.
func TestNewMultisigConfig(t *testing.T) {
	tests := []struct {
		name         string
		confMultisig []AccountMultisig
		wantCount    int
	}{
		{
			name:         "nil config returns empty map",
			confMultisig: nil,
			wantCount:    0,
		},
		{
			name: "single multisig config",
			confMultisig: []AccountMultisig{
				{
					Type:     domainAccount.AccountTypeDeposit,
					Required: 2,
					AuthUsers: []domainAccount.AuthType{
						domainAccount.AuthType1,
						domainAccount.AuthType2,
					},
				},
			},
			wantCount: 1,
		},
		{
			name: "multiple multisig configs",
			confMultisig: []AccountMultisig{
				{
					Type:     domainAccount.AccountTypeDeposit,
					Required: 2,
					AuthUsers: []domainAccount.AuthType{
						domainAccount.AuthType1,
						domainAccount.AuthType2,
					},
				},
				{
					Type:     domainAccount.AccountTypePayment,
					Required: 3,
					AuthUsers: []domainAccount.AuthType{
						domainAccount.AuthType1,
						domainAccount.AuthType2,
						domainAccount.AuthType3,
					},
				},
			},
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			multi := NewMultisigConfig(tt.confMultisig)
			require.NotNil(t, multi, "MultisigConfig should not be nil")
			require.Equal(t, tt.wantCount, len(multi.AccountMap), "AccountMap count mismatch")
		})
	}
}

// TestNewMultisigConfig_Integration tests NewMultisigConfig with real config file.
// This verifies that the conversion from TOML config to domain entity works correctly.
func TestNewMultisigConfig_Integration(t *testing.T) {
	// config
	confPath := configutil.GetConfigFilePath("account.toml")
	conf, err := NewAccount(confPath)
	require.NoError(t, err, "fail to create config")

	multi := NewMultisigConfig(conf.Multisigs)
	require.NotNil(t, multi, "MultisigConfig should not be nil")

	// Verify that deposit, payment, and stored accounts are multisig
	require.True(t, multi.IsMultisigAccount(domainAccount.AccountTypeDeposit), "deposit should be multisig")
	require.True(t, multi.IsMultisigAccount(domainAccount.AccountTypePayment), "payment should be multisig")
	require.True(t, multi.IsMultisigAccount(domainAccount.AccountTypeStored), "stored should be multisig")

	// Verify that client account is not multisig
	require.False(t, multi.IsMultisigAccount(domainAccount.AccountTypeClient), "client should not be multisig")
}
