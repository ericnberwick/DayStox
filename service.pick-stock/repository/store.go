package repository

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ericnberwick/daily-stox/service.pick-stock/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func InsertStock(stock domain.StockRecommendation) error {
	ctx := context.Background()

	// Load .env for local testing (Cloud Run will ignore this if no .env file exists)
	_ = godotenv.Load(".env")

	// 1. Get Connection String
	connStr := os.Getenv("SUPABASE_DB_URL")
	if connStr == "" {
		return fmt.Errorf("SUPABASE_DB_URL environment variable is not set")
	}

	// 2. Create a Connection Pool
	// This is better for Cloud Run as it manages connections more reliably
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %v", err)
	}
	defer pool.Close()

	// 3. Define the Query
	// Note: Using double quotes "Stocks" makes the table name case-sensitive in Postgres
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

	// 4. Execute the Insert
	// Exec returns a 'tag' which contains the count of rows affected
	tag, err := pool.Exec(ctx, query,
		stock.StockName,
		stock.StockTicker,
		stock.Currency,
		stock.CurrentPrice,
		stock.CurrentDate,
		stock.EstimatedSellDate,
		stock.ReasonForBuying,
		stock.LynchCategory,
		stock.ConvictionScore,
		stock.PegRatio,
		stock.DebtToEquityPct,
		stock.InstitutionalOwnershipPct,
		stock.SharesInFloat,
		stock.FloatPercentage,
		stock.CatalystEvent,
		stock.TargetPrice,
		stock.StopLoss,
		stock.FundamentalTrigger,
	)

	if err != nil {
		return fmt.Errorf("database insert failed: %v", err)
	}

	// 5. Explicit Log for GCP
	// If you see 'Rows affected: 0', your RLS policies in Supabase are blocking the write!
	log.Printf("Successfully inserted into DB. Rows affected: %d", tag.RowsAffected())

	return nil
}