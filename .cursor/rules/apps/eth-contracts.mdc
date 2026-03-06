---
description: ETH Contracts Directory Rules
globs: ["apps/eth-contracts/**"]
alwaysApply: false
---


# ETH Contracts Directory Rules

Rules for working in `apps/eth-contracts/` — a standalone Foundry project for Ethereum smart contracts.

## Skill

Use the `solidity-development` skill for all changes in this directory.

## Toolchain

| Tool | Command | Notes |
|------|---------|-------|
| Compile | `forge build` | Artifacts go to `out/` |
| Test | `forge test -v` | All `*.t.sol` files in `test/` |
| Lint | `bun run lint` | solhint — must exit 0, zero errors |
| Format | `bun run fmt` | dprint — must exit 0 |
| Install deps | `bun install` | npm packages (`node_modules/`) |
| Install forge libs | `forge install foundry-rs/forge-std` | Foundry libs go to `lib/` |

> `forge` is at `~/.foundry/bin/forge` if not in PATH.

## Verification (required before commit)

```bash
cd apps/eth-contracts
forge build       # must succeed with zero errors
bun run lint      # must exit 0, zero solhint errors
forge test -v     # must pass with zero failures
```

## Solidity Conventions

- **Pragma**: `^0.8.34` (or latest stable)
- **Imports**: Named imports only — `import {ERC20} from "..."` (no global imports)
- **NatSpec**: Required on all contracts (`@title`, `@author`, `@notice`) and public/external functions (`@notice`, `@param`, `@return`)
- **Visibility**: Explicit on all functions; constructors exempt in Solidity ≥0.8
- **Dependencies**: OpenZeppelin Contracts `^5.6.1` via npm (`node_modules/`)
- **Foundry libs**: forge-std via `forge install` (not npm), stored in `lib/`

## Security Rules

- **Never** hardcode private keys — use `vm.envUint("PRIVATE_KEY")` in scripts
- **Always** keep `.env` in `.gitignore`
- Do not add custom access control, pausability, or upgradeability unless explicitly required
- `_mint` should only be called in constructors (fixed supply)

## File Layout

```
apps/eth-contracts/
├── contracts/    ← production Solidity source (*.sol)
├── script/       ← Foundry deploy scripts (*.s.sol)
├── test/         ← Foundry unit tests (*.t.sol)
├── lib/          ← forge install libraries (gitignored)
├── node_modules/ ← bun install packages (gitignored)
└── out/          ← compiled artifacts (gitignored)
```

## Do NOT

- Run `npm install` or `yarn` — use `bun install`
- Use `truffle` or `hardhat` — this project uses Foundry only
- Edit files in `out/`, `lib/`, or `node_modules/` — these are generated/installed
- Commit `lib/`, `out/`, `node_modules/`, or `.env`
