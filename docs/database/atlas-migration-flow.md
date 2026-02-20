# atlas migration flow by `atlas migrate diff/apply`

以下は **Atlas v1.0.0 の versioned migrations（`atlas migrate diff/apply`）** を前提に、
`wallet / keygen / sign` の3つのDBを、**PostgreSQL / MySQL** に対して運用するフローを **mermaid** で整理したもの。
（`migrate diff` で migration を生成し、`migrate apply` で適用する流れ）

> **Note**: SQLite は Atlas による migration 管理の対象外。SQLite スキーマは `tools/sqlc/schemas/sqlite/` で手動管理し、SQLC によるコード生成のみ行う。

## 1. 初回: DBデザイン → migration生成 → DBへ適用

```mermaid
flowchart TD
  A[Start] --> B[理想スキーマ定義を作成\n例: schema.hcl または schema.sql または ORM]
  B --> C[atlas.hcl を用意\nenv を DBと方言ごとに定義\n例: pg_wallet, mysql_wallet]
  C --> D{対象DBはどれ?}

  D --> W[wallet]
  D --> K[keygen]
  D --> S[sign]

  %% wallet dialect
  W --> Wd{方言はどれ?}
  Wd --> Wpg[PostgreSQL]
  Wd --> Wmy[MySQL]

  %% keygen dialect
  K --> Kd{方言はどれ?}
  Kd --> Kpg[PostgreSQL]
  Kd --> Kmy[MySQL]

  %% sign dialect
  S --> Sd{方言はどれ?}
  Sd --> Spg[PostgreSQL]
  Sd --> Smy[MySQL]

  subgraph G[初回 migration 生成 baseline]
    G1[atlas migrate diff\n--env 対象env\n--dir file://migrations/方言/DB\n--to 理想スキーマ\n--dev-url 方言別dev DB]
  end

  subgraph AP[DBへ適用]
    AP1[atlas migrate apply\n--env 対象env\n対象DBへ pending を適用]
  end

  Wpg --> G1 --> AP1
  Wmy --> G1 --> AP1

  Kpg --> G1 --> AP1
  Kmy --> G1 --> AP1

  Spg --> G1 --> AP1
  Smy --> G1 --> AP1

  AP1 --> Z[Done]
```

ポイント:

- **`atlas migrate diff` は `dev-url のdialect` に合わせて SQL migration を生成**する
- **`atlas migrate apply` が migrations directory の未適用分をDBに適用**する
- 3 DB×2 dialectを全部カバーするなら、運用上は **migrations を dialect/DBごとに分ける**のが事故りにくい（e.g. `migrations/postgres/wallet` など）。

---

## 2. 変更: schema変更 → 更新migration作成 → DBへ適用

```mermaid
flowchart TD
  A[Start] --> B[理想スキーマ定義を更新\n例: schema.hcl または schema.sql または ORM]
  B --> C{変更対象DBはどれ?}

  C --> W[wallet]
  C --> K[keygen]
  C --> S[sign]

  %% wallet dialect
  W --> Wd{方言はどれ?}
  Wd --> Wpg[PostgreSQL]
  Wd --> Wmy[MySQL]

  %% keygen dialect
  K --> Kd{方言はどれ?}
  Kd --> Kpg[PostgreSQL]
  Kd --> Kmy[MySQL]

  %% sign dialect
  S --> Sd{方言はどれ?}
  Sd --> Spg[PostgreSQL]
  Sd --> Smy[MySQL]

  subgraph M[新しい migration を作成]
    M1[atlas migrate diff\nname を指定\n--env 対象env\n--dir migrations の場所\n--to 理想スキーマ\n--dev-url 方言別dev DB]
  end

  subgraph Q[任意: チェックとテスト]
    L1[atlas migrate lint\nmigration の整合性と危険DDLを検知]
    T1[任意: dev DB へ apply してテスト]
  end

  subgraph AP[本番または対象DBへ適用]
    A1[atlas migrate apply\n--env 対象env\n対象DBへ pending を適用]
  end

  Wpg --> M1 --> L1 --> T1 --> A1
  Wmy --> M1 --> L1 --> T1 --> A1

  Kpg --> M1 --> L1 --> T1 --> A1
  Kmy --> M1 --> L1 --> T1 --> A1

  Spg --> M1 --> L1 --> T1 --> A1
  Smy --> M1 --> L1 --> T1 --> A1

  A1 --> Z[Done]
```

補足:

- `migrate diff` は **「現在の状態 → 理想状態」への差分 migration を自動生成**するコマンド
- `migrate apply` は **未適用 migration を順に適用**する
- `atlas.hcl` の env で「どのDBに」「どの migrations dir を」「どの schema source から」扱うかをまとめて定義できる
