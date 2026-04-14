### 🧪 Testing Commands

```bash
# Unit tests
make go-test

# Integration tests (MySQL)
make integration-test

# E2E tests (SQLite)
make btc-e2e-reset P=1

# E2E tests (MySQL)
make btc-e2e-reset P=1 DB=mysql

# Verify build
make go-lint
make check-build
```
