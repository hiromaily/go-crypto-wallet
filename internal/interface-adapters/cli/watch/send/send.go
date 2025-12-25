package send

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	"github.com/hiromaily/go-crypto-wallet/internal/di"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// AddCommand creates and returns the send command
func AddCommand(wallet *wallets.Watcher, container di.Container) *cobra.Command {
	var filePath string

	cmd := &cobra.Command{
		Use:   "send",
		Short: "send signed transaction to blockchain network",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSend(container, filePath)
		},
	}
	cmd.Flags().StringVar(&filePath, "file", "", "signed transaction file path")

	return cmd
}

func runSend(container di.Container, filePath string) error {
	// validator
	if filePath == "" {
		return errors.New("file path option [-file] is required")
	}

	// Get use case from container
	useCase := container.NewWatchSendTransactionUseCase().(watchusecase.SendTransactionUseCase)

	// send signed transactions
	output, err := useCase.Execute(context.Background(), watchusecase.SendTransactionInput{
		FilePath: filePath,
	})
	if err != nil {
		return fmt.Errorf("fail to send transaction: %w", err)
	}

	// TODO: output should be json if json option is true
	fmt.Println("tx is sent!! txID: " + output.TxID)

	return nil
}
