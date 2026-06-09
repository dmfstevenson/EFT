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
