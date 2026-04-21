<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/docs/architecture/design-philosophy.tpl.md · Run `make docs` to regenerate.
-->

# Design Philosophy

## Why Clean Architecture?

This project handles sensitive financial operations including private key management and cryptocurrency transactions. Clean Architecture provides:

1. **Testability**: Core business logic can be tested without external dependencies
2. **Maintainability**: Clear boundaries make changes predictable and safe
3. **Security**: Separation of concerns helps identify and isolate security-critical code
4. **Flexibility**: Infrastructure can be replaced without affecting business logic

## Core Principles

1. **Dependency Rule**: Dependencies point inward. Outer layers depend on inner layers, never the reverse.
2. **Dependency Inversion**: High-level modules define interfaces; low-level modules implement them.
3. **Single Responsibility**: Each layer has one reason to change.
4. **Interface Segregation**: Interfaces are defined by the layer that uses them, not the layer that implements them.
