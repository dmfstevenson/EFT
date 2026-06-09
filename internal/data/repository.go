package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"etf-recommendation-api/internal/models"
)

type ETFRepository struct {
	db *sql.DB
}

func NewETFRepository(db *sql.DB) *ETFRepository {
	return &ETFRepository{db: db}
}

func (r *ETFRepository) CreateETF(etf models.ETF) error {
	query := `INSERT INTO etfs (symbol, name, description, expense_ratio) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRow(query, etf.Symbol, etf.Name, etf.Description, etf.ExpenseRatio).Scan(&etf.ID)
	if err != nil {
		return fmt.Errorf("error creating ETF: %w", err)
	}
	return nil
}

func (r *ETFRepository) GetETFBySymbol(symbol string) (*models.ETF, error) {
	query := `SELECT id, symbol, name, description, expense_ratio, created_at FROM etfs WHERE symbol = $1`
	var etf models.ETF
	err := r.db.QueryRow(query, symbol).Scan(&etf.ID, &etf.Symbol, &etf.Name, &etf.Description, &etf.ExpenseRatio, &etf.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &etf, nil
}

func (r *ETFRepository) GetAllETFs() ([]models.ETF, error) {
	query := `SELECT id, symbol, name, description, expense_ratio, created_at FROM etfs ORDER BY symbol`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var etfs []models.ETF
	for rows.Next() {
		var etf models.ETF
		if err := rows.Scan(&etf.ID, &etf.Symbol, &etf.Name, &etf.Description, &etf.ExpenseRatio, &etf.CreatedAt); err != nil {
			return nil, err
		}
		etfs = append(etfs, etf)
	}
	return etfs, nil
}

func (r *ETFRepository) CreatePrice(price models.Price) error {
	query := `INSERT INTO prices (etf_id, date, open, high, low, close, volume) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7) 
	          ON CONFLICT (etf_id, date) DO UPDATE 
	          SET open = $3, high = $4, low = $5, close = $6, volume = $7`
	_, err := r.db.Exec(query, price.ETFID, price.Date, price.Open, price.High, price.Low, price.Close, price.Volume)
	if err != nil {
		return fmt.Errorf("error creating price: %w", err)
	}
	return nil
}

func (r *ETFRepository) GetPricesByETFID(etfID int, limit int) ([]models.Price, error) {
	query := `SELECT id, etf_id, date, open, high, low, close, volume, created_at 
	          FROM prices WHERE etf_id = $1 ORDER BY date DESC LIMIT $2`
	rows, err := r.db.Query(query, etfID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prices []models.Price
	for rows.Next() {
		var price models.Price
		if err := rows.Scan(&price.ID, &price.ETFID, &price.Date, &price.Open, &price.High, &price.Low, &price.Close, &price.Volume, &price.CreatedAt); err != nil {
			return nil, err
		}
		prices = append(prices, price)
	}
	return prices, nil
}

func (r *ETFRepository) GetTopPerformers(period string, limit int) ([]models.ETF, error) {
	// Calculate performance based on period
	var days int
	switch period {
	case "1y":
		days = 365
	case "3y":
		days = 365 * 3
	case "5y":
		days = 365 * 5
	default:
		days = 365
	}

	query := `
		SELECT e.id, e.symbol, e.name, e.description, e.expense_ratio, e.created_at,
		       (p2.close - p1.close) / p1.close * 100 as performance
		FROM etfs e
		INNER JOIN prices p1 ON e.id = p1.etf_id
		INNER JOIN prices p2 ON e.id = p2.etf_id
		WHERE p1.date = (
		    SELECT date FROM prices WHERE etf_id = e.id ORDER BY date ASC LIMIT 1 OFFSET $1
		)
		AND p2.date = (
		    SELECT date FROM prices WHERE etf_id = e.id ORDER BY date DESC LIMIT 1
		)
		ORDER BY performance DESC
		LIMIT $2
	`

	rows, err := r.db.Query(query, days-1, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var etfs []models.ETF
	for rows.Next() {
		var etf models.ETF
		var performance float64
		if err := rows.Scan(&etf.ID, &etf.Symbol, &etf.Name, &etf.Description, &etf.ExpenseRatio, &etf.CreatedAt, &performance); err != nil {
			return nil, err
		}
		etfs = append(etfs, etf)
	}
	return etfs, nil
}

func (r *ETFRepository) CreatePlatform(platform models.Platform) error {
	prosJSON, _ := json.Marshal(platform.Pros)
	consJSON, _ := json.Marshal(platform.Cons)

	query := `INSERT INTO platforms (name, description, fees, pros, cons, rating, website) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := r.db.QueryRow(query, platform.Name, platform.Description, platform.Fees, prosJSON, consJSON, platform.Rating, platform.Website).Scan(&platform.ID)
	if err != nil {
		return fmt.Errorf("error creating platform: %w", err)
	}
	return nil
}

func (r *ETFRepository) GetAllPlatforms() ([]models.Platform, error) {
	query := `SELECT id, name, description, fees, pros, cons, rating, website FROM platforms ORDER BY rating DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var platforms []models.Platform
	for rows.Next() {
		var platform models.Platform
		var prosJSON, consJSON []byte
		if err := rows.Scan(&platform.ID, &platform.Name, &platform.Description, &platform.Fees, &prosJSON, &consJSON, &platform.Rating, &platform.Website); err != nil {
			return nil, err
		}
		json.Unmarshal(prosJSON, &platform.Pros)
		json.Unmarshal(consJSON, &platform.Cons)
		platforms = append(platforms, platform)
	}
	return platforms, nil
}

func (r *ETFRepository) CreateNews(news models.News) error {
	query := `INSERT INTO news (title, url, published_at, source, summary) 
	          VALUES ($1, $2, $3, $4, $5) ON CONFLICT (url) DO NOTHING`
	_, err := r.db.Exec(query, news.Title, news.URL, news.PublishedAt, news.Source, news.Summary)
	if err != nil {
		return fmt.Errorf("error creating news: %w", err)
	}
	return nil
}

func (r *ETFRepository) GetLatestNews(limit int) ([]models.News, error) {
	query := `SELECT id, title, url, published_at, source, summary FROM news ORDER BY published_at DESC LIMIT $1`
	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newsItems []models.News
	for rows.Next() {
		var news models.News
		if err := rows.Scan(&news.ID, &news.Title, &news.URL, &news.PublishedAt, &news.Source, &news.Summary); err != nil {
			return nil, err
		}
		newsItems = append(newsItems, news)
	}
	return newsItems, nil
}

func (r *ETFRepository) CreateAlert(alert models.PriceAlert) error {
	query := `INSERT INTO price_alerts (etf_id, email, threshold, direction) 
	          VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRow(query, alert.ETFID, alert.Email, alert.Threshold, alert.Direction).Scan(&alert.ID)
	if err != nil {
		return fmt.Errorf("error creating alert: %w", err)
	}
	return nil
}

func (r *ETFRepository) CheckAlerts() ([]models.PriceAlert, error) {
	query := `
		SELECT a.id, a.etf_id, a.email, a.threshold, a.direction, a.created_at, a.triggered, p.close
		FROM price_alerts a
		INNER JOIN prices p ON a.etf_id = p.etf_id
		WHERE a.triggered = false
		AND p.date = (SELECT MAX(date) FROM prices WHERE etf_id = a.etf_id)
		AND (
			(a.direction = 'above' AND p.close >= a.threshold) OR
			(a.direction = 'below' AND p.close <= a.threshold)
		)
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []models.PriceAlert
	for rows.Next() {
		var alert models.PriceAlert
		var currentPrice float64
		if err := rows.Scan(&alert.ID, &alert.ETFID, &alert.Email, &alert.Threshold, &alert.Direction, &alert.CreatedAt, &alert.Triggered, &currentPrice); err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}
	return alerts, nil
}

func (r *ETFRepository) MarkAlertTriggered(alertID int) error {
	query := `UPDATE price_alerts SET triggered = true WHERE id = $1`
	_, err := r.db.Exec(query, alertID)
	if err != nil {
		return fmt.Errorf("error marking alert as triggered: %w", err)
	}
	return nil
}

func (r *ETFRepository) InitializeSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS etfs (
		id SERIAL PRIMARY KEY,
		symbol VARCHAR(20) UNIQUE NOT NULL,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		expense_ratio DECIMAL(5, 4),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS prices (
		id SERIAL PRIMARY KEY,
		etf_id INTEGER REFERENCES etfs(id) ON DELETE CASCADE,
		date DATE NOT NULL,
		open DECIMAL(15, 4),
		high DECIMAL(15, 4),
		low DECIMAL(15, 4),
		close DECIMAL(15, 4),
		volume BIGINT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(etf_id, date)
	);

	CREATE TABLE IF NOT EXISTS platforms (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) UNIQUE NOT NULL,
		description TEXT,
		fees TEXT,
		pros JSONB,
		cons JSONB,
		rating DECIMAL(3, 2),
		website VARCHAR(255)
	);

	CREATE TABLE IF NOT EXISTS news (
		id SERIAL PRIMARY KEY,
		title VARCHAR(500) NOT NULL,
		url VARCHAR(1000) UNIQUE NOT NULL,
		published_at TIMESTAMP,
		source VARCHAR(100),
		summary TEXT
	);

	CREATE TABLE IF NOT EXISTS price_alerts (
		id SERIAL PRIMARY KEY,
		etf_id INTEGER REFERENCES etfs(id) ON DELETE CASCADE,
		email VARCHAR(255) NOT NULL,
		threshold DECIMAL(15, 4) NOT NULL,
		direction VARCHAR(10) CHECK (direction IN ('above', 'below')),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		triggered BOOLEAN DEFAULT FALSE
	);

	CREATE INDEX IF NOT EXISTS idx_prices_etf_date ON prices(etf_id, date);
	CREATE INDEX IF NOT EXISTS idx_news_published ON news(published_at DESC);
	`

	_, err := r.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("error initializing schema: %w", err)
	}
	log.Println("Database schema initialized successfully")
	return nil
}
