package data

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type YahooFinanceResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Symbol    string `json:"symbol"`
				Exchange  string `json:"exchange"`
				Currency  string `json:"currency"`
				RegularPrice float64 `json:"regularMarketPrice"`
			} `json:"meta"`
			Timestamp  []int64  `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Open  []float64 `json:"open"`
					High  []float64 `json:"high"`
					Low   []float64 `json:"low"`
					Close []float64 `json:"close"`
					Volume []int64  `json:"volume"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
	} `json:"chart"`
}

type ETFPriceData struct {
	Symbol string
	Date   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}

func FetchETFData(symbol string, startDate, endDate time.Time) ([]ETFPriceData, error) {
	start := startDate.Unix()
	end := endDate.Unix()
	
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?period1=%d&period2=%d&interval=1d", symbol, start, end)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}
	
	var data YahooFinanceResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}
	
	if len(data.Chart.Result) == 0 {
		return nil, fmt.Errorf("no data returned for symbol %s", symbol)
	}
	
	result := data.Chart.Result[0]
	var priceData []ETFPriceData
	
	for i, ts := range result.Timestamp {
		if i >= len(result.Indicators.Quote[0].Close) {
			break
		}
		
		close := result.Indicators.Quote[0].Close[i]
		if close == 0 {
			continue
		}
		
		priceData = append(priceData, ETFPriceData{
			Symbol: symbol,
			Date:   time.Unix(ts, 0),
			Open:   result.Indicators.Quote[0].Open[i],
			High:   result.Indicators.Quote[0].High[i],
			Low:    result.Indicators.Quote[0].Low[i],
			Close:  close,
			Volume: result.Indicators.Quote[0].Volume[i],
		})
	}
	
	return priceData, nil
}

func GetCurrentPrice(symbol string) (float64, error) {
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?interval=1d", symbol)
	
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error fetching data: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response: %w", err)
	}
	
	var data YahooFinanceResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, fmt.Errorf("error parsing JSON: %w", err)
	}
	
	if len(data.Chart.Result) == 0 {
		return 0, fmt.Errorf("no data returned for symbol %s", symbol)
	}
	
	return data.Chart.Result[0].Meta.RegularPrice, nil
}
