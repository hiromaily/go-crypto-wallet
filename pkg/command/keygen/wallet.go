package keygen

import (
	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/command/keygen/key"
	"github.com/hiromaily/go-bitcoin/pkg/command/keygen/signature"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
)

func WalletSubCommands(wallet wallet.Keygener, version string) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"key": func() (cli.Command, error) {
			return &key.KeyCommand{
				Name:    "key",
				Version: version,
				UI:      command.ClolorUI(),
				Wallet:  wallet,
			}, nil
		},
		"signature": func() (cli.Command, error) {
			return &signature.SignatureCommand{
				Name:   "signature",
				UI:     command.ClolorUI(),
				Wallet: wallet,
			}, nil
		},
	}
}
