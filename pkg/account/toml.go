package account

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
)

// AccountRoot account root config
type AccountRoot struct {
	Types           []domainAccount.AccountType `toml:"types" mapstructure:"types"`
	DepositReceiver domainAccount.AccountType   `toml:"deposit_receiver" mapstructure:"deposit_receiver"`
	PaymentSender   domainAccount.AccountType   `toml:"payment_sender" mapstructure:"payment_sender"`
	Multisigs       []AccountMultisig           `toml:"multisig" mapstructure:"multisig"`
}

// AccountMultisig multisig setting
type AccountMultisig struct {
	Type      domainAccount.AccountType `toml:"type" mapstructure:"type"`
	Required  int                       `toml:"required" mapstructure:"required"`
	AuthUsers []domainAccount.AuthType  `toml:"auth_users" mapstructure:"auth_users"`
}

// NewAccount creates account config
func NewAccount(file string) (*AccountRoot, error) {
	if file == "" {
		return nil, errors.New("config file should be passed")
	}

	var err error
	conf, err := loadAccount(file)
	if err != nil {
		return nil, err
	}

	// debug
	// grok.Value(conf)

	// validate
	if err = conf.validate(); err != nil {
		return nil, err
	}

	return conf, nil
}

// loadAccount load account config file
func loadAccount(path string) (*AccountRoot, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("toml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("can't read config file. %s: %w", path, err)
	}

	var config AccountRoot
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("fail to unmarshal config: %w", err)
	}

	return &config, nil
}

// validate config
func (c *AccountRoot) validate() error {
	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return err
	}

	return nil
}
