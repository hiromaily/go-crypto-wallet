# HYC Token (ERC-20) Smart Contract + Deployment Environment

## 1. Objective

Implement an ERC-20 compatible token called **HYC** and deploy it to an Ethereum node.

The goal of this task is **only**:

1. Implement the Solidity contract
2. Compile and deploy using Foundry
3. Enable ERC-20 transfer calls from the CLI

The latest tech stack in 2026

---

## 2. Target Network

Network:

* Ethereum

Local development nodes already exist via Docker Compose:

```
anvil
geth
```

These nodes expose an RPC endpoint compatible with Ethereum JSON-RPC.

Example RPC

```
http://localhost:8545
```

---

## 3. Technology Stack

Smart contract framework:

* Foundry

ERC-20 implementation:

* OpenZeppelin contracts

Package manager priority:

1. **bun**
2. fallback → **pnpm**

Rust-based tooling requirements:

Formatter

* dprint

Solidity linter

* solhint

---

## 4. Project Structure

Recommended layout:

```
hyc-token/
│
├─ contracts/
│   └─ HYC.sol
│
├─ script/
│   └─ DeployHYC.s.sol
│
├─ test/
│
├─ foundry.toml
│
├─ package.json
│
└─ bun.lockb / pnpm-lock.yaml
```

---

## 5. Install Dependencies

Prefer **bun**.

```
bun init
bun add @openzeppelin/contracts
```

If bun is not compatible with the environment:

```
pnpm init
pnpm add @openzeppelin/contracts
```

Install development tools:

```
bun add -d solhint
bun add -d dprint
```

---

## 6. HYC Token Solidity Implementation

Create file:

```
contracts/HYC.sol
```

Solidity implementation:

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.34;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract HYC is ERC20 {

    constructor(uint256 initialSupply)
        ERC20("hiromaily Coin", "HYC")
    {
        _mint(msg.sender, initialSupply);
    }

}
```

Token parameters:

```
Name: hiromaily Coin
Symbol: HYC
Decimals: 18
```

Initial supply is minted to the deployer.

Example initial supply:

```
1_000_000 * 10^18
```

---

## 7. Foundry Configuration

Create `foundry.toml`.

Example:

```
[profile.default]
src = "contracts"
out = "out"
libs = ["node_modules"]

solc_version = "0.8.34"
optimizer = true
optimizer_runs = 200
```

Compile contracts:

```
forge build
```

---

## 8. Deployment Script

Create solidity file:

Example script:

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.34;

import "forge-std/Script.sol";
import "../contracts/HYC.sol";

contract DeployHYC is Script {

    function run() external {

        uint256 deployerKey = vm.envUint("PRIVATE_KEY");

        vm.startBroadcast(deployerKey);

        new HYC(1_000_000 ether);

        vm.stopBroadcast();
    }
}
```

---

## 9. Deploy to Local Ethereum Node

Example using Anvil:

```
forge script script/{deployment script}.sol \
--rpc-url http://localhost:8545 \
--broadcast
```

Output will include:

```
contract address
transaction hash
gas usage
```

Record the deployed **token contract address**.

---

## 10. Token Transfer Flow

The CLI wallet will interact with the deployed HYC token.

Transfer flow:

```
CLI
 ↓
encode ERC20 transfer()
 ↓
create EIP-1559 transaction
 ↓
sign transaction
 ↓
eth_sendRawTransaction
 ↓
Ethereum node
```

---

## 11. ERC-20 Transfer Encoding

Function signature:

```
transfer(address,uint256)
```

Method selector:

```
0xa9059cbb
```

Transaction data format:

```
0xa9059cbb
<32 bytes recipient>
<32 bytes amount>
```

Example:

```
to = token contract address
value = 0
data = encoded transfer call
```

---

## 12. Example Raw Transaction

Example fields:

```
type: 0x2
chainId: 1
nonce: N
maxPriorityFeePerGas
maxFeePerGas
gas
to: HYC contract
value: 0
data: transfer calldata
```

Transaction type:

EIP-1559.

---

## 13. Formatter Configuration [dprint](https://github.com/dprint/dprint)

use latest version

Example `dprint.json`.

```
{
  "plugins": [
    "https://plugins.dprint.dev/typescript-0.90.0.wasm"
  ]
}
```

Run formatter:

```
dprint fmt
```

---

## 14. Solidity Linting

Example `.solhint.json`.

```
{
  "extends": "solhint:recommended"
}
```

Run lint:

```
solhint contracts/**/*.sol
```

---

## 15. Expected Outcome

After implementation:

1. HYC ERC-20 token can be deployed to Ethereum node
2. CLI can generate raw transactions
3. CLI can perform:

```
HYC transfer
```

using:

```
transfer(address,uint256)
```
