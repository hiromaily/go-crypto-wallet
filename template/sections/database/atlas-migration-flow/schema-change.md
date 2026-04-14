### 2. Schema Change: Update Schema -> Generate New Migration -> Apply to DB

```mermaid
flowchart TD
  A[Start] --> B[Update desired schema definition\ne.g. modify tools/atlas/schemas/postgres/watch.hcl]
  B --> C{Target database?}

  C --> W[watch]
  C --> K[keygen]
  C --> S[sign]

  %% watch dialect
  W --> Wd{Dialect?}
  Wd --> Wpg[PostgreSQL]
  Wd --> Wmy[MySQL]

  %% keygen dialect
  K --> Kd{Dialect?}
  Kd --> Kpg[PostgreSQL]
  Kd --> Kmy[MySQL]

  %% sign dialect
  S --> Sd{Dialect?}
  Sd --> Spg[PostgreSQL]
  Sd --> Smy[MySQL]

  subgraph M[Generate new migration]
    M1[atlas migrate diff\nwith migration name\n--env target_env\n--dir migrations directory\n--to desired schema HCL\n--dev-url dialect-specific dev DB]
  end

  subgraph Q[Optional: Validation and testing]
    L1[atlas migrate lint\nChecks integrity and\ndetects dangerous DDL]
    T1[Optional: apply to dev DB for testing]
  end

  subgraph AP[Apply to target DB]
    A1[atlas migrate apply\n--env target_env\nApply pending migrations to target DB]
  end

  Wpg --> M1 --> L1 --> T1 --> A1
  Wmy --> M1 --> L1 --> T1 --> A1

  Kpg --> M1 --> L1 --> T1 --> A1
  Kmy --> M1 --> L1 --> T1 --> A1

  Spg --> M1 --> L1 --> T1 --> A1
  Smy --> M1 --> L1 --> T1 --> A1

  A1 --> Z[Done]
```

Notes:

- `migrate diff` **auto-generates a diff migration from the current state to the desired state**
- `migrate apply` **applies unapplied migrations in order**
- The `atlas.hcl` env consolidates which DB, migrations directory, and schema source to use

---
