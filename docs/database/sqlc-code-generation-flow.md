# sqlc code generation flowchart

```mermaid
flowchart TD
  A[Start] --> B[前提: migration が各DBに完全適用済み\nwallet, keygen, sign]

  B --> C{対象のDB種別はどれ?}

  C --> PG[PostgreSQL]
  C --> MY[MySQL]
  C --> SQ[SQLite]

  %% PostgreSQL path
  PG --> PG1[Schema dump\nコマンド例\npg_dump --schema-only --no-owner --no-privileges \n  --dbname "$PG_DSN_wallet" > gen/schema/wallet.pg.sql\npg_dump --schema-only --no-owner --no-privileges \n  --dbname "$PG_DSN_keygen" > gen/schema/keygen.pg.sql\npg_dump --schema-only --no-owner --no-privileges \n  --dbname "$PG_DSN_sign" > gen/schema/sign.pg.sql]

  %% MySQL path
  MY --> MY1[Schema dump\nコマンド例\nmysqldump --no-data --skip-comments --routines=0 --triggers=0 \n  --databases wallet > gen/schema/wallet.my.sql\nmysqldump --no-data --skip-comments --routines=0 --triggers=0 \n  --databases keygen > gen/schema/keygen.my.sql\nmysqldump --no-data --skip-comments --routines=0 --triggers=0 \n  --databases sign > gen/schema/sign.my.sql]

  %% SQLite path
  SQ --> SQ1[Schema extract by Atlas\nコマンド例\natlas schema inspect -u \"sqlite://wallet.sqlite\" \n  --format \"{{ sql . }}\" > gen/schema/wallet.sqlite.sql\natlas schema inspect -u \"sqlite://keygen.sqlite\" \n  --format \"{{ sql . }}\" > gen/schema/keygen.sqlite.sql\natlas schema inspect -u \"sqlite://sign.sqlite\" \n  --format \"{{ sql . }}\" > gen/schema/sign.sqlite.sql]

  %% Common normalize
  PG1 --> N[不要な情報を除去して正規化\nコマンド例\ncat gen/schema/*.sql | ./scripts/strip_schema.sh > gen/schema/normalized.sql]
  MY1 --> N
  SQ1 --> N

  %% sqlc generation
  N --> G[sqlc generate で Go を生成\nコマンド例\nsqlc generate]

  G --> Z[Done]
```