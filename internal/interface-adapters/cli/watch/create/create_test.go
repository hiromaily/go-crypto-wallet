package create

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// TestCreateCommandsHaveQuorumFlag verifies that deposit, payment, and transfer
// commands each expose an optional --quorum flag for multisig TX file creation.
func TestCreateCommandsHaveQuorumFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cmdName string
	}{
		{"deposit", "deposit"},
		{"payment", "payment"},
		{"transfer", "transfer"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			parent := &cobra.Command{Use: "create"}
			AddCommands(parent, nil, nil)

			cmd, _, err := parent.Find([]string{tt.cmdName})
			assert.NoError(t, err)
			assert.NotNil(t, cmd, "command %s not found", tt.cmdName)

			flag := cmd.Flags().Lookup("quorum")
			assert.NotNil(t, flag, "command %q should have --quorum flag", tt.cmdName)
		})
	}
}
