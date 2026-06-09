package main

import (
	"log"
	"os"

	"etf-recommendation-api/internal/api"
	"etf-recommendation-api/internal/data"
	"etf-recommendation-api/internal/models"
	"etf-recommendation-api/internal/scheduler"

	"github.com/gin-gonic/gin"
)

func main() {
	// Get database URL from environment variable
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://localhost/etf_recommendations?sslmode=disable"
	}

	// Connect to database
	if err := data.ConnectDB(dbURL); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer data.CloseDB()

	// Initialize repository
	repo := data.NewETFRepository(data.DB)

	// Initialize database schema
	if err := repo.InitializeSchema(); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	// Seed initial data if needed
	seedData(repo)

	// Start scheduler in background
	sched := scheduler.NewScheduler(repo)
	go sched.Start()

	// Setup Gin router
	router := gin.Default()

	// Enable CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Setup handlers
	handler := api.NewHandler(repo)

	// API routes
	api := router.Group("/api")
	{
		api.GET("/health", handler.HealthCheck)
		api.GET("/etfs", handler.GetETFs)
		api.GET("/etfs/:symbol", handler.GetETFBySymbol)
		api.GET("/etfs/:symbol/prices", handler.GetETFPrices)
		api.GET("/etfs/top", handler.GetTopPerformers)
		api.GET("/platforms", handler.GetPlatforms)
		api.GET("/news", handler.GetNews)
		api.POST("/alerts", handler.CreateAlert)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func seedData(repo *data.ETFRepository) {
	// Seed Philippine ETFs
	philippineETFs := []models.ETF{
		{
			Symbol:       "FMETF.PS",
			Name:         "First Metro Phil. Equity Exchange Traded Fund",
			Description:  "Tracks the PSEi, providing exposure to the Philippine stock market",
			ExpenseRatio: 0.0050,
		},
		{
			Symbol:       "PSEETF.PS",
			Name:         "PSEi ETF",
			Description:  "Another PSEi-tracking ETF",
			ExpenseRatio: 0.0050,
		},
	}

	for _, etf := range philippineETFs {
		_, err := repo.GetETFBySymbol(etf.Symbol)
		if err != nil {
			log.Printf("Seeding ETF: %s", etf.Symbol)
			if err := repo.CreateETF(etf); err != nil {
				log.Printf("Error seeding ETF %s: %v", etf.Symbol, err)
			}
		}
	}

	// Seed trading platforms
	platforms := []models.Platform{
		{
			Name:        "COL Financial",
			Description: "One of the largest online stockbrokers in the Philippines",
			Fees:        "0.25% commission, minimum PHP 20",
			Pros:        []string{"Low minimum investment", "User-friendly platform", "Research tools available"},
			Cons:        []string{"Limited to Philippine market", "Customer service can be slow"},
			Rating:      4.5,
			Website:     "https://www.colfinancial.com",
		},
		{
			Name:        "FirstMetroSec",
			Description: "Online stockbroker backed by Metrobank",
			Fees:        "0.25% commission, minimum PHP 20",
			Pros:        []string{"Backed by major bank", "Low fees", "Easy funding via Metrobank accounts"},
			Cons:        []string{"Platform less polished than competitors", "Limited research tools"},
			Rating:      4.3,
			Website:     "https://www.firstmetrosec.com.ph",
		},
		{
			Name:        "PhilStocks",
			Description: "Online stockbroker with competitive pricing",
			Fees:        "0.20% commission, minimum PHP 20",
			Pros:        []string{"Lowest commission rates", "Fast execution", "Good mobile app"},
			Cons:        []string{"Higher minimum balance required", "Limited educational resources"},
			Rating:      4.2,
			Website:     "https://www.philstocks.com.ph",
		},
	}

	for _, platform := range platforms {
		existing, err := repo.GetAllPlatforms()
		if err != nil || len(existing) == 0 {
			log.Printf("Seeding platform: %s", platform.Name)
			if err := repo.CreatePlatform(platform); err != nil {
				log.Printf("Error seeding platform %s: %v", platform.Name, err)
			}
		}
	}
}
