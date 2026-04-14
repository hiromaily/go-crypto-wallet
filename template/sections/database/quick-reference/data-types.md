### 🔄 Data Type Mapping

| Concept | MySQL | SQLite | PostgreSQL |
|---------|-------|--------|------------|
| **Auto ID** | `BIGINT AUTO_INCREMENT` | `INTEGER AUTOINCREMENT` | `BIGSERIAL` |
| **Boolean** | `TINYINT(1)` | `INTEGER (0/1)` | `BOOLEAN` |
| **Enum** | `ENUM('a','b')` | `TEXT CHECK(...)` | `TEXT CHECK(...)` |
| **Decimal** | `DECIMAL(26,10)` | `TEXT` | `NUMERIC(26,10)` |
| **Timestamp** | `DATETIME` | `TEXT (ISO8601)` | `TIMESTAMP` |
| **Text (sized)** | `VARCHAR(255)` | `TEXT` | `VARCHAR(255)` |
| **Text (large)** | `TEXT` | `TEXT` | `TEXT` |
