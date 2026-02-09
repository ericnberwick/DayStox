package main

import (
	"fmt"
	"log" // Use log for better timestamps in GCP

	"github.com/ericnberwick/daily-stox/service.pick-stock/application"
	"github.com/ericnberwick/daily-stox/service.pick-stock/repository"
)

func main() {

	stockPick := application.GetStock()
	if stockPick == nil {
		log.Fatal("Failed to get stock pick")
	}
	fmt.Println("Stock pick is:", stockPick.StockName)

	// 2. Update price
	application.UpdateStockPriceByTicker(stockPick)

	// 3. Insert into Database
	err := repository.InsertStock(*stockPick)
	if err != nil {
		// This will stop the program and show the REAL error in GCP logs
		fmt.Println("CRITICAL DATABASE ERROR: %v", err)
	}

	// 4. Final Success Log
	fmt.Printf("SUCCESS: Picked %s and confirmed DB insertion.", stockPick.StockName)
}