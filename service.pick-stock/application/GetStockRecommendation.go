package application

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ericnberwick/daily-stox/service.pick-stock/domain"
	"github.com/ericnberwick/daily-stox/service.pick-stock/repository"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

func GetStock() *domain.StockRecommendation {
	ctx := context.Background()
	_ = godotenv.Load(".env")

    geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	cc := &genai.ClientConfig{
		APIKey: geminiAPIKey,
	}
    client, err := genai.NewClient(ctx, cc)
	
    if err != nil {
        log.Fatal(err)
    }

    result, err := client.Models.GenerateContent(
        ctx,
        "gemini-3-flash-preview",
        genai.Text(GetPrompt()),
        nil,
    )
	if err != nil {
        log.Printf("Gemini API Error: %v", err)
        return nil // or handle appropriately
    }
	stocks, err := ProcessAIResponse(result)
	if err != nil {
		log.Printf("Error processing: %v", err)
		
	}
	
	// Now you can access fields directly!
	if len(stocks) > 0 {
		fmt.Printf("Stock: %s, Currency: %s\n", stocks[0].StockName, stocks[0].Currency)
	}
	var stockPick domain.StockRecommendation = stocks[0]
	return &stockPick
}

func GetPrompt() string {
	stx,err := repository.GetStocks()
	if err != nil {
		fmt.Println("Error detected:", err) 
	}
	var stkTickers []string
	for _, stk := range stx {
		stkTickers = append(stkTickers, stk.StockTicker)
	}
    currentDate := time.Now().Format("02/01/2006")
	template := `
		Imagine that you are a financial analyst in particular Peter Lynch.
		You must pick one stock that you think will double, triple or more over the next week.
		The stock must be from U.S. Equities (Stocks & ETFs) and be from 1 of the 19 major U.S. exchanges (NYSE, NASDAQ, etc.)
		The current date is %s, dates should be in RFC3339 format.
		DO NOT SUGGEST ANY OF THE FOLLOWING COMPANY'S WITH THE TICKER : %s
		Pick the stock, show the price, and estimate a date for optimal sell.
		Output the result in JSON.
		For example:
		[
		{
			"stock_name": "Ascendis Pharma",
			"stock_ticker": "ASND",
			"currency" : US Dollar ($)
			"current_price": 145.20,
			"current_date" : "2026-02-03T00:00:00Z",
			"estimated_sell_date": "2026-02-10T00:00:00Z",
			"reason_for_buying": "Major FDA PDUFA date for TransCon CNP (Achondroplasia). Positive approval would open a multi-billion dollar market in pediatric growth disorders, likely triggering a significant valuation re-rating immediately."
			"lynch_category": "Fast Grower / Special Situation",
			"conviction_score": 8,
			"peg_ratio": 1.1,
			"debt_to_equity": 22,
			"institutional_ownership": "Low",
			"shares_in_float": "58.2M",
			"float_percentage": 98,
			"catalyst_event": "FDA Approval Decision (PDUFA)",
			"target_price": 210.00,
			"stop_loss": 115.00,
			"fundamental_trigger": "Sell immediately if FDA issues a Complete Response Letter (rejection)."
		},
		]
		`
	return fmt.Sprintf(template, currentDate, stkTickers)
}

func ProcessAIResponse(resp *genai.GenerateContentResponse) ([]domain.StockRecommendation, error) {
	// 1. Extract the raw text from the first candidate
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no content returned from AI")
	}
	
	rawText := resp.Candidates[0].Content.Parts[0].Text

	// 2. Clean Markdown backticks if they exist
	// AI often returns ```json [ { ... } ] ```
	cleanedJSON := strings.TrimSpace(rawText)
	cleanedJSON = strings.TrimPrefix(cleanedJSON, "```json")
	cleanedJSON = strings.TrimSuffix(cleanedJSON, "```")
	cleanedJSON = strings.TrimSpace(cleanedJSON)

	// 3. Unmarshal into the slice (since AI returns an array [])
	var recommendations []domain.StockRecommendation
	err := json.Unmarshal([]byte(cleanedJSON), &recommendations)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return recommendations, nil
}