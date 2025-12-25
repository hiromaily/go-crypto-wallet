package btc

import (
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin"
)

// AddCommands adds all Bitcoin API subcommands
func AddCommands(parentCmd *cobra.Command, btc bitcoin.Bitcoiner) {
	// encryptwallet command
	var encryptwalletPassphrase string
	encryptwalletCmd := &cobra.Command{
		Use:   "encryptwallet",
		Short: "encrypts the wallet with 'passphrase'",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEncryptWallet(btc, encryptwalletPassphrase)
		},
	}
	encryptwalletCmd.Flags().StringVar(&encryptwalletPassphrase, "passphrase", "", "passphrase")
	parentCmd.AddCommand(encryptwalletCmd)

	// walletpassphrase command
	var walletpassphrasePassphrase string
	walletpassphraseCmd := &cobra.Command{
		Use:   "walletpassphrase",
		Short: "stores the wallet decryption key in memory for 'timeout' seconds",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWalletPassphrase(btc, walletpassphrasePassphrase)
		},
	}
	walletpassphraseCmd.Flags().StringVar(&walletpassphrasePassphrase, "passphrase", "", "passphrase")
	parentCmd.AddCommand(walletpassphraseCmd)

	// walletpassphrasechange command
	var (
		walletpassphrasechangeOld string
		walletpassphrasechangeNew string
	)
	walletpassphrasechangeCmd := &cobra.Command{
		Use:   "walletpassphrasechange",
		Short: "changes the wallet passphrase from 'oldpassphrase' to 'newpassphrase'",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWalletPassphraseChange(btc, walletpassphrasechangeOld, walletpassphrasechangeNew)
		},
	}
	walletpassphrasechangeCmd.Flags().StringVar(&walletpassphrasechangeOld, "old", "", "old passphrase")
	walletpassphrasechangeCmd.Flags().StringVar(&walletpassphrasechangeNew, "new", "", "new passphrase")
	parentCmd.AddCommand(walletpassphrasechangeCmd)

	// walletlock command
	walletlockCmd := &cobra.Command{
		Use:   "walletlock",
		Short: "removes the wallet encryption key from memory, locking the wallet",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWalletLock(btc)
		},
	}
	parentCmd.AddCommand(walletlockCmd)

	// dumpwallet command
	var dumpwalletFile string
	dumpwalletCmd := &cobra.Command{
		Use:   "dumpwallet",
		Short: "dumps all wallet keys in a human-readable format to a server-side file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDumpWallet(btc, dumpwalletFile)
		},
	}
	dumpwalletCmd.Flags().StringVar(&dumpwalletFile, "file", "", "file name")
	parentCmd.AddCommand(dumpwalletCmd)

	// importwallet command
	var importwalletFile string
	importwalletCmd := &cobra.Command{
		Use:   "importwallet",
		Short: "Imports keys from a wallet dump file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runImportWallet(btc, importwalletFile)
		},
	}
	importwalletCmd.Flags().StringVar(&importwalletFile, "file", "", "file name")
	parentCmd.AddCommand(importwalletCmd)
}
