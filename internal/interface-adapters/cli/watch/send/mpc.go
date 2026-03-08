package send

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	wallets "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet"
	ethwallet "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet/eth"
)

// newSendMPCCommand creates the `watch send mpc` command (ETH only).
func newSendMPCCommand(wallet *wallets.Watcher) *cobra.Command {
	var (
		filePath  string
		peerAddrs string
	)

	cmd := &cobra.Command{
		Use:   "mpc",
		Short: "Initiate MPC/TSS signing session and broadcast the transaction (ETH only)",
		Long: `Read an unsigned ETH MPC transaction file, initiate a TSS signing session with
the specified MPC nodes, apply the assembled signature, and broadcast the transaction.`,
		Example: `  watch --coin eth send mpc \
    --file payment_mpc_<uuid>.json \
    --peer-addrs "localhost:9001,localhost:9002"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wallet == nil || *wallet == nil {
				return errors.New("wallet not initialized, check --coin flag")
			}
			ethWatch, ok := (*wallet).(*ethwallet.ETHWatch)
			if !ok {
				fmt.Println("this command is only supported for ETH")
				return nil
			}
			if filePath == "" {
				return errors.New("--file flag is required")
			}
			if peerAddrs == "" {
				return errors.New("--peer-addrs flag is required")
			}
			addrs := splitMPCAddrs(peerAddrs)
			if len(addrs) == 0 {
				return errors.New("--peer-addrs must contain at least one address")
			}
			return runSendMPC(cmd.Context(), ethWatch, filePath, addrs)
		},
	}

	cmd.Flags().StringVar(&filePath, "file", "", "Unsigned MPC transaction JSON file path (required)")
	cmd.Flags().StringVar(&peerAddrs, "peer-addrs", "", "Comma-separated gRPC addresses of signing nodes (required)")
	cobra.CheckErr(cmd.MarkFlagRequired("file"))
	cobra.CheckErr(cmd.MarkFlagRequired("peer-addrs"))

	return cmd
}

func runSendMPC(
	ctx context.Context,
	ethWatch *ethwallet.ETHWatch,
	filePath string,
	peerAddrs []string,
) error {
	txHash, err := ethWatch.SendMPCTx(ctx, filePath, peerAddrs)
	if err != nil {
		return fmt.Errorf("failed to send MPC transaction: %w", err)
	}

	fmt.Printf("MPC transaction broadcast successfully\n")
	fmt.Printf("  TX Hash: %s\n", txHash)

	return nil
}

func splitMPCAddrs(s string) []string {
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
