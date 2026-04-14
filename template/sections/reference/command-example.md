## Command example

- [CLI Command Reference](../../../internal/interface-adapters/cli/README.md) - Command × Chain × UseCase matrix (SSOT): which commands each chain supports, corresponding use case interfaces, and role of each command
- [docs/getting-started/commands.md](../../../docs/getting-started/commands.md) - Detailed flags, options, and usage examples per command
- [Makefile](https://github.com/hiromaily/go-crypto-wallet/blob/main/Makefile) - Main Makefile with modular includes
- Makefile modules (in `make/` directory):
  - [watch_op.mk](https://github.com/hiromaily/go-crypto-wallet/blob/main/make/watch_op.mk) - Watch wallet operations
  - [keygen_op.mk](https://github.com/hiromaily/go-crypto-wallet/blob/main/make/keygen_op.mk) - Keygen wallet operations
  - [sign_op.mk](https://github.com/hiromaily/go-crypto-wallet/blob/main/make/sign_op.mk) - Sign wallet operations
  - And other specialized modules for builds, tests, Docker, etc.
