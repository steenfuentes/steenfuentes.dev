#!/bin/bash

# Script to set up Google Cloud Storage backup for SQLite

set -e

# Check if required variables are set
if [ -z "$GCP_PROJECT_ID" ]; then
    echo "Error: GCP_PROJECT_ID environment variable is not set"
    echo "Usage: GCP_PROJECT_ID=your-project-id ./setup-backup.sh"
    exit 1
fi

BUCKET_NAME="${BUCKET_NAME:-steenfuentes-db-backups}"
REGION="${GCP_REGION:-us-central1}"

echo "Setting up backup bucket: $BUCKET_NAME"
echo "Project: $GCP_PROJECT_ID"
echo "Region: $REGION"

# Create the bucket
echo "Creating GCS bucket..."
gsutil mb -p "$GCP_PROJECT_ID" -c "standard" -l "$REGION" "gs://$BUCKET_NAME" || echo "Bucket already exists"

# Set lifecycle rule to delete old backups after 30 days
echo "Setting lifecycle rules..."
cat > /tmp/lifecycle.json <<EOF
{
  "lifecycle": {
    "rule": [
      {
        "action": {"type": "Delete"},
        "condition": {
          "age": 30,
          "matchesPrefix": ["backups/"]
        }
      }
    ]
  }
}
EOF

gsutil lifecycle set /tmp/lifecycle.json "gs://$BUCKET_NAME"
rm /tmp/lifecycle.json

# Grant service account access
SERVICE_ACCOUNT="${SERVICE_ACCOUNT:-github-actions@$GCP_PROJECT_ID.iam.gserviceaccount.com}"
echo "Granting access to service account: $SERVICE_ACCOUNT"

gsutil iam ch "serviceAccount:$SERVICE_ACCOUNT:objectAdmin" "gs://$BUCKET_NAME"

# Create initial backup directory
echo "Creating backup prefix..."
echo "Backup bucket initialized" | gsutil cp - "gs://$BUCKET_NAME/backups/README.txt"

echo ""
echo "✅ Backup setup complete!"
echo ""
echo "Add these environment variables to your Cloud Run service:"
echo "  BACKUP_BUCKET=$BUCKET_NAME"
echo "  BACKUP_INTERVAL=6h  # Optional, defaults to 6 hours"
echo ""
echo "Lifecycle policy: Backups older than 30 days will be automatically deleted"