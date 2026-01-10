package export

import (
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/internal/di"
	wallets "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet"
)

// AddCommands adds all export subcommands
func AddCommands(parentCmd *cobra.Command, wallet *wallets.Keygener, containerGetter func() di.Container) {
	// address command
	var addressAccount string
	addressCmd := &cobra.Command{
		Use:   "address",
		Short: "export generated PublicKey as csv file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAddress(containerGetter(), addressAccount)
		},
	}
	addressCmd.Flags().StringVar(&addressAccount, "account", "", "target account")
	parentCmd.AddCommand(addressCmd)

	// descriptor command
	var (
		descriptorAccount       string
		descriptorOutput        string
		descriptorFormat        string
		descriptorIncludeChange bool
	)
	descriptorCmd := &cobra.Command{
		Use:   "descriptor",
		Short: "export output descriptors to file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDescriptor(
				containerGetter(),
				descriptorAccount,
				descriptorOutput,
				descriptorFormat,
				descriptorIncludeChange,
			)
		},
	}
	descriptorCmd.Flags().StringVar(&descriptorAccount, "account", "", "target account (required)")
	descriptorCmd.Flags().StringVar(&descriptorOutput, "output", "", "output file path (required)")
	descriptorCmd.Flags().StringVar(
		&descriptorFormat,
		"format",
		"bitcoin-core",
		"output format: text, json, bitcoin-core (default)",
	)
	descriptorCmd.Flags().BoolVar(&descriptorIncludeChange, "include-change", false, "include change descriptors")
	parentCmd.AddCommand(descriptorCmd)
}
