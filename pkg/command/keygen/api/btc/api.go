package btc

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-crypto-wallet/pkg/command"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
)

// APICommand api subcommand
type APICommand struct {
	Name    string
	Version string
	UI      cli.Ui
	BTC     btcgrp.Bitcoiner
}

// Synopsis is explanation for this subcommand
func (c *APICommand) Synopsis() string {
	return "Bitcoin API functionality"
}

var (
	encryptwalletSynopsis    = "encrypts the wallet with 'passphrase'"
	walletpassphraseSynopsis = `stores the wallet decryption key in memory for 'timeout' seconds.\n
this is needed prior to performing transactions related to private keys such as sending bitcoins`
	walletpassphrasechangeSynopsis = "changes the wallet passphrase from 'oldpassphrase' to 'newpassphrase'"
	walletlockSynopsis             = "removes the wallet encryption key from memory, locking the wallet"
	dumpwalletSynopsis             = "dumps all wallet keys in a human-readable format to a server-side file"
	importwalletSynopsis           = "Imports keys from a wallet dump file"
)

// Help returns usage for this subcommand
func (c *APICommand) Help() string {
	return fmt.Sprintf(`Usage: wallet api [Subcommands...]
Subcommands:
  encryptwallet          %s
  walletpassphrase       %s
  walletpassphrasechange %s
  walletlock             %s
  dumpwallet             %s
  importwallet           %s`,
		encryptwalletSynopsis,
		walletpassphraseSynopsis,
		walletpassphrasechangeSynopsis,
		walletlockSynopsis,
		dumpwalletSynopsis,
		importwalletSynopsis)
}

// Run executes this subcommand
func (c *APICommand) Run(args []string) int {
	c.UI.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.Name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//farther subcommand import
	cmds := map[string]cli.CommandFactory{
		"encryptwallet": func() (cli.Command, error) {
			return &EncryptWalletCommand{
				name:     "encryptwallet",
				synopsis: encryptwalletSynopsis,
				ui:       command.ClolorUI(),
				btc:      c.BTC,
			}, nil
		},
		"walletpassphrase": func() (cli.Command, error) {
			return &WalletPassphraseCommand{
				name:     "walletpassphrase",
				synopsis: walletpassphraseSynopsis,
				ui:       command.ClolorUI(),
				btc:      c.BTC,
			}, nil
		},
		"walletpassphrasechange": func() (cli.Command, error) {
			return &WalletPassphraseChangeCommand{
				name:     "walletpassphrasechange",
				synopsis: walletpassphrasechangeSynopsis,
				ui:       command.ClolorUI(),
				btc:      c.BTC,
			}, nil
		},
		"walletlock": func() (cli.Command, error) {
			return &WalletLockCommand{
				name:     "walletlock",
				synopsis: walletlockSynopsis,
				ui:       command.ClolorUI(),
				btc:      c.BTC,
			}, nil
		},
		"dumpwallet": func() (cli.Command, error) {
			return &DumpWalletCommand{
				name:     "dumpwallet",
				synopsis: dumpwalletSynopsis,
				ui:       command.ClolorUI(),
				btc:      c.BTC,
			}, nil
		},
		"importwallet": func() (cli.Command, error) {
			return &ImportWalletCommand{
				name:     "importwallet",
				synopsis: importwalletSynopsis,
				ui:       command.ClolorUI(),
				btc:      c.BTC,
			}, nil
		},
	}
	cl := command.CreateSubCommand(c.Name, c.Version, args, cmds)

	code, err := cl.Run()
	if err != nil {
		c.UI.Error(fmt.Sprintf("fail to call Run() subcommand of %s: %v", c.Name, err))
	}
	return code
}
