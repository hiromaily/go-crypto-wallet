## Database Schema Change Workflow

This document provides a comprehensive guide for modifying or adding database columns and schemas in the go-crypto-wallet project, which supports multiple database backends.

### Table of Contents

- [Overview](#overview)
- [Supported Databases](#supported-databases)
- [Quick Reference](#quick-reference)
- [Step-by-Step Workflow](#step-by-step-workflow)
- [Multi-Database Considerations](#multi-database-considerations)
- [Testing Schema Changes](#testing-schema-changes)
- [Common Patterns](#common-patterns)
- [Troubleshooting](#troubleshooting)
- [Best Practices](#best-practices)

### Overview

The project uses a **declarative schema management** approach with [Atlas](https://atlasgo.io/) and type-safe code generation with [sqlc](https://sqlc.dev/).

**Key Principle**: HCL schema files are the **single source of truth**. Never edit migration SQL files or generated code directly.

#### Architecture

```
HCL Schemas (Source of Truth)
  ↓ Atlas generates
Database Migration Files (.sql)
  ↓ Applied to
Running Databases (MySQL/SQLite/PostgreSQL)
  ↓ Dumped to
SQLC Schema Files (.sql)
  ↓ sqlc generates
Type-Safe Go Code (sqlcgen/*.go)
  ↓ Used by
Repository Implementations
```

### Supported Databases

See [Database Architecture](../../../docs/database/architecture.md) for supported backends (MySQL 8.4, SQLite, PostgreSQL). Database type is selected via `database.type` in wallet config files.
