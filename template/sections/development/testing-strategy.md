## Testing Strategy

### By Layer

| Layer | Testing Approach | Dependencies |
|-------|------------------|--------------|
| Domain | Unit tests, pure functions | None (no mocks needed) |
| Application | Unit tests with mocked ports | Mock interfaces |
| Infrastructure | Integration tests | Real external systems |
| Interface Adapters | Integration tests | Full stack |

### Test Organization

- Unit tests: Same package (`*_test.go`)
- Integration tests: Build tag `//go:build integration`
- Test utilities: `pkg/testutil/` and `**/testutil/`
