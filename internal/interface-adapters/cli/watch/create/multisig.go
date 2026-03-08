package create

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	wallets "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet"
	ethwallet "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet/eth"
)

// newMultisigCommand creates the `watch create multisig` command (ETH only).
func newMultisigCommand(wallet *wallets.Watcher) *cobra.Command {
	var (
		safeAddress string
		to          string
		amount      float64
		threshold   int
		actionType  string
	)

	cmd := &cobra.Command{
		Use:   "multisig",
		Short: "Create an unsigned Safe multisig transaction proposal file (ETH only)",
		Long: `Propose a new Safe multisig transaction.
Generates a JSON file containing all Safe execution parameters and the computed
safeTxHash. The file must be signed by the required number of signers before
it can be submitted with 'watch send multisig'.`,
		Example: `  watch --coin eth create multisig \
    --safe 0xSafeAddress --to 0xRecipient --amount 0.1 --threshold 2 --action-type payment`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wallet == nil || *wallet == nil {
				return errors.New("wallet not initialized, check --coin flag")
			}
			ethWatch, ok := (*wallet).(*ethwallet.ETHWatch)
			if !ok {
				fmt.Println("this command is only supported for ETH")
				return nil
			}
			if safeAddress == "" {
				return errors.New("--safe flag is required")
			}
			if to == "" {
				return errors.New("--to flag is required")
			}
			return runCreateMultisig(ethWatch, safeAddress, to, amount, threshold, actionType)
		},
	}

	cmd.Flags().StringVar(&safeAddress, "safe", "", "Safe proxy contract address (required)")
	cmd.Flags().StringVar(&to, "to", "", "Recipient address (required)")
	cmd.Flags().Float64Var(&amount, "amount", 0, "Amount in Ether to send (required)")
	cmd.Flags().IntVar(&threshold, "threshold", 1, "Required signer count (m in m-of-n)")
	cmd.Flags().StringVar(&actionType, "action-type", "payment", "Action type (deposit, payment, transfer)")
	cobra.CheckErr(cmd.MarkFlagRequired("safe"))
	cobra.CheckErr(cmd.MarkFlagRequired("to"))
	cobra.CheckErr(cmd.MarkFlagRequired("amount"))

	return cmd
}

func runCreateMultisig(
	ethWatch *ethwallet.ETHWatch,
	safeAddress, to string,
	amount float64,
	threshold int,
	actionType string,
) error {
	filePath, uuid, err := ethWatch.CreateMultisigTx(safeAddress, to, amount, threshold, actionType)
	if err != nil {
		return fmt.Errorf("failed to create multisig transaction: %w", err)
	}

	fmt.Printf("Unsigned Safe multisig proposal created\n")
	fmt.Printf("  UUID:      %s\n", uuid)
	fmt.Printf("  File:      %s\n", filePath)
	fmt.Printf("  Threshold: %d\n", threshold)
	fmt.Printf("\nHave each authorized signer sign the file with:\n")
	fmt.Printf("  keygen --coin eth sign signature --file %s --signer-address <addr>\n", filePath)

	return nil
}
