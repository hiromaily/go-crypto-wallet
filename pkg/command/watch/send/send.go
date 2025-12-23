package send

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// AddCommand creates and returns the send command
func AddCommand(wallet *wallets.Watcher) *cobra.Command {
	var filePath string

	cmd := &cobra.Command{
		Use:   "send",
		Short: "send signed transaction to blockchain network",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSend(*wallet, filePath)
		},
	}
	cmd.Flags().StringVar(&filePath, "file", "", "signed transaction file path")

	return cmd
}

func runSend(wallet wallets.Watcher, filePath string) error {
	// validator
	if filePath == "" {
		return errors.New("file path option [-file] is required")
	}

	// send signed transactions
	txID, err := wallet.SendTx(filePath)
	if err != nil {
		return fmt.Errorf("fail to call SendTx() %w", err)
	}

	// TODO: output should be json if json option is true
	fmt.Println("tx is sent!! txID: " + txID)

	return nil
}
