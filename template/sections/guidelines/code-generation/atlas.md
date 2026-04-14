### Database Migrations (Atlas)

**Tool**: [Atlas](https://atlasgo.io/)
**Source**: `tools/atlas/schemas/{db_dialect}/*.hcl` (HCL schema files)
**Command**: `make atlas-dev-reset` (regenerate from scratch)

**Generated Files**:

- `tools/atlas/migrations/watch/*.sql` - Watch schema migrations
- `tools/atlas/migrations/keygen/*.sql` - Keygen schema migrations
- `tools/atlas/migrations/sign/*.sql` - Sign schema migrations
- `tools/atlas/migrations/*/atlas.sum` - Migration checksums

**Note**: See [Database Management Guidelines](../../../docs/database/architecture.md) for detailed workflow.
