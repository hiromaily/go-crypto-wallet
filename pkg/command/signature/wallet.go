package signature

import (
	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/command/signature/add"
	"github.com/hiromaily/go-bitcoin/pkg/command/signature/create"
	"github.com/hiromaily/go-bitcoin/pkg/command/signature/export"
	_import "github.com/hiromaily/go-bitcoin/pkg/command/signature/import"
	"github.com/hiromaily/go-bitcoin/pkg/command/signature/sign"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// WalletSubCommands
func WalletSubCommands(wallet wallets.Signer, version string) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"add": func() (cli.Command, error) {
			return &add.AddCommand{
				Name:    "add",
				Version: version,
				UI:      command.ClolorUI(),
				Wallet:  wallet,
			}, nil
		},
		"create": func() (cli.Command, error) {
			return &create.CreateCommand{
				Name:    "add",
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
			return &_import.ImportCommand{
				Name:    "import",
				Version: version,
				UI:      command.ClolorUI(),
				Wallet:  wallet,
			}, nil
		},
		"sign": func() (cli.Command, error) {
			return &sign.SignCommand{
				Name:   "signature",
				UI:     command.ClolorUI(),
				Wallet: wallet,
			}, nil
		},
	}
}
