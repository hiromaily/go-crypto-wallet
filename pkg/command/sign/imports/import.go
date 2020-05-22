package imports

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-crypto-wallet/pkg/command"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// ImportCommand import subcommand
type ImportCommand struct {
	Name    string
	Version string
	UI      cli.Ui
	Wallet  wallets.Signer
}

// Synopsis is explanation for this subcommand
func (c *ImportCommand) Synopsis() string {
	return "import resource"
}

var (
	privkeySynopsis = "import generated private key for Authorization account to database"
)

// Help returns usage for this subcommand
func (c *ImportCommand) Help() string {
	return fmt.Sprintf(`Usage: sign import [Subcommands...]
Subcommands:
  privkey   %s
`, privkeySynopsis)
}

// Run executes this subcommand
func (c *ImportCommand) Run(args []string) int {
	c.UI.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.Name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//farther subcommand import
	cmds := map[string]cli.CommandFactory{
		"privkey": func() (cli.Command, error) {
			return &PrivKeyCommand{
				name:     "privkey",
				synopsis: privkeySynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.Wallet,
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
