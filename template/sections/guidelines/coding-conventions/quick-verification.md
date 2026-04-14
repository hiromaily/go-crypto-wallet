### Quick Verification Reference

```bash
# Go files
make go-lint && make tidy && make check-build && make go-test

# Database schema (HCL)
make atlas-fmt && make atlas-lint

# SQL queries
make sqlc-validate && make sqlc
```
