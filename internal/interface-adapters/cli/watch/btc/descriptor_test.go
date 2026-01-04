package btc_test

import (
	"testing"

	"github.com/spf13/cobra"

	watchbtccmd "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/watch/btc"
)

func TestAddDescriptorCommands(t *testing.T) {
	t.Parallel()

	rootCmd := &cobra.Command{Use: "test"}
	watchbtccmd.AddDescriptorCommands(rootCmd, nil)

	descCmd, _, err := rootCmd.Find([]string{"descriptor"})
	if err != nil {
		t.Fatalf("descriptor command not found: %v", err)
	}

	if !descCmd.HasSubCommands() {
		t.Fatal("descriptor command should have subcommands")
	}

	for _, sub := range []string{"import", "validate"} {
		if cmd, _, err := descCmd.Find([]string{sub}); err != nil || cmd == nil {
			t.Fatalf("%s subcommand not found", sub)
		}
	}
}

func TestDescriptorImportFlags(t *testing.T) {
	t.Parallel()

	rootCmd := &cobra.Command{Use: "test"}
	watchbtccmd.AddDescriptorCommands(rootCmd, nil)

	importCmd, _, err := rootCmd.Find([]string{"descriptor", "import"})
	if err != nil {
		t.Fatalf("import subcommand not found: %v", err)
	}

	for _, flag := range []string{"file", "account", "start", "count"} {
		if importCmd.Flags().Lookup(flag) == nil {
			t.Fatalf("flag --%s not found on import command", flag)
		}
	}
}

func TestDescriptorValidateFlags(t *testing.T) {
	t.Parallel()

	rootCmd := &cobra.Command{Use: "test"}
	watchbtccmd.AddDescriptorCommands(rootCmd, nil)

	validateCmd, _, err := rootCmd.Find([]string{"descriptor", "validate"})
	if err != nil {
		t.Fatalf("validate subcommand not found: %v", err)
	}

	if validateCmd.Flags().Lookup("file") == nil {
		t.Fatal("flag --file not found on validate command")
	}
}
