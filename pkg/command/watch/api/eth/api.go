package eth

import (
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum"
)

// AddCommands adds all Ethereum API subcommands
func AddCommands(parentCmd *cobra.Command, eth ethereum.Ethereumer) {
	// clientversion command
	clientversionCmd := &cobra.Command{
		Use:   "clientversion",
		Short: "network version",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runClientVersion(eth)
		},
	}
	parentCmd.AddCommand(clientversionCmd)

	// nodeinfo command
	nodeinfoCmd := &cobra.Command{
		Use:   "nodeinfo",
		Short: "node info",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runNodeInfo(eth)
		},
	}
	parentCmd.AddCommand(nodeinfoCmd)

	// syncing command
	syncingCmd := &cobra.Command{
		Use:   "syncing",
		Short: "sync info",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSyncing(eth)
		},
	}
	parentCmd.AddCommand(syncingCmd)

	// netversion command
	netversionCmd := &cobra.Command{
		Use:   "netversion",
		Short: "network version",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runNetVersion(eth)
		},
	}
	parentCmd.AddCommand(netversionCmd)
}
