# eth-contracts

Ethereum smart contracts for the go-crypto-wallet project, built with [Foundry](https://book.getfoundry.sh/).

## Contracts

| Contract | File | Description |
|----------|------|-------------|
| HYC | `contracts/HYC.sol` | ERC-20 token ("hiromaily Coin", symbol `HYC`, 18 decimals) |

## Prerequisites

| Tool | Install |
|------|---------|
| [Foundry](https://book.getfoundry.sh/getting-started/installation) | `curl -L https://foundry.paradigm.xyz \| bash && foundryup` |
| [bun](https://bun.sh) | `curl -fsSL https://bun.sh/install \| bash` |

## Setup

```bash
# 1. Install forge-std (Foundry scripting/testing library)
forge install foundry-rs/forge-std

# 2. Install npm dependencies (@openzeppelin/contracts, solhint, dprint)
bun install
```

## Usage

```bash
# Compile contracts
bun run build       # forge build → artifacts in out/

# Run tests
bun run test        # forge test -v

# Lint Solidity
bun run lint        # solhint contracts/**/*.sol

# Format JS/TS files
bun run fmt         # dprint fmt
```

## Deployment

Deploy to a local Ethereum node (`anvil` or `geth`) at `http://localhost:8545`:

```bash
export PRIVATE_KEY=0x<deployer-private-key>
bun run deploy:local
# → outputs contract address, transaction hash, gas usage
```

After deployment, record the contract address in the Go wallet configuration:

```yaml
# configs/wallet.yml (example)
ethereum:
  erc20_token: hyc
  erc20s:
    hyc:
      symbol: HYC
      name: hiromaily Coin
      contract_address: "0x<deployed-address>"
      master_address: "0x<deployer-address>"
      decimals: 18
```

## Project Structure

```
eth-contracts/
├── contracts/
│   └── HYC.sol               # ERC-20 token contract
├── script/
│   └── DeployHYC.s.sol       # Foundry deployment script
├── test/
│   └── HYC.t.sol             # Foundry unit tests
├── lib/
│   └── forge-std/            # Foundry standard library (forge install)
├── node_modules/             # npm dependencies (bun install)
├── out/                      # Compiled artifacts (forge build)
├── foundry.toml              # Foundry configuration
├── package.json              # npm scripts and dependencies
├── .solhint.json             # Solidity linter configuration
└── dprint.json               # JS/TS formatter configuration
```

## Technology Stack

| Layer | Choice | Version |
|-------|--------|---------|
| Smart contract language | Solidity | `^0.8.34` |
| Build / test framework | Foundry | 1.6.0+ |
| ERC-20 base | OpenZeppelin Contracts | `^5.6.1` |
| Package manager | bun | primary (pnpm fallback) |
| Solidity linter | solhint | `^6.0.3` |
| Formatter | dprint | `^0.52.0` |

## Security

- The deployer private key is read exclusively from the `PRIVATE_KEY` environment variable — never hardcode keys in source files.
- Add `.env` to `.gitignore` to prevent accidental commits.
- HYC uses only standard OpenZeppelin ERC-20 logic with no custom transfer logic, admin functions, or upgradeability.
- The total supply is fixed at deployment time (`_mint` is called only in the constructor).
