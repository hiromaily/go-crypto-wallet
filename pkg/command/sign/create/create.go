package create

import (
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/pkg/di"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// AddCommands adds all create subcommands
func AddCommands(parentCmd *cobra.Command, wallet *wallets.Signer, container di.Container) {
	// seed command
	var seedValue string
	seedCmd := &cobra.Command{
		Use:   "seed",
		Short: "create seed",
		Long:  "create seed for wallet. If --seed is provided, it will be stored instead of generating a new one",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSeed(*wallet, seedValue)
		},
	}
	seedCmd.Flags().StringVar(&seedValue, "seed", "",
		"given seed is used to store in database instead of generating new seed (development use)")
	parentCmd.AddCommand(seedCmd)

	// hdkey command
	hdkeyCmd := &cobra.Command{
		Use:   "hdkey",
		Short: "create key for hd wallet for Authorization account",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHDKey(*wallet)
		},
	}
	parentCmd.AddCommand(hdkeyCmd)
}
