### Common Patterns

#### Pattern 1: Adding Created/Updated Timestamps

**MySQL**:

```hcl
column "created_at" {
  type = datetime
  default = sql("CURRENT_TIMESTAMP")
}

column "updated_at" {
  type = datetime
  null = true
}
```

**SQLite** (in `.sql` file):

```sql
created_at TEXT DEFAULT CURRENT_TIMESTAMP,
updated_at TEXT
```

**PostgreSQL**:

```sql
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP
```

#### Pattern 2: Nullable vs NOT NULL

**Prefer NULL for optional fields**:

```hcl
column "email" {
  type = varchar(255)
  null = true  // Optional field
}
```

**Use NOT NULL for required fields**:

```hcl
column "wallet_address" {
  type = varchar(500)
  null = false  // Required field
}
```

#### Pattern 3: Foreign Keys

**MySQL/PostgreSQL** (HCL):

```hcl
table "btc_tx_input" {
  // ...
  column "tx_id" {
    type = bigint
    null = false
  }

  foreign_key "fk_btc_tx_input_tx" {
    columns     = [column.tx_id]
    ref_columns = [table.btc_tx.column.id]
    on_delete   = CASCADE
    on_update   = CASCADE
  }
}
```

**SQLite** (limited FK support):

```sql
-- SQLite supports FK but they're not enforced by default
-- PRAGMA foreign_keys = ON; required
tx_id INTEGER NOT NULL,
FOREIGN KEY (tx_id) REFERENCES btc_tx(id) ON DELETE CASCADE
```

#### Pattern 4: Indexes for Performance

**HCL**:

```hcl
table "address" {
  // ...
  index "idx_coin_account" {
    columns = [column.coin, column.account]
  }

  index "idx_is_allocated" {
    columns = [column.is_allocated]
  }
}
```

#### Pattern 5: Unique Constraints

**HCL**:

```hcl
table "address" {
  // ...
  column "wallet_address" {
    type = varchar(500)
    null = false
  }

  index "idx_wallet_address_unique" {
    unique = true
    columns = [column.wallet_address]
  }
}
```
