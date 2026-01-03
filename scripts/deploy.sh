#!/bin/bash

# Cloud Runへのデプロイスクリプト

set -e

# カラー設定
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 設定
PROJECT_ID="${GCP_PROJECT_ID:-your-project-id}"
REGION="${GCP_REGION:-asia-northeast1}"
SERVICE_NAME="guess-title-game-api"
DB_INSTANCE_NAME="guess-title-game-db"

echo -e "${GREEN}=== Cloud Run デプロイスクリプト ===${NC}"

# プロジェクトIDの確認
if [ "$PROJECT_ID" = "your-project-id" ]; then
    echo -e "${RED}エラー: GCP_PROJECT_ID環境変数を設定してください${NC}"
    echo "例: export GCP_PROJECT_ID=your-project-id"
    exit 1
fi

echo -e "${YELLOW}プロジェクトID: $PROJECT_ID${NC}"
echo -e "${YELLOW}リージョン: $REGION${NC}"

# プロジェクトの設定
echo -e "${GREEN}[1/4] GCPプロジェクトを設定中...${NC}"
gcloud config set project $PROJECT_ID

# Cloud SQLインスタンスの接続名を取得
echo -e "${GREEN}[2/4] Cloud SQL接続情報を取得中...${NC}"
DB_CONNECTION_NAME=$(gcloud sql instances describe $DB_INSTANCE_NAME --format="value(connectionName)" 2>/dev/null || echo "")

if [ -z "$DB_CONNECTION_NAME" ]; then
    echo -e "${YELLOW}警告: Cloud SQLインスタンス '$DB_INSTANCE_NAME' が見つかりません${NC}"
    echo "Cloud SQLインスタンスを先に作成してください"
    echo ""
    echo "作成コマンド例:"
    echo "  gcloud sql instances create $DB_INSTANCE_NAME \\"
    echo "    --database-version=POSTGRES_16 \\"
    echo "    --tier=db-f1-micro \\"
    echo "    --region=$REGION"
    exit 1
fi

echo -e "${GREEN}Cloud SQL接続名: $DB_CONNECTION_NAME${NC}"

# イメージのビルドとプッシュ
echo -e "${GREEN}[3/4] Dockerイメージをビルド中...${NC}"
gcloud builds submit --tag gcr.io/$PROJECT_ID/$SERVICE_NAME

# Cloud Runにデプロイ
echo -e "${GREEN}[4/4] Cloud Runにデプロイ中...${NC}"
gcloud run deploy $SERVICE_NAME \
  --image gcr.io/$PROJECT_ID/$SERVICE_NAME \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --add-cloudsql-instances $DB_CONNECTION_NAME \
  --set-env-vars "ENV=production" \
  --set-env-vars "DB_HOST=/cloudsql/$DB_CONNECTION_NAME" \
  --set-env-vars "DB_USER=postgres" \
  --set-env-vars "DB_NAME=guess_title_game" \
  --set-env-vars "DB_SSL_MODE=disable" \
  --set-secrets "DB_PASSWORD=DB_PASSWORD:latest" \
  --min-instances 0 \
  --max-instances 10 \
  --memory 512Mi \
  --cpu 1 \
  --timeout 300

echo -e "${GREEN}デプロイ完了！${NC}"

# サービスURLを取得
SERVICE_URL=$(gcloud run services describe $SERVICE_NAME --region $REGION --format "value(status.url)")
echo -e "${GREEN}サービスURL: $SERVICE_URL${NC}"

# ヘルスチェック
echo -e "${YELLOW}ヘルスチェック実行中...${NC}"
sleep 5
HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" $SERVICE_URL/health || echo "000")

if [ "$HTTP_STATUS" = "200" ]; then
    echo -e "${GREEN}✓ ヘルスチェック成功 (HTTP $HTTP_STATUS)${NC}"
else
    echo -e "${RED}✗ ヘルスチェック失敗 (HTTP $HTTP_STATUS)${NC}"
    echo "ログを確認してください: gcloud run logs read $SERVICE_NAME --region $REGION"
fi
