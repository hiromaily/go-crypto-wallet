### ERC-20 Contract Deployment

The HYT ERC-20 contract is deployed automatically by `make eth-e2e-p2` using Foundry (`forge`). No manual deployment step is required for E2E testing.

For manual deployment (advanced use):

```bash
cd ./apps/eth-contracts
forge build
forge script script/DeployHYT.s.sol --broadcast --rpc-url http://localhost:8546
```

Requires Foundry to be installed (`curl -L https://foundry.paradigm.xyz | bash && foundryup`).
