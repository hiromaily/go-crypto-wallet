package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// New create config
func New(file string, wtype wallet.WalletType, coinTypeCode coin.CoinTypeCode) (*Config, error) {
	if file == "" {
		return nil, errors.New("config file should be passed")
	}

	var err error
	conf, err := loadConfig(file)
	if err != nil {
		return nil, err
	}

	//debug
	//grok.Value(conf)

	conf.CoinTypeCode = coinTypeCode

	//validate
	if err = conf.validate(wtype, coinTypeCode); err != nil {
		return nil, err
	}

	return conf, nil
}

// loadConfig load config file
func loadConfig(path string) (*Config, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "can't read toml file. %s", path)
	}

	var config Config
	_, err = toml.Decode(string(d), &config)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call toml.Decode()")
	}

	return &config, nil
}

// validate config
func (c *Config) validate(wtype wallet.WalletType, coinTypeCode coin.CoinTypeCode) error {
	validate := validator.New()
	//if err := validate.Struct(c); err != nil {
	//	return err
	//}

	switch coinTypeCode {
	case coin.BTC, coin.BCH:
		//flds := []string{
		//	"SubSlice[0].Test",
		//	"Sub",
		//	"SubIgnore",
		//	"Anonymous.A",
		//}
		if err := validate.StructExcept(c, "Ethereum"); err != nil {
			return err
		}
		switch wtype {
		case wallet.WalletTypeWatchOnly:
			if c.Bitcoin.Block.ConfirmationNum == 0 {
				return errors.New("Block ConfirmationNum is required in toml file")
			}
		default:
		}
	case coin.ETH:
		if err := validate.StructExcept(c, "AddressType", "Bitcoin"); err != nil {
			return err
		}
	default:
	}

	return nil
}
