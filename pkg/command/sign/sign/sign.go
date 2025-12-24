package sign

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/pkg/di"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// AddCommands adds all sign subcommands
func AddCommands(parentCmd *cobra.Command, wallet *wallets.Signer, container di.Container) {
	// signature command
	var signatureFile string
	signatureCmd := &cobra.Command{
		Use:   "signature",
		Short: "sign on signed transaction for multsig address (account would be found from file name)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSignature(*wallet, signatureFile)
		},
	}
	signatureCmd.Flags().StringVar(&signatureFile, "file", "", "import file path for signed transactions")
	parentCmd.AddCommand(signatureCmd)
}

func runSignature(wallet wallets.Signer, filePath string) error {
	fmt.Println("sign on signed transaction for multsig address")

	// validator
	if filePath == "" {
		return errors.New("file path option [-file] is required")
	}

	// sign on signed transactions
	hexTx, isSigned, generatedFileName, err := wallet.SignTx(filePath)
	if err != nil {
		return fmt.Errorf("fail to call SignTx() %w", err)
	}

	// TODO: output should be json if json option is true
	fmt.Printf("[hex]: %s\n[isCompleted]: %t\n[fileName]: %s\n", hexTx, isSigned, generatedFileName)

	return nil
}
