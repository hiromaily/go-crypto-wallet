# Dependency Injection

The `internal/di/` package wires together all dependencies:

```go
// Example: Watch wallet dependency injection
func NewWatchWallet(cfg *config.WalletConfig) (*WatchWallet, error) {
    // Infrastructure
    bitcoinClient := bitcoin.NewClient(cfg.Bitcoin)
    repository := watchrepo.NewRepository(db)

    // Application (Use Cases)
    createTxUseCase := watch.NewCreateTransactionUseCase(bitcoinClient, repository)

    // Interface Adapters (Commands)
    return &WatchWallet{
        createTxCmd: cli.NewCreateTxCommand(createTxUseCase),
    }, nil
}
```
