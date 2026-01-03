# Scripts

## migration.sh

データベースマイグレーション管理スクリプト

### 前提条件

`golang-migrate`をインストールする必要があります：

```bash
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# Windows
scoop install migrate
```

### 使い方

```bash
# マイグレーションを適用
./scripts/migration.sh up

# マイグレーションをロールバック（1つ）
./scripts/migration.sh down

# マイグレーションをロールバック（複数）
./scripts/migration.sh down 2

# 現在のマイグレーションバージョンを確認
./scripts/migration.sh version

# 新しいマイグレーションファイルを作成
./scripts/migration.sh create add_new_column

# マイグレーションバージョンを強制設定
./scripts/migration.sh force 20251223000001

# すべてのテーブルを削除（要確認）
./scripts/migration.sh drop
```

### 環境変数

`.env`ファイルまたは環境変数で設定：

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=guess_title_game
DB_SSL_MODE=disable
```

### トラブルシューティング

#### "Dirty database version" エラー

マイグレーションが途中で失敗した場合：

```bash
# 現在のバージョンを確認
./scripts/migration.sh version

# バージョンを強制的に設定（version番号は適切なものに変更）
./scripts/migration.sh force 20251223000001
```

#### データベース接続エラー

1. PostgreSQLが起動しているか確認
2. `.env`ファイルの設定を確認
3. データベースが作成されているか確認

```bash
# データベースを作成
psql -U postgres -c "CREATE DATABASE guess_title_game;"
```
