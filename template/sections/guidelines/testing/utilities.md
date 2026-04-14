### Test Utilities

Test utilities are co-located with the package they support in a `testutil/` subdirectory.
A global `pkg/testutil/` package is **not used** in this project.

**Examples:**

- `pkg/db/testutil/` — PostgreSQL connection helpers for database tests

**Rule**: Create a `testutil/` subdirectory inside the package the utilities belong to,
never a standalone top-level `pkg/testutil/` or `internal/testutil/` package.
