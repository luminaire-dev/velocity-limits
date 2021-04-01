package processor

import (
	"encoding/json"
	"strconv"
	"time"
)

const (
	dailyLimit   = 5000
	weeklyLimit  = 20000
	maxLoadCount = 3
)

var (
	approvedLoads = make(map[string][]Load)
	processedLoads = make(map[string]Response)
)

type Load struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customer_id"`
	LoadAmount string    `json:"load_amount"`
	Time       time.Time `json:"time"`
}

type Response struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	Accepted   bool   `json:"accepted"`
}

// approve or reject the incoming load based on limits validation
func Process(load []byte)[]byte {

	// parse the JSON-encoded data
	var incomingLoad Load
	if err := json.Unmarshal(load, &incomingLoad)
	err != nil {
		panic(err)
	}

	// parse load amount to float to use it mathematically
	incomingLoadAmount, parseErr := strconv.ParseFloat(incomingLoad.LoadAmount[1:], 64)
	if parseErr != nil {
	   panic(parseErr)
	 }

	// reject any incoming load amount greater than 5000, as it already fails the daily limit check
	if(incomingLoadAmount > dailyLimit){
		return GenerateResponse(incomingLoad, false)
	}

	// ignore load with duplicate ID
	if(processedLoads[incomingLoad.ID].CustomerID == incomingLoad.CustomerID) {
		return nil
	}

	// iterate over processed load to get daily and weekly data per customer
	var totalDailyAmount float64 = 0;
	var totalWeeklyAmount float64 = 0;
	var transactionCount = 0;
	for _, processedLoad := range approvedLoads[incomingLoad.CustomerID] {

		// parse processed load amount to float to use it mathematically
		processedLoadAmount, parseErr := strconv.ParseFloat(processedLoad.LoadAmount[1:], 64)
		if parseErr != nil {
			 panic(parseErr)
		 }

		// add daily total and incrment cout for same day loads
		if(DoesDayMatch(incomingLoad.Time, processedLoad.Time)){
			 totalDailyAmount += processedLoadAmount
			 transactionCount ++
		}

		// add weekly total for same week loads
	  if(DoesWeekMatch(incomingLoad.Time, processedLoad.Time)) {
			totalWeeklyAmount += processedLoadAmount
	  }
	}

	// reject load if daily transaction count is exceeded
	if(transactionCount == maxLoadCount){
		return GenerateResponse(incomingLoad, false)
	}

	// reject load if daily trasaction limit is exceeded
	if(totalDailyAmount + incomingLoadAmount > dailyLimit){
		return GenerateResponse(incomingLoad, false)
	}

	// reject load if weekly trasaction limit is exceeded
	if(totalWeeklyAmount + incomingLoadAmount > weeklyLimit){
		return GenerateResponse(incomingLoad, false)
	}

	// store load into approved loads
	approvedLoads[incomingLoad.CustomerID] = append(approvedLoads[incomingLoad.CustomerID], incomingLoad)

	return GenerateResponse(incomingLoad, true)
}

// generate the json response that will be added to the output file
func GenerateResponse(load Load, accepted bool) []byte {
	r := Response{
		ID:         load.ID,
		CustomerID: load.CustomerID,
		Accepted:   accepted,
	}

	// store load into processed Loads
	processedLoads[load.ID] = r

	res, err := json.Marshal(&r)
	if err != nil {
		panic(err)
	}
	return res
}

// check if two date are on the same day
func DoesDayMatch(date1, date2 time.Time) bool {
    y1, m1, d1 := date1.Date()
    y2, m2, d2 := date2.Date()
    return y1 == y2 && m1 == m2 && d1 == d2
}

// check if two dates are in the same week
func DoesWeekMatch(incomingLoadTime, processedLoadTime time.Time) bool {
		year, week := incomingLoadTime.ISOWeek()
		year2, week2 := processedLoadTime.ISOWeek()
		return year == year2 && week == week2
}
