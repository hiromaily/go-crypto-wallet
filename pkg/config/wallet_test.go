package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
)

func TestNewWallet(t *testing.T) {
	// Get project path
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		t.Skip("GOPATH not set, skipping integration test")
	}

	projPath := filepath.Join(gopath, "src/github.com/hiromaily/go-crypto-wallet")

	tests := []struct {
		name         string
		configFile   string
		walletType   domainWallet.WalletType
		coinTypeCode domainCoin.CoinTypeCode
		wantErr      bool
	}{
		{
			name:         "BTC Watch Wallet",
			configFile:   filepath.Join(projPath, "config/wallet/btc/watch.yaml"),
			walletType:   domainWallet.WalletTypeWatchOnly,
			coinTypeCode: domainCoin.BTC,
			wantErr:      false,
		},
		{
			name:         "BTC Keygen Wallet",
			configFile:   filepath.Join(projPath, "config/wallet/btc/keygen.yaml"),
			walletType:   domainWallet.WalletTypeKeyGen,
			coinTypeCode: domainCoin.BTC,
			wantErr:      false,
		},
		{
			name:         "BTC Sign Wallet",
			configFile:   filepath.Join(projPath, "config/wallet/btc/sign1.yaml"),
			walletType:   domainWallet.WalletTypeSign,
			coinTypeCode: domainCoin.BTC,
			wantErr:      false,
		},
		{
			name:         "ETH Watch Wallet",
			configFile:   filepath.Join(projPath, "config/wallet/eth/watch.yaml"),
			walletType:   domainWallet.WalletTypeWatchOnly,
			coinTypeCode: domainCoin.ETH,
			wantErr:      false,
		},
		{
			name:         "ETH Keygen Wallet",
			configFile:   filepath.Join(projPath, "config/wallet/eth/keygen.yaml"),
			walletType:   domainWallet.WalletTypeKeyGen,
			coinTypeCode: domainCoin.ETH,
			wantErr:      false,
		},
		{
			name:         "ETH Sign Wallet",
			configFile:   filepath.Join(projPath, "config/wallet/eth/sign.yaml"),
			walletType:   domainWallet.WalletTypeSign,
			coinTypeCode: domainCoin.ETH,
			wantErr:      false,
		},
		{
			name:         "Empty file path",
			configFile:   "",
			walletType:   domainWallet.WalletTypeWatchOnly,
			coinTypeCode: domainCoin.BTC,
			wantErr:      true,
		},
		{
			name:         "Non-existent file",
			configFile:   "/nonexistent/path/config.toml",
			walletType:   domainWallet.WalletTypeWatchOnly,
			coinTypeCode: domainCoin.BTC,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testNewWalletCase(t, tt.configFile, tt.walletType, tt.coinTypeCode, tt.wantErr)
		})
	}
}

//nolint:lll
func testNewWalletCase(
	t *testing.T,
	configFile string,
	walletType domainWallet.WalletType,
	coinTypeCode domainCoin.CoinTypeCode,
	wantErr bool,
) {
	t.Helper()
	// Skip tests if config file doesn't exist (except for error cases)
	if !wantErr && configFile != "" {
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			t.Skipf("Config file not found: %s", configFile)
		}
	}

	conf, err := NewWallet(configFile, walletType, coinTypeCode)
	if wantErr {
		require.Error(t, err, "NewWallet() should return error")
		return
	}

	require.NoError(t, err, "NewWallet() should not return error")
	require.NotNil(t, conf, "NewWallet() returned nil config without error")
	validateConfig(t, conf, coinTypeCode)
}

func validateConfig(t *testing.T, conf *WalletRoot, coinTypeCode domainCoin.CoinTypeCode) {
	t.Helper()
	// Verify that nested structures are properly loaded
	switch coinTypeCode {
	case domainCoin.BTC, domainCoin.BCH:
		validateBitcoinConfig(t, conf)
	case domainCoin.ETH, domainCoin.ERC20:
		validateEthereumConfig(t, conf)
	case domainCoin.XRP:
		validateRippleConfig(t, conf)
	case domainCoin.LTC, domainCoin.HYT:
		// Not implemented yet
	default:
		// Other coins
	}

	// Verify common fields
	validateCommonConfig(t, conf)
}

func validateBitcoinConfig(t *testing.T, conf *WalletRoot) {
	t.Helper()
	assert.NotEmpty(t, conf.Bitcoin.Host, "Bitcoin.Host should not be empty")
	assert.NotEmpty(t, conf.Bitcoin.User, "Bitcoin.User should not be empty")
}

func validateEthereumConfig(t *testing.T, conf *WalletRoot) {
	t.Helper()
	assert.NotEmpty(t, conf.Ethereum.Host, "Ethereum.Host should not be empty")
}

func validateRippleConfig(t *testing.T, conf *WalletRoot) {
	t.Helper()
	assert.NotEmpty(t, conf.Ripple.WebsocketPublicURL, "Ripple.WebsocketPublicURL should not be empty")
}

func validateCommonConfig(t *testing.T, conf *WalletRoot) {
	t.Helper()
	assert.NotEmpty(t, conf.Logger.Service, "Logger.Service should not be empty")
	assert.NotEmpty(t, conf.Database.MySQL.Host, "Database.MySQL.Host should not be empty")
}

// TestNewWallet_YAML tests the NewWallet function with YAML configuration files.
// This verifies that YAML configuration files are properly loaded and unmarshaled into WalletRoot structure.
func TestNewWallet_YAML(t *testing.T) {
	// Create temporary YAML config file
	tmpDir := t.TempDir()
	yamlContent := `
key_type: bip44
address_type: legacy
coin_type: btc

bitcoin:
  host: localhost:8332
  user: user
  pass: pass
  http_post_mode: true
  disable_tls: true
  network_type: testnet3
  block:
    confirmation_num: 6
  fee:
    adjustment_min: 0.8
    adjustment_max: 1.2

ethereum:
  host: localhost
  port: 8545
  disable_tls: true
  network_type: mainnet
  confirmation_num: 12

ripple:
  websocket_public_url: ws://localhost:6006
  websocket_admin_url: ws://localhost:6006
  network_type: testnet
  api:
    url: http://localhost:5005
    is_secure: false

logger:
  service: test-wallet
  format: json
  env: dev
  level: debug
  is_logger: true

tracer:
  type: none

database:
  type: mysql
  mysql:
    host: localhost:3306
    dbname: test_db
    user: root
    pass: root
    debug: false
  sqlite:
    path: ./data/sqlite/btc/test.db
    debug: false

file_path:
  tx: /tmp/tx
  address: /tmp/address
  full_pubkey: /tmp/pubkey
`
	yamlPath := filepath.Join(tmpDir, "test_wallet.yaml")
	err := os.WriteFile(yamlPath, []byte(yamlContent), 0o600)
	require.NoError(t, err, "Failed to create test YAML file")

	// Test loading YAML config
	conf, err := NewWallet(yamlPath, domainWallet.WalletTypeWatchOnly, domainCoin.BTC)
	require.NoError(t, err, "NewWallet() should not return error")
	require.NotNil(t, conf, "NewWallet() returned nil config")

	// Verify that viper properly loaded the YAML file
	assert.NotEmpty(t, conf.Bitcoin.Host, "Bitcoin.Host should not be empty")
	assert.NotEmpty(t, conf.Logger.Service, "Logger.Service should not be empty")
	assert.NotEmpty(t, conf.Database.MySQL.Host, "Database.MySQL.Host should not be empty")

	// Verify nested structures are loaded correctly
	assert.NotZero(t, conf.Bitcoin.Block.ConfirmationNum, "Bitcoin.Block.ConfirmationNum should not be zero")
	assert.True(t, conf.Bitcoin.Fee.AdjustmentMin != 0 || conf.Bitcoin.Fee.AdjustmentMax != 0,
		"Bitcoin fee settings should be loaded")
}

// TestNewWallet_UnsupportedExtension tests that NewWallet returns an error for unsupported file extensions.
func TestNewWallet_UnsupportedExtension(t *testing.T) {
	tmpDir := t.TempDir()
	invalidPath := filepath.Join(tmpDir, "config.json")

	// Create a dummy JSON file
	err := os.WriteFile(invalidPath, []byte(`{"test": "data"}`), 0o600)
	require.NoError(t, err, "Failed to create test file")

	conf, err := NewWallet(invalidPath, domainWallet.WalletTypeWatchOnly, domainCoin.BTC)
	require.Error(t, err, "NewWallet() should return error for unsupported extension")
	assert.Nil(t, conf, "Config should be nil on error")
	assert.Contains(t, err.Error(), "unsupported config file extension", "Error should mention unsupported extension")
}

// TestNewWallet_PostgreSQL tests the NewWallet function with PostgreSQL database configuration.
// This verifies that PostgreSQL configuration is properly validated.
func TestNewWallet_PostgreSQL(t *testing.T) {
	// Create temporary YAML config file with PostgreSQL
	tmpDir := t.TempDir()
	yamlContent := `
key_type: bip44
address_type: legacy
coin_type: btc

bitcoin:
  host: localhost:8332
  user: user
  pass: pass
  http_post_mode: true
  disable_tls: true
  network_type: testnet3
  block:
    confirmation_num: 6
  fee:
    adjustment_min: 0.8
    adjustment_max: 1.2

logger:
  service: test-wallet
  format: json
  level: debug
  is_logger: true

tracer:
  type: none

database:
  type: postgresql
  postgresql:
    host: localhost
    port: 5432
    dbname: test_db
    user: postgres
    pass: postgres
    sslmode: disable
    debug: false

file_path:
  tx: /tmp/tx
  address: /tmp/address
  full_pubkey: /tmp/pubkey
`
	yamlPath := filepath.Join(tmpDir, "test_postgresql_wallet.yaml")
	err := os.WriteFile(yamlPath, []byte(yamlContent), 0o600)
	require.NoError(t, err, "Failed to create test YAML file")

	// Test loading PostgreSQL config
	conf, err := NewWallet(yamlPath, domainWallet.WalletTypeWatchOnly, domainCoin.BTC)
	require.NoError(t, err, "NewWallet() should not return error")
	require.NotNil(t, conf, "NewWallet() returned nil config")

	// Verify PostgreSQL configuration
	assert.Equal(t, "postgresql", conf.Database.Type)
	assert.Equal(t, "localhost", conf.Database.PostgreSQL.Host)
	assert.Equal(t, 5432, conf.Database.PostgreSQL.Port)
	assert.Equal(t, "test_db", conf.Database.PostgreSQL.DB)
	assert.Equal(t, "postgres", conf.Database.PostgreSQL.User)
	assert.Equal(t, "postgres", conf.Database.PostgreSQL.Pass)
	assert.Equal(t, "disable", conf.Database.PostgreSQL.SSLMode)
	assert.False(t, conf.Database.PostgreSQL.Debug)
}

// TestValidateDatabase_PostgreSQL tests PostgreSQL database validation rules.
func TestValidateDatabase_PostgreSQL(t *testing.T) {
	tests := []struct {
		name    string
		conf    WalletRoot
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid PostgreSQL configuration with all fields",
			conf: WalletRoot{
				Database: Database{
					Type: "postgresql",
					PostgreSQL: PostgreSQL{
						Host:    "localhost",
						Port:    5432,
						DB:      "test_db",
						User:    "postgres",
						Pass:    "postgres",
						SSLMode: "disable",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Valid PostgreSQL configuration without optional fields",
			conf: WalletRoot{
				Database: Database{
					Type: "postgresql",
					PostgreSQL: PostgreSQL{
						Host: "localhost",
						// Port defaults to 0 (will use 5432 in connection factory)
						DB:   "test_db",
						User: "postgres",
						Pass: "postgres",
						// SSLMode defaults to "" (will use "prefer" in connection factory)
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Missing PostgreSQL host",
			conf: WalletRoot{
				Database: Database{
					Type: "postgresql",
					PostgreSQL: PostgreSQL{
						DB:   "test_db",
						User: "postgres",
						Pass: "postgres",
					},
				},
			},
			wantErr: true,
			errMsg:  "database.postgresql.host is required when database.type is postgresql",
		},
		{
			name: "Missing PostgreSQL dbname",
			conf: WalletRoot{
				Database: Database{
					Type: "postgresql",
					PostgreSQL: PostgreSQL{
						Host: "localhost",
						User: "postgres",
						Pass: "postgres",
					},
				},
			},
			wantErr: true,
			errMsg:  "database.postgresql.dbname is required when database.type is postgresql",
		},
		{
			name: "Missing PostgreSQL user",
			conf: WalletRoot{
				Database: Database{
					Type: "postgresql",
					PostgreSQL: PostgreSQL{
						Host: "localhost",
						DB:   "test_db",
						Pass: "postgres",
					},
				},
			},
			wantErr: true,
			errMsg:  "database.postgresql.user is required when database.type is postgresql",
		},
		{
			name: "Missing PostgreSQL password",
			conf: WalletRoot{
				Database: Database{
					Type: "postgresql",
					PostgreSQL: PostgreSQL{
						Host: "localhost",
						DB:   "test_db",
						User: "postgres",
					},
				},
			},
			wantErr: true,
			errMsg:  "database.postgresql.pass is required when database.type is postgresql",
		},
		{
			name: "Unsupported database type",
			conf: WalletRoot{
				Database: Database{
					Type: "oracle",
				},
			},
			wantErr: true,
			errMsg:  "unsupported database type: oracle (must be mysql, sqlite, or postgresql)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.conf.validateDatabase()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
