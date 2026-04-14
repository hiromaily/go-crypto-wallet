#### Scenario 2: Adding a New Table

**Example**: Add `audit_log` table to watch schema.

##### Step 1: Add Table Definition to HCL

Edit `tools/atlas/schemas/{db_dialect}/watch.hcl`:

```hcl
table "audit_log" {
  schema = schema.watch

  column "id" {
    type = bigint
    auto_increment = true
  }

  column "entity_type" {
    type = varchar(50)
    null = false
  }

  column "entity_id" {
    type = bigint
    null = false
  }

  column "action" {
    type = enum("create", "update", "delete")
    null = false
  }

  column "user_id" {
    type = varchar(100)
    null = true
  }

  column "changes" {
    type = text
    null = true
  }

  column "created_at" {
    type = datetime
    default = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  index "idx_entity" {
    columns = [column.entity_type, column.entity_id]
  }

  index "idx_created_at" {
    columns = [column.created_at]
  }
}
```

##### Step 2: Follow Same Workflow

Follow Steps 2-10 from Scenario 1.

**Additional Considerations**:

- Create corresponding queries in `tools/sqlc/queries/mysql/audit_log.sql`
- Implement repository interface and implementations for all databases
- Add integration tests for the new table

---
