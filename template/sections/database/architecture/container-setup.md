### Architecture

#### Container Setup

```yaml
services:
  # PostgreSQL (default)
  wallet-postgres:
    image: postgres:18.2
    profiles: ["postgres"]
    ports:
      - "${POSTGRESQL_PORT:-5432}:5432"
    volumes:
      - wallet-postgres:/var/lib/postgres
      - "./docker/postgres/init.d:/docker-entrypoint-initdb.d"
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres

  # MySQL (alternative)
  wallet-mysql:
    image: mysql:8.4
    profiles: ["mysql"]
    ports:
      - "${MYSQL_PORT:-3306}:3306"
    volumes:
      - wallet-mysql:/var/lib/mysql
      - "./docker/mysql/conf.d:/etc/mysql/conf.d"
      - "./docker/mysql/init.d:/docker-entrypoint-initdb.d"
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_USER: hiromaily
      MYSQL_PASSWORD: hiromaily
```

Both services include health checks and are activated via Docker Compose profiles (`--profile postgres` or `--profile mysql`).

#### Migration Services

Atlas migration services run automatically when databases start. Each schema has a dedicated migration service:

```yaml
# Example: PostgreSQL watch migration
wallet-postgres-migrate-watch:
  image: arigaio/atlas:1.1.0
  profiles: ["postgres"]
  command:
    - migrate
    - apply
    - --dir
    - "file://migrations/postgres/watch"
    - --url
    - "postgres://postgres:postgres@wallet-postgres:5432/watch?sslmode=disable"
  depends_on:
    wallet-postgres:
      condition: service_healthy
  restart: "no"
```

Migration services exist for: `watch`, `keygen`, `sign`, `sign2` (both MySQL and PostgreSQL).

#### Directory Structure

```
docker/
├── mysql/
│   ├── archive/                       # Archived SQL schema files (reference only)
│   ├── conf.d/
│   │   └── custom.cnf                 # Server-level configuration (utf8mb4)
│   ├── init.d/
│   │   └── 01_init_all_schemas_2.sql  # Creates watch, keygen, sign, sign2 databases
│   └── insert/
│       └── ganache.example.sql        # Test data for Ganache
├── postgres/
│   └── init.d/
│       └── 01_create_databases.sh     # Creates watch, keygen, sign, sign2 databases
```

**Note:** Schema definitions are managed by Atlas (HCL files in `tools/atlas/schemas/`). The init scripts only create empty databases.

#### Initialization Process

When the container starts for the first time:

1. **Database Creation**: Init scripts create four empty databases

   **PostgreSQL** (`01_create_databases.sh`):

   ```sql
   CREATE DATABASE watch;
   CREATE DATABASE keygen;
   CREATE DATABASE sign;
   CREATE DATABASE sign2;
   ```

   **MySQL** (`01_init_all_schemas_2.sql`):

   ```sql
   CREATE DATABASE `watch` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   CREATE DATABASE `keygen` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   CREATE DATABASE `sign` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   CREATE DATABASE `sign2` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```

2. **User/Permission Setup**: Grants privileges to application users

3. **Schema Migration**: Atlas migration services automatically apply migrations after the database is healthy

---
