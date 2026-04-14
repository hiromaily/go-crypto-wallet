### 5. Dialect Differences

| Feature | PostgreSQL | MySQL |
|---|---|---|
| ID generation | `identity { generated = BY_DEFAULT }` | `auto_increment = true` |
| Enum type | Named type: `enum "coin" { values = [...] }` | Inline: `enum("a", "b")` |
| Timestamp | `timestamptz` (timezone-aware) | `datetime` |
| Schema reference | `schema "public"` | `schema "watch"` / `"keygen"` / `"sign"` |
| Numeric | `numeric(26,10)` | `decimal(26,10)` |
| Binary data | `bytea` | `blob` |
| Default time | `sql("now()")` | `sql("CURRENT_TIMESTAMP")` |

---
