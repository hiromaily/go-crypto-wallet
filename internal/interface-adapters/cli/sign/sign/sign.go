package sign

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	signusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign"
	"github.com/hiromaily/go-crypto-wallet/internal/di"
	wallets "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet"
	ethwallet "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet/eth"
)

// AddCommands adds all sign subcommands
func AddCommands(parentCmd *cobra.Command, wallet *wallets.Signer, containerGetter func() di.Container) {
	// signature command
	var (
		signatureFile string
		signerAddress string
	)
	signatureCmd := &cobra.Command{
		Use:   "signature",
		Short: "sign on signed transaction for multsig address (account would be found from file name)",
		RunE: func(cmd *cobra.Command, args []string) error {
			// For ETH Safe multisig signing, route through the wallet adapter so that
			// file-format detection and signer-address routing work correctly.
			// The DI container's NewSignTransactionUseCase only handles single-sig ETH (no-op).
			if signerAddress != "" {
				if eths, ok := (*wallet).(*ethwallet.ETHSign); ok {
					eths.SetSignerAddress(signerAddress)
					return runETHMultisigSignWalletSign(eths, signatureFile)
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

// runETHMultisigSignWalletSign delegates ETH Safe multisig signing to the ETHSign wallet adapter.
func runETHMultisigSignWalletSign(eths *ethwallet.ETHSign, filePath string) error {
	fmt.Println("sign on unsigned multisig transaction (sign wallet)")
	if filePath == "" {
		return errors.New("file path option [-file] is required")
	}
	outPath, isDone, _, err := eths.SignTx(filePath)
	if err != nil {
		return fmt.Errorf("fail to sign multisig transaction: %w", err)
	}
	signedCount := extractMultisigSignedCount(outPath)
	fmt.Printf("[isCompleted]: %t\n[fileName]: %s\n[signedCount]: %d\n[unsignedCount]: %d\n",
		isDone, outPath, signedCount, 0)
	return nil
}

// extractMultisigSignedCount parses the signed-count suffix from a multisig file path.
// File naming convention: {actionType}_multisig_{uuid}_{N}.json  →  N.
func extractMultisigSignedCount(path string) int {
	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	parts := strings.Split(base, "_")
	if len(parts) > 0 {
		if n, err := strconv.Atoi(parts[len(parts)-1]); err == nil {
			return n
		}
	}
	return 0
}

func runSignature(container di.Container, filePath string) error {
	fmt.Println("sign on signed transaction for multsig address")

	// validator
	if filePath == "" {
		return errors.New("file path option [-file] is required")
	}

	// sign on signed transactions
	useCase := container.NewSignTransactionUseCase()
	output, err := useCase.Sign(context.Background(), signusecase.SignTransactionInput{
		FilePath: filePath,
	})
	if err != nil {
		return fmt.Errorf("fail to sign transaction: %w", err)
	}

	// TODO: output should be json if json option is true
	fmt.Printf("[signedData]: %s\n[isCompleted]: %t\n[fileName]: %s\n",
		output.SignedData, output.IsComplete, output.NextFilePath)

	return nil
}
