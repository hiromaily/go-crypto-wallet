package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// NewWallet creates wallet config
func NewWallet(file string, wtype wallet.WalletType, coinTypeCode coin.CoinTypeCode) (*WalletRoot, error) {
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
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "can't read toml file. %s", path)
	}

	var config WalletRoot
	_, err = toml.Decode(string(d), &config)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call toml.Decode()")
	}

	return &config, nil
}

// validate config
func (c *WalletRoot) validate(wtype wallet.WalletType, coinTypeCode coin.CoinTypeCode) error {
	validate := validator.New()

	switch coinTypeCode {
	case coin.BTC, coin.BCH:
		if err := validate.StructExcept(c, "Ethereum", "Ripple"); err != nil {
			return err
		}
		switch wtype {
		case wallet.WalletTypeWatchOnly:
			if c.Bitcoin.Block.ConfirmationNum == 0 {
				return errors.New("Block ConfirmationNum is required in toml file")
			}
		default:
		}
	case coin.ETH, coin.ERC20:
		if err := validate.StructExcept(c, "AddressType", "Bitcoin", "Ripple"); err != nil {
			return err
		}
	case coin.XRP:
		if err := validate.StructExcept(c, "AddressType", "Bitcoin", "Ethereum"); err != nil {
			return err
		}
	default:
	}

	return nil
}

func (c *WalletRoot) ValidateERC20(token coin.ERC20Token) error {
	if _, ok := c.Ethereum.ERC20s[token]; !ok {
		errors.Errorf("erc20 token information for [%s] is required", token.String())
	}
	return nil
}
