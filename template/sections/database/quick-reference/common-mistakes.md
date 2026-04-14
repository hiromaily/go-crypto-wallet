### ❌ Common Mistakes to Avoid

| ❌ Don't Do This | ✅ Do This Instead |
|------------------|-------------------|
| Edit migration SQL files | Edit HCL schemas, regenerate migrations |
| Edit generated SQLC code | Modify queries or schemas, regenerate code |
| Create MySQL-only schemas | Ensure SQLite/PostgreSQL equivalents exist |
| Skip `atlas-fmt` and `atlas-lint` | Always format and validate before regenerating |
| Commit without testing | Run full test cycle before commit |
| Use different column names | Maintain identical names across all databases |
| Modify only one database | Update all three databases (MySQL, SQLite, PostgreSQL) |
