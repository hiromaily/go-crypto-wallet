package api

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/wallets"
)

//api subcommand
type APICommand struct {
	Name    string
	Version string
	UI      cli.Ui
	Wallet  wallets.Walleter
}

func (c *APICommand) Synopsis() string {
	return "Bitcoin API functionality"
}

var (
	unlocktxSynopsis        = "unlock locked transaction for unspent transaction"
	estimatefeeSynopsis     = "estimate fee"
	loggingSynopsis         = "logging"
	getnetworkinfoSynopsis  = "call getnetworkinfo"
	validateaddressSynopsis = "validate address"
	listunspentSynopsis     = "call listunspent"
	balanceSynopsis         = "get balance for account"
)

func (c *APICommand) Help() string {
	return fmt.Sprintf(`Usage: wallet api [Subcommands...]
Subcommands:
  unlocktx         %s
  estimatefee      %s
  logging          %s
  getnetworkinfo   %s
  validateaddress  %s
  listunspent      %s
`, unlocktxSynopsis, estimatefeeSynopsis, loggingSynopsis, getnetworkinfoSynopsis, validateaddressSynopsis, listunspentSynopsis)
}

func (c *APICommand) Run(args []string) int {
	c.UI.Output(c.Synopsis())

	flags := flag.NewFlagSet(c.Name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//farther subcommand import
	cmds := map[string]cli.CommandFactory{
		"unlocktx": func() (cli.Command, error) {
			return &UnLockTxCommand{
				name:     "unlocktx",
				synopsis: unlocktxSynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.Wallet,
			}, nil
		},
		"estimatefee": func() (cli.Command, error) {
			return &EstimateFeeCommand{
				name:     "estimatefee",
				synopsis: estimatefeeSynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.Wallet,
			}, nil
		},
		"logging": func() (cli.Command, error) {
			return &LoggingCommand{
				name:     "logging",
				synopsis: loggingSynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.Wallet,
			}, nil
		},
		"getnetworkinfo": func() (cli.Command, error) {
			return &GetnetworkInfoCommand{
				name:     "getnetworkinfo",
				synopsis: getnetworkinfoSynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.Wallet,
			}, nil
		},
		"validateaddress": func() (cli.Command, error) {
			return &ValidateAddressCommand{
				name:     "validateaddress",
				synopsis: validateaddressSynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.Wallet,
			}, nil
		},
		"listunspent": func() (cli.Command, error) {
			return &ListUnspentCommand{
				name:     "listunspent",
				synopsis: listunspentSynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.Wallet,
			}, nil
		},
		"balance": func() (cli.Command, error) {
			return &BalanceCommand{
				name:     "balance",
				synopsis: balanceSynopsis,
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
