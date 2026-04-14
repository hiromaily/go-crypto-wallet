## Requirements

### Core Dependencies

| Tool | Version | Description |
|------|---------|-------------|
| Go | 1.26.2 | Programming language |
| Atlas | 1.1.0 | Database schema migration |
| sqlc | 1.30.0 | SQL code generator |
| Docker | latest | Container runtime |
| Docker Compose | latest | Container orchestration |
| [golangci-lint](https://github.com/golangci/golangci-lint) | v2.10.0+ | Linter (for development) |
| [protoc](https://grpc.io/docs/protoc-installation/) | 33.0+ | Protocol buffer compiler (**Edition 2024**) |
| [buf](https://buf.build/) | latest | Protocol buffer management (lint, format) |

### Database

Supported databases (choose one)

| Tool | Version | Description |
|------|---------|-------------|
| PostgreSQL | 18.2+ | Database (via Docker) |
| MySQL | 8.4+ | Database (via Docker) |
| SQLite | 3.0+ | Database |

### Blockchain Nodes

| Chain | Node | Version | Notes |
|-------|------|---------|-------|
| BTC | [Bitcoin Core](https://bitcoin.org/en/bitcoin-core/) | 28.0+ | supports v28-v30, [Docker image](https://hub.docker.com/r/bitcoin/bitcoin) |
| BCH | [Bitcoin ABC](https://www.bitcoinabc.org/) | 0.21+ | Bitcoin Cash node |
| ETH | [go-ethereum](https://github.com/ethereum/go-ethereum) | latest | Geth client |
| ETH | [Anvil](https://getfoundry.sh/anvil/overview/) | latest | For local development (Foundry) |
| XRP | [rippled](https://xrpl.org/manage-the-rippled-server.html) | latest | Ripple node |

### Major Go Dependencies

| Package | Version | Description |
|---------|---------|-------------|
| btcsuite/btcd | v0.25.0 | Bitcoin library |
| ethereum/go-ethereum | v1.16.7 | Ethereum library |
| spf13/cobra | v1.10.2 | CLI framework |
| spf13/viper | v1.21.0 | Configuration management |
| google.golang.org/grpc | v1.78.0 | gRPC for XRP communication |
| golang.org/x/crypto | v0.46.0 | Cryptographic functions |
