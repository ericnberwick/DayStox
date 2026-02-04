package application

import (
	"context"
	"fmt"
	"log"
	"time"

	"os"

	"github.com/ericnberwick/daily-stox/domain"
	"github.com/joho/godotenv"
	massive "github.com/massive-com/client-go/v2/rest"
	"github.com/massive-com/client-go/v2/rest/models"
)

func UpdateStockPriceByTicker(stock *domain.StockRecommendation) error {
	godotenv.Load()
	massiveApiKey := os.Getenv("MASSIVE_API_KEY")
	c := massive.New(massiveApiKey)

	params := &models.GetDailyOpenCloseAggParams{
        Ticker: stock.StockTicker,
        // Format: YYYY-MM-DD (Previous business day)
        Date:   models.Date(time.Now().AddDate(0, 0, -1)), 
    }

    res, err := c.GetDailyOpenCloseAgg(context.Background(), params)
    if err != nil {
        log.Fatal("Error fetching daily data: ", err)
    }

    fmt.Printf("The stock: %s closed at price: %.2f on %s\n", res.Symbol, res.Close, res.From)

	stock.CurrentPrice = res.Close

	return nil
}
