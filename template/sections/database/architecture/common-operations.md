### Common Operations

#### Database Access

**PostgreSQL**:

```bash
# Access watch database
docker exec -it wallet-postgres psql -U postgres -d watch

# Access keygen database
docker exec -it wallet-postgres psql -U postgres -d keygen

# Access sign database
docker exec -it wallet-postgres psql -U postgres -d sign
```

**MySQL**:

```bash
# Access watch schema
docker exec -it wallet-mysql mysql -uroot -proot watch

# Access keygen schema
docker exec -it wallet-mysql mysql -uroot -proot keygen

# Access sign schema
docker exec -it wallet-mysql mysql -uroot -proot sign
```

#### Schema Export (Backup)

Export schema structure without data:

```bash
# Export all schemas (uses DB_DIALECT, default: postgres)
make dump-schema-all

# Export for specific dialect
make dump-schema-all DB_DIALECT=mysql
make dump-schema-all DB_DIALECT=postgres

# Export individual schemas
make dump-schema-watch
make dump-schema-keygen
make dump-schema-sign
```

Output location: `data/dump/sql/dump_*.{dialect}.sql`

#### Data Export (Full Backup)

**PostgreSQL**:

```bash
docker exec wallet-postgres pg_dump -U postgres watch > backups/watch_$(date +%Y%m%d).sql
docker exec wallet-postgres pg_dump -U postgres keygen > backups/keygen_$(date +%Y%m%d).sql
docker exec wallet-postgres pg_dump -U postgres sign > backups/sign_$(date +%Y%m%d).sql
```

**MySQL**:

```bash
docker exec wallet-mysql mysqldump -uroot -proot watch > backups/watch_$(date +%Y%m%d).sql
docker exec wallet-mysql mysqldump -uroot -proot keygen > backups/keygen_$(date +%Y%m%d).sql
docker exec wallet-mysql mysqldump -uroot -proot sign > backups/sign_$(date +%Y%m%d).sql
```

#### Reset Database

Complete database reset (WARNING: deletes all data):

```bash
# Reset specific dialect
make reset-docker                    # PostgreSQL (default)
make reset-docker DB_DIALECT=mysql   # MySQL

# Remove all database containers
make remove-all-dbs
```

---
