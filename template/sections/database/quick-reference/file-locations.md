### 📁 File Locations

```
tools/atlas/
├── schemas/              # ✏️  EDIT HERE - Source of truth
│   ├── watch.hcl
│   ├── keygen.hcl
│   └── sign.hcl
└── migrations/           # 🔒 AUTO-GENERATED - Do not edit
    ├── watch/*.sql
    ├── keygen/*.sql
    └── sign/*.sql

tools/sqlc/
├── queries/
│   ├── mysql/            # ✏️  EDIT HERE - MySQL queries (? placeholders)
│   │   ├── address.sql
│   │   ├── btc_tx.sql
│   │   └── *.sql
│   └── postgres/       # ✏️  EDIT HERE - PostgreSQL queries ($1,$2 placeholders)
│       └── *.sql
├── schemas/
│   ├── mysql/            # 🔄 EXTRACTED - From MySQL dump
│   │   ├── 01_watch.sql
│   │   ├── 02_keygen.sql
│   │   └── 03_sign.sql
│   ├── postgres/       # 🔄 EXTRACTED - From PostgreSQL dump
│   │   └── *.sql
│   └── sqlite/           # ✏️  CONVERTED - Manual type mapping
│       └── *.sql

internal/infrastructure/database/
├── mysql/sqlcgen/        # 🔒 AUTO-GENERATED
├── sqlite/sqlcgen/       # 🔒 AUTO-GENERATED
└── postgres/sqlcgen/   # 🔒 AUTO-GENERATED (coming soon)
```

**Legend**:

- ✏️  Manual editing allowed/required
- 🔒 Auto-generated - Do not edit
- 🔄 Extracted from database
