# SQLBoiler to SQLC Migration Status

## Issue #32 - Progress Report

### Completed (Phase 1)

✅ **Foundation Setup**
- Installed sqlc v1.30.0
- Created `tools/sqlc/` directory structure
- Configured `sqlc.yml` for watch database
- Created clean SQL schema files for all 8 watch database tables:
  - `01_btc_tx.sql` - BTC transaction tables (btc_tx, btc_tx_input, btc_tx_output)
  - `02_tx.sql` - Generic transaction table (ETH/XRP/HYT)
  - `03_eth_detail_tx.sql` - Ethereum transaction details
  - `04_xrp_detail_tx.sql` - XRP transaction details
  - `05_payment_request.sql` - Payment requests
  - `06_address.sql` - Address management

✅ **Pattern Establishment (tx table)**
- Created basic CRUD queries for `tx` table in `tools/sqlc/queries/tx.sql`
- Generated sqlc code successfully to `pkg/db/rdb/sqlcgen/`
- Verified generated code compiles

### Decision Made

📌 **Scope**: Focus on **watch database only** (8 tables)
- Keygen database (4 tables) - Future phase
- Sign database (2 tables) - Future phase

This decision was made because:
1. Watch database is the most complex and frequently used
2. Keygen and Sign databases have overlapping table names (seed) requiring separate configs
3. Incremental migration is safer and more manageable

### Remaining Work

✅ **Phase 2 - SQL Queries for All Watch Repositories** (COMPLETED)
1. ✅ tx - Pattern established (6 queries)
2. ✅ btc_tx - Complex queries (11 queries including joins, counts)
3. ✅ btc_tx_input - Insert and query operations (3 queries)
4. ✅ btc_tx_output - Insert and query operations (3 queries)
5. ✅ eth_detail_tx - With JOIN to tx table (8 queries)
6. ✅ xrp_detail_tx - XRP-specific fields with JOIN (8 queries)
7. ✅ payment_request - CRUD operations (6 queries)
8. ✅ address - Address allocation logic (5 queries)

**Total: 50 SQL queries created and generated successfully**

⏳ **Phase 3 - Repository Migration**
- Migrate each repository from sqlboiler to sqlc
- Update interfaces to use sqlc types
- Ensure backward compatibility

⏳ **Phase 4 - Testing**
- Update integration tests
- Verify all existing tests pass
- Add new tests for sqlc-specific functionality

⏳ **Phase 5 - Cleanup**
- Remove sqlboiler dependencies from `go.mod`
- Remove `pkg/models/rdb/` directory
- Update Makefile (remove sqlboiler, add sqlc targets)
- Update documentation

### Files Created/Modified

**Created (Phase 1 & 2):**
- `tools/sqlc/sqlc.yml` - sqlc configuration
- `tools/sqlc/schemas/*.sql` - 6 schema files defining 8 tables
- `tools/sqlc/queries/*.sql` - 9 query files with 50 total queries:
  - `tx.sql`, `address.sql`, `payment_request.sql`
  - `btc_tx.sql`, `btc_tx_input.sql`, `btc_tx_output.sql`
  - `eth_detail_tx.sql`, `xrp_detail_tx.sql`
- `pkg/db/rdb/sqlcgen/*.go` - 10 generated Go files:
  - `db.go` - Queries struct
  - `models.go` - Type definitions for all 8 tables
  - 8 query files (*.sql.go)
- `make/codegen.mk` - Added `make sqlc` target
- `SQLC_MIGRATION_STATUS.md` - Progress tracking document

**Modified:**
- None yet (no breaking changes to existing code)

### Next Steps

The migration requires completing the remaining work outlined above. This is a substantial undertaking that affects:
- 8 repository implementations (~1,865 lines of code)
- Multiple integration tests
- Database interaction patterns throughout the codebase

### Recommendations

1. **Continue incrementally**: Complete one repository at a time
2. **Test thoroughly**: Run integration tests after each repository migration
3. **Maintain compatibility**: Keep both sqlboiler and sqlc working during transition
4. **Document patterns**: Use tx repository as the reference pattern

### Estimated Complexity

- ✅ **SQL queries written**: 50 queries (Phase 2 complete!)
- ⏳ **Code migration remaining**: ~1,500 lines to refactor
- ⏳ **Test updates remaining**: ~200 lines

**Progress: Phase 1 & 2 Complete (40% done)**

This is a **medium-to-large refactoring task** being tackled systematically.

---

## How to Continue

To continue this migration:

1. Pick next repository (recommend: address - simplest)
2. Analyze existing repository methods
3. Write SQL queries for all methods
4. Regenerate sqlc code
5. Create new repository using sqlc
6. Run tests
7. Repeat for remaining repositories

## References

- Issue #32: https://github.com/hiromaily/go-crypto-wallet/issues/32
- sqlc Documentation: https://docs.sqlc.dev/
- Existing repositories: `pkg/repository/watchrepo/`
- Generated models: `pkg/db/rdb/sqlcgen/`
