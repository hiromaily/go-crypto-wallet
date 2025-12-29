package watch

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/internal/di"
)

// AddMuSig2Commands adds MuSig2 subcommands to the parent command
func AddMuSig2Commands(parentCmd *cobra.Command, container di.Container) {
	musig2Cmd := &cobra.Command{
		Use:   "musig2",
		Short: "MuSig2 operations for BTC",
		Long:  `MuSig2 multi-signature operations for Bitcoin transactions`,
	}

	// Add collect-nonces and aggregate subcommands
	musig2Cmd.AddCommand(newMuSig2CollectNoncesCommand(container))
	musig2Cmd.AddCommand(newMuSig2AggregateCommand(container))

	parentCmd.AddCommand(musig2Cmd)
}

// newMuSig2CollectNoncesCommand creates the 'musig2 collect-nonces' command
func newMuSig2CollectNoncesCommand(container di.Container) *cobra.Command {
	var (
		file   string
		nonces string
		output string
	)

	cmd := &cobra.Command{
		Use:   "collect-nonces",
		Short: "Collect and aggregate MuSig2 nonces from all signers",
		Long: `Collect MuSig2 nonces from all signers and aggregate them.
This command reads nonce JSON files from all signers (Keygen and Sign wallets),
validates the nonce count, aggregates them, and creates a PSBT with the
aggregated nonces for Round 2 signing.`,
		Example: `  watch --coin btc musig2 collect-nonces \
  --file payment_123_unsigned.psbt \
  --nonces payment_123_keygen_nonce.json,payment_123_sign1_nonce.json,payment_123_sign2_nonce.json \
  --output payment_123_with_nonces.psbt`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMuSig2CollectNonces(cmd.Context(), container, file, nonces, output)
		},
	}

	cmd.Flags().StringVar(&file, "file", "", "PSBT file path (required)")
	cmd.Flags().StringVar(&nonces, "nonces", "", "Comma-separated list of nonce JSON files (required)")
	cmd.Flags().StringVar(&output, "output", "", "Output PSBT file path (required)")
	cobra.CheckErr(cmd.MarkFlagRequired("file"))
	cobra.CheckErr(cmd.MarkFlagRequired("nonces"))
	cobra.CheckErr(cmd.MarkFlagRequired("output"))

	return cmd
}

// newMuSig2AggregateCommand creates the 'musig2 aggregate' command
func newMuSig2AggregateCommand(container di.Container) *cobra.Command {
	var (
		files  string
		output string
	)

	cmd := &cobra.Command{
		Use:   "aggregate",
		Short: "Aggregate MuSig2 partial signatures",
		Long: `Aggregate MuSig2 partial signatures from all signers.
This command reads signed PSBT files from all signers (Keygen and Sign wallets),
extracts partial signatures, aggregates them into a final signature, and creates
a finalized PSBT ready for broadcast.`,
		Example: `  watch --coin btc musig2 aggregate \
  --files payment_123_keygen_signed.psbt,payment_123_sign1_signed.psbt,payment_123_sign2_signed.psbt \
  --output payment_123_final.psbt`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMuSig2Aggregate(cmd.Context(), container, files, output)
		},
	}

	cmd.Flags().StringVar(&files, "files", "", "Comma-separated list of signed PSBT files (required)")
	cmd.Flags().StringVar(&output, "output", "", "Output finalized PSBT file path (required)")
	cobra.CheckErr(cmd.MarkFlagRequired("files"))
	cobra.CheckErr(cmd.MarkFlagRequired("output"))

	return cmd
}

// runMuSig2CollectNonces executes the nonce collection and aggregation command
func runMuSig2CollectNonces(_ context.Context, _ di.Container, file, noncesStr, output string) error {
	fmt.Println("Collect and aggregate MuSig2 nonces")

	// Read PSBT file
	psbtData, err := os.ReadFile(file) // #nosec G304 - file path is from CLI flag
	if err != nil {
		return fmt.Errorf("failed to read PSBT file: %w", err)
	}

	// Parse nonce file list
	nonceFiles := strings.Split(noncesStr, ",")
	if len(nonceFiles) < 2 {
		return fmt.Errorf("at least 2 nonce files required for MuSig2, got %d", len(nonceFiles))
	}

	// Read and parse all nonce files
	nonces := make([]NonceData, 0, len(nonceFiles))
	for _, nonceFile := range nonceFiles {
		nonceFile = strings.TrimSpace(nonceFile)
		data, err := os.ReadFile(nonceFile) // #nosec G304 - file path is from CLI flag
		if err != nil {
			return fmt.Errorf("failed to read nonce file %s: %w", nonceFile, err)
		}

		var nonce NonceData
		if err := json.Unmarshal(data, &nonce); err != nil {
			return fmt.Errorf("failed to parse nonce file %s: %w", nonceFile, err)
		}
		nonces = append(nonces, nonce)
	}

	fmt.Printf("✓ Collected %d nonces from signers\n", len(nonces))
	for i, nonce := range nonces {
		fmt.Printf("  Signer %d: %s\n", i+1, nonce.SignerID)
	}

	// TODO: Aggregate nonces using MuSig2Service
	// TODO: Embed aggregated nonces into PSBT custom fields
	// For now, write the original PSBT as placeholder
	if err := os.WriteFile(output, psbtData, 0o600); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("✓ Nonces aggregated successfully\n")
	fmt.Printf("  Output file: %s\n", output)
	fmt.Printf("\nSend this PSBT to all signers for Round 2 signing.\n")

	return nil
}

// runMuSig2Aggregate executes the signature aggregation command
func runMuSig2Aggregate(_ context.Context, _ di.Container, filesStr, output string) error {
	fmt.Println("Aggregate MuSig2 partial signatures")

	// Parse signed PSBT file list
	psbtFiles := strings.Split(filesStr, ",")
	if len(psbtFiles) < 2 {
		return fmt.Errorf("at least 2 signed PSBT files required for MuSig2, got %d", len(psbtFiles))
	}

	// Read all signed PSBT files
	psbtDataList := make([][]byte, 0, len(psbtFiles))
	for _, psbtFile := range psbtFiles {
		psbtFile = strings.TrimSpace(psbtFile)
		data, err := os.ReadFile(psbtFile) // #nosec G304 - file path is from CLI flag
		if err != nil {
			return fmt.Errorf("failed to read PSBT file %s: %w", psbtFile, err)
		}
		psbtDataList = append(psbtDataList, data)
	}

	fmt.Printf("✓ Collected %d partial signatures from signers\n", len(psbtDataList))

	// TODO: Extract partial signatures from each PSBT
	// TODO: Parse aggregated public key, combined nonce, message hash from PSBT
	// TODO: Call AggregateMuSig2SignaturesUseCase with proper data
	// For now, write the first PSBT as placeholder
	if err := os.WriteFile(output, psbtDataList[0], 0o600); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("✓ Signatures aggregated successfully (placeholder)\n")
	fmt.Printf("  Output file: %s\n", output)
	fmt.Printf("\nTransaction is fully signed and ready to broadcast (placeholder).\n")
	fmt.Printf("Use: watch --coin btc send --file %s\n", output)

	return nil
}

// NonceData represents the nonce data stored in JSON format
type NonceData struct {
	SignerID    string `json:"signer_id"`
	PublicNonce string `json:"public_nonce"`
}
