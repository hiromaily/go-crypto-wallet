package key

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/command/signature/key/add"
	"github.com/hiromaily/go-bitcoin/pkg/command/signature/key/create"
	"github.com/hiromaily/go-bitcoin/pkg/command/signature/key/export"
	_import "github.com/hiromaily/go-bitcoin/pkg/command/signature/key/import"
	"github.com/hiromaily/go-bitcoin/pkg/wallets"
)

//key subcommand
type KeyCommand struct {
	Name    string
	Version string
	UI      cli.Ui
	Wallet  wallets.Signer
}

func (c *KeyCommand) Synopsis() string {
	return "key importing functionality"
}

var (
	addSynopsis    = "add key for multisig address"
	createSynopsis = "create resources"
	exportSynopsis = "export resources"
	importSynopsis = "import resource"
)

func (c *KeyCommand) Help() string {
	return fmt.Sprintf(`Usage: sign key [Subcommands...]
Subcommands:
  add     %s
  create  %s
  export  %s
  import  %s
`, addSynopsis, createSynopsis, exportSynopsis, importSynopsis)
}

func (c *KeyCommand) Run(args []string) int {
	c.UI.Output(c.Synopsis())

	flags := flag.NewFlagSet(c.Name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//farther subcommand import
	cmds := map[string]cli.CommandFactory{
		"add": func() (cli.Command, error) {
			return &add.AddCommand{
				Name:        "add",
				Version:     c.Version,
				SynopsisExp: addSynopsis,
				UI:          command.ClolorUI(),
				Wallet:      c.Wallet,
			}, nil
		},
		"create": func() (cli.Command, error) {
			return &create.CreateCommand{
				Name:        "create",
				Version:     c.Version,
				SynopsisExp: createSynopsis,
				UI:          command.ClolorUI(),
				Wallet:      c.Wallet,
			}, nil
		},
		"export": func() (cli.Command, error) {
			return &export.ExportCommand{
				Name:        "export",
				Version:     c.Version,
				SynopsisExp: exportSynopsis,
				UI:          command.ClolorUI(),
				Wallet:      c.Wallet,
			}, nil
		},
		"import": func() (cli.Command, error) {
			return &_import.ImportCommand{
				Name:        "import",
				Version:     c.Version,
				SynopsisExp: importSynopsis,
				UI:          command.ClolorUI(),
				Wallet:      c.Wallet,
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
