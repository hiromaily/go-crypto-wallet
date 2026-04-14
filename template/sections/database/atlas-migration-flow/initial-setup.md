### 1. Initial Setup: Schema Design -> Generate Baseline Migration -> Apply to DB

```mermaid
flowchart TD
  A[Start] --> B[Create desired schema definitions\nin HCL format\ne.g. tools/atlas/schemas/postgres/watch.hcl]
  B --> C[Configure atlas.hcl\nDefine env per DB and dialect\ne.g. local_postgres_watch, local_mysql_watch]
  C --> D{Target database?}

  D --> W[watch]
  D --> K[keygen]
  D --> S[sign]

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

  subgraph G[Generate baseline migration]
    G1[atlas migrate diff\n--env target_env\n--dir file://migrations/dialect/db\n--to desired schema HCL\n--dev-url dialect-specific dev DB]
  end

  subgraph AP[Apply to DB]
    AP1[atlas migrate apply\n--env target_env\nApply pending migrations to target DB]
  end

  Wpg --> G1 --> AP1
  Wmy --> G1 --> AP1

  Kpg --> G1 --> AP1
  Kmy --> G1 --> AP1

  Spg --> G1 --> AP1
  Smy --> G1 --> AP1

  AP1 --> Z[Done]
```

Key points:

- **`atlas migrate diff`** generates SQL migration files matching the `dev-url` dialect
- **`atlas migrate apply`** applies pending (unapplied) migrations from the migrations directory to the target DB
- Migrations are **separated by dialect and database** to prevent cross-dialect issues (e.g. `migrations/postgres/watch/`, `migrations/mysql/keygen/`)

---
