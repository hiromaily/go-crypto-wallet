package sign_test

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/sign"
)

func TestAddMuSig2Commands(t *testing.T) {
	t.Parallel()

	// Create a root command
	rootCmd := &cobra.Command{
		Use: "test",
	}

	// Add MuSig2 commands with nil container (testing command registration only)
	sign.AddMuSig2Commands(rootCmd, nil)

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

	// Verify nonce subcommand
	nonceCmd, _, err := musig2Cmd.Find([]string{"nonce"})
	if err != nil {
		t.Errorf("nonce subcommand not found: %v", err)
	}
	if nonceCmd == nil {
		t.Error("nonce subcommand is nil")
	}

	// Verify sign subcommand
	signCmd, _, err := musig2Cmd.Find([]string{"sign"})
	if err != nil {
		t.Errorf("sign subcommand not found: %v", err)
	}
	if signCmd == nil {
		t.Error("sign subcommand is nil")
	}
}

func TestMuSig2NonceCommandFlags(t *testing.T) {
	t.Parallel()

	rootCmd := &cobra.Command{Use: "test"}
	sign.AddMuSig2Commands(rootCmd, nil)

	musig2Cmd, _, err := rootCmd.Find([]string{"musig2"})
	if err != nil {
		t.Fatalf("musig2 command not found: %v", err)
	}

	nonceCmd, _, err := musig2Cmd.Find([]string{"nonce"})
	if err != nil {
		t.Fatalf("nonce subcommand not found: %v", err)
	}

	// Verify required flags exist
	fileFlag := nonceCmd.Flags().Lookup("file")
	if fileFlag == nil {
		t.Error("--file flag not found")
	}

	outputFlag := nonceCmd.Flags().Lookup("output")
	if outputFlag == nil {
		t.Error("--output flag not found")
	}
}

func TestMuSig2SignCommandFlags(t *testing.T) {
	t.Parallel()

	rootCmd := &cobra.Command{Use: "test"}
	sign.AddMuSig2Commands(rootCmd, nil)

	musig2Cmd, _, err := rootCmd.Find([]string{"musig2"})
	if err != nil {
		t.Fatalf("musig2 command not found: %v", err)
	}

	signCmd, _, err := musig2Cmd.Find([]string{"sign"})
	if err != nil {
		t.Fatalf("sign subcommand not found: %v", err)
	}

	// Verify required flags exist
	fileFlag := signCmd.Flags().Lookup("file")
	if fileFlag == nil {
		t.Error("--file flag not found")
	}

	outputFlag := signCmd.Flags().Lookup("output")
	if outputFlag == nil {
		t.Error("--output flag not found")
	}
}
