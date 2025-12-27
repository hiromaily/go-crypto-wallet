# Atlas Configuration for go-crypto-wallet
# Manages three separate schemas: watch, keygen, sign
# Uses HCL schema files for declarative schema management

# Variable to control destructive changes (drop schema/table/column)
variable "destructive" {
  type    = bool
  default = false
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

# Local development environment - Watch schema
env "local_watch" {
  url = "mysql://root:root@127.0.0.1:3306/watch?charset=utf8mb4&parseTime=True&loc=Local"
  src = "file://schemas/watch.hcl"
  migration {
    dir = "file://migrations/watch"
  }
  # Optional: Use a dev database for testing schema changes
  # dev = "docker://mysql/8/dev"
}

# Local development environment - Keygen schema
env "local_keygen" {
  url = "mysql://root:root@127.0.0.1:3306/keygen?charset=utf8mb4&parseTime=True&loc=Local"
  src = "file://schemas/keygen.hcl"
  migration {
    dir = "file://migrations/keygen"
  }
}

# Local development environment - Sign schema
env "local_sign" {
  url = "mysql://root:root@127.0.0.1:3306/sign?charset=utf8mb4&parseTime=True&loc=Local"
  src = "file://schemas/sign.hcl"
  migration {
    dir = "file://migrations/sign"
  }
}

# Usage examples:
#
# Apply HCL schema directly:
#   atlas schema apply --env local_watch
#   atlas schema apply --env local_keygen
#   atlas schema apply --env local_sign
#
# Show diff between database and HCL schema:
#   atlas schema diff --env local_watch
#   atlas schema diff --env local_keygen
#   atlas schema diff --env local_sign
#
# Generate migration from HCL schema diff:
#   atlas migrate diff --env local_watch --name add_new_table
#   atlas migrate diff --env local_keygen --name update_account_key
#   atlas migrate diff --env local_sign --name add_index
#
# Apply migrations:
#   atlas migrate apply --env local_watch
#   atlas migrate apply --env local_keygen
#   atlas migrate apply --env local_sign
#
# Check migration status:
#   atlas migrate status --env local_watch
#   atlas migrate status --env local_keygen
#   atlas migrate status --env local_sign
#
# With destructive changes enabled:
#   atlas schema apply --env local_watch -var destructive=true

