package sign

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

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
			// For ETH Safe multisig signing, route through the wallet adapter so that
			// file-format detection and signer-address routing work correctly.
			// The DI container's NewKeygenSignTransactionUseCase only handles single-sig ETH.
			if signerAddress != "" {
				if ethk, ok := (*wallet).(*ethwallet.ETHKeygen); ok {
					ethk.SetSignerAddress(signerAddress)
					return runETHMultisigKeygenSign(ethk, signatureFile)
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

// runETHMultisigKeygenSign delegates ETH Safe multisig signing to the wallet adapter.
// The adapter detects the multisig file format and routes to the multisig use case.
func runETHMultisigKeygenSign(ethk *ethwallet.ETHKeygen, filePath string) error {
	fmt.Println("sign on unsigned multisig transaction")
	if filePath == "" {
		return errors.New("file path option [-file] is required")
	}
	outPath, isDone, _, err := ethk.SignTx(filePath)
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
