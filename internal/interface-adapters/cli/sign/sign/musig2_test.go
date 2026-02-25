package sign_test

import (
	"testing"

	"github.com/spf13/cobra"

	signcli "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/sign/sign"
)

func TestAddMuSig2Commands(t *testing.T) {
	t.Parallel()

	rootCmd := &cobra.Command{Use: "test"}
	signCmd := &cobra.Command{Use: "sign"}
	rootCmd.AddCommand(signCmd)
	signcli.AddCommands(signCmd, nil, nil)

	musig2Cmd, _, err := signCmd.Find([]string{"musig2"})
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

	if !musig2Cmd.HasSubCommands() {
		t.Error("musig2 command should have subcommands")
	}

	nonceCmd, _, err := musig2Cmd.Find([]string{"nonce"})
	if err != nil {
		t.Errorf("nonce subcommand not found: %v", err)
	}
	if nonceCmd == nil {
		t.Error("nonce subcommand is nil")
	}

	signSubCmd, _, err := musig2Cmd.Find([]string{"sign"})
	if err != nil {
		t.Errorf("sign subcommand not found: %v", err)
	}
	if signSubCmd == nil {
		t.Error("sign subcommand is nil")
	}
}

func TestMuSig2NonceCommandFlags(t *testing.T) {
	t.Parallel()

	signCmd := &cobra.Command{Use: "sign"}
	signcli.AddCommands(signCmd, nil, nil)

	musig2Cmd, _, err := signCmd.Find([]string{"musig2"})
	if err != nil {
		t.Fatalf("musig2 command not found: %v", err)
	}

	nonceCmd, _, err := musig2Cmd.Find([]string{"nonce"})
	if err != nil {
		t.Fatalf("nonce subcommand not found: %v", err)
	}

	if nonceCmd.Flags().Lookup("file") == nil {
		t.Error("--file flag not found")
	}
	if nonceCmd.Flags().Lookup("output") == nil {
		t.Error("--output flag not found")
	}
}

func TestMuSig2SignCommandFlags(t *testing.T) {
	t.Parallel()

	signCmd := &cobra.Command{Use: "sign"}
	signcli.AddCommands(signCmd, nil, nil)

	musig2Cmd, _, err := signCmd.Find([]string{"musig2"})
	if err != nil {
		t.Fatalf("musig2 command not found: %v", err)
	}

	signSubCmd, _, err := musig2Cmd.Find([]string{"sign"})
	if err != nil {
		t.Fatalf("sign subcommand not found: %v", err)
	}

	if signSubCmd.Flags().Lookup("file") == nil {
		t.Error("--file flag not found")
	}
	if signSubCmd.Flags().Lookup("output") == nil {
		t.Error("--output flag not found")
	}
}
