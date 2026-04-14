### 🗄️ Database Configuration

#### MySQL (Production)

```toml
[database]
type = "mysql"

[database.mysql]
host = "127.0.0.1:3306"
dbname = "watch"  # or "keygen", "sign"
user = "hiromaily"
pass = "hiromaily"
```

#### SQLite (E2E Testing)

```toml
[database]
type = "sqlite"

[database.sqlite]
path = "./data/sqlite/btc/e2e.db"
debug = true
```

#### PostgreSQL (Coming Soon)

```toml
[database]
type = "postgres"

[database.postgres]
host = "127.0.0.1"
port = 5432
dbname = "watch"  # or "keygen", "sign"
user = "hiromaily"
pass = "hiromaily"
sslmode = "prefer"
```
