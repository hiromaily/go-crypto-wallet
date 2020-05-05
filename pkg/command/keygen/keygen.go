package keygen

import (
	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/command/keygen/api"
	"github.com/hiromaily/go-bitcoin/pkg/command/keygen/create"
	"github.com/hiromaily/go-bitcoin/pkg/command/keygen/export"
	"github.com/hiromaily/go-bitcoin/pkg/command/keygen/imports"
	"github.com/hiromaily/go-bitcoin/pkg/command/keygen/sign"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// WalletSubCommands returns subcommand for keygen wallet
func WalletSubCommands(wallet wallets.Keygener, version string) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"api": func() (cli.Command, error) {
			return &api.APICommand{
				Name:    "api",
				Version: version,
				UI:      command.ClolorUI(),
				BTC:     wallet.GetBTC(),
			}, nil
		},
		"create": func() (cli.Command, error) {
			return &create.CreateCommand{
				Name:    "create",
				Version: version,
				UI:      command.ClolorUI(),
				Wallet:  wallet,
			}, nil
		},
		"export": func() (cli.Command, error) {
			return &export.ExportCommand{
				Name:    "export",
				Version: version,
				UI:      command.ClolorUI(),
				Wallet:  wallet,
			}, nil
		},
		"import": func() (cli.Command, error) {
			return &imports.ImportCommand{
				Name:    "import",
				Version: version,
				UI:      command.ClolorUI(),
				Wallet:  wallet,
			}, nil
		},
		"sign": func() (cli.Command, error) {
			return &sign.SignatureCommand{
				Name:   "sign",
				UI:     command.ClolorUI(),
				Wallet: wallet,
			}, nil
		},
	}
}
