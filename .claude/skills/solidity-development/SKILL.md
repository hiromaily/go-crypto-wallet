---
name: solidity-development
description: Workflow for Solidity smart contract development. Use when modifying smart contracts in apps/erc20-token/contracts/.
---

# Solidity Development Workflow

Standard workflow for Solidity smart contract changes.

## Applicable Directories

| Path | Description |
|------|-------------|
| `apps/erc20-token/contracts/` | Smart contract source files |
| `contracts/` | ABI files |

## Branch Management

Same as other development:

```bash
git fetch origin
git checkout main
git reset --hard origin/main
git checkout -b {branch-type}/issue-{number}-{brief-description}
```

## Verification Commands

```bash
cd apps/erc20-token
npm install           # Install dependencies
truffle compile       # Compile contracts
truffle test          # Run tests
npm run lint          # Lint Solidity code
```

## Self-Review Checklist

### Code Quality

- [ ] Follows Solidity best practices
- [ ] Gas optimization considered
- [ ] Proper visibility modifiers (public, private, internal, external)
- [ ] Events emitted for state changes

### Security

- [ ] No reentrancy vulnerabilities
- [ ] Integer overflow/underflow protection
- [ ] Access control properly implemented
- [ ] No hardcoded addresses (use constructor parameters)

### Testing

- [ ] Unit tests cover all functions
- [ ] Edge cases tested
- [ ] Gas consumption verified

## ABI Generation

After modifying contracts, regenerate ABI:

```bash
# Compile contracts
cd apps/erc20-token
truffle compile

# Update ABI in main contracts directory
cp build/contracts/Token.json ../../contracts/token.abi
```

**Note**: After ABI update, regenerate Go bindings:

```bash
# From project root
make gen-abi  # If available, or manual abigen command
```

## Commit Message Format

```
feat(contract): {brief description}

- {detail 1}
- {detail 2}

Closes #{issue_number}
```

## Related Chain Context

- ETH (Ethereum mainnet/testnet)
- ERC20 (Token standard)
