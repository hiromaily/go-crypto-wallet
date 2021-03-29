package create

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-crypto-wallet/pkg/command"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// CreateCommand create subcommand
type CreateCommand struct {
	Name    string
	Version string
	UI      cli.Ui
	Wallet  wallets.Keygener
}

// Synopsis is explanation for this subcommand
func (c *CreateCommand) Synopsis() string {
	return "create resources"
}

var (
	keySynopsis      = "create one key for debug use"
	hdkeySynopsis    = "create HD key"
	seedSynopsis     = "create seed"
	multisigSynopsis = "create multisig address"
)

// Help returns usage for this subcommand
func (c *CreateCommand) Help() string {
	return fmt.Sprintf(`Usage: keygen create [Subcommands...]
Subcommands:
  key       %s
  hdkey     %s
  seed      %s
  multisig  %s
`, keySynopsis, hdkeySynopsis, seedSynopsis, multisigSynopsis)
}

// Run executes this subcommand
func (c *CreateCommand) Run(args []string) int {
	c.UI.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.Name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// farther subcommand import
	cmds := map[string]cli.CommandFactory{
		"key": func() (cli.Command, error) {
			return &HDKeyCommand{
				name:     "key",
				synopsis: keySynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.Wallet,
			}, nil
		},
		"hdkey": func() (cli.Command, error) {
			return &HDKeyCommand{
				name:     "hdkey",
				synopsis: hdkeySynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.Wallet,
			}, nil
		},
		"seed": func() (cli.Command, error) {
			return &SeedCommand{
				name:     "seed",
				synopsis: seedSynopsis,
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
