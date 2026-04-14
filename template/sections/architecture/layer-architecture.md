## Layer Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Interface Adapters Layer                      │
│                 (CLI Commands, HTTP Handlers)                    │
│                  internal/interface-adapters/                    │
└────────────────────────────┬────────────────────────────────────┘
                             │ depends on
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Application Layer                            │
│            (Use Cases, Ports, DTOs)                              │
│                  internal/application/                           │
└──────────────┬─────────────────────────────────┬────────────────┘
               │ depends on                       │ depends on
               ▼                                  ▼
┌──────────────────────────────┐    ┌─────────────────────────────┐
│        Domain Layer          │    │    Infrastructure Layer     │
│    (Pure Business Logic)     │    │   (External Dependencies)   │
│      internal/domain/        │    │   internal/infrastructure/  │
└──────────────────────────────┘    └──────────────┬──────────────┘
                                                   │ implements
                                                   ▼
                                    ┌─────────────────────────────┐
                                    │   Application Ports         │
                                    │ (Interface Definitions)     │
                                    │  internal/application/ports/│
                                    └─────────────────────────────┘
```
