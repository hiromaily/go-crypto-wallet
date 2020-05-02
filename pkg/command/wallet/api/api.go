package api

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// APICommand api subcommand
type APICommand struct {
	Name    string
	Version string
	UI      cli.Ui
	Wallet  wallets.Watcher
}

// Synopsis is explanation for this subcommand
func (c *APICommand) Synopsis() string {
	return "Bitcoin API functionality"
}

var (
	balanceSynopsis         = "get balance for account"
	estimatefeeSynopsis     = "estimate fee"
	getnetworkinfoSynopsis  = "call getnetworkinfo"
	getaddressinfoSynopsis  = "call getaddressinfo"
	listunspentSynopsis     = "call listunspent"
	loggingSynopsis         = "logging"
	unlocktxSynopsis        = "unlock locked transaction for unspent transaction"
	validateaddressSynopsis = "validate address"
)

// Help returns usage for this subcommand
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

// Run executes this subcommand
func (c *APICommand) Run(args []string) int {
	c.UI.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.Name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//farther subcommand import
	cmds := map[string]cli.CommandFactory{
		"balance": func() (cli.Command, error) {
			return &BalanceCommand{
				name:     "balance",
				synopsis: balanceSynopsis,
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
		"getaddressinfo": func() (cli.Command, error) {
			return &GetAddressInfoCommand{
				name:     "getnetworkinfo",
				synopsis: getaddressinfoSynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.Wallet,
			}, nil
		},
		"getnetworkinfo": func() (cli.Command, error) {
			return &GetNetworkInfoCommand{
				name:     "getnetworkinfo",
				synopsis: getnetworkinfoSynopsis,
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
		"logging": func() (cli.Command, error) {
			return &LoggingCommand{
				name:     "logging",
				synopsis: loggingSynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.Wallet,
			}, nil
		},
		"unlocktx": func() (cli.Command, error) {
			return &UnLockTxCommand{
				name:     "unlocktx",
				synopsis: unlocktxSynopsis,
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
	}
	cl := command.CreateSubCommand(c.Name, c.Version, args, cmds)

	code, err := cl.Run()
	if err != nil {
		c.UI.Error(fmt.Sprintf("fail to call Run() subcommand of %s: %v", c.Name, err))
	}
	return code
}
