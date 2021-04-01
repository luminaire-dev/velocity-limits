package processor

import (
	"encoding/json"
	"strconv"
	"fmt"
	"time"
)

const (
	dailyLimit   = 500000
	weeklyLimit  = 2000000
	maxLoadCount = 3
)

type Load struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customer_id"`
	LoadAmount string    `json:"load_amount"`
	Time       time.Time `json:"time"`
}

// Approves or rejects the incoming load based on limits
func Process(load []byte) bool {

	// parse the JSON-encoded data
	var incomingLoad Load
	if err := json.Unmarshal(load, &incomingLoad)
	err != nil {
		panic(err)
	}
	fmt.Println(incomingLoad)

	// parse load amount to float so we can use it mathematically
	loadAmount, parseErr := strconv.ParseFloat(incomingLoad.LoadAmount[1:], 64)
	if parseErr != nil {
	   panic(parseErr)
	 }

	if(loadAmount > 5000){
	  return false
	}

	return true
}
