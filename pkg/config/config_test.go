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
			configFile:   filepath.Join(projPath, "data/config/btc_watch.toml"),
			walletType:   domainWallet.WalletTypeWatchOnly,
			coinTypeCode: domainCoin.BTC,
			wantErr:      false,
		},
		{
			name:         "BTC Keygen Wallet",
			configFile:   filepath.Join(projPath, "data/config/btc_keygen.toml"),
			walletType:   domainWallet.WalletTypeKeyGen,
			coinTypeCode: domainCoin.BTC,
			wantErr:      false,
		},
		{
			name:         "BTC Sign Wallet",
			configFile:   filepath.Join(projPath, "data/config/btc_sign.toml"),
			walletType:   domainWallet.WalletTypeSign,
			coinTypeCode: domainCoin.BTC,
			wantErr:      false,
		},
		{
			name:         "ETH Watch Wallet",
			configFile:   filepath.Join(projPath, "data/config/eth_watch.toml"),
			walletType:   domainWallet.WalletTypeWatchOnly,
			coinTypeCode: domainCoin.ETH,
			wantErr:      false,
		},
		{
			name:         "ETH Keygen Wallet",
			configFile:   filepath.Join(projPath, "data/config/eth_keygen.toml"),
			walletType:   domainWallet.WalletTypeKeyGen,
			coinTypeCode: domainCoin.ETH,
			wantErr:      false,
		},
		{
			name:         "ETH Sign Wallet",
			configFile:   filepath.Join(projPath, "data/config/eth_sign.toml"),
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
	assert.NotEmpty(t, conf.MySQL.Host, "MySQL.Host should not be empty")
}

func TestLoadWallet(t *testing.T) {
	// Get project path
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		t.Skip("GOPATH not set, skipping integration test")
	}

	projPath := filepath.Join(gopath, "src/github.com/hiromaily/go-crypto-wallet")
	configPath := filepath.Join(projPath, "data/config/btc_watch.toml")

	// Skip if config file doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Skipf("Config file not found: %s", configPath)
	}

	conf, err := loadWallet(configPath)
	require.NoError(t, err, "loadWallet() should not return error")
	require.NotNil(t, conf, "loadWallet() returned nil config")

	// Verify that viper properly loaded the TOML file
	assert.NotEmpty(t, conf.Bitcoin.Host, "Bitcoin.Host should not be empty")

	// Verify nested structures are loaded correctly
	assert.True(t, conf.Bitcoin.Fee.AdjustmentMin != 0 || conf.Bitcoin.Fee.AdjustmentMax != 0,
		"Bitcoin fee settings should be loaded")

	// Verify map structures (if any ERC20s configured)
	if len(conf.Ethereum.ERC20s) > 0 {
		for token, erc20 := range conf.Ethereum.ERC20s {
			assert.NotEmpty(t, erc20.Symbol, "ERC20 token %v should have symbol", token)
		}
	}
}
