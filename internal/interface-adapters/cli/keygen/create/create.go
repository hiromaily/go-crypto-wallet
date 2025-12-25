package create

import (
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/internal/di"
	wallets "github.com/hiromaily/go-crypto-wallet/internal/wallet"
)

// AddCommands adds all create subcommands
func AddCommands(parentCmd *cobra.Command, wallet *wallets.Keygener, container di.Container) {
	// key command
	keyCmd := &cobra.Command{
		Use:   "key",
		Short: "create one key for debug use",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runKey(*wallet)
		},
	}
	parentCmd.AddCommand(keyCmd)

	// hdkey command
	var (
		hdkeyKeyNum    uint64
		hdkeyAccount   string
		hdkeyIsKeyPair bool
	)
	hdkeyCmd := &cobra.Command{
		Use:   "hdkey",
		Short: "create HD key",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHDKeyWithFlags(container, hdkeyKeyNum, hdkeyAccount, hdkeyIsKeyPair)
		},
	}
	hdkeyCmd.Flags().Uint64Var(&hdkeyKeyNum, "keynum", 0, "number of generating hd key")
	hdkeyCmd.Flags().StringVar(&hdkeyAccount, "account", "", "target account")
	hdkeyCmd.Flags().BoolVar(&hdkeyIsKeyPair, "keypair", false, "keypair for XRP")
	parentCmd.AddCommand(hdkeyCmd)

	// seed command
	var seedValue string
	seedCmd := &cobra.Command{
		Use:   "seed",
		Short: "create seed",
		Long:  "create seed for wallet. If --seed is provided, it will be stored instead of generating a new one",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSeed(container, seedValue)
		},
	}
	seedCmd.Flags().StringVar(&seedValue, "seed", "",
		"given seed is used to store in database instead of generating new seed (development use)")
	parentCmd.AddCommand(seedCmd)

	// multisig command
	var multisigAccount string
	multisigCmd := &cobra.Command{
		Use:   "multisig",
		Short: "create multisig address",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMultisigWithAccount(container, multisigAccount)
		},
	}
	multisigCmd.Flags().StringVar(&multisigAccount, "account", "", "target account")
	parentCmd.AddCommand(multisigCmd)
}
