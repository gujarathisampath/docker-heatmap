# 🐳 Docker Heatmap

Generate beautiful GitHub-style contribution heatmaps for your Docker Hub activity. Embed in your README and showcase your container contributions.

![Docker Heatmap Preview](https://via.placeholder.com/800x200/161b22/39d353?text=Docker+Activity+Heatmap)

## ✨ Features

- **🔐 GitHub OAuth** - Secure authentication with your GitHub account
- **📊 Beautiful Heatmaps** - GitHub-style SVG contribution graphs
- **🔗 Easy Embedding** - Copy-paste URLs for README or any website
- **🔒 Secure Storage** - AES-256 encrypted token storage (zero plaintext)
- **⚡ Auto Refresh** - Background jobs keep data up-to-date
- **📡 Public API** - JSON endpoints for custom integrations

## 🛠 Tech Stack

| Layer              | Technology                                                           |
| ------------------ | -------------------------------------------------------------------- |
| **Frontend**       | Next.js 14, TypeScript, Tailwind CSS, shadcn/ui, TanStack Query, Zod |
| **Backend**        | Go, GoFiber, GORM                                                    |
| **Database**       | PostgreSQL                                                           |
| **Auth**           | GitHub OAuth                                                         |
| **Infrastructure** | Docker, Docker Compose                                               |

## 📁 Project Structure

```
docker-heatmap/
├── frontend/           # Next.js frontend
│   ├── app/           # App router pages
│   ├── components/    # UI components
│   ├── lib/           # Utilities & API client
│   ├── context/       # React contexts
│   └── hooks/         # Custom hooks
├── backend/           # Go backend
│   ├── cmd/           # Entry point
│   └── internal/
│       ├── config/    # Configuration
│       ├── database/  # Database connection
│       ├── handlers/  # HTTP handlers
│       ├── middleware/# Auth & rate limiting
│       ├── models/    # GORM models
│       ├── services/  # Business logic
│       ├── utils/     # Utilities
│       └── worker/    # Background jobs
├── infra/             # Infrastructure
│   ├── docker-compose.yml
│   └── nginx/
└── .env.example       # Environment template
```

## 🚀 Quick Start

### Prerequisites

- Node.js 20+ (or Bun)
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL (or use Docker)

### 1. Clone & Setup

```bash
git clone https://github.com/yourusername/docker-heatmap.git
cd docker-heatmap

# Copy environment template
cp .env.example .env
```

### 2. Configure GitHub OAuth

1. Go to [GitHub Developer Settings](https://github.com/settings/developers)
2. Create a new OAuth App:
   - **Homepage URL:** `http://localhost:3000`
   - **Callback URL:** `http://localhost:8080/api/auth/github/callback`
3. Copy Client ID and Secret to `.env`

### 3. Start Database & Backend (Docker)

```bash
cd infra
docker-compose up -d
```

### 4. Start Frontend

```bash
cd frontend
npm install  # or: bun install
npm run dev  # or: bun dev
```

### 5. Open App

Visit [http://localhost:3000](http://localhost:3000)

## 🔧 Configuration

### Environment Variables

| Variable               | Description                  | Required |
| ---------------------- | ---------------------------- | -------- |
| `GITHUB_CLIENT_ID`     | GitHub OAuth Client ID       | ✅       |
| `GITHUB_CLIENT_SECRET` | GitHub OAuth Secret          | ✅       |
| `JWT_SECRET`           | Secret for JWT signing       | ✅       |
| `ENCRYPTION_KEY`       | 32-char key for AES-256      | ✅       |
| `DATABASE_URL`         | PostgreSQL connection string | ✅       |
| `FRONTEND_URL`         | Frontend URL for CORS        | ✅       |
| `PORT`                 | Backend port (default: 8080) | ❌       |

### Generating Secrets

```bash
# JWT Secret
openssl rand -hex 32

# Encryption Key (exactly 32 characters)
openssl rand -base64 24 | head -c 32
```

## 📡 API Endpoints

### Authentication

| Method | Endpoint                    | Description        |
| ------ | --------------------------- | ------------------ |
| GET    | `/api/auth/github`          | Start GitHub OAuth |
| GET    | `/api/auth/github/callback` | OAuth callback     |
| POST   | `/api/auth/logout`          | Logout             |

### User

| Method | Endpoint          | Description      |
| ------ | ----------------- | ---------------- |
| GET    | `/api/user/me`    | Get current user |
| PUT    | `/api/user/me`    | Update profile   |
| GET    | `/api/user/embed` | Get embed codes  |

### Docker

| Method | Endpoint                 | Description           |
| ------ | ------------------------ | --------------------- |
| POST   | `/api/docker/connect`    | Connect Docker Hub    |
| GET    | `/api/docker/account`    | Get connected account |
| DELETE | `/api/docker/disconnect` | Disconnect account    |
| POST   | `/api/docker/sync`       | Trigger sync          |

### Public (Embeddable)

| Method | Endpoint                       | Description   |
| ------ | ------------------------------ | ------------- |
| GET    | `/api/heatmap/:username.svg`   | SVG heatmap   |
| GET    | `/api/activity/:username.json` | Activity JSON |
| GET    | `/api/profile/:username`       | Profile data  |

## 🎨 Embedding Your Heatmap

### Markdown (GitHub README)

```markdown
![Docker Activity](https://your-domain.com/api/heatmap/your-docker-username.svg)
```

### HTML

```html
<img
  src="https://your-domain.com/api/heatmap/your-docker-username.svg"
  alt="Docker Activity"
/>
```

### With Link

```html
<a href="https://your-domain.com/profile/your-docker-username">
  <img
    src="https://your-domain.com/api/heatmap/your-docker-username.svg"
    alt="Docker Activity"
  />
</a>
```

## 🏗 Development

### Backend Only

```bash
cd backend
go mod download
go run cmd/main.go
```

### Frontend Only

```bash
cd frontend
npm install
npm run dev
```

### Full Stack (Docker)

```bash
cd infra
docker-compose up --build
```

## 🚢 Production Deployment

### VPS with Docker

1. Set production environment variables
2. Configure Nginx with SSL
3. Run Docker Compose with production profile

```bash
docker-compose --profile production up -d
```

### Environment for Production

```env
ENVIRONMENT=production
FRONTEND_URL=https://your-domain.com
GITHUB_CALLBACK_URL=https://your-domain.com/api/auth/github/callback
```

## 🔐 Security

- **Token Encryption:** Docker Hub tokens are encrypted with AES-256-GCM
- **OAuth State:** CSRF protection with state tokens
- **Rate Limiting:** Different tiers for API, auth, and public endpoints
- **JWT Auth:** Stateless authentication with 7-day expiry

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🤝 Contributing

Contributions are welcome! Please open an issue or submit a PR.

---

<p align="center">
  Made with ❤️ for the Docker community
</p>
