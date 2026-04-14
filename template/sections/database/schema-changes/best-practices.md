### Best Practices

#### 1. Always Use HCL as Source of Truth

✅ **DO**:

```bash
# Edit HCL file
vim tools/atlas/schemas/watch.hcl

# Regenerate migrations
make atlas-dev-reset
```

❌ **DON'T**:

```bash
# Never edit migration files directly
vim tools/atlas/migrations/watch/20240215120000.sql  # WRONG!
```

#### 2. Test Locally Before Committing

```bash
# Complete local test cycle
docker compose down -v
docker compose up -d wallet-mysql
make atlas-migrate-docker
make sqlc-all
make go-lint
make check-build
make go-test
```

#### 3. Keep Migrations Small and Focused

✅ **Good**: One logical change per commit

- Add `email` column to `address` table
- Create `audit_log` table
- Add index on `created_at`

❌ **Bad**: Multiple unrelated changes

- Add `email` column, create `audit_log` table, modify `btc_tx` structure

#### 4. Document Schema Changes

Add comments to HCL files:

```hcl
// Added 2024-02-15: Email notifications feature (#566)
column "email" {
  type = varchar(255)
  null = true
}
```

#### 5. Backward Compatibility

When modifying schemas:

- **Adding columns**: Always make them nullable or provide default values
- **Removing columns**: Deprecated first, remove in next major version
- **Changing types**: Ensure data can be migrated without loss

#### 6. Security for Sensitive Columns

For sensitive data (seeds, private keys):

- Use appropriate encryption
- Never log sensitive column values
- Consider separate schemas for hot/cold storage

```hcl
// Keygen schema (OFFLINE - sensitive data)
table "seed" {
  schema = schema.keygen

  column "encrypted_seed" {
    type = text
    null = false
  }

  // Never store unencrypted
}
```

#### 7. Index Strategy

Add indexes for:

- Foreign key columns
- Frequently queried columns
- Columns used in WHERE clauses
- Columns used in ORDER BY

```hcl
index "idx_coin_account" {
  columns = [column.coin, column.account]
}
```

#### 8. Version Control

Commit in this order:

1. HCL schema changes
2. Generated migration files
3. Updated SQLC schema files
4. Generated SQLC code
5. Repository implementations
6. Tests

#### 9. CI/CD Integration

Ensure CI tests schema changes:

```yaml
# .github/workflows/test.yml
- name: Test schema migrations
  run: |
    docker compose up -d wallet-mysql
    make atlas-migrate-docker
    make sqlc-all
    make check-build
    make go-test
```

#### 10. Documentation

Update these files when schema changes:

- `docs/database/architecture.md` - Database architecture
- This file - If new patterns emerge
- `tools/atlas/README.md` - Atlas-specific details
- Schema-specific docs (if applicable)
