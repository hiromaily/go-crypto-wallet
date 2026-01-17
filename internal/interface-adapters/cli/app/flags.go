package app

import (
	"github.com/spf13/cobra"
)

// AddGlobalFlags adds common persistent flags to a cobra command.
// These flags are shared across keygen, sign, and watch wallets.
// It accepts an optional coinHelp string to customize the help text for the --coin flag.
func AddGlobalFlags(cmd *cobra.Command, opts *Options, coinHelp ...string) {
	cmd.PersistentFlags().StringVarP(
		&opts.ConfigPath,
		"config", "c", "",
		"config file path (required)",
	)
	cmd.PersistentFlags().StringVar(
		&opts.AccountConfigPath,
		"account-config", "",
		"account config file path for multisig settings",
	)

	helpText := "coin type code: btc, bch, eth, xrp, hyt"
	if len(coinHelp) > 0 && coinHelp[0] != "" {
		helpText = coinHelp[0]
	}
	cmd.PersistentFlags().StringVar(
		&opts.CoinTypeCode,
		"coin", "btc",
		helpText,
	)
	cmd.PersistentFlags().StringVarP(
		&opts.BTCWallet,
		"wallet", "w", "",
		"specify wallet.dat in bitcoin core",
	)
}
