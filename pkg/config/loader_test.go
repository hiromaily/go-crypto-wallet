package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoadConfig tests the loadConfig function with various file formats.
// This verifies that both TOML and YAML files can be loaded automatically.
func TestLoadConfig(t *testing.T) {
	// Create temporary test config files
	tmpDir := t.TempDir()

	// Create TOML test file
	tomlContent := `
types = ["client", "deposit"]
deposit_receiver = "deposit"
payment_sender = "payment"
`
	tomlPath := filepath.Join(tmpDir, "test.toml")
	err := os.WriteFile(tomlPath, []byte(tomlContent), 0o600)
	require.NoError(t, err, "Failed to create test TOML file")

	// Create YAML test file
	yamlContent := `
types:
  - client
  - deposit
deposit_receiver: deposit
payment_sender: payment
`
	yamlPath := filepath.Join(tmpDir, "test.yaml")
	err = os.WriteFile(yamlPath, []byte(yamlContent), 0o600)
	require.NoError(t, err, "Failed to create test YAML file")

	// Create YML test file
	ymlPath := filepath.Join(tmpDir, "test.yml")
	err = os.WriteFile(ymlPath, []byte(yamlContent), 0o600)
	require.NoError(t, err, "Failed to create test YML file")

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "Load TOML file",
			path:    tomlPath,
			wantErr: false,
		},
		{
			name:    "Load YAML file",
			path:    yamlPath,
			wantErr: false,
		},
		{
			name:    "Load YML file",
			path:    ymlPath,
			wantErr: false,
		},
		{
			name:    "Load non-existent file",
			path:    filepath.Join(tmpDir, "nonexistent.toml"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var conf AccountRoot
			err := loadConfig(tt.path, &conf)

			if tt.wantErr {
				require.Error(t, err, "loadConfig() should return error")
				return
			}

			require.NoError(t, err, "loadConfig() should not return error")
			assert.NotEmpty(t, conf.Types, "Types should not be empty")
			assert.Equal(t, "deposit", string(conf.DepositReceiver), "DepositReceiver mismatch")
			assert.Equal(t, "payment", string(conf.PaymentSender), "PaymentSender mismatch")
		})
	}
}

// TestLoadConfigWithUnsupportedExtension tests loadConfig with unsupported file extensions.
func TestLoadConfigWithUnsupportedExtension(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name string
		path string
	}{
		{
			name: "JSON extension",
			path: filepath.Join(tmpDir, "test.json"),
		},
		{
			name: "XML extension",
			path: filepath.Join(tmpDir, "test.xml"),
		},
		{
			name: "No extension",
			path: filepath.Join(tmpDir, "test"),
		},
		{
			name: "Uppercase unsupported extension",
			path: filepath.Join(tmpDir, "test.JSON"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var conf AccountRoot
			err := loadConfig(tt.path, &conf)

			require.Error(t, err, "loadConfig() should return error for unsupported extension")
			assert.Contains(t, err.Error(), "unsupported config file extension",
				"Error message should mention unsupported extension")
		})
	}
}

// TestLoadConfigWithEnvOverride tests that environment variables override config file values.
func TestLoadConfigWithEnvOverride(t *testing.T) {
	tmpDir := t.TempDir()

	// Create YAML test file with initial values
	yamlContent := `
deposit_receiver: deposit
payment_sender: payment
`
	yamlPath := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(yamlPath, []byte(yamlContent), 0o600)
	require.NoError(t, err, "Failed to create test YAML file")

	tests := []struct {
		name             string
		envKey           string
		envValue         string
		expectedReceiver string
		expectedSender   string
	}{
		{
			name:             "Override deposit_receiver with env var",
			envKey:           "WALLET_DEPOSIT_RECEIVER",
			envValue:         "custom_deposit",
			expectedReceiver: "custom_deposit",
			expectedSender:   "payment",
		},
		{
			name:             "Override payment_sender with env var",
			envKey:           "WALLET_PAYMENT_SENDER",
			envValue:         "custom_payment",
			expectedReceiver: "deposit",
			expectedSender:   "custom_payment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			t.Setenv(tt.envKey, tt.envValue)

			var conf AccountRoot
			err := loadConfig(yamlPath, &conf)

			require.NoError(t, err, "loadConfig() should not return error")
			assert.Equal(t, tt.expectedReceiver, string(conf.DepositReceiver),
				"DepositReceiver should match expected value")
			assert.Equal(t, tt.expectedSender, string(conf.PaymentSender),
				"PaymentSender should match expected value")
		})
	}
}

// TestLoadConfigEnvPrefixAndNesting tests environment variable prefix and nested key handling.
func TestLoadConfigEnvPrefixAndNesting(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a simple YAML config for testing nested override
	yamlContent := `
deposit_receiver: original_deposit
payment_sender: original_payment
`
	yamlPath := filepath.Join(tmpDir, "nested_test.yaml")
	err := os.WriteFile(yamlPath, []byte(yamlContent), 0o600)
	require.NoError(t, err, "Failed to create test YAML file")

	t.Run("Environment variable without prefix should not override", func(t *testing.T) {
		// Set env var without WALLET_ prefix
		t.Setenv("DEPOSIT_RECEIVER", "should_not_work")

		var conf AccountRoot
		err := loadConfig(yamlPath, &conf)

		require.NoError(t, err)
		assert.Equal(t, "original_deposit", string(conf.DepositReceiver),
			"Value should not be overridden by env var without prefix")
	})

	t.Run("Environment variable with WALLET_ prefix should override", func(t *testing.T) {
		// Set env var with WALLET_ prefix
		t.Setenv("WALLET_DEPOSIT_RECEIVER", "overridden_deposit")

		var conf AccountRoot
		err := loadConfig(yamlPath, &conf)

		require.NoError(t, err)
		assert.Equal(t, "overridden_deposit", string(conf.DepositReceiver),
			"Value should be overridden by env var with WALLET_ prefix")
	})
}

// TestLoadConfigNestedKeyOverride tests that nested configuration keys can be overridden
// using environment variables with underscore as separator (e.g., bitcoin.host -> WALLET_BITCOIN_HOST).
func TestLoadConfigNestedKeyOverride(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a YAML config with nested structure matching WalletRoot
	yamlContent := `
address_type: "legacy"
bitcoin:
  host: "127.0.0.1:18332"
  user: "testuser"
  pass: "testpass"
  http_post_mode: true
  disable_tls: true
  network_type: "regtest"
database:
  type: "mysql"
  mysql:
    host: "localhost:3306"
    dbname: "testdb"
    user: "dbuser"
    pass: "dbpass"
  sqlite:
    path: "./data/sqlite/btc/test.db"
    debug: false
logger:
  service: "test-service"
  format: "json"
  level: "debug"
file_path:
  tx: "./data/tx/"
  address: "./data/address/"
  full_pubkey: "./data/fullpubkey/"
`
	yamlPath := filepath.Join(tmpDir, "wallet_test.yaml")
	err := os.WriteFile(yamlPath, []byte(yamlContent), 0o600)
	require.NoError(t, err, "Failed to create test YAML file")

	t.Run("Override nested bitcoin.host with WALLET_BITCOIN_HOST", func(t *testing.T) {
		// Set env var for nested key: bitcoin.host -> WALLET_BITCOIN_HOST
		t.Setenv("WALLET_BITCOIN_HOST", "overridden-host:19999")

		var conf WalletRoot
		err := loadConfig(yamlPath, &conf)

		require.NoError(t, err, "loadConfig() should not return error")
		assert.Equal(t, "overridden-host:19999", conf.Bitcoin.Host,
			"Nested key bitcoin.host should be overridden by WALLET_BITCOIN_HOST")
		// Verify other nested values are not affected
		assert.Equal(t, "testuser", conf.Bitcoin.User,
			"Other nested values should remain unchanged")
	})

	t.Run("Override nested database.mysql.host with WALLET_DATABASE_MYSQL_HOST", func(t *testing.T) {
		// Set env var for nested key: database.mysql.host -> WALLET_DATABASE_MYSQL_HOST
		t.Setenv("WALLET_DATABASE_MYSQL_HOST", "mysql-server:3307")

		var conf WalletRoot
		err := loadConfig(yamlPath, &conf)

		require.NoError(t, err, "loadConfig() should not return error")
		assert.Equal(t, "mysql-server:3307", conf.Database.MySQL.Host,
			"Nested key database.mysql.host should be overridden by WALLET_DATABASE_MYSQL_HOST")
		// Verify other nested values are not affected
		assert.Equal(t, "testdb", conf.Database.MySQL.DB,
			"Other nested values should remain unchanged")
	})

	t.Run("Override top-level address_type with WALLET_ADDRESS_TYPE", func(t *testing.T) {
		// Set env var for top-level key
		t.Setenv("WALLET_ADDRESS_TYPE", "bech32")

		var conf WalletRoot
		err := loadConfig(yamlPath, &conf)

		require.NoError(t, err, "loadConfig() should not return error")
		assert.Equal(t, "bech32", string(conf.AddressType),
			"Top-level key address_type should be overridden by WALLET_ADDRESS_TYPE")
	})

	t.Run("Override deeply nested logger.level with WALLET_LOGGER_LEVEL", func(t *testing.T) {
		// Set env var for nested key: logger.level -> WALLET_LOGGER_LEVEL
		t.Setenv("WALLET_LOGGER_LEVEL", "error")

		var conf WalletRoot
		err := loadConfig(yamlPath, &conf)

		require.NoError(t, err, "loadConfig() should not return error")
		assert.Equal(t, "error", conf.Logger.Level,
			"Nested key logger.level should be overridden by WALLET_LOGGER_LEVEL")
	})
}

// TestLoadConfigCaseInsensitive tests that file extension detection is case-insensitive.
func TestLoadConfigCaseInsensitive(t *testing.T) {
	tmpDir := t.TempDir()

	tomlContent := `
types = ["client", "deposit"]
deposit_receiver = "deposit"
payment_sender = "payment"
`

	yamlContent := `
types:
  - client
  - deposit
deposit_receiver: deposit
payment_sender: payment
`

	tests := []struct {
		name     string
		filename string
		content  string
	}{
		{name: "Uppercase TOML", filename: "test.TOML", content: tomlContent},
		{name: "Mixed case TOML", filename: "test.Toml", content: tomlContent},
		{name: "Uppercase YAML", filename: "test.YAML", content: yamlContent},
		{name: "Mixed case YAML", filename: "test.Yaml", content: yamlContent},
		{name: "Uppercase YML", filename: "test.YML", content: yamlContent},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tt.filename)
			err := os.WriteFile(path, []byte(tt.content), 0o600)
			require.NoError(t, err, "Failed to create test file")

			var conf AccountRoot
			err = loadConfig(path, &conf)

			require.NoError(t, err, "loadConfig() should handle case-insensitive extensions")
			assert.NotEmpty(t, conf.Types, "Types should not be empty")
		})
	}
}
