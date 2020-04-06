package key

import (
	"flag"
	"log"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/service"
)

const keyName = "key"

//key subcommand
type KeyCommand struct {
	version string
	ui      cli.Ui
	wallet  *service.Wallet
}

func (c *KeyCommand) Synopsis() string {
	return "key importing functionality"
}

func (c *KeyCommand) Help() string {
	return `Usage: wallet key [Subcommands...]
Subcommands:
  import  import file 
`
}

func (c *KeyCommand) Run(args []string) int {
	flags := flag.NewFlagSet(keyName, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//farther subcommand import
	cmds := map[string]cli.CommandFactory{
		"import": func() (cli.Command, error) {
			return &ImportCommand{
				ui:     command.ClolorUI(),
				wallet: c.wallet,
			}, nil
		},
	}
	cl := command.CreateSubCommand(keyName, c.version, args, cmds)

	code, err := cl.Run()
	if err != nil {
		log.Printf("fail to call Run() subcommand of key: %v\n", err)
	}
	return code
}
