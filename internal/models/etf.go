package models

import "time"

type ETF struct {
	ID           int     `json:"id" db:"id"`
	Symbol       string  `json:"symbol" db:"symbol"`
	Name         string  `json:"name" db:"name"`
	Description  string  `json:"description" db:"description"`
	ExpenseRatio float64 `json:"expense_ratio" db:"expense_ratio"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type Price struct {
	ID        int       `json:"id" db:"id"`
	ETFID     int       `json:"etf_id" db:"etf_id"`
	Date      time.Time `json:"date" db:"date"`
	Open      float64   `json:"open" db:"open"`
	High      float64   `json:"high" db:"high"`
	Low       float64   `json:"low" db:"low"`
	Close     float64   `json:"close" db:"close"`
	Volume    int64     `json:"volume" db:"volume"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Platform struct {
	ID          int      `json:"id" db:"id"`
	Name        string   `json:"name" db:"name"`
	Description string   `json:"description" db:"description"`
	Fees        string   `json:"fees" db:"fees"`
	Pros        []string `json:"pros" db:"pros"`
	Cons        []string `json:"cons" db:"cons"`
	Rating      float64  `json:"rating" db:"rating"`
	Website     string   `json:"website" db:"website"`
}

type News struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	URL         string    `json:"url" db:"url"`
	PublishedAt time.Time `json:"published_at" db:"published_at"`
	Source      string    `json:"source" db:"source"`
	Summary     string    `json:"summary" db:"summary"`
}

type PriceAlert struct {
	ID        int       `json:"id" db:"id"`
	ETFID     int       `json:"etf_id" db:"etf_id"`
	Email     string    `json:"email" db:"email"`
	Threshold float64   `json:"threshold" db:"threshold"`
	Direction string    `json:"direction" db:"direction"` // "above" or "below"
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Triggered bool      `json:"triggered" db:"triggered"`
}
