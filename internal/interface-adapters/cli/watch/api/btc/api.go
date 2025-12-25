package btc

import (
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin"
)

// AddCommands adds all Bitcoin API subcommands
func AddCommands(parentCmd *cobra.Command, btc bitcoin.Bitcoiner) {
	// balance command
	var balanceAccount string
	balanceCmd := &cobra.Command{
		Use:   "balance",
		Short: "get balance for account",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBalance(btc, balanceAccount)
		},
	}
	balanceCmd.Flags().StringVar(&balanceAccount, "account", "", "account")
	parentCmd.AddCommand(balanceCmd)

	// estimatefee command
	estimatefeeCmd := &cobra.Command{
		Use:   "estimatefee",
		Short: "estimate fee",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEstimateFee(btc)
		},
	}
	parentCmd.AddCommand(estimatefeeCmd)

	// getnetworkinfo command
	getnetworkinfoCmd := &cobra.Command{
		Use:   "getnetworkinfo",
		Short: "call getnetworkinfo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetNetworkInfo(btc)
		},
	}
	parentCmd.AddCommand(getnetworkinfoCmd)

	// getaddressinfo command
	var getaddressinfoAddress string
	getaddressinfoCmd := &cobra.Command{
		Use:   "getaddressinfo",
		Short: "call getaddressinfo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetAddressInfo(btc, getaddressinfoAddress)
		},
	}
	getaddressinfoCmd.Flags().StringVar(&getaddressinfoAddress, "address", "", "address")
	parentCmd.AddCommand(getaddressinfoCmd)

	// listunspent command
	var (
		listunspentAccount string
		listunspentNum     int64
	)
	listunspentCmd := &cobra.Command{
		Use:   "listunspent",
		Short: "call listunspent",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListUnspent(btc, listunspentAccount, listunspentNum)
		},
	}
	listunspentCmd.Flags().StringVar(&listunspentAccount, "account", "", "account")
	listunspentCmd.Flags().Int64Var(&listunspentNum, "num", -1, "confirmation number")
	parentCmd.AddCommand(listunspentCmd)

	// logging command
	loggingCmd := &cobra.Command{
		Use:   "logging",
		Short: "logging",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogging(btc)
		},
	}
	parentCmd.AddCommand(loggingCmd)

	// unlocktx command
	unlocktxCmd := &cobra.Command{
		Use:   "unlocktx",
		Short: "unlock locked transaction for unspent transaction",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUnlockTx(btc)
		},
	}
	parentCmd.AddCommand(unlocktxCmd)

	// validateaddress command
	var validateaddressAddress string
	validateaddressCmd := &cobra.Command{
		Use:   "validateaddress",
		Short: "validate address",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runValidateAddress(btc, validateaddressAddress)
		},
	}
	validateaddressCmd.Flags().StringVar(
		&validateaddressAddress, "address", "", "address like '2NFXSXxw8Fa6P6CSovkdjXE6UF4hupcTHtr'")
	parentCmd.AddCommand(validateaddressCmd)
}
