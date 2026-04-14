### Important Rules

1. **Never manually edit auto-generated files** - Changes will be overwritten on next generation
2. **Edit source files instead**:
   - Atlas: Edit `tools/atlas/schemas/{db_dialect}/*.hcl` (HCL schema files)
   - SQLC Schemas: **DO NOT EDIT** `tools/sqlc/schemas/{db_dialect}/*.sql` - these are auto-generated from database dumps. Edit `tools/atlas/schemas/{db_dialect}/*.hcl` instead.
   - SQLC Queries: Edit `tools/sqlc/queries/{db_dialect}/*.sql` (manually edited)
   - Mockery: Edit `.mockery.yaml` to add new interfaces, then run `make mockery`
   - Protocol Buffers: Edit `proto/rippleapi/*.proto`
   - ABI: Edit `contracts/token.abi` (or regenerate from Solidity source)
3. **Regenerate after source changes** - Run the appropriate make command after modifying source files
4. **Verify generation** - Run `make check-build` after regenerating to ensure code compiles
