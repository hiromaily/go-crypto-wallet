package account

import (
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
)

// AccountRoot account root config
type AccountRoot struct {
	Types           []AccountType     `toml:"types"`
	DepositReceiver AccountType       `toml:"deposit_receiver"`
	PaymentSender   AccountType       `toml:"payment_sender"`
	Multisigs       []AccountMultisig `toml:"multisig"`
}

// AccountMultisig multisig setting
type AccountMultisig struct {
	Type      AccountType `toml:"type"`
	Required  int         `toml:"required"`
	AuthUsers []AuthType  `toml:"auth_users"`
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
	d, err := os.ReadFile(path) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("can't read toml file. %s: %w", path, err)
	}

	var config AccountRoot
	_, err = toml.Decode(string(d), &config)
	if err != nil {
		return nil, fmt.Errorf("fail to call toml.Decode(): %w", err)
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
