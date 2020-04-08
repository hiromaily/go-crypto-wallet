package wallet

import (
	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/command/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/service"
)

func WalletSubCommands(wallet *service.Wallet, version string) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"api": func() (cli.Command, error) {
			return &api.APICommand{
				Name:    "api",
				Version: version,
				UI:      command.ClolorUI(),
				Wallet:  wallet,
			}, nil
		},
	}
}
