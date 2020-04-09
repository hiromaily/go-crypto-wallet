package _import

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
)

//privkey subcommand
type PrivKeyCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallet.Signer
}

func (c *PrivKeyCommand) Synopsis() string {
	return c.synopsis
}

func (c *PrivKeyCommand) Help() string {
	return `Usage: sign key import privkey
`
}

func (c *PrivKeyCommand) Run(args []string) int {
	c.ui.Output(c.Synopsis())

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// import generated private key for Authorization account to database
	err := c.wallet.ImportPrivateKey(account.AccountTypeAuthorization)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call ImportPrivateKey() %+v", err))
	}
	c.ui.Output("Done!")

	return 0
}
