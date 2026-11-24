# Root Makefile for steenfuentes.dev

# Variables
GCP_PROJECT_ID ?= steenfuentes-dev
GCP_REGION ?= us-west1
DOMAIN ?= steenfuentes.dev

# Default target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make setup-gcp      - Set up GCP resources"
	@echo "  make setup-backup   - Set up database backup bucket"
	@echo "  make deploy         - Deploy to production (via GitHub Actions)"
	@echo "  make deploy-frontend - Manually deploy frontend to Cloud Run"
	@echo "  make dev            - Run development environment"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build artifacts"

# Set up GCP resources
.PHONY: setup-gcp
setup-gcp:
	@echo "Setting up GCP resources..."
	# Enable APIs
	gcloud services enable run.googleapis.com
	gcloud services enable artifactregistry.googleapis.com
	gcloud services enable secretmanager.googleapis.com

	# Create Artifact Registry
	gcloud artifacts repositories create steenfuentes \
		--repository-format=docker \
		--location=$(GCP_REGION) || echo "Repository already exists"

	# Create service account
	gcloud iam service-accounts create github-actions \
		--display-name="GitHub Actions Service Account" || echo "Service account already exists"

	# Grant permissions
	gcloud projects add-iam-policy-binding $(GCP_PROJECT_ID) \
		--member="serviceAccount:github-actions@$(GCP_PROJECT_ID).iam.gserviceaccount.com" \
		--role="roles/run.admin"
	gcloud projects add-iam-policy-binding $(GCP_PROJECT_ID) \
		--member="serviceAccount:github-actions@$(GCP_PROJECT_ID).iam.gserviceaccount.com" \
		--role="roles/artifactregistry.writer"
	gcloud projects add-iam-policy-binding $(GCP_PROJECT_ID) \
		--member="serviceAccount:github-actions@$(GCP_PROJECT_ID).iam.gserviceaccount.com" \
		--role="roles/secretmanager.secretAccessor"

	@echo "✅ GCP setup complete!"
	@echo "Next: Create service account key and add to GitHub secrets"

# Set up backup bucket
.PHONY: setup-backup
setup-backup:
	GCP_PROJECT_ID=$(GCP_PROJECT_ID) ./scripts/setup-backup.sh

# Development environment
.PHONY: dev
dev:
	docker-compose up

# Deploy to production (manual deployment)
.PHONY: deploy
deploy:
	@echo "Deployment should be done via GitHub Actions"
	@echo "Push to main branch to trigger deployment"

# Deploy frontend manually to Cloud Run
.PHONY: deploy-frontend
deploy-frontend:
	@echo "Deploying frontend to Cloud Run..."
	./scripts/deploy-frontend.sh

# Run frontend locally
.PHONY: frontend-dev
frontend-dev:
	cd frontend && npm run dev

# Run backend locally
.PHONY: backend-dev
backend-dev:
	cd backend && go run cmd/server/main.go

# Build all Docker images
.PHONY: docker-build
docker-build:
	cd backend && docker build -t steenfuentes-backend .
	cd frontend && docker build -t steenfuentes-frontend .

# Clean build artifacts
.PHONY: clean
clean:
	cd frontend && make clean
	cd backend && go clean
	docker system prune -f

# Check deployment status
.PHONY: status
status:
	@echo "Backend service:"
	gcloud run services describe steenfuentes-backend --region=$(GCP_REGION) --format="value(status.url)" || echo "Not deployed"
	@echo "\nFrontend service:"
	gcloud run services describe steenfuentes-frontend --region=$(GCP_REGION) --format="value(status.url)" || echo "Not deployed"

# View logs
.PHONY: logs-backend
logs-backend:
	gcloud run logs read steenfuentes-backend --region=$(GCP_REGION) --limit=50

.PHONY: logs-frontend
logs-frontend:
	gcloud run logs read steenfuentes-frontend --region=$(GCP_REGION) --limit=50
