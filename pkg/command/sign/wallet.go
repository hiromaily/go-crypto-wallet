package sign

import (
	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-crypto-wallet/pkg/command"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/keygen/api/btc"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/sign/create"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/sign/export"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/sign/imports"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/sign/sign"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/wallet/api/eth"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/btcwallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/ethwallet"
)

// WalletSubCommands returns subcommand for signature
func WalletSubCommands(wallet wallets.Signer, version string) map[string]cli.CommandFactory {
	cmds := map[string]cli.CommandFactory{
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
			return &imports.ImportCommand{
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
	switch v := wallet.(type) {
	case *btcwallet.BTCSign:
		cmds["api"] = func() (cli.Command, error) {
			return &btc.APICommand{
				Name:    "api",
				Version: version,
				UI:      command.ClolorUI(),
				BTC:     v.BTC,
			}, nil
		}
	case *ethwallet.ETHSign:
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
