package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/ericnberwick/daily-stox/service.pick-stock/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)



func GetStocks() ([]domain.StockRecommendation, error) {
	ctx := context.Background()
	_ = godotenv.Load(".env")

	// 1. Get Connection String and Pool
	connStr := os.Getenv("SUPABASE_DB_URL")
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	defer pool.Close()

	// 2. Define the Query
	// We select all columns that match your struct fields
	query := `
        SELECT 
            stock_name, stock_ticker, currency, current_price, "current_date", 
            estimated_sell_date, reason_for_buying, lynch_category, 
            conviction_score, peg_ratio, debt_to_equity_pct, 
            institutional_ownership_pct, shares_in_float, float_percentage, 
            catalyst_event, target_price, stop_loss, fundamental_trigger
        FROM "Stocks"`

	// 3. Execute the Query
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close() // Important: Always close rows to release the connection

	// 4. Iterate through the results
	var stocks []domain.StockRecommendation

	for rows.Next() {
		var s domain.StockRecommendation
		
		// The order of arguments in Scan must match the order in your SELECT statement
		err := rows.Scan(
			&s.StockName, &s.StockTicker, &s.Currency, &s.CurrentPrice, &s.CurrentDate,
			&s.EstimatedSellDate, &s.ReasonForBuying, &s.LynchCategory,
			&s.ConvictionScore, &s.PegRatio, &s.DebtToEquityPct,
			&s.InstitutionalOwnershipPct, &s.SharesInFloat, &s.FloatPercentage,
			&s.CatalystEvent, &s.TargetPrice, &s.StopLoss, &s.FundamentalTrigger,
		)
		if err != nil {
			fmt.Println("We got an error")
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		stocks = append(stocks, s)
	}

	// 5. Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}

	return stocks, nil
}