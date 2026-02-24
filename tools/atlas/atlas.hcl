# Atlas Configuration for go-crypto-wallet
# Manages three separate schemas: watch, keygen, sign
# Uses HCL schema files for declarative schema management
# Atlas version: 1.1.0

# Variable to control destructive changes (drop schema/table/column)
variable "destructive" {
  type        = bool
  default     = false
  description = "Allow destructive changes (drop schema/table/column)"
}

# Diff configuration to prevent accidental data loss
diff {
  skip {
    drop_schema = !var.destructive
    drop_table  = !var.destructive
    drop_column = !var.destructive
  }
}

# Lint configuration for schema validation
lint {
  # Error on destructive changes (e.g., dropping tables/columns)
  destructive {
    error = true
  }

  # Naming conventions
  naming {
    match   = "^[a-z][a-z0-9_]*$"
    message = "Table and column names must be lowercase with underscores"
  }
}

# Local development environment - Watch schema
# Used for migrations and schema apply (schema-scoped operations)
env "local_mysql_watch" {
  url     = "mysql://root:root@127.0.0.1:3306/watch?charset=utf8mb4&parseTime=True&loc=Local"
  src     = "file://schemas/mysql/watch.hcl"
  schemas = ["watch"]
  migration {
    dir = "file://migrations/mysql/watch"
  }
  # Dev database with schema name - prevents CREATE DATABASE in migrations
  dev = "docker://mysql/8/watch"
}

# Local development environment - Keygen schema
# Used for migrations and schema apply (schema-scoped operations)
env "local_mysql_keygen" {
  url     = "mysql://root:root@127.0.0.1:3306/keygen?charset=utf8mb4&parseTime=True&loc=Local"
  src     = "file://schemas/mysql/keygen.hcl"
  schemas = ["keygen"]
  migration {
    dir = "file://migrations/mysql/keygen"
  }
  # Dev database with schema name - prevents CREATE DATABASE in migrations
  dev = "docker://mysql/8/keygen"
}

# Local development environment - Sign schema
# Used for migrations and schema apply (schema-scoped operations)
env "local_mysql_sign" {
  url     = "mysql://root:root@127.0.0.1:3306/sign?charset=utf8mb4&parseTime=True&loc=Local"
  src     = "file://schemas/mysql/sign.hcl"
  schemas = ["sign"]
  migration {
    dir = "file://migrations/mysql/sign"
  }
  # Dev database with schema name - prevents CREATE DATABASE in migrations
  dev = "docker://mysql/8/sign"
}

###############################################################################
# Admin environments for schema-level operations (clean, drop schema, etc.)
# These environments use URLs without schema name to allow ModifySchema operations
###############################################################################

# Admin environment - Watch schema (for atlas schema clean)
env "admin_mysql_watch" {
  url     = "mysql://root:root@127.0.0.1:3306/?charset=utf8mb4&parseTime=True&loc=Local"
  src     = "file://schemas/mysql/watch.hcl"
  schemas = ["watch"]
}

# Admin environment - Keygen schema (for atlas schema clean)
env "admin_mysql_keygen" {
  url     = "mysql://root:root@127.0.0.1:3306/?charset=utf8mb4&parseTime=True&loc=Local"
  src     = "file://schemas/mysql/keygen.hcl"
  schemas = ["keygen"]
}

# Admin environment - Sign schema (for atlas schema clean)
env "admin_mysql_sign" {
  url     = "mysql://root:root@127.0.0.1:3306/?charset=utf8mb4&parseTime=True&loc=Local"
  src     = "file://schemas/mysql/sign.hcl"
  schemas = ["sign"]
}

###############################################################################
# PostgreSQL environments for local development
# Uses PostgreSQL 18.2 Docker container
###############################################################################

# Local development environment - PostgreSQL Watch schema
env "local_postgres_watch" {
  url     = "postgres://postgres:postgres@localhost:5432/watch?sslmode=disable"
  src     = "file://schemas/postgres/watch.hcl"
  schemas = ["public"]
  migration {
    dir = "file://migrations/postgres/watch"
  }
  dev = "docker://postgres/18/watch"
}

# Local development environment - PostgreSQL Keygen schema
env "local_postgres_keygen" {
  url     = "postgres://postgres:postgres@localhost:5432/keygen?sslmode=disable"
  src     = "file://schemas/postgres/keygen.hcl"
  schemas = ["public"]
  migration {
    dir = "file://migrations/postgres/keygen"
  }
  dev = "docker://postgres/18/keygen"
}

# Local development environment - PostgreSQL Sign schema
env "local_postgres_sign" {
  url     = "postgres://postgres:postgres@localhost:5432/sign?sslmode=disable"
  src     = "file://schemas/postgres/sign.hcl"
  schemas = ["public"]
  migration {
    dir = "file://migrations/postgres/sign"
  }
  dev = "docker://postgres/18/sign"
}

###############################################################################
# PostgreSQL Admin environments for schema-level operations
# These environments use URLs without database name for ModifySchema operations
###############################################################################

# Admin environment - PostgreSQL Watch schema
env "admin_postgres_watch" {
  url     = "postgres://postgres:postgres@localhost:5432/?sslmode=disable"
  src     = "file://schemas/postgres/watch.hcl"
  schemas = ["public"]
}

# Admin environment - PostgreSQL Keygen schema
env "admin_postgres_keygen" {
  url     = "postgres://postgres:postgres@localhost:5432/?sslmode=disable"
  src     = "file://schemas/postgres/keygen.hcl"
  schemas = ["public"]
}

# Admin environment - PostgreSQL Sign schema
env "admin_postgres_sign" {
  url     = "postgres://postgres:postgres@localhost:5432/?sslmode=disable"
  src     = "file://schemas/postgres/sign.hcl"
  schemas = ["public"]
}

# Usage examples:
#
# Apply HCL schema directly:
#   atlas schema apply --env local_mysql_watch
#   atlas schema apply --env local_mysql_keygen
#   atlas schema apply --env local_mysql_sign
#
# Show diff between database and HCL schema:
#   atlas schema diff --env local_mysql_watch
#   atlas schema diff --env local_mysql_keygen
#   atlas schema diff --env local_mysql_sign
#
# Generate migration from HCL schema diff:
#   atlas migrate diff --env local_mysql_watch --name add_new_table
#   atlas migrate diff --env local_mysql_keygen --name update_account_key
#   atlas migrate diff --env local_mysql_sign --name add_index
#
# Apply migrations:
#   atlas migrate apply --env local_mysql_watch
#   atlas migrate apply --env local_mysql_keygen
#   atlas migrate apply --env local_mysql_sign
#
# Check migration status:
#   atlas migrate status --env local_mysql_watch
#   atlas migrate status --env local_mysql_keygen
#   atlas migrate status --env local_mysql_sign
#
# With destructive changes enabled:
#   atlas schema apply --env local_mysql_watch -var destructive=true

