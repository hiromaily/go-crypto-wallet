# XRPL gRPC Server

Modern XRP Ledger gRPC server built with Bun runtime.

## Overview

This server provides gRPC endpoints for interacting with the XRP Ledger, implementation with modern tooling.

## Tech Stack

| Component | Technology |
|-----------|------------|
| Runtime | [Bun](https://bun.sh/) >= 1.3.6 |
| Language | TypeScript 5.9.3 |
| XRP Library | xrpl.js 4.5.0 |
| gRPC Framework | [ConnectRPC](https://connectrpc.com/) (@connectrpc/connect) |
| Protobuf | @bufbuild/protobuf |
| Linter/Formatter | [Biome](https://biomejs.dev/) |

## Prerequisites

- [Bun](https://bun.sh/) >= 1.3.6
- protoc >= 33.4 (for proto generation)
- Access to XRP Ledger node (testnet by default)

## Setup

### Install dependencies

```bash
bun install
```

### Generate proto files

```bash
make proto
```

### Start development server

```bash
make dev
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `GRPC_PORT` | `50051` | gRPC server port |
| `GRPC_HOST` | `0.0.0.0` | gRPC server host |
| `XRP_WS_URL` | - | Direct WebSocket URL (overrides XRP_NETWORK) |
| `XRP_NETWORK` | `testnet` | Network selection: `mainnet`, `testnet`, `devnet` |
| `XRP_CONNECTION_TIMEOUT` | `20000` | Connection timeout in milliseconds |
| `XRP_MAX_RECONNECT_ATTEMPTS` | `5` | Maximum reconnection attempts |
| `XRP_RECONNECT_DELAY` | `1000` | Delay between reconnection attempts (ms) |

### Network WebSocket URLs

| Network | WebSocket URL |
|---------|---------------|
| mainnet | `wss://xrplcluster.com` |
| testnet | `wss://s.altnet.rippletest.net:51233` |
| devnet | `wss://s.devnet.rippletest.net:51233` |

## Commands

| Command | Description |
|---------|-------------|
| `make install` | Install dependencies |
| `make dev` | Start dev server with hot reload |
| `make build` | Build for production |
| `make lint` | Run Biome linter |
| `make lint-fix` | Fix lint issues |
| `make format` | Format code with Biome |
| `make typecheck` | Run TypeScript type checking |
| `make proto` | Generate proto files |
| `make clean` | Clean build artifacts |

### npm scripts

| Script | Description |
|--------|-------------|
| `bun run start` | Start production server |
| `bun run dev` | Start dev server with watch mode |
| `bun run typecheck` | Run TypeScript type checking |
| `bun run lint` | Run Biome linter |
| `bun run lint:fix` | Fix lint issues |
| `bun run format` | Format code |

## Project Structure

```
apps/xrpl-grpc-server/
├── src/                # Source code
│   ├── gen/            # Generated protobuf/connect code
│   ├── index.ts        # Entry point
│   ├── server.ts       # ConnectRPC server setup
│   ├── config.ts       # Environment configuration
│   ├── services/       # gRPC service implementations
│   │   ├── account.ts  # RippleAccountAPI
│   │   ├── address.ts  # RippleAddressAPI
│   │   └── transaction.ts # RippleTransactionAPI
│   └── xrpl/           # XRPL client wrapper
│       └── client.ts
├── docs/               # Documentation
│   ├── MIGRATION-GUIDE.md
│   └── PROTOBUF-EDITION-2024.md
├── biome.json          # Biome configuration
├── package.json        # Dependencies
├── tsconfig.json       # TypeScript configuration
├── bun.lock            # Bun lockfile
├── Makefile            # Build commands
└── README.md           # This file
```

## References

- [xrpl.js Documentation](https://js.xrpl.org/)
- [xrpl.org Migration Guide](https://xrpl.org/docs/references/xrpljs2-migration-guide)
- [ConnectRPC Documentation](https://connectrpc.com/docs/node/getting-started)
- [Biome Documentation](https://biomejs.dev/)
- [Bun Documentation](https://bun.sh/docs)

## Related

- Part of the [go-crypto-wallet](../../README.md) multi-signature wallet system
- See issue [#470](https://github.com/hiromaily/go-crypto-wallet/issues/470) for migration plan
