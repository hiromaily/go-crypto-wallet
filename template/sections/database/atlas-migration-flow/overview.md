## Atlas Migration Flow (`atlas migrate diff/apply`)

This document describes the versioned migration workflow using **Atlas v1.1.0** (`atlas migrate diff/apply`) for
the three databases (`watch / keygen / sign`) across **PostgreSQL** and **MySQL** dialects.

> **Note**: SQLite is NOT managed by Atlas migrations. SQLite schemas are manually maintained in `tools/sqlc/schemas/sqlite/` and used only for SQLC code generation.

### Project Structure

```
tools/atlas/
├── atlas.hcl                          # Atlas configuration (env definitions)
├── schemas/
│   ├── postgres/
│   │   ├── watch.hcl                  # Desired schema (HCL) for watch DB
│   │   ├── keygen.hcl                 # Desired schema (HCL) for keygen DB
│   │   └── sign.hcl                   # Desired schema (HCL) for sign DB
│   └── mysql/
│       ├── watch.hcl
│       ├── keygen.hcl
│       └── sign.hcl
└── migrations/
    ├── postgres/
    │   ├── watch/                     # Migration files + atlas.sum
    │   ├── keygen/
    │   └── sign/
    └── mysql/
        ├── watch/
        ├── keygen/
        └── sign/
```
