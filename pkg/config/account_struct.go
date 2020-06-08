package config

import (
	"github.com/hiromaily/go-crypto-wallet/pkg/account"
)

// AccountRoot account root config
type AccountRoot struct {
	Types     []account.AccountType `toml:"types"`
	Multisigs []AccountMultisig     `toml:"multisig"`
}

// AccountMultisig account.multisig
type AccountMultisig struct {
	Type      account.AccountType `toml:"type"`
	Required  int                 `toml:"required"`
	AuthUsers account.AuthType    `toml:"auth_users"`
}
