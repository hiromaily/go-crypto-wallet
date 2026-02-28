# Implementation Plan

- [x] 1. Set up the Foundry dependency layer
- [x] 1.1 Install the forge-std library
  - Run `forge install foundry-rs/forge-std` inside `apps/eth-contracts/` to populate `lib/forge-std/`
  - Confirm `lib/forge-std/src/Script.sol` and `lib/forge-std/src/Test.sol` are present
  - Add `lib/` to `.gitignore` (or commit the submodule if preferred)
  - _Requirements: 2.1, 2.2, 2.3_

- [x] 1.2 Install npm packages
  - Run `bun install` inside `apps/eth-contracts/` to install `@openzeppelin/contracts`, `solhint`, and `dprint`
  - Confirm `node_modules/@openzeppelin/contracts/token/ERC20/ERC20.sol` is present
  - _Requirements: 2.5, 2.6, 7.7_

- [x] 1.3 Verify the Foundry build compiles without errors
  - Run `forge build` and confirm compiled artifacts appear in `out/`
  - Confirm `out/HYC.sol/HYC.json` exists and contains the ABI
  - _Requirements: 2.4_

- [x] 2. (P) Implement the HYC ERC-20 token contract
  - Write `contracts/HYC.sol` inheriting from OpenZeppelin `ERC20`
  - Set token name to `"hiromaily Coin"` and symbol to `"HYC"` in the `ERC20` constructor call
  - Accept `uint256 initialSupply` as a constructor parameter and mint it to `msg.sender`
  - Use `pragma solidity ^0.8.34`; decimals (18) are inherited from the base contract
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6_

- [x] 3. (P) Implement the HYC deployment script
  - Write `script/DeployHYC.s.sol` inheriting from forge-std `Script`
  - Read the deployer private key from the `PRIVATE_KEY` environment variable via `vm.envUint`
  - Wrap `new HYC(1_000_000 ether)` inside `vm.startBroadcast` / `vm.stopBroadcast`
  - Confirm the script works with both anvil and geth RPC endpoints at `http://localhost:8545`
  - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6_

- [x] 4. (P) Configure development tooling
- [x] 4.1 (P) Configure the Solidity linter
  - Ensure `solhint` is listed as a dev dependency in `package.json`
  - Ensure `.solhint.json` extends `solhint:recommended`
  - Run `bun run lint` and confirm zero errors on `contracts/**/*.sol`
  - _Requirements: 7.1, 7.2, 7.3_

- [x] 4.2 (P) Configure the JS/TS formatter
  - Ensure `dprint` is listed as a dev dependency in `package.json`
  - Ensure `dprint.json` specifies the TypeScript plugin (latest wasm URL)
  - Run `bun run fmt` and confirm it exits cleanly with no violations
  - _Requirements: 7.4, 7.5, 7.6, 7.7_

- [x] 5. Write and run contract unit tests
  - Write `test/HYC.t.sol` as a Foundry test contract inheriting from `forge-std/Test.sol`
  - Verify constructor mints `1_000_000 ether` to `msg.sender` and `totalSupply` matches
  - Verify `name()` returns `"hiromaily Coin"`, `symbol()` returns `"HYC"`, `decimals()` returns `18`
  - Verify `transfer(address, uint256)` correctly updates sender and recipient balances
  - Verify `transfer` to the zero address reverts (OpenZeppelin invariant)
  - Run `forge test` and confirm all tests pass with zero failures
  - _Requirements: 1.1, 1.2, 1.3, 1.4_
