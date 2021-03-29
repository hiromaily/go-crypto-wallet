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
	Wallet  wallets.Watcher
}

// Synopsis is explanation for this subcommand
func (c *CreateCommand) Synopsis() string {
	return "creating functionality"
}

var (
	depositSynopsis  = "create a deposit unsigned transaction file for client account"
	paymentSynopsis  = "create a payment unsigned transaction file for payment account"
	transferSynopsis = "create a transfer unsigned transaction file between accounts"
	dbSynopsis       = "create payment_request table with dummy data for development use"
)

// Help returns usage for this subcommand
func (c *CreateCommand) Help() string {
	return fmt.Sprintf(`Usage: wallet create [Subcommands...]
Subcommands:
  deposit  %s
  payment  %s
  transfer %s
  db       %s
`, depositSynopsis, paymentSynopsis, transferSynopsis, dbSynopsis)
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
		"deposit": func() (cli.Command, error) {
			return &DepositCommand{
				name:     "deposit",
				synopsis: depositSynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.Wallet,
			}, nil
		},
		"payment": func() (cli.Command, error) {
			return &PaymentCommand{
				name:     "payment",
				synopsis: paymentSynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.Wallet,
			}, nil
		},
		"transfer": func() (cli.Command, error) {
			return &TransferCommand{
				name:     "transfer",
				synopsis: transferSynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.Wallet,
			}, nil
		},
		"db": func() (cli.Command, error) {
			return &DBCommand{
				name:     "db",
				synopsis: dbSynopsis,
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
