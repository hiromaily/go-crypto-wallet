package file

import (
	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
)

// MPCFileRepositorier handles ETHMPCTransactionFile JSON file I/O for the MPC/TSS flow.
//
// It is defined as a separate interface from MultisigFileRepositorier and
// TransactionFileRepositorier to satisfy the Interface Segregation Principle:
// the Safe multisig and single-sig use cases are unaffected by MPC file operations.
//
// File naming convention:
//
//	Unsigned: {action_type}_mpc_{uuid}.json
//	Signed:   {action_type}_mpc_{uuid}_signed.json
type MPCFileRepositorier interface {
	// ReadETHMPCJSONFile reads the JSON file at filePath, deserialises it, and calls
	// Validate() before returning. Returns an error if the file is missing, malformed,
	// or fails validation.
	ReadETHMPCJSONFile(filePath string) (*dtoeth.ETHMPCTransactionFile, error)

	// WriteETHMPCJSONFile serialises file to indented JSON and writes it to the canonical
	// path derived from file.ActionType, file.UUID, and isSigned. Returns the written file
	// path on success.
	WriteETHMPCJSONFile(file *dtoeth.ETHMPCTransactionFile, isSigned bool) (string, error)

	// CreateMPCFilePath composes the canonical file path from actionType, uuid, and isSigned
	// without performing any I/O. Useful for pre-computing the output path before writing.
	CreateMPCFilePath(actionType string, uuid string, isSigned bool) string
}
