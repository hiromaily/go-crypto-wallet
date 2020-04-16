package _import

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// import subcommand
type ImportCommand struct {
	Name    string
	Version string
	UI      cli.Ui
	Wallet  wallets.Keygener
}

func (c *ImportCommand) Synopsis() string {
	return "import resources"
}

var (
	privkeySynopsis  = "import generated private key in database to keygen wallet"
	multisigSynopsis = "import multisig addresses exported by signature wallet from csv file to database"
)

func (c *ImportCommand) Help() string {
	return fmt.Sprintf(`Usage: keygen import [Subcommands...]
Subcommands:
  privkey   %s
  multisig  %s
`, privkeySynopsis, multisigSynopsis)
}

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
		"multisig": func() (cli.Command, error) {
			return &MultisigCommand{
				name:     "multisig",
				synopsis: multisigSynopsis,
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
