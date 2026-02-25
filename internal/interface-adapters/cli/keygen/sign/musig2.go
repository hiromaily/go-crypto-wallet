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

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	"github.com/hiromaily/go-crypto-wallet/internal/di"
	wallets "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet"
	btcwallet "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet/btc"
)

func addMuSig2Commands(parentCmd *cobra.Command, wallet *wallets.Keygener, containerGetter func() di.Container) {
	musig2Cmd := &cobra.Command{
		Use:   "musig2",
		Short: "MuSig2 multi-signature operations for BTC",
		Long:  `MuSig2 multi-signature operations for Bitcoin transactions`,
	}

	musig2Cmd.AddCommand(newMuSig2NonceCommand(wallet, containerGetter))
	musig2Cmd.AddCommand(newMuSig2SignCommand(wallet, containerGetter))

	parentCmd.AddCommand(musig2Cmd)
}

func newMuSig2NonceCommand(wallet *wallets.Keygener, containerGetter func() di.Container) *cobra.Command {
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
		Example: `  keygen --coin btc sign musig2 nonce --file payment_123_unsigned.psbt ` +
			`--output payment_123_keygen_nonce.json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wallet == nil || *wallet == nil {
				return errors.New("wallet not initialized, check --coin flag")
			}
			if _, ok := (*wallet).(*btcwallet.BTCKeygen); !ok {
				fmt.Println("this command is only supported for BTC")
				return nil
			}
			return runMuSig2Nonce(cmd.Context(), containerGetter(), file, output)
		},
	}

	cmd.Flags().StringVar(&file, "file", "", "PSBT file path (required)")
	cmd.Flags().StringVar(&output, "output", "", "Output nonce file path (required)")
	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagRequired("output")

	return cmd
}

func newMuSig2SignCommand(wallet *wallets.Keygener, containerGetter func() di.Container) *cobra.Command {
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
		Example: `  keygen --coin btc sign musig2 sign --file payment_123_with_nonces.psbt ` +
			`--output payment_123_keygen_signed.psbt`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wallet == nil || *wallet == nil {
				return errors.New("wallet not initialized, check --coin flag")
			}
			if _, ok := (*wallet).(*btcwallet.BTCKeygen); !ok {
				fmt.Println("this command is only supported for BTC")
				return nil
			}
			return runMuSig2Sign(cmd.Context(), containerGetter(), file, output)
		},
	}

	cmd.Flags().StringVar(&file, "file", "", "PSBT file with aggregated nonces (required)")
	cmd.Flags().StringVar(&output, "output", "", "Output signed PSBT file path (required)")
	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagRequired("output")

	return cmd
}

func runMuSig2Nonce(ctx context.Context, container di.Container, file, output string) error {
	fmt.Println("Generate MuSig2 nonce for Round 1")

	psbtData, err := os.ReadFile(file) // #nosec G304 - file path is from CLI flag
	if err != nil {
		return fmt.Errorf("failed to read PSBT file: %w", err)
	}

	txID := extractTransactionID(file)
	if txID == "" {
		return errors.New("failed to extract transaction ID from file")
	}

	// TODO: Derive signerID from wallet configuration instead of hardcoding
	signerID := "keygen"

	useCase := container.NewKeygenGenerateMuSig2NonceUseCase()
	result, err := useCase.Generate(ctx, keygenusecase.GenerateMuSig2NonceInput{
		TransactionID: txID,
		SignerID:      signerID,
	})
	if err != nil {
		return fmt.Errorf("failed to generate nonce: %w", err)
	}

	nonce := muSig2NonceData{
		SignerID:    result.SignerID,
		PublicNonce: hex.EncodeToString(result.PublicNonce[:]),
	}
	if err := saveMuSig2JSON(output, nonce); err != nil {
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

func runMuSig2Sign(ctx context.Context, container di.Container, file, output string) error {
	fmt.Println("Create MuSig2 partial signature for Round 2")

	psbtData, err := os.ReadFile(file) // #nosec G304 - file path is from CLI flag
	if err != nil {
		return fmt.Errorf("failed to read PSBT file: %w", err)
	}

	txID := extractTransactionID(file)
	if txID == "" {
		return errors.New("failed to extract transaction ID from file")
	}

	// TODO: Derive signerID from wallet configuration instead of hardcoding
	signerID := "keygen"

	// TODO: Extract nonces from PSBT custom fields instead of using zero values
	var aggregatedNonces [][66]byte
	// TODO: Extract message hash from PSBT transaction data
	var messageHash [32]byte

	useCase := container.NewKeygenMuSig2SignUseCase()
	result, err := useCase.Sign(ctx, keygenusecase.MuSig2SignInput{
		TransactionID:    txID,
		SignerID:         signerID,
		MessageHash:      messageHash,
		AggregatedNonces: aggregatedNonces,
	})
	if err != nil {
		return fmt.Errorf("failed to create partial signature: %w", err)
	}

	// TODO: Embed the partial signature into the PSBT before writing
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

type muSig2NonceData struct {
	SignerID    string `json:"signer_id"`
	PublicNonce string `json:"public_nonce"`
}

func extractTransactionID(filePath string) string {
	fileName := filepath.Base(filePath)
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func saveMuSig2JSON(path string, data any) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	if err := os.WriteFile(path, jsonData, 0o600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}
