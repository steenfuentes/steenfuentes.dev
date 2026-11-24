# SQLite Backup Strategy

This document explains the backup strategy for the SQLite database on Cloud Run.

## Overview

Since Cloud Run containers are ephemeral, we use Google Cloud Storage (GCS) to persist SQLite backups. This approach is:
- ✅ Cost-effective (~$0.50/month)
- ✅ Simple to implement
- ✅ Good enough for a personal blog
- ⚠️ Risk of data loss between backups

## How It Works

1. **On Startup**: 
   - Check GCS for latest backup
   - Restore if found
   - Start with empty DB if no backup exists

2. **During Runtime**:
   - Create backup every 6 hours (configurable)
   - Upload to GCS bucket
   - Keep last 30 days of backups

3. **On Shutdown/Crash**:
   - Data since last backup may be lost
   - Acceptable for low-traffic personal blog

## Setup Instructions

### 1. Create Backup Bucket

```bash
# Set your project ID
export GCP_PROJECT_ID=your-project-id

# Run setup script
make setup-backup GCP_PROJECT_ID=$GCP_PROJECT_ID
```

### 2. Add GitHub Secret

Add `BACKUP_BUCKET` secret with value `steenfuentes-db-backups` (or your custom bucket name).

### 3. Deploy

The application will automatically start using backups on next deployment.

## Configuration

Environment variables:
- `BACKUP_BUCKET`: GCS bucket name (required to enable backups)
- `BACKUP_INTERVAL`: Backup frequency (default: "6h")
- `DATA_DIR`: Local SQLite directory (default: "/data")

## Manual Backup/Restore

### Create Manual Backup

```bash
# SSH into Cloud Run instance
gcloud run services exec steenfuentes-backend --region=us-central1

# Inside container
sqlite3 /data/portfolio.db ".backup /tmp/manual-backup.db"
gsutil cp /tmp/manual-backup.db gs://steenfuentes-db-backups/manual/$(date +%Y%m%d-%H%M%S).db
```

### Restore Specific Backup

```bash
# List backups
gsutil ls gs://steenfuentes-db-backups/backups/

# Download and restore
gsutil cp gs://steenfuentes-db-backups/backups/portfolio-20240115-120000.db /tmp/
# Then redeploy service
```

## Monitoring

### Check Backup Status

```bash
# List recent backups
gsutil ls -l gs://steenfuentes-db-backups/backups/ | tail -10

# Check backup size
gsutil du -s gs://steenfuentes-db-backups/
```

### View Logs

```bash
# Check for backup-related logs
gcloud run logs read steenfuentes-backend \
  --region=us-central1 \
  --filter="backup" \
  --limit=50
```

## Limitations

1. **Data Loss Window**: Up to 6 hours of data could be lost
2. **Single Instance**: Only works with min-instances=1
3. **No Point-in-Time Recovery**: Only discrete backup points
4. **Manual Intervention**: Need to manually restore if corruption occurs

## Cost Breakdown

- Storage: ~$0.02/GB/month (SQLite is tiny)
- Operations: ~$0.005 per 10,000 operations
- Network: Free (same region)
- **Total**: Usually under $0.50/month

## Migration Path

If you outgrow this solution:

1. **Cloud SQL**: Full managed database (~$15/month)
2. **Firestore**: NoSQL option with better scaling
3. **Persistent Disk**: Keep SQLite with true persistence (~$5/month)

## Troubleshooting

### No Backups Created
- Check `BACKUP_BUCKET` env var is set
- Verify service account has Storage Admin role
- Check logs for permission errors

### Restore Failed
- Ensure bucket exists and has backups
- Check service account permissions
- Verify backup file isn't corrupted

### Database Corrupted
- Download latest backup locally
- Verify with `sqlite3 backup.db "PRAGMA integrity_check;"`
- Restore from older backup if needed