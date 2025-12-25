package account

import (
	"os"
	"testing"

	"github.com/bookerzzz/grok"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	configutil "github.com/hiromaily/go-crypto-wallet/pkg/config/testutil"
)

// TestNewAccount is test for NewAccount
func TestNewAccount(t *testing.T) {
	// t.SkipNow()
	confPath := configutil.GetConfigFilePath("account.toml")
	conf, err := NewAccount(confPath)
	require.NoError(t, err, "fail to create config")
	grok.Value(conf)
}

// TestNewAccountWithViper tests account loading using viper
func TestNewAccountWithViper(t *testing.T) {
	confPath := configutil.GetConfigFilePath("account.toml")

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

// TestLoadAccount tests the loadAccount function
func TestLoadAccount(t *testing.T) {
	confPath := configutil.GetConfigFilePath("account.toml")

	// Skip if config file doesn't exist
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		t.Skipf("Config file not found: %s", confPath)
	}

	conf, err := loadAccount(confPath)
	require.NoError(t, err, "loadAccount() should not return error")
	require.NotNil(t, conf, "loadAccount() returned nil config")
}
