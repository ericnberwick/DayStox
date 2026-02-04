package repository

import (
	"context"
	"fmt"

	"os"

	"github.com/ericnberwick/daily-stox/domain"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func InsertStock(stock domain.StockRecommendation) error {
	ctx := context.Background()
	godotenv.Load() 

    supabaseDbURLapiKey := os.Getenv("SUPABASE_DB_URL")
	// Get this from your Supabase Dashboard
	dbURL := supabaseDbURLapiKey
	connStr := dbURL
	
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	query := `
		INSERT INTO "Stocks" (
			stock_name, stock_ticker, currency, current_price, "current_date", 
			estimated_sell_date, reason_for_buying, lynch_category, 
			conviction_score, peg_ratio, debt_to_equity_pct, 
			institutional_ownership_pct, shares_in_float, float_percentage, 
			catalyst_event, target_price, stop_loss, fundamental_trigger
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		)`

	_, err = conn.Exec(ctx, query,
		stock.StockName, stock.StockTicker, stock.Currency, stock.CurrentPrice, stock.CurrentDate,
		stock.EstimatedSellDate, stock.ReasonForBuying, stock.LynchCategory,
		stock.ConvictionScore, stock.PegRatio, stock.DebtToEquityPct,
		stock.InstitutionalOwnershipPct, stock.SharesInFloat, stock.FloatPercentage,
		stock.CatalystEvent, stock.TargetPrice, stock.StopLoss, stock.FundamentalTrigger,
	)

	return err
}

