## Installation

This guide covers setting up the development environment on macOS.

### Requirements

| Tool | Version | Purpose |
|------|---------|---------|
| [Go](https://go.dev/dl/) | 1.26.2+ | Build the wallet binaries |
| [Docker](https://www.docker.com/get-started) | latest | Blockchain nodes and databases |
| Docker Compose | latest | Container orchestration |
| [Foundry](https://getfoundry.sh/) | latest | ETH E2E only: deploy ERC-20 and Safe contracts (P2, P3), cast for MPC-TSS (P4) |

> **For E2E tests (recommended entry point):** Go, Docker, and Docker Compose are sufficient for BTC, BCH, and XRP. ETH E2E patterns P2, P3, and P4 additionally require Foundry.

Install Foundry (macOS/Linux):

```bash
curl -L https://foundry.paradigm.xyz | bash
foundryup
```
