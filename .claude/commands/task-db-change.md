# Database Change Task - {description}

## Task Type

Database Change (データベース変更)

## Required Context

このタスクを開始する前に、以下のドキュメントを読み込んでください：

### 必須ドキュメント

1. **DB Change Context**: `docs/ai-agents/task-contexts/db-change.md`
2. **Database Management**: `agents/database.md`
3. **Code Generation**: `agents/code-generation.md`

### 確認が必要なファイル

1. **Schema Files**: `tools/atlas/schemas/*.hcl`
2. **SQLC Queries**: `tools/sqlc/*.sql`
3. **SQLC Config**: `tools/sqlc/sqlc.yml`

## Parameters

- `{description}`: 変更の概要（例: "labelカラム追加", "新しいテーブル作成"）

## Required Tools

作業開始前にツールバージョンを確認:

```bash
atlas version    # v1.0.0 必須
sqlc version     # 最新推奨
```

## Process

### Step 1: コンテキストのロード

上記の必須ドキュメントを読み込む

### Step 2: 現状確認

```bash
# 現在のスキーマを確認
cat tools/atlas/schemas/*.hcl

# 既存のマイグレーションを確認
ls tools/atlas/migrations/
```

### Step 3: スキーマ編集

`tools/atlas/schemas/` 内のHCLファイルを編集

### Step 4: フォーマットとリント

```bash
make atlas-fmt
make atlas-lint
```

### Step 5: マイグレーション生成

```bash
make atlas-dev-reset
```

### Step 6: Docker環境で検証

```bash
docker compose down -v
docker compose up -d
# ログを確認
docker compose logs mysql
```

### Step 7: SQLCクエリ追加/更新

`tools/sqlc/` にクエリを追加/更新

### Step 8: SQLCコード生成

```bash
make sqlc
```

### Step 9: Repository実装（必要な場合）

`internal/infrastructure/repository/` に実装を追加

### Step 10: 最終検証

```bash
make go-lint && make tidy && make check-build && make gotest
```

### Step 11: コミット & PR

```bash
git add tools/atlas/schemas/
git add tools/atlas/migrations/
git add tools/sqlc/
git add internal/infrastructure/database/sqlc/
git add internal/infrastructure/repository/  # 変更がある場合

git commit -m "feat(db): {description}

- Schema changes
- Migration generated
- SQLC queries updated"

gh pr create --title "DB: {description}"
```

## ⚠️ Important Notes

1. **自動生成ファイルを手動編集しない**: `DO NOT EDIT` コメントがあるファイル
2. **Atlas v1.0.0 必須**: バージョンが異なると互換性問題が発生
3. **Docker環境で検証**: マイグレーションが正常に適用されることを確認

## Examples

```
/task-db-change addressテーブルにlabelカラム追加
/task-db-change nonce_commitmentsテーブル作成
/task-db-change txテーブルにインデックス追加
```

## Related Documents

- [DB Change Context](../../docs/ai-agents/task-contexts/db-change.md)
- [Database Management](../../agents/database.md)
- [Code Generation](../../agents/code-generation.md)

