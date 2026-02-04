package domain

import (
	"time"
)
type StockRecommendation struct {
	ID                     int64   `db:"id"`
	StockName             string  `json:"stock_name"`
	StockTicker           string  `json:"stock_ticker"`
	Currency			  string 	`json:"currency"`
	CurrentPrice          float64  `json:"current_price"`
	CurrentDate 		  time.Time `json:"current_date"`
	EstimatedSellDate    time.Time  `json:"estimated_sell_date"`
	ReasonForBuying       string  `json:"reason_for_buying"`
	LynchCategory         string  `json:"lynch_category"`
	ConvictionScore       int     `json:"conviction_score"`
	PegRatio              float64 `json:"peg_ratio"`
	DebtToEquityPct          float32  `json:"debt_to_equity_pct"`
	InstitutionalOwnershipPct float32  `json:"institutional_ownership_pct"`
	SharesInFloat         string  `json:"shares_in_float"`
	FloatPercentage       float32  `json:"float_percentage"`
	CatalystEvent         string  `json:"catalyst_event"`
	TargetPrice           float64  `json:"target_price"`
	StopLoss              float64  `json:"stop_loss"`
	FundamentalTrigger    string  `json:"fundamental_trigger"`
}