package watch_test

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/watch"
)

func TestAddMuSig2Commands(t *testing.T) {
	t.Parallel()

	// Create a root command
	rootCmd := &cobra.Command{
		Use: "test",
	}

	// Add MuSig2 commands with nil container (testing command registration only)
	watch.AddMuSig2Commands(rootCmd, nil)

	// Verify musig2 command was added
	musig2Cmd, _, err := rootCmd.Find([]string{"musig2"})
	if err != nil {
		t.Fatalf("musig2 command not found: %v", err)
	}

	if musig2Cmd == nil {
		t.Fatal("musig2 command is nil")
		return
	}

	if musig2Cmd.Use != "musig2" {
		t.Errorf("expected Use='musig2', got '%s'", musig2Cmd.Use)
	}

	// Verify subcommands exist
	if !musig2Cmd.HasSubCommands() {
		t.Error("musig2 command should have subcommands")
	}

	// Verify collect-nonces subcommand
	collectNoncesCmd, _, err := musig2Cmd.Find([]string{"collect-nonces"})
	if err != nil {
		t.Errorf("collect-nonces subcommand not found: %v", err)
	}
	if collectNoncesCmd == nil {
		t.Error("collect-nonces subcommand is nil")
	}

	// Verify aggregate subcommand
	aggregateCmd, _, err := musig2Cmd.Find([]string{"aggregate"})
	if err != nil {
		t.Errorf("aggregate subcommand not found: %v", err)
	}
	if aggregateCmd == nil {
		t.Error("aggregate subcommand is nil")
	}
}

func TestMuSig2CollectNoncesCommandFlags(t *testing.T) {
	t.Parallel()

	rootCmd := &cobra.Command{Use: "test"}
	watch.AddMuSig2Commands(rootCmd, nil)

	musig2Cmd, _, err := rootCmd.Find([]string{"musig2"})
	if err != nil {
		t.Fatalf("musig2 command not found: %v", err)
	}

	collectNoncesCmd, _, err := musig2Cmd.Find([]string{"collect-nonces"})
	if err != nil {
		t.Fatalf("collect-nonces subcommand not found: %v", err)
	}

	// Verify required flags exist
	fileFlag := collectNoncesCmd.Flags().Lookup("file")
	if fileFlag == nil {
		t.Error("--file flag not found")
	}

	noncesFlag := collectNoncesCmd.Flags().Lookup("nonces")
	if noncesFlag == nil {
		t.Error("--nonces flag not found")
	}

	outputFlag := collectNoncesCmd.Flags().Lookup("output")
	if outputFlag == nil {
		t.Error("--output flag not found")
	}
}

func TestMuSig2AggregateCommandFlags(t *testing.T) {
	t.Parallel()

	rootCmd := &cobra.Command{Use: "test"}
	watch.AddMuSig2Commands(rootCmd, nil)

	musig2Cmd, _, err := rootCmd.Find([]string{"musig2"})
	if err != nil {
		t.Fatalf("musig2 command not found: %v", err)
	}

	aggregateCmd, _, err := musig2Cmd.Find([]string{"aggregate"})
	if err != nil {
		t.Fatalf("aggregate subcommand not found: %v", err)
	}

	// Verify required flags exist
	filesFlag := aggregateCmd.Flags().Lookup("files")
	if filesFlag == nil {
		t.Error("--files flag not found")
	}

	outputFlag := aggregateCmd.Flags().Lookup("output")
	if outputFlag == nil {
		t.Error("--output flag not found")
	}
}
