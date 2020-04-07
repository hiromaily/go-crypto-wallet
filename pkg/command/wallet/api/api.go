package api

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/service"
)

const (
	apiName = "api"
)

//api subcommand
type APICommand struct {
	name    string
	version string
	ui      cli.Ui
	wallet  *service.Wallet
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
	c.ui.Output(c.Synopsis())

	flags := flag.NewFlagSet(apiName, flag.ContinueOnError)
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
				wallet:   c.wallet,
			}, nil
		},
		"estimatefee": func() (cli.Command, error) {
			return &EstimateFeeCommand{
				name:     "estimatefee",
				synopsis: estimatefeeSynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.wallet,
			}, nil
		},
		"logging": func() (cli.Command, error) {
			return &LoggingCommand{
				name:     "logging",
				synopsis: loggingSynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.wallet,
			}, nil
		},
		"getnetworkinfo": func() (cli.Command, error) {
			return &GetnetworkInfoCommand{
				name:     "getnetworkinfo",
				synopsis: getnetworkinfoSynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.wallet,
			}, nil
		},
		"validateaddress": func() (cli.Command, error) {
			return &ValidateAddressCommand{
				name:     "validateaddress",
				synopsis: validateaddressSynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.wallet,
			}, nil
		},
		"listunspent": func() (cli.Command, error) {
			return &ListUnspentCommand{
				name:     "listunspent",
				synopsis: listunspentSynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.wallet,
			}, nil
		},
	}
	cl := command.CreateSubCommand(apiName, c.version, args, cmds)

	code, err := cl.Run()
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call Run() subcommand of %s: %v", apiName, err))
	}
	return code
}
