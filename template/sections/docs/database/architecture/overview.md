## Database Architecture

This document describes the database architecture and operations for the go-crypto-wallet project.

---

### Table of Contents

- [Overview](#overview)
- [Supported Databases](#supported-databases)
- [Architecture](#architecture)
- [Schema Design](#schema-design)
- [Setup and Configuration](#setup-and-configuration)
- [Common Operations](#common-operations)
- [Database Management](#database-management)
- [Schema Migrations with Atlas](#schema-migrations-with-atlas)
- [SQLC Code Generation](#sqlc-code-generation)
- [SQLite for E2E Testing](#sqlite-for-e2e-testing)
- [Troubleshooting](#troubleshooting)
- [Migration Guide](#migration-guide)

---

### Overview

The project supports **three database backends**:

| Database | Version | Use Case | Features |
|----------|---------|----------|----------|
| **PostgreSQL** | 18.2 | Production (default) | Docker container, schema separation, `identity` columns |
| **MySQL** | 8.4 | Production (alternative) | Docker container, schema separation, `auto_increment` |
| **SQLite** | - | E2E testing, CI/CD | Local file, fast startup, no Docker required |

All backends use **four separate databases/schemas** to manage wallet data:

- **`watch`**: Online wallet data (addresses, transactions, payment requests)
- **`keygen`**: Key generation data (seeds, account keys, full public keys)
- **`sign`**: Signing wallet data (auth account keys, seeds)
- **`sign2`**: Second signing wallet (same schema as `sign`, separate database)

This approach provides:

- Reduced resource usage (single DB instance per dialect)
- Simplified deployment and maintenance
- Data isolation through database/schema separation
- Easier backup and restore operations
- Single point of configuration

---
