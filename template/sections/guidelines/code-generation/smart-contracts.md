### Smart Contract ABI Code

**Tool**: [abigen](https://geth.ethereum.org/docs/tools/abigen) (from go-ethereum)
**Source**: `contracts/token.abi`
**Command**: `make generate-abi` (or `abigen --abi ./contracts/token.abi --pkg contract --type Token --out ./internal/infrastructure/contract/token-abi.go`)

**Generated Files**:

- `internal/infrastructure/contract/token-abi.go` - ERC-20 token contract bindings

**Note**: ABI code is generated from Ethereum smart contract ABI JSON files.
