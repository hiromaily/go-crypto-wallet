# Research & Design Decisions

---

## Summary

- **Feature**: `erc20-token`
- **Discovery Scope**: New Feature (greenfield Foundry project in `apps/eth-contracts/`)
- **Key Findings**:
  - The HYC token is a standard OpenZeppelin ERC-20 with no custom logic beyond minting — no new ABI patterns needed.
  - Foundry's `forge script` with `vm.envUint("PRIVATE_KEY")` is the standard secure deployment approach for local dev nodes.
  - The project lives entirely in `apps/eth-contracts/`; Go wallet integration (ERC-20 transfer encoding) is handled by existing infrastructure and is out of scope for this spec.

---

## Research Log

### OpenZeppelin ERC-20 (^5.x) vs older versions

- **Context**: Spec references `@openzeppelin/contracts` without pinning a version. Need to confirm constructor API.
- **Findings**:
  - OZ v5.x removed `_setupRole` and other deprecated patterns but `ERC20(name, symbol)` + `_mint` constructor pattern is unchanged.
  - `pragma solidity ^0.8.34` is fully compatible with OZ v5.x.
  - No breaking changes affect the HYC contract as written.
- **Implications**: Use `@openzeppelin/contracts": "^5.6.1"` in `package.json`.

### Foundry `libs` path with bun/npm `node_modules`

- **Context**: Foundry default `libs` is `lib/` (forge install). Spec requires npm-style deps via bun.
- **Findings**:
  - Setting `libs = ["node_modules"]` in `foundry.toml` allows `import "@openzeppelin/..."` to resolve from `node_modules/@openzeppelin/`.
  - This is the standard pattern when using npm/bun instead of `forge install`.
- **Implications**: No `forge install` step; `bun install` is the sole dependency management step.

### Deployment key handling

- **Context**: `vm.envUint("PRIVATE_KEY")` reads a private key from environment. Anvil/geth local nodes provide well-known funded test keys.
- **Findings**:
  - Anvil exposes 10 pre-funded accounts with known private keys (printed on startup).
  - `vm.envUint` is the Foundry-idiomatic approach; avoids hardcoding keys.
- **Implications**: `.env` file with `PRIVATE_KEY=0x...` should be gitignored.

### dprint TypeScript plugin versioning

- **Context**: Spec says "use latest version". As of 2026-02 the latest stable TypeScript plugin is `0.93.x`.
- **Findings**: `https://plugins.dprint.dev/typescript-0.93.0.wasm` is current stable.
- **Implications**: Pin to `0.93.0` in `dprint.json`; update when new versions release.

---

## Architecture Pattern Evaluation

| Option | Description | Strengths | Risks | Notes |
|--------|-------------|-----------|-------|-------|
| Foundry-only project | Pure Foundry with OZ via npm | Simple, no forge install quirks | Requires bun/npm | Matches spec exactly |
| Hardhat | JS-centric tooling | Wide ecosystem | Not in spec | Excluded by spec |
| Foundry + forge install | Foundry native deps | No npm needed | Different import paths than spec | Conflicts with spec |

**Selected**: Foundry + bun (npm-style), per spec.

---

## Design Decisions

### Decision: `apps/eth-contracts/` as project root (not `apps/hyc-token/`)

- **Context**: Spec uses `hyc-token/` as example name; project is ERC-20 token feature.
- **Selected Approach**: `apps/eth-contracts/` to match the spec feature name.
- **Rationale**: Consistent with `apps/` naming convention (feature-based, not token-name-based).

### Decision: Go wallet integration is out of scope

- **Context**: Requirements 4–6 describe ERC-20 encoding and CLI transfer. The existing Go infrastructure (`internal/infrastructure/api/eth/erc20/erc20.go`) already implements this.
- **Selected Approach**: Design covers only the Foundry project. The contract address produced by deployment is the integration point with Go config.
- **Rationale**: User confirmed Go context is not relevant to this spec.

### Decision: No EIP-1559 in the Solidity project itself

- **Context**: EIP-1559 is a transaction format concern, not a contract concern.
- **Selected Approach**: HYC.sol is a pure ERC-20 contract. Transaction type is handled by the deployer/caller.
- **Rationale**: Smart contract is network-agnostic regarding transaction types.

---

## Risks & Mitigations

- **Anvil vs geth RPC differences** — `forge script --broadcast` works with both; test against both nodes in CI.
- **Private key exposure** — `.env` must be gitignored; never commit `PRIVATE_KEY`.
- **OZ v5 breaking changes in future** — pin major version (`^5.0.0`) to avoid unexpected upgrades.

---

## References

- [OpenZeppelin ERC20 v5 docs](https://docs.openzeppelin.com/contracts/5.x/erc20)
- [Foundry Book — Deploying with Scripts](https://book.getfoundry.sh/tutorials/solidity-scripting)
- [Foundry — npm dependencies](https://book.getfoundry.sh/config/dependencies#node-modules-layout)
- [dprint TypeScript plugin](https://dprint.dev/plugins/typescript/)
