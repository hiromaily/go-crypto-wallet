### 6. Destructive Change Protection

The `atlas.hcl` configuration includes a `destructive` variable (default: `false`) that controls whether destructive changes (e.g. `DROP TABLE`, `DROP COLUMN`) are allowed during migration generation. Lint rules also enforce naming conventions (lowercase with underscores) and flag dangerous DDL operations.

Admin environments (`admin_*`) are configured without a specific schema/database name, allowing `ModifySchema` operations such as dropping and recreating entire databases.
