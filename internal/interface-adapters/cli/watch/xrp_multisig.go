package watch

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	"github.com/hiromaily/go-crypto-wallet/internal/di"
)

// AddXRPMultisigCommands adds XRP multi-signature and regular key subcommands
func AddXRPMultisigCommands(parentCmd *cobra.Command, containerGetter func() di.Container) {
	xrpCmd := &cobra.Command{
		Use:   "xrp",
		Short: "XRP advanced operations",
		Long:  `XRP multi-signature and regular key operations`,
	}

	// Add subcommands
	xrpCmd.AddCommand(newSetRegularKeyCommand(containerGetter))
	xrpCmd.AddCommand(newSetSignerListCommand(containerGetter))
	xrpCmd.AddCommand(newCreateMultisigTxCommand(containerGetter))
	xrpCmd.AddCommand(newAddMultisigSignatureCommand(containerGetter))
	xrpCmd.AddCommand(newSubmitMultisigTxCommand(containerGetter))

	parentCmd.AddCommand(xrpCmd)
}

// newSetRegularKeyCommand creates the 'xrp set-regular-key' command
func newSetRegularKeyCommand(containerGetter func() di.Container) *cobra.Command {
	var (
		account    string
		regularKey string
	)

	cmd := &cobra.Command{
		Use:   "set-regular-key",
		Short: "Set or remove a regular key for an XRP account",
		Long: `Configure a regular key for an XRP account.
A regular key allows signing transactions without exposing the master key.
Leave --regular-key empty to remove the current regular key.`,
		Example: "  watch --coin xrp set-regular-key --account rAddr --regular-key rKeyAddr",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSetRegularKey(cmd.Context(), containerGetter(), account, regularKey)
		},
	}

	cmd.Flags().StringVar(&account, "account", "", "XRP account address (required)")
	cmd.Flags().StringVar(&regularKey, "regular-key", "", "Regular key address (empty to remove)")
	cobra.CheckErr(cmd.MarkFlagRequired("account"))

	return cmd
}

// newSetSignerListCommand creates the 'xrp set-signer-list' command
func newSetSignerListCommand(containerGetter func() di.Container) *cobra.Command {
	var (
		account string
		signers string
		quorum  uint32
	)

	cmd := &cobra.Command{
		Use:   "set-signer-list",
		Short: "Configure multi-signature signer list for an XRP account",
		Long: `Configure a signer list for multi-signature transactions.
The signer list specifies which accounts can sign and their weights.
The quorum is the minimum total weight required to authorize a transaction.`,
		Example: "  watch --coin xrp set-signer-list --account rAddr --signers \"r1:1,r2:1\" --quorum",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSetSignerList(cmd.Context(), containerGetter(), account, signers, quorum)
		},
	}

	cmd.Flags().StringVar(&account, "account", "", "XRP account address (required)")
	cmd.Flags().StringVar(&signers, "signers", "", "Signer list as 'address:weight,...' (required)")
	cmd.Flags().Uint32Var(&quorum, "quorum", 0, "Required quorum weight (required)")
	cobra.CheckErr(cmd.MarkFlagRequired("account"))
	cobra.CheckErr(cmd.MarkFlagRequired("signers"))
	cobra.CheckErr(cmd.MarkFlagRequired("quorum"))

	return cmd
}

// newCreateMultisigTxCommand creates the 'xrp create-multisig-tx' command
func newCreateMultisigTxCommand(containerGetter func() di.Container) *cobra.Command {
	var (
		account  string
		receiver string
		amount   float64
		txType   string
	)

	cmd := &cobra.Command{
		Use:   "create-multisig-tx",
		Short: "Create a pending multi-signature transaction",
		Long: `Create a pending multi-signature transaction that requires signatures
from authorized signers. The transaction will be stored and assigned a UUID
for tracking signatures.`,
		Example: "  watch --coin xrp create-multisig-tx --account rAddr --receiver rRecv --amount 100",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreateMultisigTx(cmd.Context(), containerGetter(), account, receiver, amount, txType)
		},
	}

	cmd.Flags().StringVar(&account, "account", "", "Multi-sig account address (required)")
	cmd.Flags().StringVar(&receiver, "receiver", "", "Receiver address (required)")
	cmd.Flags().Float64Var(&amount, "amount", 0, "Amount in XRP (required)")
	cmd.Flags().StringVar(&txType, "type", "Payment", "Transaction type (default: Payment)")
	cobra.CheckErr(cmd.MarkFlagRequired("account"))
	cobra.CheckErr(cmd.MarkFlagRequired("receiver"))
	cobra.CheckErr(cmd.MarkFlagRequired("amount"))

	return cmd
}

// newAddMultisigSignatureCommand creates the 'xrp add-multisig-signature' command
func newAddMultisigSignatureCommand(containerGetter func() di.Container) *cobra.Command {
	var (
		txUUID       string
		signer       string
		signedTxBlob string
	)

	cmd := &cobra.Command{
		Use:   "add-multisig-signature",
		Short: "Add a signature to a pending multi-signature transaction",
		Long: `Add a signature from an authorized signer to a pending multi-signature
transaction. When the quorum is met, signatures will be automatically combined.`,
		Example: "  watch --coin xrp add-multisig-signature --tx-uuid UUID --signer rAddr --signed-tx-blob BLOB",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAddMultisigSignature(cmd.Context(), containerGetter(), txUUID, signer, signedTxBlob)
		},
	}

	cmd.Flags().StringVar(&txUUID, "tx-uuid", "", "Transaction UUID (required)")
	cmd.Flags().StringVar(&signer, "signer", "", "Signer account address (required)")
	cmd.Flags().StringVar(&signedTxBlob, "signed-tx-blob", "", "Signed transaction blob (required)")
	cobra.CheckErr(cmd.MarkFlagRequired("tx-uuid"))
	cobra.CheckErr(cmd.MarkFlagRequired("signer"))
	cobra.CheckErr(cmd.MarkFlagRequired("signed-tx-blob"))

	return cmd
}

// newSubmitMultisigTxCommand creates the 'xrp submit-multisig-tx' command
func newSubmitMultisigTxCommand(containerGetter func() di.Container) *cobra.Command {
	var txUUID string

	cmd := &cobra.Command{
		Use:   "submit-multisig-tx",
		Short: "Submit a ready multi-signature transaction to the network",
		Long: `Submit a multi-signature transaction that has met its quorum requirement.
The combined transaction will be submitted to the XRP Ledger.`,
		Example: "  watch --coin xrp submit-multisig-tx \\\n    --tx-uuid \"550e8400-e29b-41d4-a716-446655440000\"",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSubmitMultisigTx(cmd.Context(), containerGetter(), txUUID)
		},
	}

	cmd.Flags().StringVar(&txUUID, "tx-uuid", "", "Transaction UUID (required)")
	cobra.CheckErr(cmd.MarkFlagRequired("tx-uuid"))

	return cmd
}

// runSetRegularKey executes the set-regular-key command
func runSetRegularKey(ctx context.Context, container di.Container, account, regularKey string) error {
	fmt.Println("Setting regular key for XRP account")

	useCase := container.NewXRPWatchSetRegularKeyUseCase()
	output, err := useCase.Execute(ctx, watchusecase.SetRegularKeyInput{
		AccountAddress:    account,
		RegularKeyAddress: regularKey,
	})
	if err != nil {
		return fmt.Errorf("failed to set regular key: %w", err)
	}

	if regularKey == "" {
		fmt.Println("Regular key removal transaction created")
	} else {
		fmt.Println("Regular key set transaction created")
		fmt.Printf("  Regular Key: %s\n", regularKey)
	}
	fmt.Printf("  Account: %s\n", output.AccountID)
	fmt.Printf("  Transaction File: %s\n", output.FileName)
	fmt.Println("\nSign this transaction with the master key and broadcast it.")

	return nil
}

// runSetSignerList executes the set-signer-list command
func runSetSignerList(
	ctx context.Context,
	container di.Container,
	account, signersStr string,
	quorum uint32,
) error {
	fmt.Println("Setting signer list for XRP account")

	// Parse signers from "address:weight,..." format
	signerEntries, err := parseSignerEntries(signersStr)
	if err != nil {
		return fmt.Errorf("invalid signers format: %w", err)
	}

	useCase := container.NewXRPWatchSetSignerListUseCase()
	output, err := useCase.Execute(ctx, watchusecase.SetSignerListInput{
		AccountAddress: account,
		SignerQuorum:   quorum,
		SignerEntries:  signerEntries,
	})
	if err != nil {
		return fmt.Errorf("failed to set signer list: %w", err)
	}

	fmt.Println("SignerListSet transaction created")
	fmt.Printf("  Account: %s\n", account)
	fmt.Printf("  Quorum: %d\n", quorum)
	fmt.Printf("  Signer List ID: %d\n", output.SignerListID)
	fmt.Printf("  Transaction File: %s\n", output.FileName)
	fmt.Println("\nSigners:")
	for _, entry := range signerEntries {
		fmt.Printf("  - %s (weight: %d)\n", entry.Account, entry.Weight)
	}
	fmt.Println("\nSign this transaction with the master key and broadcast it.")

	return nil
}

// runCreateMultisigTx executes the create-multisig-tx command
func runCreateMultisigTx(
	ctx context.Context,
	container di.Container,
	account, receiver string,
	amount float64,
	txType string,
) error {
	fmt.Println("Creating multi-signature transaction")

	useCase := container.NewXRPWatchCreateMultisigTxUseCase()
	output, err := useCase.Execute(ctx, watchusecase.CreateMultisigTxInput{
		AccountAddress:  account,
		ReceiverAddress: receiver,
		Amount:          amount,
		TxType:          txType,
	})
	if err != nil {
		return fmt.Errorf("failed to create multisig transaction: %w", err)
	}

	fmt.Println("Multi-signature transaction created")
	fmt.Printf("  TX UUID: %s\n", output.TxUUID)
	fmt.Printf("  Pending ID: %d\n", output.PendingID)
	fmt.Printf("  Required Quorum: %d\n", output.RequiredQuorum)
	fmt.Println("\nUnsigned Transaction JSON:")
	fmt.Println(output.UnsignedTxJSON)
	fmt.Println("\nHave each authorized signer sign this transaction and submit signatures.")

	return nil
}

// runAddMultisigSignature executes the add-multisig-signature command
func runAddMultisigSignature(
	ctx context.Context,
	container di.Container,
	txUUID, signer, signedTxBlob string,
) error {
	fmt.Println("Adding signature to multi-signature transaction")

	useCase := container.NewXRPWatchAddMultisigSignatureUseCase()
	output, err := useCase.Execute(ctx, watchusecase.AddMultisigSignatureInput{
		TxUUID:        txUUID,
		SignerAccount: signer,
		SignedTxBlob:  signedTxBlob,
	})
	if err != nil {
		return fmt.Errorf("failed to add signature: %w", err)
	}

	fmt.Println("Signature added successfully")
	fmt.Printf("  Current Weight: %d / %d (required)\n", output.CurrentWeight, output.RequiredQuorum)

	if output.IsReady {
		fmt.Println("\nQuorum met! Transaction is ready for submission.")
		fmt.Printf("  Combined TX Blob: %s\n", output.CombinedTxBlob)
		fmt.Printf("\nSubmit with: watch --coin xrp xrp submit-multisig-tx --tx-uuid %s\n", txUUID)
	} else {
		remaining := output.RequiredQuorum - output.CurrentWeight
		fmt.Printf("\nNeed %d more weight to meet quorum.\n", remaining)
	}

	return nil
}

// runSubmitMultisigTx executes the submit-multisig-tx command
func runSubmitMultisigTx(ctx context.Context, container di.Container, txUUID string) error {
	fmt.Println("Submitting multi-signature transaction")

	useCase := container.NewXRPWatchSubmitMultisigTxUseCase()
	output, err := useCase.Execute(ctx, watchusecase.SubmitMultisigTxInput{
		TxUUID: txUUID,
	})
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	fmt.Println("Transaction submitted successfully")
	fmt.Printf("  TX Hash: %s\n", output.TxHash)
	if output.IsQueued {
		fmt.Println("  Status: Queued for validation")
	} else {
		fmt.Println("  Status: Submitted (check ledger for confirmation)")
	}

	return nil
}

// parseSignerEntries parses "address:weight,address:weight,..." format
func parseSignerEntries(signersStr string) ([]watchusecase.SignerEntry, error) {
	parts := strings.Split(signersStr, ",")
	entries := make([]watchusecase.SignerEntry, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		kv := strings.Split(part, ":")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid signer format '%s', expected 'address:weight'", part)
		}

		address := strings.TrimSpace(kv[0])
		weightStr := strings.TrimSpace(kv[1])

		weight, err := strconv.ParseUint(weightStr, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid weight '%s' for signer %s: %w", weightStr, address, err)
		}

		entries = append(entries, watchusecase.SignerEntry{
			Account: address,
			Weight:  uint32(weight),
		})
	}

	if len(entries) == 0 {
		return nil, errors.New("no valid signers provided")
	}

	return entries, nil
}
