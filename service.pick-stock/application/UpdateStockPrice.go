package application

import (
	"context"
	"fmt"
	"log"
	"time"

	"os"

	"github.com/ericnberwick/daily-stox/service.pick-stock/domain"
	"github.com/joho/godotenv"
	massive "github.com/massive-com/client-go/v2/rest"
	"github.com/massive-com/client-go/v2/rest/models"
)

func UpdateStockPriceByTicker(stock *domain.StockRecommendation) error {
	_ = godotenv.Load(".env")
	massiveApiKey := os.Getenv("MASSIVE_API_KEY")
	c := massive.New(massiveApiKey)

	params := &models.GetDailyOpenCloseAggParams{
        Ticker: stock.StockTicker,
        // Format: YYYY-MM-DD (Previous business day)
		Date: models.Date(lastTradingDay()), 
    }

    res, err := c.GetDailyOpenCloseAgg(context.Background(), params)
    if err != nil {
		fmt.Println("Searching for stock: ", stock.StockTicker)
        log.Fatal("Error fetching daily data: ", err)
    }

    fmt.Printf("The stock: %s closed at price: %.2f on %s\n", res.Symbol, res.Close, res.From)

	stock.CurrentPrice = res.Close

	return nil
}


func lastTradingDay() time.Time {
    t := time.Now()
    switch t.Weekday() {
    case time.Sunday:
        return t.AddDate(0, 0, -2) // Friday
    case time.Monday:
        return t.AddDate(0, 0, -3) // Friday
    case time.Saturday:
        return t.AddDate(0, 0, -1) // Friday
    default:
        return t.AddDate(0, 0, -1) // Previous day
    }
}