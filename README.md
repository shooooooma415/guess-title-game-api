# Guess Title Game API

バックエンドAPIサーバー（Go + Echo + PostgreSQL + WebSocket）

## 機能

- ルーム管理（作成、参加、ゲーム進行）
- リアルタイム通信（WebSocket）
- タイマー機能（5分間のカウントダウン）
- PostgreSQLによる永続化

## 必要要件

- Docker & Docker Compose
- Go 1.24+（ローカル開発の場合）
- PostgreSQL 16+（ローカル開発の場合）

## クイックスタート

### Docker Composeで起動

```bash
# 1. すべてのサービスを起動（DB + Migration + API）
docker compose up -d

# 2. ログを確認
docker compose logs -f api

# 3. ヘルスチェック
curl http://localhost:8080/health

# 4. 停止
docker compose down

# 5. データも削除して停止
docker compose down -v
```

### ローカル開発

```bash
# 1. データベースのみ起動
docker compose -f docker-compose.dev.yml up -d

# 2. 環境変数を設定
cp .env.example .env

# 3. マイグレーション実行
./scripts/migration.sh up

# 4. アプリケーション起動
go run cmd/main.go
```

## API エンドポイント

### HTTP API

| Method | Endpoint | 説明 |
|--------|----------|------|
| GET | `/health` | ヘルスチェック |
| POST | `/api/rooms` | ルーム作成 |
| POST | `/api/user` | ユーザー参加 |
| POST | `/api/rooms/:room_id/start` | ゲーム開始 |
| POST | `/api/rooms/:room_id/topic` | トピック設定 |
| POST | `/api/rooms/:room_id/answer` | 回答送信 |
| POST | `/api/rooms/:room_id/skip-discussion` | 議論スキップ |
| POST | `/api/rooms/:room_id/finish` | ゲーム終了 |

### WebSocket

```
ws://localhost:8080/ws?room_id={room_id}
```

#### クライアント → サーバー

- `CLIENT_CONNECTED` - クライアント接続通知
- `FETCH_PARTICIPANTS` - 参加者リスト取得
- `SUBMIT_TOPIC` - トピック情報送信
- `ANSWERING` - 回答情報送信

#### サーバー → クライアント

- `STATE_UPDATE` - 状態遷移通知
- `PARTICIPANT_UPDATE` - 参加者リスト更新
- `TIMER_TICK` - タイマー更新（毎秒）
- `ERROR` - エラー通知

## データベース

### マイグレーション

```bash
# マイグレーション実行
./scripts/migration.sh up

# ロールバック
./scripts/migration.sh down

# バージョン確認
./scripts/migration.sh version

# 新しいマイグレーション作成
./scripts/migration.sh create add_new_table
```

### テーブル構造

- `users` - ユーザー情報
- `themes` - テーマ情報
- `rooms` - ルーム情報（ゲームデータ含む）
- `participants` - 参加者情報
- `room_emojis` - ルームの絵文字情報

## 開発

### ディレクトリ構造

```
.
├── cmd/                    # エントリーポイント
├── config/                 # 設定管理
├── internal/
│   ├── domain/            # ドメイン層
│   ├── usecase/           # ユースケース層
│   ├── interface/         # インターフェース層
│   └── infrastructure/    # インフラ層
├── db/migrations/         # マイグレーションファイル
├── scripts/               # スクリプト
└── utils/                 # ユーティリティ

```

### ビルド

```bash
# ローカルビルド
go build -o bin/api ./cmd/main.go

# Dockerビルド
docker build -t guess-title-game-api .
```

## トラブルシューティング

### データベース接続エラー

```bash
# データベースのログを確認
docker compose logs db

# データベースに直接接続
docker compose exec db psql -U postgres -d guess_title_game
```

### マイグレーションエラー

```bash
# dirty database エラーの場合
./scripts/migration.sh force 20251223000001

# 最初からやり直す
docker compose down -v
docker compose up -d
```

## 環境変数

| 変数名 | 説明 | デフォルト値 |
|--------|------|--------------|
| SERVER_PORT | サーバーポート | 8080 |
| DB_HOST | データベースホスト | localhost |
| DB_PORT | データベースポート | 5432 |
| DB_USER | データベースユーザー | postgres |
| DB_PASSWORD | データベースパスワード | postgres |
| DB_NAME | データベース名 | guess_title_game |
| DB_SSL_MODE | SSL モード | disable |

## ライセンス

MIT
