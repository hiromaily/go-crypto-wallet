package app

import (
	"github.com/spf13/cobra"
)

// AddGlobalFlags adds common persistent flags to a cobra command.
// These flags are shared across keygen, sign, and watch wallets.
func AddGlobalFlags(cmd *cobra.Command, opts *Options) {
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
	cmd.PersistentFlags().StringVar(
		&opts.CoinTypeCode,
		"coin", "btc",
		"coin type code: btc, bch, eth, xrp, hyt",
	)
	cmd.PersistentFlags().StringVarP(
		&opts.BTCWallet,
		"wallet", "w", "",
		"specify wallet.dat in bitcoin core",
	)
}

// AddGlobalFlagsWithCoinOptions adds common persistent flags with customized coin help text.
func AddGlobalFlagsWithCoinOptions(cmd *cobra.Command, opts *Options, coinHelp string) {
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
	cmd.PersistentFlags().StringVar(
		&opts.CoinTypeCode,
		"coin", "btc",
		coinHelp,
	)
	cmd.PersistentFlags().StringVarP(
		&opts.BTCWallet,
		"wallet", "w", "",
		"specify wallet.dat in bitcoin core",
	)
}
