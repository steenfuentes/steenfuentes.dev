# Portfolio Website

A modern full-stack portfolio website built with SvelteKit and Golang, featuring a blog/journal system with GitHub OAuth authentication.

## Tech Stack

- **Frontend**: SvelteKit 5, Tailwind CSS
- **Backend**: Golang, SQLite
- **Authentication**: GitHub OAuth
- **Deployment**: Google Cloud Run

## Features

- 📝 Blog/Journal management system
- 🔐 GitHub OAuth for admin access
- 🎨 Elegant, minimalist design
- 📱 Fully responsive
- ⚡ Fast and lightweight

## Getting Started

### Prerequisites

- Docker and Docker Compose
- GitHub OAuth App (for authentication)
- Go 1.21+ (for local development)
- Node.js 20+ (for local development)

### Setup

1. Clone the repository:
```bash
git clone https://github.com/steenfuentes/portfolio.git
cd portfolio
```

2. Create environment file:
```bash
cp backend/.env.example backend/.env
```

#### Option A: Development Mode (No Authentication)

For local development, you can disable authentication:

```bash
# Edit backend/.env and set:
AUTH_ENABLED=false
```

Then run:
```bash
docker-compose up --build
```

You can now access the admin panel directly at http://localhost:3000/admin without logging in.

#### Option B: Production Mode (With GitHub OAuth)

1. Create a GitHub OAuth App:
   - Go to GitHub Settings > Developer settings > OAuth Apps
   - Create a new OAuth App with:
     - Homepage URL: `http://localhost:3000`
     - Authorization callback URL: `http://localhost:8080/auth/callback`

2. Update the `.env` file with your GitHub OAuth credentials:
```env
AUTH_ENABLED=true
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret
GITHUB_ALLOWED_USER=your_github_username
```

### Development

Using Docker Compose (recommended):
```bash
docker-compose up
```

Or run services separately:

Backend:
```bash
cd backend
go run cmd/server/main.go
```

Frontend:
```bash
cd frontend
npm install
npm run dev
```

The application will be available at:
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080

## API Endpoints

### Public Endpoints
- `GET /api/posts` - List published posts
- `GET /api/posts/:slug` - Get post by slug

### Protected Endpoints (requires authentication)
- `POST /api/posts` - Create new post
- `PUT /api/posts/:id` - Update post
- `DELETE /api/posts/:id` - Delete post

### Authentication
- `GET /api/auth/github` - Initiate GitHub OAuth
- `GET /api/auth/github/callback` - OAuth callback
- `GET /api/auth/me` - Get current user
- `POST /api/auth/logout` - Logout

## Deployment to Google Cloud

1. Build and push Docker images to Google Container Registry
2. Deploy to Cloud Run with persistent disk for SQLite
3. Configure custom domain (steenfuentes.dev)
4. Set environment variables in Cloud Run

See deployment guide in `docs/deployment.md` for detailed instructions.

## Project Structure

```
portfolio/
├── frontend/          # SvelteKit application
├── backend/           # Golang API server
├── docker-compose.yml # Development environment
└── README.md
```

## License

MIT