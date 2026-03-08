package sign

import (
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/internal/di"
	wallets "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet"
	ethwallet "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet/eth"
)

// AddCommands adds all sign subcommands
func AddCommands(parentCmd *cobra.Command, wallet *wallets.Keygener, containerGetter func() di.Container) {
	// signature command
	var (
		signatureFile string
		signerAddress string
	)
	signatureCmd := &cobra.Command{
		Use:   "signature",
		Short: "sign on unsigned transaction (account would be found from file name)",
		RunE: func(cmd *cobra.Command, args []string) error {
			// For ETH keygen, pass the signer address before delegating to SignTx.
			if signerAddress != "" {
				if ethk, ok := (*wallet).(*ethwallet.ETHKeygen); ok {
					ethk.SetSignerAddress(signerAddress)
				}
			}
			return runSignature(containerGetter(), signatureFile)
		},
	}
	signatureCmd.Flags().StringVar(&signatureFile, "file", "", "import file path for signed transactions")
	signatureCmd.Flags().StringVar(&signerAddress, "signer-address", "",
		"signer ETH address for Safe multisig signing (ETH only)")
	parentCmd.AddCommand(signatureCmd)

	// musig2 commands (BTC only)
	addMuSig2Commands(parentCmd, wallet, containerGetter)
}
