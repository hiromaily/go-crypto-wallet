package watch

import (
	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-crypto-wallet/pkg/command"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/watch/api/btc"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/watch/api/eth"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/watch/api/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/watch/create"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/watch/db"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/watch/imports"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/watch/monitor"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/watch/send"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/btcwallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/ethwallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/xrpwallet"
)

// WatchSubCommands returns subcommand for wallet
// nolint: golint
func WatchSubCommands(wallet wallets.Watcher, version string, conf *config.WalletRoot) map[string]cli.CommandFactory {
	cmds := map[string]cli.CommandFactory{
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
		"db": func() (cli.Command, error) {
			return &db.DBCommand{
				Name:    "db",
				Version: version,
				UI:      command.ClolorUI(),
				Wallet:  wallet,
			}, nil
		},
	}
	switch v := wallet.(type) {
	case *btcwallet.BTCWatch:
		cmds["api"] = func() (cli.Command, error) {
			return &btc.APICommand{
				Name:    "api",
				Version: version,
				UI:      command.ClolorUI(),
				BTC:     v.BTC,
			}, nil
		}
	case *ethwallet.ETHWatch:
		cmds["api"] = func() (cli.Command, error) {
			return &eth.APICommand{
				Name:    "api",
				Version: version,
				UI:      command.ClolorUI(),
				ETH:     v.ETH,
			}, nil
		}
	case *xrpwallet.XRPWatch:
		cmds["api"] = func() (cli.Command, error) {
			return &xrp.APICommand{
				Name:    "api",
				Version: version,
				UI:      command.ClolorUI(),
				XRP:     v.XRP,
				TxData:  &conf.Ripple.API.TxData,
			}, nil
		}
	}
	return cmds
}
