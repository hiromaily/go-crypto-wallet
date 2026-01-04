package btc_test

import (
	"testing"

	"github.com/spf13/cobra"

	btckeygen "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/keygen/btc"
)

func TestAddDescriptorCommands(t *testing.T) {
	t.Parallel()

	rootCmd := &cobra.Command{Use: "test"}
	btckeygen.AddDescriptorCommands(rootCmd, nil)

	descCmd, _, err := rootCmd.Find([]string{"descriptor"})
	if err != nil {
		t.Fatalf("descriptor command not found: %v", err)
	}

	if !descCmd.HasSubCommands() {
		t.Fatal("descriptor command should have subcommands")
	}

	for _, sub := range []string{"generate", "export", "export-all"} {
		if cmd, _, err := descCmd.Find([]string{sub}); err != nil || cmd == nil {
			t.Fatalf("%s subcommand not found", sub)
		}
	}
}

func TestDescriptorGenerateFlags(t *testing.T) {
	t.Parallel()

	rootCmd := &cobra.Command{Use: "test"}
	btckeygen.AddDescriptorCommands(rootCmd, nil)

	generateCmd, _, err := rootCmd.Find([]string{"descriptor", "generate"})
	if err != nil {
		t.Fatalf("generate subcommand not found: %v", err)
	}

	for _, flag := range []string{"account", "address-type", "change", "required-sigs"} {
		if generateCmd.Flags().Lookup(flag) == nil {
			t.Fatalf("flag --%s not found on generate command", flag)
		}
	}
}

func TestDescriptorExportFlags(t *testing.T) {
	t.Parallel()

	rootCmd := &cobra.Command{Use: "test"}
	btckeygen.AddDescriptorCommands(rootCmd, nil)

	exportCmd, _, err := rootCmd.Find([]string{"descriptor", "export"})
	if err != nil {
		t.Fatalf("export subcommand not found: %v", err)
	}

	for _, flag := range []string{"account", "output", "format", "include-change"} {
		if exportCmd.Flags().Lookup(flag) == nil {
			t.Fatalf("flag --%s not found on export command", flag)
		}
	}
}

func TestDescriptorExportAllFlags(t *testing.T) {
	t.Parallel()

	rootCmd := &cobra.Command{Use: "test"}
	btckeygen.AddDescriptorCommands(rootCmd, nil)

	exportAllCmd, _, err := rootCmd.Find([]string{"descriptor", "export-all"})
	if err != nil {
		t.Fatalf("export-all subcommand not found: %v", err)
	}

	for _, flag := range []string{"account", "output", "format", "include-change"} {
		if exportAllCmd.Flags().Lookup(flag) == nil {
			t.Fatalf("flag --%s not found on export-all command", flag)
		}
	}
}
