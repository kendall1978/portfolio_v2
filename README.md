# Kendall Roberts — Portfolio & Salary Dashboard

A personal portfolio website with an integrated salary benchmarking dashboard. The public-facing side showcases a profile, blog, and social links. Behind authentication, a salary dashboard lets you compare your compensation against regional BLS data.

## Features

### Public (no login required)
- **Home Page** — Profile photo, bio sections, and social links. All content is editable from the admin panel.
- **Blog** — Blog posts with images, pulled from the database.
- **Dark / Light Mode** — Toggle in the nav bar. Defaults to dark, remembers your preference.

### Authenticated
- **Salary Dashboard** — Interactive charts comparing your salary against regional benchmarks using BLS OES data. Filter by occupation, state, and area type.
- **Admin Panel** — Manage blog posts, home page sections, social links, and site settings (profile image, title, subtitle).
- **Settings** — Update your profile and configure MFA (TOTP).
- **Data Management** — Trigger BLS data imports (current year or historical backfill).

### Security
- JWT authentication with optional TOTP multi-factor authentication
- Recovery codes for MFA backup
- Argon2 password hashing, AES-256-GCM encryption for TOTP secrets

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.25+, Gin, GORM |
| Frontend | Vue 3 (Composition API), TypeScript, Vite |
| Styling | Tailwind CSS 4 |
| State | Pinia |
| Charts | Apache ECharts |
| Database | PostgreSQL 16 |
| Infra | Docker Compose |

## Project Structure

```
├── backend/
│   ├── cmd/
│   │   ├── server/             # Main application entry point
│   │   └── migrate-firebase/   # One-time Firebase migration script
│   └── internal/
│       ├── config/             # Environment configuration
│       ├── crypto/             # AES-256-GCM encryption
│       ├── database/           # GORM setup & auto-migration
│       ├── handlers/           # API endpoint handlers
│       │   ├── auth.go         # Register, login, MFA
│       │   ├── blog.go         # Blog post CRUD
│       │   ├── compensation.go # Salary tracking CRUD
│       │   ├── home.go         # Home section CRUD
│       │   ├── social.go       # Social link CRUD
│       │   ├── sitesettings.go # Key-value site settings
│       │   ├── upload.go       # Image upload (10MB, JPEG/PNG/WebP/GIF)
│       │   └── ...
│       ├── middleware/         # JWT auth middleware
│       ├── models/             # Database models
│       ├── routes/             # Route registration
│       └── scraper/            # BLS Excel data parser
├── frontend/
│   └── src/
│       ├── api/client.ts       # Axios client with JWT interceptors
│       ├── components/         # ThemeToggle, SocialIcon, charts
│       ├── stores/             # Pinia stores (auth, salary, portfolio, theme)
│       ├── views/              # Page components
│       │   ├── HomeView.vue    # Public portfolio home
│       │   ├── BlogView.vue    # Public blog listing
│       │   ├── AdminView.vue   # Content management (auth)
│       │   ├── DashboardView.vue # Salary dashboard (auth)
│       │   └── ...
│       └── router/index.ts     # Route definitions & auth guard
├── docs/plans/                 # Design & implementation plans
└── docker-compose.yml          # PostgreSQL service
```

## Getting Started

### Prerequisites
- Go 1.25+
- Node.js 20+
- Docker & Docker Compose

### Setup

1. **Start the database:**
   ```bash
   docker compose up -d
   ```

2. **Configure the backend:**
   ```bash
   cd backend
   cp .env.example .env
   ```
   Edit `.env` and set `ENCRYPTION_KEY` (generate with `openssl rand -hex 32`). The other defaults work for local development.

3. **Install frontend dependencies:**
   ```bash
   cd frontend
   npm install
   ```

4. **Run both servers:**
   ```bash
   # From the project root
   npm run dev
   ```
   This starts the Go backend on `http://localhost:8080` and the Vite dev server on `http://localhost:5173`.

5. **Create your account:**
   Visit `http://localhost:5173/login` and register. The first user registered gets access — registration closes after that.

### Seeding Content

The home page, blog, and social links are managed through the admin panel at `/admin` after logging in. You can:

- Add/edit/delete **blog posts** (with image upload)
- Add/edit/delete **home page sections** (bio cards)
- Add/edit/delete **social links** (GitHub, LinkedIn, etc.)
- Update **site settings** (title, subtitle, profile image URL)

### Firebase Migration (optional)

If migrating from the old Firebase-hosted portfolio:

1. Download your Firebase service account key from the Firebase Console
2. Save it as `backend/firebase-service-account.json`
3. Run:
   ```bash
   cd backend
   go run ./cmd/migrate-firebase/
   ```

### Loading Salary Data

1. Log in and go to **Data Management** (`/scrape`)
2. Click **Pull Latest Data** for the current year, or enable **Backfill** to import historical data (2018–present)
3. Optionally set a `BLS_API_KEY` in `.env` for higher rate limits (register at [data.bls.gov](https://data.bls.gov/registrationEngine/))

## API Overview

### Public Endpoints
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/blog` | List all blog posts |
| GET | `/api/blog/:id` | Get a single blog post |
| GET | `/api/home-sections` | List home page sections |
| GET | `/api/social-links` | List social links |
| GET | `/api/site-settings` | Get site settings |
| GET | `/api/health` | Health check |

### Authenticated Endpoints
All require `Authorization: Bearer <token>` header.

| Method | Path | Description |
|--------|------|-------------|
| POST/PUT/DELETE | `/api/blog[/:id]` | Blog post management |
| POST/PUT/DELETE | `/api/home-sections[/:id]` | Home section management |
| POST/PUT/DELETE | `/api/social-links[/:id]` | Social link management |
| PUT | `/api/site-settings` | Update site settings |
| POST | `/api/upload` | Upload an image |
| GET/POST/PUT/DELETE | `/api/compensation[/:id]` | Salary tracking |
| GET | `/api/salaries` | Query salary data |
| GET | `/api/salaries/trends` | Multi-year trends |
| GET | `/api/salaries/compare` | Compare against benchmarks |

## Theming

The site supports dark and light modes with a configurable accent color. The accent color is defined as CSS custom properties in `frontend/src/assets/main.css`:

```css
:root {
    --accent-color: #14b8a6;  /* teal — change this to swap accent */
    --accent-hover: #0d9488;
    --accent-light: #ccfbf1;
    --accent-muted: #5eead4;
}
```

Candidate accent colors:
- Teal: `#14b8a6`
- Violet: `#8b5cf6`
- Amber: `#f59e0b`
- Coral: `#E07A5F`

## License

MIT
