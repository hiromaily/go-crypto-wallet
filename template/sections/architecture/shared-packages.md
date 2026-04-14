## Shared Packages (`pkg/`)

The `pkg/` directory contains reusable utilities that can be used across the application:

```
pkg/
├── config/       # Configuration loading utilities
├── logger/       # Structured logging (slog-based)
├── db/mysql/     # Database connection utilities
├── grpc/         # gRPC client utilities
├── websocket/    # WebSocket client utilities
├── uuid/         # UUID generation
├── decimal/      # Decimal number handling
└── retry/        # Retry utilities
```

**Critical Rule**: Packages in `pkg/` MUST NOT import from `internal/`.
