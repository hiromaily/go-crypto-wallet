package signature

import (
	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/command/signature/key"
	"github.com/hiromaily/go-bitcoin/pkg/command/signature/signature"
	"github.com/hiromaily/go-bitcoin/pkg/wallets"
)

func WalletSubCommands(wallet wallets.Signer, version string) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"key": func() (cli.Command, error) {
			return &key.KeyCommand{
				Name:    "key",
				Version: version,
				UI:      command.ClolorUI(),
				Wallet:  wallet,
			}, nil
		},
		"signature": func() (cli.Command, error) {
			return &signature.SignatureCommand{
				Name:   "signature",
				UI:     command.ClolorUI(),
				Wallet: wallet,
			}, nil
		},
	}
}
