package config

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
)

// NewWallet creates wallet config
func NewWallet(file string, wtype domainWallet.WalletType, coinTypeCode domainCoin.CoinTypeCode) (*WalletRoot, error) {
	if file == "" {
		return nil, errors.New("config file should be passed")
	}

	var err error
	conf, err := loadWallet(file)
	if err != nil {
		return nil, err
	}

	// debug
	// debug.Debug(conf)

	// validate
	if err = conf.validate(wtype, coinTypeCode); err != nil {
		return nil, err
	}

	return conf, nil
}

// loadWallet load config file
func loadWallet(path string) (*WalletRoot, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("toml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("can't read config file. %s: %w", path, err)
	}

	var config WalletRoot
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("fail to unmarshal config: %w", err)
	}

	return &config, nil
}

// validate config
func (c *WalletRoot) validate(wtype domainWallet.WalletType, coinTypeCode domainCoin.CoinTypeCode) error {
	validate := validator.New()

	switch coinTypeCode {
	case domainCoin.BTC, domainCoin.BCH:
		if err := validate.StructExcept(c, "Ethereum", "Ripple"); err != nil {
			return err
		}
		switch wtype {
		case domainWallet.WalletTypeWatchOnly:
			if c.Bitcoin.Block.ConfirmationNum == 0 {
				return errors.New("block ConfirmationNum is required in toml file")
			}
		case domainWallet.WalletTypeKeyGen, domainWallet.WalletTypeSign:
			// No additional validation needed
		default:
		}
	case domainCoin.ETH, domainCoin.ERC20:
		if err := validate.StructExcept(c, "AddressType", "Bitcoin", "Ripple"); err != nil {
			return err
		}
	case domainCoin.XRP:
		if err := validate.StructExcept(c, "AddressType", "Bitcoin", "Ethereum"); err != nil {
			return err
		}
	case domainCoin.LTC, domainCoin.HYC:
		// Not implemented yet
	default:
	}

	return nil
}

func (c *WalletRoot) ValidateERC20(token domainCoin.ERC20Token) error {
	if _, ok := c.Ethereum.ERC20s[token]; !ok {
		return fmt.Errorf("erc20 token information for [%s] is required", token.String())
	}
	return nil
}
