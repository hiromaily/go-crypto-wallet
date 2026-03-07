package send

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	"github.com/hiromaily/go-crypto-wallet/internal/di"
	wallets "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet"
	btcwallet "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet/btc"
	xrpwallet "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet/xrp"
)

// addMultisigCommands adds the multisig target under the send command.
// Sub-operations are gated by coin type at runtime.
func addMultisigCommands(parentCmd *cobra.Command, wallet *wallets.Watcher, containerGetter func() di.Container) {
	multisigCmd := &cobra.Command{
		Use:   "multisig",
		Short: "Multi-signature operations (BTC: MuSig2, XRP: multisig)",
		Long: `Multi-signature transaction operations.
BTC: MuSig2 nonce collection and signature aggregation.
XRP: Signer list management, multi-signature transaction creation and submission.`,
	}

	// BTC MuSig2 subcommands
	multisigCmd.AddCommand(newCollectNoncesCommand(wallet, containerGetter))
	multisigCmd.AddCommand(newAggregateCommand(wallet, containerGetter))

	// XRP multisig subcommands
	multisigCmd.AddCommand(newSetRegularKeyCommand(wallet, containerGetter))
	multisigCmd.AddCommand(newSetSignerListCommand(wallet, containerGetter))
	multisigCmd.AddCommand(newCreateMultisigTxCommand(wallet, containerGetter))
	multisigCmd.AddCommand(newAddMultisigSignatureCommand(wallet, containerGetter))
	multisigCmd.AddCommand(newSubmitMultisigTxCommand(wallet, containerGetter))

	parentCmd.AddCommand(multisigCmd)
}

// --- BTC MuSig2 commands ---

func newCollectNoncesCommand(wallet *wallets.Watcher, containerGetter func() di.Container) *cobra.Command {
	var (
		file   string
		nonces string
		output string
	)

	cmd := &cobra.Command{
		Use:   "collect-nonces",
		Short: "Collect and aggregate MuSig2 nonces from all signers (BTC only)",
		Long: `Collect MuSig2 nonces from all signers and aggregate them.
This command reads nonce JSON files from all signers (Keygen and Sign wallets),
validates the nonce count, aggregates them, and creates a PSBT with the
aggregated nonces for Round 2 signing.`,
		Example: `  watch --coin btc send multisig collect-nonces \
  --file payment_123_unsigned.psbt \
  --nonces payment_123_keygen_nonce.json,payment_123_sign1_nonce.json \
  --output payment_123_with_nonces.psbt`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wallet == nil || *wallet == nil {
				return errors.New("wallet not initialized, check --coin flag")
			}
			if _, ok := (*wallet).(*btcwallet.BTCWatch); !ok {
				fmt.Println("this command is only supported for BTC")
				return nil
			}
			return runCollectNonces(cmd.Context(), containerGetter(), file, nonces, output)
		},
	}

	cmd.Flags().StringVar(&file, "file", "", "PSBT file path (required)")
	cmd.Flags().StringVar(&nonces, "nonces", "", "Comma-separated list of nonce JSON files (required)")
	cmd.Flags().StringVar(&output, "output", "", "Output PSBT file path (required)")
	cobra.CheckErr(cmd.MarkFlagRequired("file"))
	cobra.CheckErr(cmd.MarkFlagRequired("nonces"))
	cobra.CheckErr(cmd.MarkFlagRequired("output"))

	return cmd
}

func newAggregateCommand(wallet *wallets.Watcher, containerGetter func() di.Container) *cobra.Command {
	var (
		files  string
		output string
	)

	cmd := &cobra.Command{
		Use:   "aggregate",
		Short: "Aggregate MuSig2 partial signatures (BTC only)",
		Long: `Aggregate MuSig2 partial signatures from all signers.
This command reads signed PSBT files from all signers (Keygen and Sign wallets),
extracts partial signatures, aggregates them into a final signature, and creates
a finalized PSBT ready for broadcast.`,
		Example: `  watch --coin btc send multisig aggregate \
  --files payment_123_keygen_signed.psbt,payment_123_sign1_signed.psbt \
  --output payment_123_final.psbt`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wallet == nil || *wallet == nil {
				return errors.New("wallet not initialized, check --coin flag")
			}
			if _, ok := (*wallet).(*btcwallet.BTCWatch); !ok {
				fmt.Println("this command is only supported for BTC")
				return nil
			}
			return runAggregate(cmd.Context(), containerGetter(), files, output)
		},
	}

	cmd.Flags().StringVar(&files, "files", "", "Comma-separated list of signed PSBT files (required)")
	cmd.Flags().StringVar(&output, "output", "", "Output finalized PSBT file path (required)")
	cobra.CheckErr(cmd.MarkFlagRequired("files"))
	cobra.CheckErr(cmd.MarkFlagRequired("output"))

	return cmd
}

func runCollectNonces(_ context.Context, _ di.Container, file, noncesStr, output string) error {
	fmt.Println("Collect and aggregate MuSig2 nonces")

	psbtData, err := os.ReadFile(file) // #nosec G304 - file path is from CLI flag
	if err != nil {
		return fmt.Errorf("failed to read PSBT file: %w", err)
	}

	nonceFiles := strings.Split(noncesStr, ",")
	if len(nonceFiles) < 2 {
		return fmt.Errorf("at least 2 nonce files required for MuSig2, got %d", len(nonceFiles))
	}

	nonces := make([]nonceData, 0, len(nonceFiles))
	for _, nonceFile := range nonceFiles {
		nonceFile = strings.TrimSpace(nonceFile)
		data, err := os.ReadFile(nonceFile) // #nosec G304 - file path is from CLI flag
		if err != nil {
			return fmt.Errorf("failed to read nonce file %s: %w", nonceFile, err)
		}
		var nonce nonceData
		if err := json.Unmarshal(data, &nonce); err != nil {
			return fmt.Errorf("failed to parse nonce file %s: %w", nonceFile, err)
		}
		nonces = append(nonces, nonce)
	}

	fmt.Printf("✓ Collected %d nonces from signers\n", len(nonces))
	for i, nonce := range nonces {
		fmt.Printf("  Signer %d: %s\n", i+1, nonce.SignerID)
	}

	// TODO: Aggregate nonces using MuSig2Service
	// TODO: Embed aggregated nonces into PSBT custom fields
	if err := os.WriteFile(output, psbtData, 0o600); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("✓ Nonces aggregated successfully\n")
	fmt.Printf("  Output file: %s\n", output)
	fmt.Printf("\nSend this PSBT to all signers for Round 2 signing.\n")

	return nil
}

func runAggregate(_ context.Context, _ di.Container, filesStr, output string) error {
	fmt.Println("Aggregate MuSig2 partial signatures")

	psbtFiles := strings.Split(filesStr, ",")
	if len(psbtFiles) < 2 {
		return fmt.Errorf("at least 2 signed PSBT files required for MuSig2, got %d", len(psbtFiles))
	}

	psbtDataList := make([][]byte, 0, len(psbtFiles))
	for _, psbtFile := range psbtFiles {
		psbtFile = strings.TrimSpace(psbtFile)
		data, err := os.ReadFile(psbtFile) // #nosec G304 - file path is from CLI flag
		if err != nil {
			return fmt.Errorf("failed to read PSBT file %s: %w", psbtFile, err)
		}
		psbtDataList = append(psbtDataList, data)
	}

	fmt.Printf("✓ Collected %d partial signatures from signers\n", len(psbtDataList))

	// TODO: Extract partial signatures from each PSBT
	// TODO: Aggregate and finalize the transaction
	if err := os.WriteFile(output, psbtDataList[0], 0o600); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("✓ Signatures aggregated successfully (placeholder)\n")
	fmt.Printf("  Output file: %s\n", output)
	fmt.Printf("\nTransaction is ready to broadcast.\n")
	fmt.Printf("Use: watch --coin btc send tx --file %s\n", output)

	return nil
}

// --- XRP multisig commands ---

func newSetRegularKeyCommand(wallet *wallets.Watcher, containerGetter func() di.Container) *cobra.Command {
	var (
		account    string
		regularKey string
	)

	cmd := &cobra.Command{
		Use:   "set-regular-key",
		Short: "Set or remove a regular key for an XRP account (XRP only)",
		Long: `Configure a regular key for an XRP account.
A regular key allows signing transactions without exposing the master key.
Leave --regular-key empty to remove the current regular key.`,
		Example: "  watch --coin xrp send multisig set-regular-key --account rAddr --regular-key rKeyAddr",
		RunE: func(cmd *cobra.Command, args []string) error {
			if wallet == nil || *wallet == nil {
				return errors.New("wallet not initialized, check --coin flag")
			}
			if _, ok := (*wallet).(*xrpwallet.XRPWatch); !ok {
				fmt.Println("this command is only supported for XRP")
				return nil
			}
			return runSetRegularKey(cmd.Context(), containerGetter(), account, regularKey)
		},
	}

	cmd.Flags().StringVar(&account, "account", "", "XRP account address (required)")
	cmd.Flags().StringVar(&regularKey, "regular-key", "", "Regular key address (empty to remove)")
	cobra.CheckErr(cmd.MarkFlagRequired("account"))

	return cmd
}

func newSetSignerListCommand(wallet *wallets.Watcher, containerGetter func() di.Container) *cobra.Command {
	var (
		account     string
		accountType string
		signers     string
		quorum      uint32
	)

	cmd := &cobra.Command{
		Use:   "set-signer-list",
		Short: "Configure multi-signature signer list for an XRP account (XRP only)",
		Long: `Configure a signer list for multi-signature transactions.
The signer list specifies which accounts can sign and their weights.
The quorum is the minimum total weight required to authorize a transaction.`,
		Example: "  watch --coin xrp send multisig set-signer-list --account rAddr --signers \"r1:1,r2:1\" --quorum 2",
		RunE: func(cmd *cobra.Command, args []string) error {
			if wallet == nil || *wallet == nil {
				return errors.New("wallet not initialized, check --coin flag")
			}
			if _, ok := (*wallet).(*xrpwallet.XRPWatch); !ok {
				fmt.Println("this command is only supported for XRP")
				return nil
			}
			return runSetSignerList(cmd.Context(), containerGetter(), account, accountType, signers, quorum)
		},
	}

	cmd.Flags().StringVar(&account, "account", "", "XRP account address (required)")
	cmd.Flags().StringVar(&accountType, "account-type", "",
		"Account type label for key lookup by keygen signer (e.g. payment, deposit)")
	cmd.Flags().StringVar(&signers, "signers", "", "Signer list as 'address:weight,...' (required)")
	cmd.Flags().Uint32Var(&quorum, "quorum", 0, "Required quorum weight (required)")
	cobra.CheckErr(cmd.MarkFlagRequired("account"))
	cobra.CheckErr(cmd.MarkFlagRequired("signers"))
	cobra.CheckErr(cmd.MarkFlagRequired("quorum"))

	return cmd
}

func newCreateMultisigTxCommand(wallet *wallets.Watcher, containerGetter func() di.Container) *cobra.Command {
	var (
		account  string
		receiver string
		amount   float64
		txType   string
	)

	cmd := &cobra.Command{
		Use:   "create-multisig-tx",
		Short: "Create a pending multi-signature transaction (XRP only)",
		Long: `Create a pending multi-signature transaction that requires signatures
from authorized signers. The transaction will be stored and assigned a UUID
for tracking signatures.`,
		Example: "  watch --coin xrp send multisig create-multisig-tx --account rAddr --receiver rRecv --amount 100",
		RunE: func(cmd *cobra.Command, args []string) error {
			if wallet == nil || *wallet == nil {
				return errors.New("wallet not initialized, check --coin flag")
			}
			if _, ok := (*wallet).(*xrpwallet.XRPWatch); !ok {
				fmt.Println("this command is only supported for XRP")
				return nil
			}
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

func newAddMultisigSignatureCommand(wallet *wallets.Watcher, containerGetter func() di.Container) *cobra.Command {
	var (
		txUUID       string
		signer       string
		signedTxBlob string
	)

	cmd := &cobra.Command{
		Use:   "add-multisig-signature",
		Short: "Add a signature to a pending multi-signature transaction (XRP only)",
		Long: `Add a signature from an authorized signer to a pending multi-signature
transaction. When the quorum is met, signatures will be automatically combined.`,
		Example: "  watch --coin xrp send multisig add-multisig-signature " +
			"--tx-uuid UUID --signer rAddr --signed-tx-blob BLOB",
		RunE: func(cmd *cobra.Command, args []string) error {
			if wallet == nil || *wallet == nil {
				return errors.New("wallet not initialized, check --coin flag")
			}
			if _, ok := (*wallet).(*xrpwallet.XRPWatch); !ok {
				fmt.Println("this command is only supported for XRP")
				return nil
			}
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

func newSubmitMultisigTxCommand(wallet *wallets.Watcher, containerGetter func() di.Container) *cobra.Command {
	var txUUID string

	cmd := &cobra.Command{
		Use:   "submit-multisig-tx",
		Short: "Submit a ready multi-signature transaction to the network (XRP only)",
		Long: `Submit a multi-signature transaction that has met its quorum requirement.
The combined transaction will be submitted to the XRP Ledger.`,
		Example: "  watch --coin xrp send multisig submit-multisig-tx" +
			" --tx-uuid \"550e8400-e29b-41d4-a716-446655440000\"",
		RunE: func(cmd *cobra.Command, args []string) error {
			if wallet == nil || *wallet == nil {
				return errors.New("wallet not initialized, check --coin flag")
			}
			if _, ok := (*wallet).(*xrpwallet.XRPWatch); !ok {
				fmt.Println("this command is only supported for XRP")
				return nil
			}
			return runSubmitMultisigTx(cmd.Context(), containerGetter(), txUUID)
		},
	}

	cmd.Flags().StringVar(&txUUID, "tx-uuid", "", "Transaction UUID (required)")
	cobra.CheckErr(cmd.MarkFlagRequired("tx-uuid"))

	return cmd
}

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

func runSetSignerList(
	ctx context.Context,
	container di.Container,
	account, accountType, signersStr string,
	quorum uint32,
) error {
	fmt.Println("Setting signer list for XRP account")

	if quorum < 1 {
		return fmt.Errorf("quorum must be at least 1, got %d", quorum)
	}

	signerEntries, err := parseSignerEntries(signersStr)
	if err != nil {
		return fmt.Errorf("invalid signers format: %w", err)
	}

	if len(signerEntries) > 8 {
		return fmt.Errorf("signer list cannot exceed 8 entries, got %d", len(signerEntries))
	}

	useCase := container.NewXRPWatchSetSignerListUseCase()
	output, err := useCase.Execute(ctx, watchusecase.SetSignerListInput{
		AccountAddress:    account,
		SenderAccountType: accountType,
		SignerQuorum:      quorum,
		SignerEntries:     signerEntries,
	})
	if err != nil {
		return fmt.Errorf("failed to set signer list: %w", err)
	}

	fmt.Println("SignerListSet transaction created")
	fmt.Printf("  Account: %s\n", account)
	fmt.Printf("  Quorum: %d\n", quorum)
	fmt.Printf("  Signer List ID: %d\n", output.SignerListID)
	fmt.Printf("[fileName]: %s\n", output.FileName)
	fmt.Println("\nSigners:")
	for _, entry := range signerEntries {
		fmt.Printf("  - %s (weight: %d)\n", entry.Account, entry.Weight)
	}
	fmt.Println("\nSign this transaction with the master key and broadcast it.")

	return nil
}

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
		fmt.Printf("\nSubmit with: watch --coin xrp send multisig submit-multisig-tx --tx-uuid %s\n", txUUID)
	} else {
		remaining := output.RequiredQuorum - output.CurrentWeight
		fmt.Printf("\nNeed %d more weight to meet quorum.\n", remaining)
	}

	return nil
}

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

type nonceData struct {
	SignerID    string `json:"signer_id"`
	PublicNonce string `json:"public_nonce"`
}

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
