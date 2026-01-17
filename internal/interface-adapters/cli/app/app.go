// Package app provides shared CLI application logic for wallet commands.
// This package centralizes common initialization, configuration loading,
// and flag handling across keygen, sign, and watch wallets.
package app

import (
	"errors"
	"fmt"
	"log"

	"github.com/hiromaily/go-crypto-wallet/internal/di"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
)

// App holds the initialized wallet application state.
type App struct {
	WalletType  domainWallet.WalletType
	Config      *config.WalletRoot
	AccountConf *config.AccountRoot
	Container   di.Container
}

// Options holds CLI flag values for wallet initialization.
type Options struct {
	ConfigPath        string
	AccountConfigPath string
	CoinTypeCode      string
	BTCWallet         string
}

// NewApp creates a new App instance by loading configuration and initializing the DI container.
func NewApp(walletType domainWallet.WalletType, opts Options) (*App, error) {
	// Validate coin type
	if err := ValidateCoinType(opts.CoinTypeCode); err != nil {
		return nil, err
	}

	// Validate config path is provided
	if opts.ConfigPath == "" {
		return nil, errors.New("--config flag is required")
	}

	// Load wallet config
	conf, err := config.NewWallet(opts.ConfigPath, walletType, domainCoin.CoinTypeCode(opts.CoinTypeCode))
	if err != nil {
		return nil, fmt.Errorf("failed to load wallet config: %w", err)
	}

	// Load account config (optional)
	accountConf := &config.AccountRoot{}
	if opts.AccountConfigPath != "" {
		accountConf, err = config.NewAccount(opts.AccountConfigPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load account config: %w", err)
		}
	}

	// Override config with CLI values
	conf.CoinTypeCode = domainCoin.CoinTypeCode(opts.CoinTypeCode)

	// Handle ERC20 token configuration
	if domainCoin.IsERC20Token(opts.CoinTypeCode) {
		if err := conf.HasERC20Config(domainCoin.ERC20Token(opts.CoinTypeCode)); err != nil {
			return nil, fmt.Errorf("failed to validate ERC20 token: %w", err)
		}
		conf.Ethereum.ERC20Token = domainCoin.ERC20Token(opts.CoinTypeCode)
	}

	// Override Bitcoin host for specific wallet
	if opts.BTCWallet != "" {
		conf.Bitcoin.Host = fmt.Sprintf("%s/wallet/%s", conf.Bitcoin.Host, opts.BTCWallet)
		log.Println("conf.Bitcoin.Host:", conf.Bitcoin.Host)
	}

	// Create DI container
	container := di.NewContainer(conf, accountConf, walletType)

	return &App{
		WalletType:  walletType,
		Config:      conf,
		AccountConf: accountConf,
		Container:   container,
	}, nil
}
