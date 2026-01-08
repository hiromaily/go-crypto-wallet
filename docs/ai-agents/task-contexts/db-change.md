---
task_type: db-change
description: データベーススキーマ変更タスク用のコンテキスト
version: 1.0.0
---

# Task: Database Change (データベース変更)

## When to Use

- データベーススキーマを変更する場合
- 新しいテーブルを追加する場合
- カラムの追加/削除/変更
- インデックスの追加/削除
- マイグレーションの作成
- SQLCクエリの追加/変更

## Required Context Documents

### Must Read (必須)

以下のドキュメントは**必ず**読み込んでください：

| Document | Path | Purpose |
|----------|------|---------|
| Database Management | `agents/database.md` | Atlas/SQLC手順、スキーマ管理 |
| Code Generation | `agents/code-generation.md` | 自動生成ファイルの扱い |
| Workflow Guidelines | `agents/workflow.md` | 検証ステップ |

### Must Check (必須確認)

| Resource | Path | Purpose |
|----------|------|---------|
| Schema Files | `tools/atlas/schemas/*.hcl` | 現在のスキーマ定義 |
| SQLC Queries | `tools/sqlc/*.sql` | 既存クエリ |
| SQLC Config | `tools/sqlc/sqlc.yml` | SQLC設定 |

### Conditional Read (条件付き)

| Condition | Document | Path |
|-----------|----------|------|
| 新しいエンティティ追加 | Architecture | `agents/architecture.md` |
| Repository実装変更 | Internal Guidelines | `internal/AGENTS.md` |

## Task-Specific Rules

### 1. スキーマ変更の手順

```
❌ 直接SQLを実行してスキーマ変更
✅ Atlas HCLファイルを編集 → マイグレーション生成
```

**必須の手順:**
1. `tools/atlas/schemas/*.hcl` を編集
2. `make atlas-fmt` でフォーマット
3. `make atlas-lint` でリント
4. `make atlas-dev-reset` でマイグレーション生成
5. Docker環境で検証
6. `make sqlc` でSQLCコード生成

### 2. 自動生成ファイルの扱い

```
❌ 自動生成ファイルを手動編集
✅ 生成元（HCL/SQL）を編集して再生成
```

**DO NOT EDIT** コメントがあるファイル:
- `internal/infrastructure/database/sqlc/*.go`
- `tools/atlas/migrations/*.sql` (生成後)

### 3. マイグレーションの原則

```
❌ 破壊的な変更（データ損失）
✅ 段階的で安全な変更
```

- カラム削除は段階的に行う
- デフォルト値を設定してからNOT NULL追加
- 大きなテーブルの変更は本番影響を考慮

### 4. SQLCクエリの規則

```sql
-- ❌ Bad: 曖昧な命名
-- name: Get :one
SELECT * FROM users WHERE id = ?;

-- ✅ Good: 明確な命名
-- name: GetUserByID :one
SELECT * FROM users WHERE id = ?;
```

## Pre-Task Checklist

- [ ] 現在のスキーマを確認したか（`tools/atlas/schemas/`）
- [ ] 変更がデータ損失を引き起こさないか確認したか
- [ ] 関連するSQLCクエリへの影響を確認したか
- [ ] Docker環境が利用可能か確認したか
- [ ] Atlas v1.0.0がインストールされているか確認したか

## Required Tools

```bash
# ツールバージョン確認
atlas version    # v1.0.0 必須
sqlc version     # 最新推奨

# インストール（必要な場合）
make install-atlas
make install-sqlc
```

## Verification Commands

```bash
# スキーマ変更の検証
make atlas-fmt       # HCLフォーマット
make atlas-lint      # スキーマリント
make atlas-dev-reset # マイグレーション生成

# Docker環境での検証
docker compose down -v
docker compose up -d
# マイグレーションが正常に適用されることを確認

# SQLCコード生成
make sqlc

# ビルドとテスト
make check-build
make gotest
```

## Database Change Workflow

### Step 1: 現状確認

```bash
# 現在のスキーマを確認
cat tools/atlas/schemas/*.hcl

# 既存のマイグレーションを確認
ls tools/atlas/migrations/
```

### Step 2: スキーマ編集

`tools/atlas/schemas/` 内のHCLファイルを編集:

```hcl
# 例: 新しいテーブル追加
table "new_table" {
  schema = schema.wallet

  column "id" {
    type = bigint
    auto_increment = true
  }

  column "name" {
    type = varchar(255)
  }

  column "created_at" {
    type = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }
}
```

### Step 3: フォーマットとリント

```bash
make atlas-fmt    # HCLファイルをフォーマット
make atlas-lint   # スキーマの問題をチェック
```

### Step 4: マイグレーション生成

```bash
make atlas-dev-reset
```

生成されたマイグレーションファイルを確認:
- `tools/atlas/migrations/` に新しいSQLファイルが作成される

### Step 5: Docker環境で検証

```bash
# 環境をリセット
docker compose down -v

# 新しいマイグレーションで起動
docker compose up -d

# ログを確認
docker compose logs mysql
```

### Step 6: SQLCクエリ追加/更新

`tools/sqlc/` にクエリを追加:

```sql
-- name: GetNewTableByID :one
SELECT * FROM new_table WHERE id = ?;

-- name: CreateNewTable :execresult
INSERT INTO new_table (name) VALUES (?);
```

### Step 7: SQLCコード生成

```bash
make sqlc
```

生成されるファイル:
- `internal/infrastructure/database/sqlc/models.go`
- `internal/infrastructure/database/sqlc/new_table.sql.go`

### Step 8: Repository実装（必要な場合）

新しいテーブルに対応するRepositoryを実装:

```
internal/infrastructure/repository/{entity}/
├── repository.go       # Repository実装
└── repository_test.go  # テスト
```

### Step 9: 最終検証

```bash
make go-lint && make tidy && make check-build && make gotest
```

### Step 10: コミット

```bash
git add tools/atlas/schemas/
git add tools/atlas/migrations/
git add tools/sqlc/
git add internal/infrastructure/database/sqlc/
git add internal/infrastructure/repository/  # 変更がある場合

git commit -m "feat(db): add new_table schema and queries

- Add new_table schema in Atlas HCL
- Generate migration for new_table
- Add SQLC queries for new_table
- Implement repository (if applicable)"
```

## Common Scenarios

### Scenario 1: 新しいテーブル追加

```
1. tools/atlas/schemas/ にテーブル定義追加
2. make atlas-fmt && make atlas-lint
3. make atlas-dev-reset
4. docker compose で検証
5. tools/sqlc/ にクエリ追加
6. make sqlc
7. Repository実装
8. make gotest
```

### Scenario 2: カラム追加

```
1. HCLファイルでカラム追加
2. make atlas-fmt && make atlas-lint
3. make atlas-dev-reset
4. 既存データへの影響確認
5. 関連SQLCクエリの更新
6. make sqlc
7. make gotest
```

### Scenario 3: カラム削除（注意が必要）

```
Phase 1: カラムを使用しないように変更
  - アプリケーションコードからカラム参照を削除
  - SQLCクエリから除外
  - デプロイ

Phase 2: カラム削除
  - HCLからカラム削除
  - マイグレーション生成
  - デプロイ
```

## Examples

### Example 1: 新しいエンティティのテーブル追加

```
User: "XRPトランザクション履歴を保存するテーブルを追加して"

Agent Actions:
1. agents/database.md を読み込み
2. 既存のスキーマを確認（tools/atlas/schemas/）
3. 新しいテーブル定義をHCLに追加
4. make atlas-fmt && make atlas-lint
5. make atlas-dev-reset
6. Docker環境で検証
7. SQLCクエリを追加
8. make sqlc
9. 必要に応じてRepository実装
10. 検証コマンド実行
```

### Example 2: 既存テーブルにカラム追加

```
User: "addressテーブルにlabelカラムを追加して"

Agent Actions:
1. agents/database.md を読み込み
2. 既存のaddressテーブル定義を確認
3. HCLにlabelカラムを追加
4. make atlas-fmt && make atlas-lint
5. make atlas-dev-reset
6. 影響を受けるSQLCクエリを更新
7. make sqlc
8. 関連するコードを更新
9. 検証コマンド実行
```

## Testing Requirements

### 必須

- [ ] マイグレーションがDocker環境で正常に適用される
- [ ] 生成されたSQLCコードでビルドが通る
- [ ] 既存のテストがパスする

### 推奨

- [ ] 新しいクエリのテスト追加
- [ ] Repository実装のテスト
- [ ] 統合テストでDB操作を確認

## Troubleshooting

### Atlas関連

```bash
# スキーマの差分確認
atlas schema diff --from "file://tools/atlas/schemas" --to "mysql://..."

# マイグレーションの適用状態確認
atlas migrate status --dir "file://tools/atlas/migrations"
```

### SQLC関連

```bash
# SQLCの検証
sqlc compile

# 詳細なエラー出力
sqlc generate --experimental
```

### Docker関連

```bash
# MySQLコンテナのログ確認
docker compose logs mysql

# MySQLに直接接続
docker compose exec mysql mysql -u wallet -p wallet
```

## Related Documents

- [Database Management](../../../agents/database.md) - Atlas/SQLC詳細手順
- [Code Generation](../../../agents/code-generation.md) - 自動生成ファイルの扱い
- [Architecture Guidelines](../../../agents/architecture.md) - Repository層の設計

