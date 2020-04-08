package wallet

import (
	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/command/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/command/wallet/db"
	"github.com/hiromaily/go-bitcoin/pkg/command/wallet/key"
	"github.com/hiromaily/go-bitcoin/pkg/command/wallet/monitoring"
	"github.com/hiromaily/go-bitcoin/pkg/command/wallet/payment"
	"github.com/hiromaily/go-bitcoin/pkg/command/wallet/receipt"
	"github.com/hiromaily/go-bitcoin/pkg/command/wallet/sending"
	"github.com/hiromaily/go-bitcoin/pkg/command/wallet/transfer"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/service"
)

func WalletSubCommands(wallet *service.Wallet, version string) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"api": func() (cli.Command, error) {
			return &api.APICommand{
				Name:    "api",
				Version: version,
				UI:      command.ClolorUI(),
				Wallet:  wallet,
			}, nil
		},
		"db": func() (cli.Command, error) {
			return &db.DBCommand{
				Name:    "db",
				Version: version,
				UI:      command.ClolorUI(),
				Wallet:  wallet,
			}, nil
		},
		"key": func() (cli.Command, error) {
			return &key.KeyCommand{
				Name:    "key",
				Version: version,
				UI:      command.ClolorUI(),
				Wallet:  wallet,
			}, nil
		},
		"monitoring": func() (cli.Command, error) {
			return &monitoring.MonitoringCommand{
				Name:    "monitoring",
				Version: version,
				UI:      command.ClolorUI(),
				Wallet:  wallet,
			}, nil
		},
		"payment": func() (cli.Command, error) {
			return &payment.PaymentCommand{
				Name:    "payment",
				Version: version,
				UI:      command.ClolorUI(),
				Wallet:  wallet,
			}, nil
		},
		"receipt": func() (cli.Command, error) {
			return &receipt.ReceiptCommand{
				Name:    "receipt",
				Version: version,
				UI:      command.ClolorUI(),
				Wallet:  wallet,
			}, nil
		},
		"sending": func() (cli.Command, error) {
			return &sending.SendingCommand{
				Name:   "sending",
				UI:     command.ClolorUI(),
				Wallet: wallet,
			}, nil
		},
		"transfer": func() (cli.Command, error) {
			return &transfer.TransferCommand{
				Name:   "transfer",
				UI:     command.ClolorUI(),
				Wallet: wallet,
			}, nil
		},
	}
}
