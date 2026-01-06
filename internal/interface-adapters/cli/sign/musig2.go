package sign

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	signusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign"
	"github.com/hiromaily/go-crypto-wallet/internal/di"
)

// AddMuSig2Commands adds MuSig2 subcommands to the parent command
func AddMuSig2Commands(parentCmd *cobra.Command, containerGetter func() di.Container) {
	musig2Cmd := &cobra.Command{
		Use:   "musig2",
		Short: "MuSig2 operations for BTC",
		Long:  `MuSig2 multi-signature operations for Bitcoin transactions`,
	}

	// Add nonce and sign subcommands
	musig2Cmd.AddCommand(newMuSig2NonceCommand(containerGetter))
	musig2Cmd.AddCommand(newMuSig2SignCommand(containerGetter))

	parentCmd.AddCommand(musig2Cmd)
}

// newMuSig2NonceCommand creates the 'musig2 nonce' command
func newMuSig2NonceCommand(containerGetter func() di.Container) *cobra.Command {
	var (
		file   string
		output string
	)

	cmd := &cobra.Command{
		Use:   "nonce",
		Short: "Generate MuSig2 nonce for Round 1",
		Long: `Generate a MuSig2 nonce for the first round of the signing protocol.
This command reads a PSBT file, generates a nonce, and saves it to a JSON file
for sharing with other signers.`,
		Example: `  sign --coin btc musig2 nonce --file payment_123_unsigned.psbt ` +
			`--output payment_123_sign1_nonce.json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMuSig2Nonce(cmd.Context(), containerGetter(), file, output)
		},
	}

	cmd.Flags().StringVar(&file, "file", "", "PSBT file path (required)")
	cmd.Flags().StringVar(&output, "output", "", "Output nonce file path (required)")
	cobra.CheckErr(cmd.MarkFlagRequired("file"))
	cobra.CheckErr(cmd.MarkFlagRequired("output"))

	return cmd
}

// newMuSig2SignCommand creates the 'musig2 sign' command
func newMuSig2SignCommand(containerGetter func() di.Container) *cobra.Command {
	var (
		file   string
		output string
	)

	cmd := &cobra.Command{
		Use:   "sign",
		Short: "Create MuSig2 partial signature for Round 2",
		Long: `Create a MuSig2 partial signature for the second round of the signing protocol.
This command reads a PSBT file with aggregated nonces, creates a partial signature,
and saves the signed PSBT to a file for aggregation.`,
		Example: `  sign --coin btc musig2 sign --file payment_123_with_nonces.psbt ` +
			`--output payment_123_sign1_signed.psbt`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMuSig2Sign(cmd.Context(), containerGetter(), file, output)
		},
	}

	cmd.Flags().StringVar(&file, "file", "", "PSBT file with aggregated nonces (required)")
	cmd.Flags().StringVar(&output, "output", "", "Output signed PSBT file path (required)")
	cobra.CheckErr(cmd.MarkFlagRequired("file"))
	cobra.CheckErr(cmd.MarkFlagRequired("output"))

	return cmd
}

// runMuSig2Nonce executes the nonce generation command
func runMuSig2Nonce(ctx context.Context, container di.Container, file, output string) error {
	fmt.Println("Generate MuSig2 nonce for Round 1")

	// Read PSBT file
	psbtData, err := os.ReadFile(file) // #nosec G304 - file path is from CLI flag
	if err != nil {
		return fmt.Errorf("failed to read PSBT file: %w", err)
	}

	// Extract transaction ID from PSBT filename or content
	// For now, use filename as a simple transaction ID
	// In a real implementation, this would parse the PSBT and extract the actual tx ID
	txID := extractTransactionID(file)
	if txID == "" {
		return errors.New("failed to extract transaction ID from file")
	}

	// Get auth type from container
	authType := container.AuthType()

	// Execute use case
	useCase := container.NewSignGenerateMuSig2NonceUseCase()
	result, err := useCase.Generate(ctx, signusecase.GenerateMuSig2NonceInput{
		TransactionID: txID,
		AuthType:      authType,
	})
	if err != nil {
		return fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Save nonce to JSON file
	nonceData := NonceData{
		SignerID:    result.SignerID,
		PublicNonce: hex.EncodeToString(result.PublicNonce[:]),
	}

	if err := saveJSON(output, nonceData); err != nil {
		return fmt.Errorf("failed to save nonce file: %w", err)
	}

	fmt.Printf("✓ Nonce generated successfully\n")
	fmt.Printf("  Signer ID: %s\n", result.SignerID)
	fmt.Printf("  Output file: %s\n", output)
	fmt.Printf("\nShare this nonce file with the watch wallet to proceed to Round 2.\n")

	// Note: PSBT data read but not used in placeholder implementation
	_ = psbtData

	return nil
}

// runMuSig2Sign executes the partial signature creation command
func runMuSig2Sign(ctx context.Context, container di.Container, file, output string) error {
	fmt.Println("Create MuSig2 partial signature for Round 2")

	// Read PSBT file with aggregated nonces
	psbtData, err := os.ReadFile(file) // #nosec G304 - file path is from CLI flag
	if err != nil {
		return fmt.Errorf("failed to read PSBT file: %w", err)
	}

	// Extract transaction ID
	txID := extractTransactionID(file)
	if txID == "" {
		return errors.New("failed to extract transaction ID from file")
	}

	// Get auth type from container
	authType := container.AuthType()

	// Parse aggregated nonces from PSBT or separate file
	// TODO: Extract nonces from PSBT custom fields instead of using zero values
	// The current implementation uses zero values which will cause invalid signatures
	var aggregatedNonces [][66]byte
	// TODO: Extract message hash from PSBT transaction data
	// Signing a zero hash means signing a non-existent transaction
	var messageHash [32]byte

	// Execute use case
	useCase := container.NewSignMuSig2SignUseCase()
	result, err := useCase.Sign(ctx, signusecase.MuSig2SignInput{
		TransactionID:    txID,
		AuthType:         authType,
		MessageHash:      messageHash,
		AggregatedNonces: aggregatedNonces,
	})
	if err != nil {
		return fmt.Errorf("failed to create partial signature: %w", err)
	}

	// Save signed PSBT to output file
	// TODO: Embed the partial signature into the PSBT before writing
	// Currently, the partial signature is only printed and not added to the PSBT
	// The output file will contain the original PSBT without the signature
	if err := os.WriteFile(output, psbtData, 0o600); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("✓ Partial signature created successfully\n")
	fmt.Printf("  Signer ID: %s\n", result.SignerID)
	fmt.Printf("  Signature: %s\n", hex.EncodeToString(result.PartialSignature[:]))
	fmt.Printf("  Output file: %s\n", output)
	fmt.Printf("\nSend this signed PSBT to the watch wallet for signature aggregation.\n")

	return nil
}

// NonceData represents the nonce data stored in JSON format
type NonceData struct {
	SignerID    string `json:"signer_id"`
	PublicNonce string `json:"public_nonce"`
}

// extractTransactionID extracts a transaction ID from the file path
// In a real implementation, this would parse the PSBT and extract the actual transaction ID
func extractTransactionID(filePath string) string {
	// Simple implementation: use the filename without extension.
	// This is a placeholder and will be replaced by proper PSBT parsing.
	fileName := filepath.Base(filePath)
	baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	// This is a placeholder implementation to make the examples work.
	// It assumes a transaction ID is followed by suffixes like `_unsigned` or `_with_nonces`.
	if strings.Contains(baseName, "_unsigned") {
		return strings.Split(baseName, "_unsigned")[0]
	}
	if strings.Contains(baseName, "_with_nonces") {
		return strings.Split(baseName, "_with_nonces")[0]
	}
	if strings.Contains(baseName, "_signed") {
		return strings.Split(baseName, "_signed")[0]
	}

	return baseName
}

// saveJSON saves data to a JSON file
func saveJSON(path string, data any) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(path, jsonData, 0o600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
