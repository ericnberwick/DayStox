package main

import (
	"fmt"

	"github.com/ericnberwick/daily-stox/application"
	"github.com/ericnberwick/daily-stox/repository"
)

func main(){
	
	stockPick := application.GetStock()
	fmt.Println("Stock pick is: ",stockPick.StockName)
	application.UpdateStockPriceByTicker(stockPick)
	err := repository.InsertStock(*stockPick)
	if err != nil {
		fmt.Errorf("error inserting stock: %w", err)
	}
	fmt.Printf("Sucessfully picked stock %s today and inserted into the database", stockPick.StockName)
}