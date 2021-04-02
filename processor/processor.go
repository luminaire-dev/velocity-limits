package processor

import (
	"strconv"
	"time"
)

const (
	dailyLimit   = 5000
	weeklyLimit  = 20000
	maxLoadCount = 3
)

type Processor struct {
	approvedLoads  map[string][]Load
	processedLoads map[string]Response
}

func NewProcessor() *Processor {
	return &Processor{
		approvedLoads:  make(map[string][]Load),
		processedLoads: make(map[string]Response),
	}
}

// Process approves or rejects the incoming load based on limits validation
func (p *Processor) Load(incomingLoad Load) *Response {
	// approve status set to true be default
	approved := true

	// parse load amount to float to use it mathematically
	incomingLoadAmount, err := strconv.ParseFloat(incomingLoad.LoadAmount[1:], 64)
	if err != nil {
		panic(err)
	}

	// reject any incoming load amount greater than 5000, as it already fails the daily limit check
	if incomingLoadAmount > dailyLimit {
		return p.newResponse(incomingLoad, false)
	}

	// ignore load with duplicate ID
	if p.processedLoads[incomingLoad.ID].CustomerID == incomingLoad.CustomerID {
		return nil
	}

	// iterate over processed load to get daily and weekly data per customer
	var totalDailyAmount float64 = 0
	var totalWeeklyAmount float64 = 0
	var transactionCount = 0
	for _, processedLoad := range p.approvedLoads[incomingLoad.CustomerID] {

		// parse processed load amount to float to use it mathematically
		processedLoadAmount, err := strconv.ParseFloat(processedLoad.LoadAmount[1:], 64)
		if err != nil {
			panic(err)
		}

		// add daily total and incrment cout for same day loads
		if IsSameDay(incomingLoad.Time, processedLoad.Time) {
			totalDailyAmount += processedLoadAmount
			transactionCount++
		}

		// add weekly total for same week loads
		if IsSameWeek(incomingLoad.Time, processedLoad.Time) {
			totalWeeklyAmount += processedLoadAmount
		}
	}

	// reject load if daily transaction count is exceeded
	if transactionCount == maxLoadCount {
		approved = false
		return p.newResponse(incomingLoad, approved)
	}

	// reject load if daily trasaction limit is exceeded
	if totalDailyAmount+incomingLoadAmount > dailyLimit {
		approved = false
		return p.newResponse(incomingLoad, approved)
	}

	// reject load if weekly trasaction limit is exceeded
	if totalWeeklyAmount+incomingLoadAmount > weeklyLimit {
		approved = false
		return p.newResponse(incomingLoad, approved)
	}

	// store load into approved loads
	p.approvedLoads[incomingLoad.CustomerID] = append(p.approvedLoads[incomingLoad.CustomerID], incomingLoad)

	return p.newResponse(incomingLoad, approved)
}

// newResponse retuns a response that will be added to the output file
func (p *Processor) newResponse(load Load, accepted bool) *Response {
	r := Response{
		ID:         load.ID,
		CustomerID: load.CustomerID,
		Accepted:   accepted,
	}

	// store load into processed Loads
	p.processedLoads[load.ID] = r

	return &r
}

// check if two date are on the same day
func IsSameDay(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// check if two dates are in the same week
func IsSameWeek(incomingLoadTime, processedLoadTime time.Time) bool {
	year, week := incomingLoadTime.ISOWeek()
	year2, week2 := processedLoadTime.ISOWeek()
	return year == year2 && week == week2
}
