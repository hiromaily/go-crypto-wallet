package imports

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// FullPubKeyCommand fullpubkey subcommand
type FullPubKeyCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Keygener
}

// Synopsis is explanation for this subcommand
func (c *FullPubKeyCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *FullPubKeyCommand) Help() string {
	return `Usage: keygen key import fullpubkey [options...]
Options:
  -file  full-pubkey file path
`
}

// Run executes this subcommand
func (c *FullPubKeyCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	var (
		fileName string
	)
	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.StringVar(&fileName, "file", "", "full-pubkey file path")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//validator
	if fileName == "" {
		c.ui.Error("file option [-file] is required")
		return 1
	}

	//import generated private key to keygen wallet
	err := c.wallet.FullPubKeyImport().ImportFullPubKey(fileName)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call FullPubKeyImport().ImportFullPubKey() %+v", err))
		return 1
	}
	c.ui.Output("Done!")

	return 0
}
