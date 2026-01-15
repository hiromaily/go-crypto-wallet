package config

import (
	"os"
	"testing"

	"github.com/bookerzzz/grok"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	configutil "github.com/hiromaily/go-crypto-wallet/pkg/config/testutil"
)

// TestNewAccount tests the NewAccount function for loading account configuration from TOML files.
func TestNewAccount(t *testing.T) {
	// t.SkipNow()
	confPath := configutil.GetConfigFilePath("account.yaml")
	// Skip if config file path is empty (file not found)
	if confPath == "" {
		t.Skip("Config file not found: account.yaml")
	}
	conf, err := NewAccount(confPath)
	require.NoError(t, err, "fail to create config")
	grok.Value(conf)
}

// TestNewAccountWithViper tests account loading using viper.
// This verifies that the TOML file is properly loaded and unmarshaled into AccountRoot structure.
func TestNewAccountWithViper(t *testing.T) {
	confPath := configutil.GetConfigFilePath("account.yaml")

	// Skip if config file doesn't exist
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		t.Skipf("Config file not found: %s", confPath)
	}

	conf, err := NewAccount(confPath)
	require.NoError(t, err, "NewAccount() should not return error")
	require.NotNil(t, conf, "NewAccount() returned nil config")

	// Verify that viper properly loaded the TOML file
	assert.NotEmpty(t, conf.Types, "Account types should not be empty")
	assert.NotEmpty(t, conf.DepositReceiver, "DepositReceiver should not be empty")
	assert.NotEmpty(t, conf.PaymentSender, "PaymentSender should not be empty")
	assert.NotEmpty(t, conf.Multisigs, "Multisigs should not be empty")

	// Verify multisig structure loaded correctly
	for i, ms := range conf.Multisigs {
		assert.NotEmpty(t, ms.Type, "Multisig[%d].Type should not be empty", i)
		assert.NotZero(t, ms.Required, "Multisig[%d].Required should not be zero", i)
		assert.NotEmpty(t, ms.AuthUsers, "Multisig[%d].AuthUsers should not be empty", i)
	}
}

// TestLoadAccount tests the TOML loading functionality indirectly through NewAccount.
// This test verifies that NewAccount correctly loads and validates account configuration
// from a TOML file using the common loadTOML function.
func TestLoadAccount(t *testing.T) {
	confPath := configutil.GetConfigFilePath("account.yaml")

	// Skip if config file doesn't exist
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		t.Skipf("Config file not found: %s", confPath)
	}

	conf, err := NewAccount(confPath)
	require.NoError(t, err, "NewAccount() should not return error")
	require.NotNil(t, conf, "NewAccount() returned nil config")
}

// TestNewMultisigConfig tests the conversion from AccountMultisig config to domain MultisigConfig.
// This verifies that the conversion logic works correctly.
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
	confPath := configutil.GetConfigFilePath("account.yaml")
	// Skip if config file path is empty (file not found)
	if confPath == "" {
		t.Skip("Config file not found: account.yaml")
	}
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

// TestNewAccount_YAML tests the NewAccount function for loading account configuration from YAML files.
// This verifies that YAML configuration files are properly loaded and unmarshaled into AccountRoot structure.
func TestNewAccount_YAML(t *testing.T) {
	confPath := configutil.GetConfigFilePath("account.yaml")

	// Skip if config file doesn't exist
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		t.Skipf("Config file not found: %s", confPath)
	}

	conf, err := NewAccount(confPath)
	require.NoError(t, err, "NewAccount() should not return error")
	require.NotNil(t, conf, "NewAccount() returned nil config")

	// Verify that viper properly loaded the YAML file
	assert.NotEmpty(t, conf.Types, "Account types should not be empty")
	assert.NotEmpty(t, conf.DepositReceiver, "DepositReceiver should not be empty")
	assert.NotEmpty(t, conf.PaymentSender, "PaymentSender should not be empty")
	assert.NotEmpty(t, conf.Multisigs, "Multisigs should not be empty")

	// Verify multisig structure loaded correctly
	for i, ms := range conf.Multisigs {
		assert.NotEmpty(t, ms.Type, "Multisig[%d].Type should not be empty", i)
		assert.NotZero(t, ms.Required, "Multisig[%d].Required should not be zero", i)
		assert.NotEmpty(t, ms.AuthUsers, "Multisig[%d].AuthUsers should not be empty", i)
	}
}

// TestNewAccount_TOMLvsYAML tests that TOML and YAML configs produce identical results.
// This verifies backward compatibility and ensures both formats are equivalent.
func TestNewAccount_TOMLvsYAML(t *testing.T) {
	tomlPath := configutil.GetConfigFilePath("account.toml")
	yamlPath := configutil.GetConfigFilePath("account.yaml")

	// Skip if either config file doesn't exist
	if _, err := os.Stat(tomlPath); os.IsNotExist(err) {
		t.Skipf("Config file not found: %s", tomlPath)
	}
	if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
		t.Skipf("Config file not found: %s", yamlPath)
	}

	tomlConf, err := NewAccount(tomlPath)
	require.NoError(t, err, "NewAccount(TOML) should not return error")

	yamlConf, err := NewAccount(yamlPath)
	require.NoError(t, err, "NewAccount(YAML) should not return error")

	// Compare configurations
	assert.Equal(t, tomlConf.Types, yamlConf.Types, "Types should match")
	assert.Equal(t, tomlConf.DepositReceiver, yamlConf.DepositReceiver, "DepositReceiver should match")
	assert.Equal(t, tomlConf.PaymentSender, yamlConf.PaymentSender, "PaymentSender should match")
	assert.Equal(t, len(tomlConf.Multisigs), len(yamlConf.Multisigs), "Multisigs count should match")

	// Compare multisig configurations
	for i := range tomlConf.Multisigs {
		assert.Equal(t, tomlConf.Multisigs[i].Type, yamlConf.Multisigs[i].Type,
			"Multisig[%d].Type should match", i)
		assert.Equal(t, tomlConf.Multisigs[i].Required, yamlConf.Multisigs[i].Required,
			"Multisig[%d].Required should match", i)
		assert.Equal(t, tomlConf.Multisigs[i].AuthUsers, yamlConf.Multisigs[i].AuthUsers,
			"Multisig[%d].AuthUsers should match", i)
	}
}
