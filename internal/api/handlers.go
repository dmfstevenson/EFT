package api

import (
	"net/http"
	"strconv"

	"etf-recommendation-api/internal/data"
	"etf-recommendation-api/internal/models"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo *data.ETFRepository
}

func NewHandler(repo *data.ETFRepository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) GetETFs(c *gin.Context) {
	etfs, err := h.repo.GetAllETFs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, etfs)
}

func (h *Handler) GetETFBySymbol(c *gin.Context) {
	symbol := c.Param("symbol")
	etf, err := h.repo.GetETFBySymbol(symbol)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ETF not found"})
		return
	}
	c.JSON(http.StatusOK, etf)
}

func (h *Handler) GetETFPrices(c *gin.Context) {
	symbol := c.Param("symbol")
	etf, err := h.repo.GetETFBySymbol(symbol)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ETF not found"})
		return
	}

	limit := 365
	if l := c.Query("limit"); l != "" {
		if num, err := strconv.Atoi(l); err == nil && num > 0 {
			limit = num
		}
	}

	prices, err := h.repo.GetPricesByETFID(etf.ID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prices)
}

func (h *Handler) GetTopPerformers(c *gin.Context) {
	period := c.DefaultQuery("period", "1y")
	limit := 10
	if l := c.Query("limit"); l != "" {
		if num, err := strconv.Atoi(l); err == nil && num > 0 {
			limit = num
		}
	}

	etfs, err := h.repo.GetTopPerformers(period, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, etfs)
}

func (h *Handler) GetPlatforms(c *gin.Context) {
	platforms, err := h.repo.GetAllPlatforms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, platforms)
}

func (h *Handler) GetNews(c *gin.Context) {
	limit := 20
	if l := c.Query("limit"); l != "" {
		if num, err := strconv.Atoi(l); err == nil && num > 0 {
			limit = num
		}
	}

	news, err := h.repo.GetLatestNews(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, news)
}

func (h *Handler) CreateAlert(c *gin.Context) {
	var alert models.PriceAlert
	if err := c.ShouldBindJSON(&alert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.CreateAlert(alert); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, alert)
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}
