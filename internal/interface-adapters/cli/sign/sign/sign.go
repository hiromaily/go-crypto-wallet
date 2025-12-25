package sign

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	signusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign"
	"github.com/hiromaily/go-crypto-wallet/internal/di"
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
			return runSignature(container, signatureFile)
		},
	}
	signatureCmd.Flags().StringVar(&signatureFile, "file", "", "import file path for signed transactions")
	parentCmd.AddCommand(signatureCmd)
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
	fmt.Printf("[hex]: %s\n[isCompleted]: %t\n[fileName]: %s\n",
		output.SignedHex, output.IsComplete, output.NextFilePath)

	return nil
}
