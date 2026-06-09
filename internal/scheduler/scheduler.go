package scheduler

import (
	"log"
	"time"

	"etf-recommendation-api/internal/data"
	"etf-recommendation-api/internal/models"
)

type Scheduler struct {
	repo *data.ETFRepository
}

func NewScheduler(repo *data.ETFRepository) *Scheduler {
	return &Scheduler{repo: repo}
}

func (s *Scheduler) UpdateDailyPrices() {
	log.Println("Starting daily price update...")

	etfs, err := s.repo.GetAllETFs()
	if err != nil {
		log.Printf("Error fetching ETFs: %v", err)
		return
	}

	for _, etf := range etfs {
		endDate := time.Now()
		startDate := endDate.AddDate(0, 0, -7) // Fetch last 7 days

		priceData, err := data.FetchETFData(etf.Symbol, startDate, endDate)
		if err != nil {
			log.Printf("Error fetching data for %s: %v", etf.Symbol, err)
			continue
		}

		for _, pd := range priceData {
			price := models.Price{
				ETFID:  etf.ID,
				Date:   pd.Date,
				Open:   pd.Open,
				High:   pd.High,
				Low:    pd.Low,
				Close:  pd.Close,
				Volume: pd.Volume,
			}
			if err := s.repo.CreatePrice(price); err != nil {
				log.Printf("Error saving price for %s on %s: %v", etf.Symbol, pd.Date.Format("2006-01-02"), err)
			}
		}
		log.Printf("Updated prices for %s", etf.Symbol)
	}

	log.Println("Daily price update completed")
}

func (s *Scheduler) CheckPriceAlerts() {
	log.Println("Checking price alerts...")

	alerts, err := s.repo.CheckAlerts()
	if err != nil {
		log.Printf("Error checking alerts: %v", err)
		return
	}

	for _, alert := range alerts {
		// In production, send email notification here
		log.Printf("Alert triggered: ETF ID %d, Email %s, Threshold %.2f, Direction %s",
			alert.ETFID, alert.Email, alert.Threshold, alert.Direction)

		if err := s.repo.MarkAlertTriggered(alert.ID); err != nil {
			log.Printf("Error marking alert %d as triggered: %v", alert.ID, err)
		}
	}

	log.Printf("Checked %d price alerts", len(alerts))
}

func (s *Scheduler) Start() {
	// Run initial update
	s.UpdateDailyPrices()

	// Schedule daily updates at 6:00 PM UTC (2:00 AM PHT)
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		s.UpdateDailyPrices()
		s.CheckPriceAlerts()
	}
}
