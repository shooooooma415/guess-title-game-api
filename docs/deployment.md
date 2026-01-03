# Cloud Run デプロイガイド

このガイドでは、guess-title-game-apiをGoogle Cloud Runにデプロイする手順を説明します。

## 前提条件

1. Google Cloud Platform (GCP) アカウント
2. gcloud CLI がインストール済み
3. プロジェクトが作成済み

## 1. 初期セットアップ

### 1.1 gcloud CLIのインストール（未インストールの場合）

```bash
# macOS
brew install --cask google-cloud-sdk

# 認証
gcloud auth login
```

### 1.2 必要なAPIの有効化

```bash
# プロジェクトIDを設定
export GCP_PROJECT_ID="your-project-id"
gcloud config set project $GCP_PROJECT_ID

# 必要なAPIを有効化
gcloud services enable \
  cloudbuild.googleapis.com \
  run.googleapis.com \
  sqladmin.googleapis.com \
  secretmanager.googleapis.com
```

## 2. Cloud SQL (PostgreSQL) のセットアップ

### 2.1 Cloud SQLインスタンスの作成

個人利用向けの最小構成（月額 $5-7程度）：

```bash
# 最小構成（個人利用・開発環境向け）
gcloud sql instances create guess-title-game-db \
  --database-version=POSTGRES_16 \
  --tier=db-f1-micro \
  --region=asia-northeast1 \
  --storage-type=HDD \
  --storage-size=10GB \
  --no-backup \
  --activation-policy=ALWAYS
```

**コスト削減のヒント:**
```bash
# 使わない時はインスタンスを停止（ストレージ料金のみ $1-2/月）
gcloud sql instances patch guess-title-game-db --activation-policy NEVER

# 再起動する時
gcloud sql instances patch guess-title-game-db --activation-policy ALWAYS
```

**本番環境用（推奨）:**
```bash
gcloud sql instances create guess-title-game-db \
  --database-version=POSTGRES_16 \
  --tier=db-g1-small \
  --region=asia-northeast1 \
  --storage-type=SSD \
  --storage-size=10GB \
  --backup \
  --backup-start-time=03:00
```

### 2.2 データベースの作成

```bash
gcloud sql databases create guess_title_game \
  --instance=guess-title-game-db
```

### 2.3 データベースパスワードの設定

```bash
# PostgreSQLユーザーのパスワードを設定
gcloud sql users set-password postgres \
  --instance=guess-title-game-db \
  --password="YOUR_SECURE_PASSWORD"
```

### 2.4 Secret Managerにパスワードを保存

```bash
# Secretを作成
echo -n "YOUR_SECURE_PASSWORD" | gcloud secrets create DB_PASSWORD \
  --data-file=-

# Cloud Runサービスアカウントに権限付与
PROJECT_NUMBER=$(gcloud projects describe $GCP_PROJECT_ID --format="value(projectNumber)")
gcloud secrets add-iam-policy-binding DB_PASSWORD \
  --member="serviceAccount:${PROJECT_NUMBER}-compute@developer.gserviceaccount.com" \
  --role="roles/secretmanager.secretAccessor"
```

## 3. マイグレーションの実行

Cloud SQLインスタンスに接続してマイグレーションを実行します。

### 3.1 Cloud SQL Proxyを使用する方法

```bash
# Cloud SQL Proxyのダウンロード
curl -o cloud-sql-proxy https://storage.googleapis.com/cloud-sql-connectors/cloud-sql-proxy/v2.14.2/cloud-sql-proxy.darwin.amd64
chmod +x cloud-sql-proxy

# Proxyの起動（別ターミナルで実行）
./cloud-sql-proxy $GCP_PROJECT_ID:asia-northeast1:guess-title-game-db --port 5432

# マイグレーションツールのインストール
brew install golang-migrate

# マイグレーション実行（元のターミナルで）
migrate -path ./db/migrations \
  -database "postgresql://postgres:YOUR_SECURE_PASSWORD@localhost:5432/guess_title_game?sslmode=disable" \
  up
```

### 3.2 Cloud Shellを使用する方法

```bash
# Cloud Shellで実行
gcloud cloud-shell ssh

# リポジトリをクローン
git clone YOUR_REPOSITORY_URL
cd guess-title-game-api

# マイグレーションツールのインストール
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# マイグレーション実行
~/go/bin/migrate -path ./db/migrations \
  -database "postgresql://postgres:YOUR_SECURE_PASSWORD@/guess_title_game?host=/cloudsql/$GCP_PROJECT_ID:asia-northeast1:guess-title-game-db" \
  up
```

## 4. Cloud Runへのデプロイ

### 4.1 自動デプロイスクリプトを使用

```bash
export GCP_PROJECT_ID="your-project-id"
export GCP_REGION="asia-northeast1"

./deploy.sh
```

### 4.2 手動デプロイ

```bash
# イメージのビルドとプッシュ
gcloud builds submit --tag gcr.io/$GCP_PROJECT_ID/guess-title-game-api

# Cloud Runにデプロイ
gcloud run deploy guess-title-game-api \
  --image gcr.io/$GCP_PROJECT_ID/guess-title-game-api \
  --platform managed \
  --region asia-northeast1 \
  --allow-unauthenticated \
  --add-cloudsql-instances $GCP_PROJECT_ID:asia-northeast1:guess-title-game-db \
  --set-env-vars "ENV=production" \
  --set-env-vars "DB_HOST=/cloudsql/$GCP_PROJECT_ID:asia-northeast1:guess-title-game-db" \
  --set-env-vars "DB_USER=postgres" \
  --set-env-vars "DB_NAME=guess_title_game" \
  --set-env-vars "DB_SSL_MODE=disable" \
  --set-secrets "DB_PASSWORD=DB_PASSWORD:latest" \
  --min-instances 0 \
  --max-instances 10 \
  --memory 512Mi \
  --cpu 1 \
  --timeout 300
```

## 5. 動作確認

```bash
# サービスURLを取得
SERVICE_URL=$(gcloud run services describe guess-title-game-api \
  --region asia-northeast1 \
  --format "value(status.url)")

# ヘルスチェック
curl $SERVICE_URL/health

# 期待されるレスポンス
# {"status":"ok"}
```

## 6. CI/CDセットアップ（オプション）

GitHub ActionsまたはCloud Buildトリガーを設定して、自動デプロイを有効にします。

### 6.1 Cloud Build トリガーの作成

```bash
gcloud builds triggers create github \
  --repo-name=guess-title-game-api \
  --repo-owner=YOUR_GITHUB_USERNAME \
  --branch-pattern="^main$" \
  --build-config=cloudbuild.yaml
```

## 7. モニタリングとログ

### ログの確認

```bash
# リアルタイムログ
gcloud run logs tail guess-title-game-api --region asia-northeast1

# 最新のログ
gcloud run logs read guess-title-game-api --region asia-northeast1 --limit 50
```

### メトリクスの確認

```bash
# Cloud Consoleでメトリクスを確認
echo "https://console.cloud.google.com/run/detail/asia-northeast1/guess-title-game-api/metrics?project=$GCP_PROJECT_ID"
```

## 8. コスト最適化

### 8.1 最小インスタンス数を0に設定

```bash
gcloud run services update guess-title-game-api \
  --region asia-northeast1 \
  --min-instances 0
```

### 8.2 Cloud SQLの自動シャットダウン設定（開発環境のみ）

```bash
# Cloud Consoleから手動で設定
# または、使用しない時間帯にインスタンスを停止
gcloud sql instances patch guess-title-game-db --activation-policy NEVER
```

## 9. トラブルシューティング

### デプロイが失敗する場合

```bash
# ビルドログを確認
gcloud builds list --limit 5

# 特定のビルドの詳細を確認
gcloud builds log BUILD_ID
```

### データベース接続エラー

```bash
# Cloud SQLインスタンスの状態を確認
gcloud sql instances describe guess-title-game-db

# Cloud Runサービスの環境変数を確認
gcloud run services describe guess-title-game-api \
  --region asia-northeast1 \
  --format yaml
```

### WebSocketの動作確認

Cloud Runは WebSocket をサポートしていますが、以下の点に注意してください：

- タイムアウト: 最大60分
- 接続数: サービスのスケーリング設定による

## 10. セキュリティベストプラクティス

1. **認証の追加**（本番環境推奨）
```bash
gcloud run services update guess-title-game-api \
  --region asia-northeast1 \
  --no-allow-unauthenticated
```

2. **Cloud Armorの設定**（DDoS対策）

3. **VPCコネクタの使用**（プライベートネットワーク）

## 参考リンク

- [Cloud Run ドキュメント](https://cloud.google.com/run/docs)
- [Cloud SQL for PostgreSQL](https://cloud.google.com/sql/docs/postgres)
- [Secret Manager](https://cloud.google.com/secret-manager/docs)
