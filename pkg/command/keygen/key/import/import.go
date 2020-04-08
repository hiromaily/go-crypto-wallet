package _import

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
)

// import subcommand
type ImportCommand struct {
	Name        string
	Version     string
	SynopsisExp string
	UI          cli.Ui
	Wallet      wallet.Keygener
}

func (c *ImportCommand) Synopsis() string {
	return c.SynopsisExp
}

var (
	privkeySynopsis  = "import private key"
	multisigSynopsis = "import multisig address"
)

func (c *ImportCommand) Help() string {
	return fmt.Sprintf(`Usage: keygen import [Subcommands...]
Subcommands:
  privkey   %s
  multisig  %s
`, privkeySynopsis, multisigSynopsis)
}

func (c *ImportCommand) Run(args []string) int {
	c.UI.Output(c.Synopsis())

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
