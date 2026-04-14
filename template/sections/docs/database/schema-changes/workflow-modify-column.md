#### Scenario 3: Modifying Existing Column

**Example**: Change `wallet_address` from `VARCHAR(255)` to `VARCHAR(500)`.

##### Step 1: Modify HCL Schema

Edit `tools/atlas/schemas/{db_dialect}/watch.hcl`:

```hcl
table "address" {
  // ...
  column "wallet_address" {
    type = varchar(500)  // Changed from 255
    null = false
  }
  // ...
}
```

##### Step 2-10: Follow Same Workflow

Same as Scenario 1, Steps 2-10.

**Important**: Atlas will generate an `ALTER TABLE` migration.

---
