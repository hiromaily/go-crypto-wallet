package wallet

import (
	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/command/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/command/wallet/create"
	"github.com/hiromaily/go-bitcoin/pkg/command/wallet/db"
	"github.com/hiromaily/go-bitcoin/pkg/command/wallet/imports"
	"github.com/hiromaily/go-bitcoin/pkg/command/wallet/monitor"
	"github.com/hiromaily/go-bitcoin/pkg/command/wallet/send"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// WalletSubCommands returns subcommand for wallet
// nolint: golint
func WalletSubCommands(wallet wallets.Watcher, version string) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"import": func() (cli.Command, error) {
			return &imports.ImportCommand{
				Name:    "import",
				Version: version,
				UI:      command.ClolorUI(),
				Wallet:  wallet,
			}, nil
		},
		"create": func() (cli.Command, error) {
			return &create.CreateCommand{
				Name:   "transfer",
				UI:     command.ClolorUI(),
				Wallet: wallet,
			}, nil
		},
		"send": func() (cli.Command, error) {
			return &send.SendCommand{
				Name:   "send",
				UI:     command.ClolorUI(),
				Wallet: wallet,
			}, nil
		},
		"monitor": func() (cli.Command, error) {
			return &monitor.MonitorCommand{
				Name:    "monitor",
				Version: version,
				UI:      command.ClolorUI(),
				Wallet:  wallet,
			}, nil
		},
		"api": func() (cli.Command, error) {
			return &api.APICommand{
				Name:    "api",
				Version: version,
				UI:      command.ClolorUI(),
				Wallet:  wallet,
			}, nil
		},
		"db": func() (cli.Command, error) {
			return &db.DBCommand{
				Name:    "db",
				Version: version,
				UI:      command.ClolorUI(),
				Wallet:  wallet,
			}, nil
		},
	}
}
