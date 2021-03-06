package create

import (
	"flag"
	"fmt"

	"github.com/bookerzzz/grok"
	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// HDKeyCommand hdkey subcommand
type HDKeyCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Keygener
}

// Synopsis is explanation for this subcommand
func (c *HDKeyCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *HDKeyCommand) Help() string {
	return `Usage: keygen create hdkey [options...]
Options:
  -keynum   number of generating hd key
  -account  target account
  -keypair  keypair for XRP
`
}

// Run executes this subcommand
func (c *HDKeyCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	var (
		keyNum    uint64
		acnt      string
		isKeyPair bool
	)
	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.Uint64Var(&keyNum, "keynum", 0, "number of generating hd key")
	flags.StringVar(&acnt, "account", "", "target account")
	flags.BoolVar(&isKeyPair, "keypair", false, "keypair for XRP")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//validator
	if keyNum == 0 {
		c.ui.Error("number of key option [-keynum] is required")
		return 1
	}
	if !account.ValidateAccountType(acnt) {
		c.ui.Error("account option [-account] is invalid")
		return 1
	}
	if !account.NotAllow(acnt, []account.AccountType{account.AccountTypeAuthorization}) {
		c.ui.Error(fmt.Sprintf("account: %s is not allowed", account.AccountTypeAuthorization))
		return 1
	}

	// create seed
	bSeed, err := c.wallet.GenerateSeed()
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call GenerateSeed() %+v", err))
		return 1
	}

	//generate key for hd wallet
	keys, err := c.wallet.GenerateAccountKey(account.AccountType(acnt), bSeed, uint32(keyNum), isKeyPair)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call GenerateAccountKey() %+v", err))
		return 1
	}
	grok.Value(keys)

	return 0
}
