package create

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	wallets "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet"
	ethwallet "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet/eth"
)

// newMPCCommand creates the `watch create mpc` command (ETH only).
func newMPCCommand(wallet *wallets.Watcher) *cobra.Command {
	var (
		from       string
		to         string
		amount     float64
		threshold  int
		partyIDs   string
		actionType string
	)

	cmd := &cobra.Command{
		Use:   "mpc",
		Short: "Create an unsigned MPC/TSS transaction file (ETH only)",
		Long: `Create an unsigned EIP-1559 transaction file for MPC/TSS threshold signing.
The file can then be signed cooperatively by the required number of MPC nodes
using 'keygen serve mpc' and submitted with 'watch send mpc'.`,
		Example: `  watch --coin eth create mpc \
    --from 0xFromAddress --to 0xRecipient --amount 0.1 \
    --threshold 2 --party-ids "node1,node2,node3" --action-type payment`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wallet == nil || *wallet == nil {
				return errors.New("wallet not initialized, check --coin flag")
			}
			ethWatch, ok := (*wallet).(*ethwallet.ETHWatch)
			if !ok {
				fmt.Println("this command is only supported for ETH")
				return nil
			}
			if from == "" {
				return errors.New("--from flag is required")
			}
			if to == "" {
				return errors.New("--to flag is required")
			}
			if partyIDs == "" {
				return errors.New("--party-ids flag is required")
			}
			ids := splitCSV(partyIDs)
			if len(ids) < threshold {
				return fmt.Errorf("number of party IDs (%d) must be >= threshold (%d)", len(ids), threshold)
			}
			return runCreateMPC(cmd.Context(), ethWatch, from, to, amount, threshold, ids, actionType)
		},
	}

	cmd.Flags().StringVar(&from, "from", "", "Sender address (required)")
	cmd.Flags().StringVar(&to, "to", "", "Recipient address (required)")
	cmd.Flags().Float64Var(&amount, "amount", 0, "Amount in Ether to send (required)")
	cmd.Flags().IntVar(&threshold, "threshold", 2, "Required signer count (m in m-of-n)")
	cmd.Flags().StringVar(&partyIDs, "party-ids", "", "Comma-separated party IDs for all MPC nodes (required)")
	cmd.Flags().StringVar(&actionType, "action-type", "payment", "Action type (deposit, payment, transfer)")
	cobra.CheckErr(cmd.MarkFlagRequired("from"))
	cobra.CheckErr(cmd.MarkFlagRequired("to"))
	cobra.CheckErr(cmd.MarkFlagRequired("amount"))
	cobra.CheckErr(cmd.MarkFlagRequired("party-ids"))

	return cmd
}

func runCreateMPC(
	ctx context.Context,
	ethWatch *ethwallet.ETHWatch,
	from, to string,
	amount float64,
	threshold int,
	partyIDs []string,
	actionType string,
) error {
	filePath, txHash, err := ethWatch.CreateMPCTx(ctx, from, to, amount, threshold, partyIDs, actionType)
	if err != nil {
		return fmt.Errorf("failed to create MPC transaction: %w", err)
	}

	fmt.Printf("Unsigned MPC transaction file created\n")
	fmt.Printf("  File:      %s\n", filePath)
	fmt.Printf("  TX Hash:   %s\n", txHash)
	fmt.Printf("  Threshold: %d-of-%d\n", threshold, len(partyIDs))
	fmt.Printf("\nStart MPC nodes and sign with:\n")
	fmt.Printf("  keygen --coin eth serve mpc --shard-path <path> --passphrase <pass>\n")
	fmt.Printf("  watch --coin eth send mpc --file %s --peer-addrs <addr1,addr2,...>\n", filePath)

	return nil
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
