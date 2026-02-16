# PostgreSQL Watch Schema Migrations

## Status

Migration directory created, awaiting Atlas migration generation.

## How to Generate Migrations

The initial PostgreSQL migration must be generated using Atlas CLI with a running PostgreSQL Docker container:

```bash
# 1. Ensure Docker is running
docker ps

# 2. Start PostgreSQL container (if using Docker Compose)
docker compose up -d wallet-db-postgresql

# 3. Generate migration from HCL schema
cd tools/atlas
atlas migrate diff initial_schema \
  --env local_postgresql_watch \
  --config file://atlas.hcl

# 4. Verify migration was created
ls -la migrations/postgresql_watch/
```

## Expected Output

After running the atlas migrate diff command, this directory should contain:
- `<timestamp>_initial_schema.sql` - PostgreSQL DDL migration file
- `atlas.sum` - Migration checksum file

## Verification

The generated migration should:
- Use PostgreSQL syntax (BIGSERIAL, CHECK constraints, TIMESTAMP)
- Include all tables from watch.hcl schema
- Have proper COMMENT statements for documentation
- Not include MySQL-specific syntax (enum, AUTO_INCREMENT, backticks)

## Configuration

Atlas environment is configured in `../atlas.hcl`:
- Environment: `local_postgresql_watch`
- Database URL: `postgres://postgres:postgres@localhost:5432/watch?sslmode=disable`
- Schema source: `file://schemas/watch.hcl`
- Dev database: `docker://postgres/18/watch`
