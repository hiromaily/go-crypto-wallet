# XRPL gRPC Server

Modern XRP Ledger gRPC server built with Bun runtime.

## Overview

This server provides gRPC endpoints for interacting with the XRP Ledger, replacing the older `ripple-lib-server` implementation with modern tooling.

## Prerequisites

- [Bun](https://bun.sh/) >= 1.3.6
- Buf CLI >= 1.64.0 (not yet available - see [SETUP.md](SETUP.md) for details)

## Installation

```bash
bun install
```

## Code Generation

Generate TypeScript code from proto files:

```bash
bun run proto
```

This generates TypeScript types and Connect-ES service clients from the proto definitions in `../../proto/rippleapi/`.

⚠️ **Important**: This **currently does not work** due to buf CLI not supporting Protobuf Edition 2024 yet. The proto files use `edition = "2024"` (see `../../docs/proto.md`). Code generation will work once buf CLI >= 1.64.0 is released. See [SETUP.md](SETUP.md) for full details.

## Development

```bash
# Start server in development mode with hot reload
bun run dev

# Run type checking
bun run typecheck

# Run linter
bun run lint

# Fix linting issues
bun run lint:fix

# Format code
bun run format
```

## Production

```bash
bun run start
```

## Technology Stack

- **Runtime**: Bun
- **Language**: TypeScript 5.9.3
- **gRPC Framework**: Connect RPC (@connectrpc/connect)
- **Protobuf**: Buf (@bufbuild/protobuf)
- **XRP Library**: xrpl.js (xrpl)
- **Linter/Formatter**: Biome

## Project Structure

```
apps/xrpl-grpc-server/
├── src/           # Source code
├── biome.json     # Biome configuration
├── package.json   # Dependencies
├── tsconfig.json  # TypeScript configuration
└── README.md      # This file
```

## Related

- Part of the [go-crypto-wallet](../../README.md) multi-signature wallet system
- See issue [#470](https://github.com/hiromaily/go-crypto-wallet/issues/470) for migration plan
