#!/bin/bash
# Deploy frontend to GCP Cloud Run

set -e  # Exit on error

# Configuration
export GCP_PROJECT_ID=steenfuentes-dev
export GCP_REGION=us-west1
export SERVICE_NAME=steenfuentes-frontend
export IMAGE_TAG=$(date +%Y%m%d-%H%M%S)
export IMAGE_URL=$GCP_REGION-docker.pkg.dev/$GCP_PROJECT_ID/steenfuentes/$SERVICE_NAME:$IMAGE_TAG

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== Frontend Deployment to Cloud Run ===${NC}"
echo "Project: $GCP_PROJECT_ID"
echo "Region: $GCP_REGION"
echo "Service: $SERVICE_NAME"
echo "Image: $IMAGE_URL"
echo ""

# Prompt to ensure user has ran `gcloud auth login`
read -p "Have you run 'gcloud auth login'? (y/n): " gcloud_login
if [ "$gcloud_login" != "y" ]; then
    echo -e "${RED}Please run 'gcloud auth login' and try again.${NC}"
    exit 1
fi

# Set the project
echo -e "${GREEN}Setting GCP project...${NC}"
gcloud config set project $GCP_PROJECT_ID

# Configure docker authentication for Artifact Registry
echo -e "${GREEN}Configuring Docker authentication...${NC}"
gcloud auth configure-docker $GCP_REGION-docker.pkg.dev --quiet

# Build the frontend image
echo -e "${GREEN}Building Docker image...${NC}"
docker buildx build --platform=linux/amd64 -t $IMAGE_URL ./frontend

# Push the image to Artifact Registry
echo -e "${GREEN}Pushing image to Artifact Registry...${NC}"
docker push $IMAGE_URL

# Deploy to Cloud Run
echo -e "${GREEN}Deploying to Cloud Run...${NC}"
gcloud run deploy $SERVICE_NAME \
    --image=$IMAGE_URL \
    --region=$GCP_REGION \
    --platform=managed \
    --allow-unauthenticated \
    --port=3000 \
    --memory=512Mi \
    --cpu=1 \
    --min-instances=0 \
    --max-instances=100 \
    --set-env-vars="HOST=0.0.0.0,NODE_ENV=production"

# Get the service URL
SERVICE_URL=$(gcloud run services describe $SERVICE_NAME --region=$GCP_REGION --format='value(status.url)')

echo ""
echo -e "${GREEN}=== Deployment Complete! ===${NC}"
echo -e "Service URL: ${YELLOW}$SERVICE_URL${NC}"
echo ""

# Optional: Tag the image as latest
echo -e "${GREEN}Tagging image as latest...${NC}"
docker tag $IMAGE_URL $GCP_REGION-docker.pkg.dev/$GCP_PROJECT_ID/steenfuentes/$SERVICE_NAME:latest
docker push $GCP_REGION-docker.pkg.dev/$GCP_PROJECT_ID/steenfuentes/$SERVICE_NAME:latest
