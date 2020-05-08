package keygen

import (
	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/command/keygen/api/btc"
	"github.com/hiromaily/go-bitcoin/pkg/command/keygen/api/eth"
	"github.com/hiromaily/go-bitcoin/pkg/command/keygen/create"
	"github.com/hiromaily/go-bitcoin/pkg/command/keygen/export"
	"github.com/hiromaily/go-bitcoin/pkg/command/keygen/imports"
	"github.com/hiromaily/go-bitcoin/pkg/command/keygen/sign"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets/btcwallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets/ethwallet"
)

// WalletSubCommands returns subcommand for keygen wallet
func WalletSubCommands(wallet wallets.Keygener, version string) map[string]cli.CommandFactory {
	cmds := map[string]cli.CommandFactory{
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
	switch v := wallet.(type) {
	case *btcwallet.BTCKeygen:
		cmds["api"] = func() (cli.Command, error) {
			return &btc.APICommand{
				Name:    "api",
				Version: version,
				UI:      command.ClolorUI(),
				BTC:     v.BTC,
			}, nil
		}
	case *ethwallet.ETHKeygen:
		cmds["api"] = func() (cli.Command, error) {
			return &eth.APICommand{
				Name:    "api",
				Version: version,
				UI:      command.ClolorUI(),
				ETH:     v.ETH,
			}, nil
		}
	}
	return cmds
}
