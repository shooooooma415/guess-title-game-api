# Cloud Run クイックスタート

最小限の手順でCloud Runにデプロイする方法です。

## 必要な環境変数の設定

```bash
export GCP_PROJECT_ID="your-project-id"
export GCP_REGION="asia-northeast1"
export DB_PASSWORD="your-secure-password"
```

## ワンライナーデプロイ

```bash
# 1. APIの有効化
gcloud services enable cloudbuild.googleapis.com run.googleapis.com sqladmin.googleapis.com secretmanager.googleapis.com --project=$GCP_PROJECT_ID

# 2. Cloud SQLインスタンスの作成
gcloud sql instances create guess-title-game-db \
  --database-version=POSTGRES_16 \
  --tier=db-f1-micro \
  --region=$GCP_REGION \
  --project=$GCP_PROJECT_ID

# 3. データベースとユーザーの設定
gcloud sql databases create guess_title_game --instance=guess-title-game-db --project=$GCP_PROJECT_ID && \
gcloud sql users set-password postgres --instance=guess-title-game-db --password="$DB_PASSWORD" --project=$GCP_PROJECT_ID

# 4. Secretの作成
echo -n "$DB_PASSWORD" | gcloud secrets create DB_PASSWORD --data-file=- --project=$GCP_PROJECT_ID

# 5. Secretへのアクセス権限付与
PROJECT_NUMBER=$(gcloud projects describe $GCP_PROJECT_ID --format="value(projectNumber)") && \
gcloud secrets add-iam-policy-binding DB_PASSWORD \
  --member="serviceAccount:${PROJECT_NUMBER}-compute@developer.gserviceaccount.com" \
  --role="roles/secretmanager.secretAccessor" \
  --project=$GCP_PROJECT_ID

# 6. マイグレーション実行（Cloud SQL Proxy経由）
# 別ターミナルでProxyを起動してから実行
# ./cloud-sql-proxy $GCP_PROJECT_ID:$GCP_REGION:guess-title-game-db --port 5432
# migrate -path ./db/migrations -database "postgresql://postgres:$DB_PASSWORD@localhost:5432/guess_title_game?sslmode=disable" up

# 7. デプロイ
./deploy.sh
```

## または自動スクリプトで一括実行

```bash
export GCP_PROJECT_ID="your-project-id"
./deploy.sh
```

## 確認

```bash
# サービスURLを取得
gcloud run services describe guess-title-game-api --region $GCP_REGION --format "value(status.url)"

# ヘルスチェック
curl $(gcloud run services describe guess-title-game-api --region $GCP_REGION --format "value(status.url)")/health
```

## よくあるコマンド

```bash
# ログの確認
gcloud run logs tail guess-title-game-api --region $GCP_REGION

# 再デプロイ
gcloud builds submit --tag gcr.io/$GCP_PROJECT_ID/guess-title-game-api && \
gcloud run deploy guess-title-game-api --image gcr.io/$GCP_PROJECT_ID/guess-title-game-api --region $GCP_REGION

# サービスの削除
gcloud run services delete guess-title-game-api --region $GCP_REGION

# Cloud SQLインスタンスの削除
gcloud sql instances delete guess-title-game-db
```
