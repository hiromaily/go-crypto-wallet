package eth

import (
	"encoding/json"
	"os"
)

// isETHMultisigFile returns true if the file at filePath contains a "safe_address" key,
// indicating it is an ETH Safe multisig transaction file (ETHMultisigTransactionFile).
func isETHMultisigFile(filePath string) bool {
	data, err := os.ReadFile(filePath) // #nosec G304 -- file path comes from trusted CLI input
	if err != nil {
		return false
	}
	var peek struct {
		SafeAddress string `json:"safe_address"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return false
	}
	return peek.SafeAddress != ""
}
