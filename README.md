# Personal Salary Dashboard

A web application to track and compare regional salary data using annual OES (Occupational Employment and Wage Statistics) data from the U.S. Bureau of Labor Statistics (BLS).

## Features

- **Dynamic Data Scraping**: Downloads and parses BLS annual Excel spreadsheets (2018-present).
- **Interactive Dashboard**:
  - Regional Map (USA states/metros).
  - State Comparison horizontal bar charts (25th/Median/75th percentiles).
  - Multi-year Salary Trends line charts.
- **Granular Filtering**:
  - Filter by occupation (SOC codes).
  - Filter by state.
  - Toggle between State and Metro area aggregations.
- **Personalized Insights**: Compare your own salary against regional benchmarks.
- **Secure Authentication**: JWT-based auth with multi-factor authentication (TOTP) support.

## Tech Stack

### Backend
- **Language**: Go 1.25+
- **Framework**: Gin Gonic
- **ORM**: GORM (Postgres)
- **Data Parsing**: `excelize` for high-performance Excel streaming.
- **Security**: Argon2 password hashing, JWT, MFA via `pquerna/otp`.

### Frontend
- **Framework**: Vue 3 (Composition API)
- **State Management**: Pinia
- **Styling**: Tailwind CSS 4
- **Charts**: Apache ECharts (via `vue-echarts`)
- **Build Tool**: Vite

## Project Structure

```
├── backend/             # Go source code
│   ├── cmd/server/      # Application entry point
│   ├── internal/        # Private library code
│   │   ├── handlers/    # API endpoints
│   │   ├── models/      # Database schemas
│   │   ├── scraper/     # BLS Excel parsing logic
│   │   └── middleware/  # Auth & security
│   └── testdata/        # Test fixtures (Excel)
├── frontend/            # Vue.js source code
│   ├── src/
│   │   ├── api/         # Axios client
│   │   ├── components/  # Reusable UI & charts
│   │   ├── stores/      # Pinia state
│   │   └── views/       # Page components
│   └── public/          # Static assets
└── docker-compose.yml   # Infrastructure (PostgreSQL)
```

## Getting Started

### Prerequisites
- Go 1.25+
- Node.js 20.10+
- Docker & Docker Compose (for database)

### Installation

1. **Clone the repository**:
   ```bash
   git clone <repo-url>
   cd project_salary
   ```

2. **Start the database**:
   ```bash
   docker-compose up -d
   ```

3. **Backend Setup**:
   ```bash
   cd backend
   cp .env.example .env  # Update with your DB credentials
   go mod download
   go run cmd/server/main.go
   ```

4. **Frontend Setup**:
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

5. **Initialize Data**:
   - Register an account at `http://localhost:5173/`.
   - Go to "Data Management" and click "Pull Latest Data" (or use the backfill option for historical data).

## License
MIT
