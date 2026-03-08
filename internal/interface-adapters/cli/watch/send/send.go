package send

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	"github.com/hiromaily/go-crypto-wallet/internal/di"
	wallets "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet"
)

// AddCommand creates and returns the send command with tx and multisig subcommands
func AddCommand(wallet *wallets.Watcher, containerGetter func() di.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "send transactions to the blockchain network",
	}

	// tx target
	cmd.AddCommand(newTxCommand(containerGetter))

	// mpc target (ETH only: MPC/TSS threshold signing)
	cmd.AddCommand(newSendMPCCommand(wallet))

	// multisig target (BTC: MuSig2, XRP: multisig)
	addMultisigCommands(cmd, wallet, containerGetter)

	return cmd
}

func newTxCommand(containerGetter func() di.Container) *cobra.Command {
	var filePath string

	cmd := &cobra.Command{
		Use:   "tx",
		Short: "send signed transaction to blockchain network",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSend(containerGetter(), filePath)
		},
	}
	cmd.Flags().StringVar(&filePath, "file", "", "signed transaction file path")

	return cmd
}

func runSend(container di.Container, filePath string) error {
	if filePath == "" {
		return errors.New("file path option [-file] is required")
	}

	useCase := container.NewWatchSendTransactionUseCase().(watchusecase.SendTransactionUseCase)

	output, err := useCase.Execute(context.Background(), watchusecase.SendTransactionInput{
		FilePath: filePath,
	})
	if err != nil {
		return fmt.Errorf("fail to send transaction: %w", err)
	}

	fmt.Println("tx is sent!! txID: " + output.TxID)

	return nil
}
