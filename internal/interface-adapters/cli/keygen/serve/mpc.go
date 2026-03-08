// Package serve provides the `keygen serve mpc` command for running the MPC node daemon.
package serve

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	wallets "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet"
	ethwallet "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet/eth"
)

// AddCommands adds the serve mpc command under the given parent command.
func AddCommands(parentCmd *cobra.Command, wallet *wallets.Keygener) {
	mpcCmd := &cobra.Command{
		Use:   "mpc",
		Short: "Start the MPC node gRPC signing daemon (ETH only)",
		Long: `Start the MPC node daemon that accepts threshold signing sessions from the
Watch wallet coordinator. Blocks until SIGINT or SIGTERM is received.`,
		Example: `  keygen --coin eth serve mpc \
    --listen-addr localhost:9001 \
    --shard-path shard.json --passphrase secret \
    --party-id node1 --all-party-ids "node1,node2,node3"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wallet == nil || *wallet == nil {
				return errors.New("wallet not initialized, check --coin flag")
			}
			ethKeygen, ok := (*wallet).(*ethwallet.ETHKeygen)
			if !ok {
				fmt.Println("this command is only supported for ETH")
				return nil
			}

			listenAddr, _ := cmd.Flags().GetString("listen-addr")
			shardPath, _ := cmd.Flags().GetString("shard-path")
			passphrase, _ := cmd.Flags().GetString("passphrase")
			partyID, _ := cmd.Flags().GetString("party-id")
			allPartyIDsStr, _ := cmd.Flags().GetString("all-party-ids")

			if listenAddr == "" {
				return errors.New("--listen-addr flag is required")
			}
			if shardPath == "" {
				return errors.New("--shard-path flag is required")
			}
			if partyID == "" {
				return errors.New("--party-id flag is required")
			}

			allPartyIDs := splitServeCSV(allPartyIDsStr)

			// Set up context that cancels on SIGINT/SIGTERM for graceful shutdown.
			ctx := cmd.Context()
			ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
			defer cancel()

			input := keygenusecase.ServeMPCInput{
				ListenAddr:  listenAddr,
				ShardPath:   shardPath,
				Passphrase:  passphrase,
				PartyID:     partyID,
				AllPartyIDs: allPartyIDs,
			}

			fmt.Printf("Starting MPC node daemon on %s (party: %s)\n", listenAddr, partyID)
			fmt.Printf("Press Ctrl+C to stop.\n")

			if err := ethKeygen.ServeMPC(ctx, input); err != nil {
				return fmt.Errorf("MPC node daemon stopped with error: %w", err)
			}

			fmt.Printf("MPC node daemon stopped gracefully.\n")
			return nil
		},
	}

	mpcCmd.Flags().String("listen-addr", "", "gRPC listen address (e.g. localhost:9001) (required)")
	mpcCmd.Flags().String("shard-path", "", "Path to the encrypted key shard file (required)")
	mpcCmd.Flags().String("passphrase", "", "Passphrase for decrypting the key shard")
	mpcCmd.Flags().String("party-id", "", "This node's party ID (required)")
	mpcCmd.Flags().String("all-party-ids", "", "Comma-separated IDs of all ceremony participants")
	cobra.CheckErr(mpcCmd.MarkFlagRequired("listen-addr"))
	cobra.CheckErr(mpcCmd.MarkFlagRequired("shard-path"))
	cobra.CheckErr(mpcCmd.MarkFlagRequired("party-id"))

	parentCmd.AddCommand(mpcCmd)
}

func splitServeCSV(s string) []string {
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
