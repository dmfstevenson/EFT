# Philippine ETF Recommendations

A real-time ETF recommendation website for the Philippines, built with Go and React.

## Features

- **ETF Analysis**: View Philippine ETFs with performance metrics and historical data
- **Trading Platform Reviews**: Compare popular trading platforms in the Philippines
- **News Feed**: Latest ETF and financial news
- **Price Alerts**: Set email alerts when ETF prices reach your target thresholds
- **Interactive Charts**: Visualize ETF price history with Recharts

## Tech Stack

### Backend
- **Go 1.21+**: Backend API
- **Gin**: Web framework
- **PostgreSQL**: Database for historical data
- **Yahoo Finance API**: ETF price data source
- **Cloudflare Workers**: Deployment platform

### Frontend
- **React 18**: UI framework
- **TypeScript**: Type safety
- **TailwindCSS**: Styling
- **Recharts**: Data visualization
- **Lucide React**: Icons
- **Cloudflare Pages**: Deployment platform

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go          # Backend entry point
├── internal/
│   ├── api/                 # API handlers
│   ├── data/                # Data fetching and repository
│   ├── models/              # Database models
│   └── scheduler/           # Daily data update scheduler
├── frontend/
│   ├── src/
│   │   ├── components/      # React components
│   │   ├── services/        # API client
│   │   └── App.tsx          # Main application
│   └── package.json
├── wrangler.toml            # Cloudflare Workers config
└── README.md
```

## Local Development

### Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL 14+

### Backend Setup

1. Set up PostgreSQL database:
```bash
createdb etf_recommendations
```

2. Set environment variables:
```bash
export DATABASE_URL="postgres://localhost/etf_recommendations?sslmode=disable"
export PORT=8080
```

3. Run the backend:
```bash
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`

### Frontend Setup

1. Navigate to frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Set environment variables:
```bash
cp .env.example .env
# Edit .env to set VITE_API_URL=http://localhost:8080/api
```

4. Run the development server:
```bash
npm run dev
```

The frontend will be available at `http://localhost:5173`

## Deployment

### Backend (Cloudflare Workers)

1. Install Wrangler CLI:
```bash
npm install -g wrangler
```

2. Login to Cloudflare:
```bash
wrangler login
```

3. Create D1 database:
```bash
wrangler d1 create etf_recommendations
```

4. Update `wrangler.toml` with your database ID

5. Deploy:
```bash
wrangler deploy
```

### Frontend (Cloudflare Pages)

**Option 1: Manual Deployment**
1. Build the frontend:
```bash
cd frontend
npm run build
```

2. Deploy to Cloudflare Pages:
```bash
npx wrangler pages deploy dist
```

**Option 2: GitHub Integration (Recommended)**
1. Connect your GitHub repository to Cloudflare Pages
2. In Cloudflare Pages settings, configure:
   - **Build command**: `cd frontend && npm install && npm run build`
   - **Build output directory**: `frontend/dist`
   - **Root directory**: `/` (repository root)

**Important**: Cloudflare Pages should only deploy the React frontend. The `.cloudflare/pages.json` file in the repository root configures this. Do not configure Cloudflare Pages to deploy the Go backend - that should be deployed separately to Cloudflare Workers.

## API Endpoints

- `GET /api/health` - Health check
- `GET /api/etfs` - List all ETFs
- `GET /api/etfs/:symbol` - Get ETF by symbol
- `GET /api/etfs/:symbol/prices` - Get ETF price history
- `GET /api/etfs/top?period=1y&limit=10` - Get top performers
- `GET /api/platforms` - Get trading platform reviews
- `GET /api/news` - Get latest news
- `POST /api/alerts` - Create price alert

## Data Sources

- **ETF Prices**: Yahoo Finance API (free tier)
- **News**: Public RSS feeds from financial news sites
- **Platform Reviews**: Curated content based on research

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details
