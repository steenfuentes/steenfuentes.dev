# GitHub Actions CI/CD Setup

This directory contains GitHub Actions workflows for deploying the application to Google Cloud Run.

## Workflows

### 1. `deploy.yml` - Production Deployment
- Triggers on pushes to `main` branch
- Deploys both backend and frontend to Cloud Run
- Uses Artifact Registry for container storage
- Configures production environment variables

### 2. `preview.yml` - Preview Deployments
- Triggers on pull requests
- Creates temporary preview environments
- Posts preview URLs as PR comments
- Disables authentication for easy testing

### 3. `cleanup-preview.yml` - Preview Cleanup
- Triggers when PRs are closed
- Automatically deletes preview services
- Keeps Cloud Run costs under control

## Setup Instructions

### 1. Create GCP Resources

```bash
# Enable required APIs
gcloud services enable run.googleapis.com
gcloud services enable artifactregistry.googleapis.com
gcloud services enable secretmanager.googleapis.com

# Create Artifact Registry repository
gcloud artifacts repositories create steenfuentes \
  --repository-format=docker \
  --location=us-central1

# Create service account for GitHub Actions
gcloud iam service-accounts create github-actions \
  --display-name="GitHub Actions Service Account"

# Grant necessary permissions
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
  --member="serviceAccount:github-actions@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/run.admin"

gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
  --member="serviceAccount:github-actions@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/artifactregistry.writer"

gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
  --member="serviceAccount:github-actions@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/secretmanager.secretAccessor"

# Create and download service account key
gcloud iam service-accounts keys create github-actions-key.json \
  --iam-account=github-actions@YOUR_PROJECT_ID.iam.gserviceaccount.com
```

### 2. Configure GitHub Secrets

Add these secrets to your GitHub repository (Settings → Secrets and variables → Actions):

- `GCP_PROJECT_ID`: Your Google Cloud project ID
- `GCP_SA_KEY`: Contents of the service account key JSON file
- `GITHUB_CLIENT_ID`: GitHub OAuth app client ID
- `GITHUB_CLIENT_SECRET`: GitHub OAuth app client secret
- `GITHUB_ALLOWED_USER`: Your GitHub username (for admin access)
- `FRONTEND_URL`: Production frontend URL (e.g., https://steenfuentes.dev)

### 3. Create GitHub OAuth App

1. Go to GitHub Settings → Developer settings → OAuth Apps
2. Create a new OAuth App with:
   - Homepage URL: `https://steenfuentes.dev`
   - Authorization callback URL: `https://steenfuentes-backend-YOUR_PROJECT_ID.a.run.app/api/auth/github/callback`

### 4. (Optional) Set up custom domain

```bash
# Map custom domain to Cloud Run service
gcloud run domain-mappings create \
  --service=steenfuentes-frontend \
  --domain=steenfuentes.dev \
  --region=us-central1
```

## Environment Variables

### Backend
- `ENVIRONMENT`: Set to "production"
- `PORT`: Cloud Run sets this automatically
- `DATA_DIR`: Directory for SQLite database
- `MIGRATIONS_PATH`: Path to migration files
- `GITHUB_CLIENT_ID`: OAuth client ID
- `GITHUB_CLIENT_SECRET`: OAuth client secret (stored in Secret Manager)
- `GITHUB_ALLOWED_USER`: GitHub username allowed to admin
- `FRONTEND_URL`: Frontend URL for CORS

### Frontend
- `VITE_API_URL`: Backend API URL
- `PORT`: Cloud Run sets this automatically
- `NODE_ENV`: Set to "production"

## Cost Optimization

- Minimum instances set to 1 for production (prevents cold starts)
- Maximum instances set to 10 (adjust based on traffic)
- Preview environments have max 2 instances
- Preview environments are automatically deleted when PRs close

## Security Notes

- Consider using Workload Identity Federation instead of service account keys
- Store sensitive values in Google Secret Manager
- Review and rotate service account keys regularly
- Use least-privilege IAM roles