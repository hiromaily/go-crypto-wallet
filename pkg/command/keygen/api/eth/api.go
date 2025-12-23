package eth

import (
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp"
)

// AddCommands adds all Ethereum API subcommands
func AddCommands(parentCmd *cobra.Command, eth ethgrp.Ethereumer) {
	// importrawkey command
	var (
		importrawkeyPrivKey    string
		importrawkeyPassPhrase string
	)
	importrawkeyCmd := &cobra.Command{
		Use:   "importrawkey",
		Short: "import raw key",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runImportRawKey(eth, importrawkeyPrivKey, importrawkeyPassPhrase)
		},
	}
	importrawkeyCmd.Flags().StringVar(&importrawkeyPrivKey, "key", "", "private key")
	importrawkeyCmd.Flags().StringVar(&importrawkeyPassPhrase, "pass", "", "passphrase")
	parentCmd.AddCommand(importrawkeyCmd)
}
