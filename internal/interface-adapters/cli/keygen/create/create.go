package create

import (
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/internal/di"
	wallets "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet"
)

// AddCommands adds all create subcommands
func AddCommands(parentCmd *cobra.Command, wallet *wallets.Keygener, containerGetter func() di.Container) {
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
			return runHDKeyWithFlags(containerGetter(), hdkeyKeyNum, hdkeyAccount, hdkeyIsKeyPair)
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
			return runSeed(containerGetter(), seedValue)
		},
	}
	seedCmd.Flags().StringVar(&seedValue, "seed", "",
		"given seed is used to store in database instead of generating new seed (development use)")
	parentCmd.AddCommand(seedCmd)

	// multisig command
	var (
		multisigAccount string
		multisigType    string
	)
	multisigCmd := &cobra.Command{
		Use:   "multisig",
		Short: "create multisig address (traditional or MuSig2)",
		Long: `Create multisig addresses using either traditional P2SH/P2WSH multisig or MuSig2 Taproot.

Traditional multisig (default):
  - Uses Bitcoin Core's addmultisigaddress RPC
  - Creates P2SH or P2WSH addresses depending on --address-type flag
  - Compatible with all Bitcoin address types

MuSig2 multisig:
  - Creates Taproot (P2TR) addresses with aggregated public keys
  - Uses MuSig2 Schnorr signature aggregation protocol
  - More efficient and private than traditional multisig
  - Requires BIP86 Taproot support`,
		Example: `  # Create traditional multisig address (default)
  keygen --coin btc create multisig --account payment

  # Create MuSig2 Taproot address
  keygen --coin btc create multisig --account payment --multisig-type musig2`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMultisigWithFlags(containerGetter(), multisigAccount, multisigType)
		},
	}
	multisigCmd.Flags().StringVar(&multisigAccount, "account", "", "target account")
	multisigCmd.Flags().StringVar(&multisigType, "multisig-type", "traditional",
		"multisig type: 'traditional' (P2SH/P2WSH) or 'musig2' (Taproot)")
	parentCmd.AddCommand(multisigCmd)
}
