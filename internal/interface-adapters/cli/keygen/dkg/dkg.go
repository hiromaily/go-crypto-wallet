// Package dkg provides CLI commands for MPC/TSS Distributed Key Generation.
package dkg

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	wallets "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet"
	ethwallet "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet/eth"
)

// AddCommands adds DKG-related commands under the given parent command.
func AddCommands(parentCmd *cobra.Command, wallet *wallets.Keygener) {
	parentCmd.AddCommand(newDKGCommand(wallet))
	parentCmd.AddCommand(newPreParamsCommand(wallet))
}

// newDKGCommand creates the `keygen dkg` command.
func newDKGCommand(wallet *wallets.Keygener) *cobra.Command {
	var (
		partyID         string
		allPartyIDs     string
		threshold       int
		peerAddrs       string
		preParamsPath   string
		shardOutputPath string
		passphrase      string
	)

	cmd := &cobra.Command{
		Use:   "dkg",
		Short: "Run the DKG ceremony for this MPC node (ETH only)",
		Long: `Run the Distributed Key Generation ceremony for this MPC node.
All participating nodes must run this command simultaneously.
On success, an encrypted key shard is saved to --shard-output.`,
		Example: `  keygen --coin eth dkg \
    --party-id node1 --all-party-ids "node1,node2,node3" \
    --threshold 2 --peers "localhost:9001,localhost:9002" \
    --pre-params-path pre_params.json --shard-output shard.json --passphrase secret`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wallet == nil || *wallet == nil {
				return errors.New("wallet not initialized, check --coin flag")
			}
			ethKeygen, ok := (*wallet).(*ethwallet.ETHKeygen)
			if !ok {
				fmt.Println("this command is only supported for ETH")
				return nil
			}
			if partyID == "" {
				return errors.New("--party-id flag is required")
			}
			if allPartyIDs == "" {
				return errors.New("--all-party-ids flag is required")
			}
			ids := splitCSV(allPartyIDs)
			peers := splitCSV(peerAddrs)
			input := keygenusecase.RunDKGInput{
				PartyID:         partyID,
				AllPartyIDs:     ids,
				Threshold:       threshold,
				PeerAddrs:       peers,
				PreParamsPath:   preParamsPath,
				ShardOutputPath: shardOutputPath,
				Passphrase:      passphrase,
			}
			output, err := ethKeygen.RunDKG(cmd.Context(), input)
			if err != nil {
				return fmt.Errorf("DKG ceremony failed: %w", err)
			}
			fmt.Printf("DKG ceremony completed successfully\n")
			fmt.Printf("  Joint ETH Address: %s\n", output.EthAddress)
			fmt.Printf("  Shard saved to:    %s\n", output.ShardPath)
			return nil
		},
	}

	cmd.Flags().StringVar(&partyID, "party-id", "", "This node's party ID (required)")
	cmd.Flags().StringVar(&allPartyIDs, "all-party-ids", "",
		"Comma-separated IDs of all ceremony participants (required)")
	cmd.Flags().IntVar(&threshold, "threshold", 2, "Required signer count (m in m-of-n)")
	cmd.Flags().StringVar(&peerAddrs, "peers", "", "Comma-separated gRPC addresses of other nodes")
	cmd.Flags().StringVar(&preParamsPath, "pre-params-path", "", "Path to pre-computed Paillier parameters JSON file")
	cmd.Flags().StringVar(&shardOutputPath, "shard-output", "", "Output path for the encrypted key shard")
	cmd.Flags().StringVar(&passphrase, "passphrase", "", "Passphrase for encrypting the key shard")
	cobra.CheckErr(cmd.MarkFlagRequired("party-id"))
	cobra.CheckErr(cmd.MarkFlagRequired("all-party-ids"))

	return cmd
}

// newPreParamsCommand creates the `keygen pre-params` command.
func newPreParamsCommand(wallet *wallets.Keygener) *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:   "pre-params",
		Short: "Generate Paillier pre-computation parameters for MPC DKG (ETH only)",
		Long: `Generate the Paillier key pre-computation parameters required for the DKG ceremony.
This is a CPU-intensive operation that may take several minutes.
The output file is required for the 'keygen dkg' command.`,
		Example: "  keygen --coin eth pre-params --output pre_params.json",
		RunE: func(cmd *cobra.Command, args []string) error {
			if wallet == nil || *wallet == nil {
				return errors.New("wallet not initialized, check --coin flag")
			}
			ethKeygen, ok := (*wallet).(*ethwallet.ETHKeygen)
			if !ok {
				fmt.Println("this command is only supported for ETH")
				return nil
			}
			if output == "" {
				return errors.New("--output flag is required")
			}
			if err := ethKeygen.GeneratePreParams(cmd.Context(), output); err != nil {
				return fmt.Errorf("failed to generate pre-params: %w", err)
			}
			fmt.Printf("Pre-computation parameters generated\n")
			fmt.Printf("  Output: %s\n", output)
			return nil
		},
	}

	cmd.Flags().StringVar(&output, "output", "", "Output path for the pre-params JSON file (required)")
	cobra.CheckErr(cmd.MarkFlagRequired("output"))

	return cmd
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
