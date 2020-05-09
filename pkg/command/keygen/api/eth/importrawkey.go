package eth

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/ethgrp"
)

// ImportRawKeyCommand syncing subcommand
type ImportRawKeyCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	eth      ethgrp.Ethereumer
}

// Synopsis is explanation for this subcommand
func (c *ImportRawKeyCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *ImportRawKeyCommand) Help() string {
	return `Usage: keygen api importrawkey [options...]
Options:
  -key   private key
  -pass  passphrase
`
}

// Run executes this subcommand
func (c *ImportRawKeyCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	var privKey, passPhrase string

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.StringVar(&privKey, "key", "", "private key")
	flags.StringVar(&passPhrase, "pass", "", "passphrase")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// validation
	if privKey == "" {
		c.ui.Error("key option [-key] is invalid")
		return 1
	}
	if passPhrase == "" {
		c.ui.Error("pass option [-pass] is invalid")
		return 1
	}

	//if strings.HasPrefix(privKey, "0x") {
	//	privKey = strings.TrimLeft(privKey, "0x")
	//}

	addr, err := c.eth.ImportRawKey(privKey, passPhrase)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call eth.ImportRawKey() %+v", err))
		return 1
	}

	c.ui.Info(fmt.Sprintf("new address: %s", addr))

	return 0
}
