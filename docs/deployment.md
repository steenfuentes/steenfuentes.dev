# Deployment Guide

This guide covers deploying steenfuentes.dev to Google Cloud Run using GitHub Actions.

## Architecture Overview

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   GitHub    │────▶│ GitHub       │────▶│ Google      │
│   Push      │     │ Actions      │     │ Cloud Run   │
└─────────────┘     └──────────────┘     └─────────────┘
                            │
                            ▼
                    ┌──────────────┐
                    │  Artifact    │
                    │  Registry    │
                    └──────────────┘
```

## Prerequisites

1. Google Cloud Project with billing enabled
2. GitHub repository
3. Domain name (optional)
4. GitHub OAuth App

## Initial Setup

### 1. Fork/Clone Repository

```bash
git clone https://github.com/steenfuentes/steenfuentes.git
cd steenfuentes
```

### 2. Set Up Google Cloud

```bash
# Set your project ID
export GCP_PROJECT_ID=your-project-id

# Configure gcloud
gcloud config set project $GCP_PROJECT_ID

# Run setup
make setup-gcp GCP_PROJECT_ID=$GCP_PROJECT_ID
```

### 3. Create Service Account Key

```bash
# Create key
gcloud iam service-accounts keys create github-actions-key.json \
  --iam-account=github-actions@$GCP_PROJECT_ID.iam.gserviceaccount.com

# Display key (to copy to GitHub)
cat github-actions-key.json
```

### 4. Configure GitHub Secrets

Go to your repository Settings → Secrets and variables → Actions, and add:

| Secret Name | Description | Example |
|-------------|-------------|---------|
| `GCP_PROJECT_ID` | Your GCP project ID | `my-project-123` |
| `GCP_SA_KEY` | Service account key JSON | `{...}` (entire JSON) |
| `GITHUB_CLIENT_ID` | OAuth app client ID | `abc123` |
| `GITHUB_CLIENT_SECRET` | OAuth app secret | `xyz789` |
| `GITHUB_ALLOWED_USER` | Your GitHub username | `steenfuentes` |
| `FRONTEND_URL` | Production frontend URL | `https://steenfuentes.dev` |

### 5. Create GitHub OAuth App

1. Go to GitHub Settings → Developer settings → OAuth Apps
2. Click "New OAuth App"
3. Fill in:
   - **Application name**: steenfuentes.dev Admin
   - **Homepage URL**: https://steenfuentes.dev
   - **Authorization callback URL**: https://steenfuentes-backend-PROJECT_ID.a.run.app/api/auth/github/callback
4. Save Client ID and Client Secret

### 6. First Deployment

```bash
# Push to main branch
git add .
git commit -m "Initial deployment"
git push origin main
```

GitHub Actions will automatically:
1. Build Docker images
2. Push to Artifact Registry
3. Deploy to Cloud Run
4. Output service URLs

## Post-Deployment Steps

### 1. Update OAuth Redirect URI

After first deployment, update your GitHub OAuth App callback URL with the actual backend URL from Cloud Run.

### 2. Set Up Custom Domain (Optional)

```bash
# For frontend
gcloud run domain-mappings create \
  --service=steenfuentes-frontend \
  --domain=steenfuentes.dev \
  --region=us-central1

# Add DNS records as instructed by Cloud Run
```

### 3. Configure SQLite Persistence

For production, consider:
- Using Cloud SQL instead of SQLite
- Or mounting a persistent disk to Cloud Run
- Or using Google Cloud Storage with SQLite backups

## Monitoring

### View Logs

```bash
# Backend logs
make logs-backend

# Frontend logs  
make logs-frontend
```

### Check Status

```bash
make status
```

## Troubleshooting

### Build Failures

1. Check GitHub Actions logs
2. Ensure all secrets are set correctly
3. Verify service account permissions

### Runtime Errors

1. Check Cloud Run logs
2. Verify environment variables
3. Check CORS settings if frontend can't reach backend

### Database Issues

- SQLite data is ephemeral by default on Cloud Run
- Consider implementing regular backups
- For production, migrate to Cloud SQL

## Cost Management

- Set maximum instance limits
- Use minimum instances = 0 for dev/staging
- Monitor usage in GCP Console
- Preview environments auto-delete on PR close

## Security Best Practices

1. Rotate service account keys regularly
2. Use Secret Manager for sensitive values
3. Enable Cloud Run authentication for admin endpoints
4. Regular security updates for dependencies
5. Consider Workload Identity Federation instead of keys

## Next Steps

- [ ] Set up monitoring and alerts
- [ ] Configure backup strategy
- [ ] Implement health checks
- [ ] Set up staging environment
- [ ] Configure CI testing