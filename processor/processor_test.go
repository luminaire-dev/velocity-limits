package processor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsSameDay(t *testing.T) {
	result := false

	// dates should be on same day
	result = IsSameDay(time.Now(), time.Now())
	assert.Equal(t, result, true)

	result = IsSameDay(time.Date(2000, 1, 2, 1, 0, 0, 0, time.UTC), time.Date(2000, 1, 2, 4, 30, 0, 0, time.UTC))
	assert.Equal(t, result, true)

	result = IsSameDay(time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC), time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC))
	assert.Equal(t, result, true)

	// dates should not be on same day
	result = IsSameDay(time.Now(), time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC))
	assert.Equal(t, result, false)

	result = IsSameDay(time.Date(2000, 1, 2, 1, 0, 0, 0, time.UTC), time.Date(2000, 1, 6, 4, 30, 0, 0, time.UTC))
	assert.Equal(t, result, false)

	result = IsSameDay(time.Date(2002, 1, 2, 1, 0, 0, 0, time.UTC), time.Date(2000, 1, 2, 4, 30, 0, 0, time.UTC))
	assert.Equal(t, result, false)
}

func TestIsWeekDay(t *testing.T) {
	result := false

	// dates should be in same week
	result = IsSameWeek(time.Now(), time.Now())
	assert.Equal(t, result, true)

	result = IsSameWeek(time.Date(2000, 1, 4, 0, 0, 0, 0, time.UTC), time.Date(2000, 1, 8, 0, 0, 0, 0, time.UTC))
	assert.Equal(t, result, true)

	result = IsSameWeek(time.Date(2000, 6, 30, 1, 0, 0, 0, time.UTC), time.Date(2000, 7, 1, 1, 0, 0, 0, time.UTC))
	assert.Equal(t, result, true)

	// dates should not be in same week
	result = IsSameWeek(time.Now(), time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC))
	assert.Equal(t, result, false)

	result = IsSameWeek(time.Date(2000, 6, 1, 1, 0, 0, 0, time.UTC), time.Date(2000, 7, 1, 1, 0, 0, 0, time.UTC))
	assert.Equal(t, result, false)

	result = IsSameWeek(time.Date(2002, 1, 2, 1, 0, 0, 0, time.UTC), time.Date(2000, 1, 2, 1, 0, 0, 0, time.UTC))
	assert.Equal(t, result, false)
}

func TestSuccessfulLoadAttempt(t *testing.T) {
	// load attempt should be approved - no other loads deposited for customer
	p := NewProcessor()
	result := p.Load(Load{
		ID:         "1",
		CustomerID: "800",
		LoadAmount: "$3382.87",
		Time:       time.Now(),
	})

	resp := Response{
		ID:         "1",
		CustomerID: "800",
		Accepted:   true,
	}

	assert.Equal(t, resp, *result)

	// load attempt should be approved  - daily and weekly limit not reached
	result = p.Load(Load{
		ID:         "2",
		CustomerID: "800",
		LoadAmount: "$100.00",
		Time:       time.Date(2021, 4, 1, 1, 0, 0, 0, time.UTC),
	})

	resp = Response{
		ID:         "2",
		CustomerID: "800",
		Accepted:   true,
	}

	assert.Equal(t, resp, *result)

	// load attempt should be approved - daily and weekly limit not reached
	result = p.Load(Load{
		ID:         "3",
		CustomerID: "800",
		LoadAmount: "$100.00",
		Time:       time.Date(2021, 4, 9, 1, 0, 0, 0, time.UTC),
	})

	resp = Response{
		ID:         "3",
		CustomerID: "800",
		Accepted:   true,
	}

	assert.Equal(t, resp, *result)
}

func TestRejectLoadDueToDailyTotalLimit(t *testing.T) {

	// deposit load of $4000 for customer 123
	p := NewProcessor()
	p.Load(Load{
		ID:         "1",
		CustomerID: "123",
		LoadAmount: "$4000.00",
		Time:       time.Now(),
	})

	// load attempts should be rejected - daily limit exceeded
	proccesedLoad := p.Load(Load{
		ID:         "2",
		CustomerID: "123",
		LoadAmount: "$2000.00",
		Time:       time.Now(),
	})

	resp := Response{
		ID:         "2",
		CustomerID: "123",
		Accepted:   false,
	}

	assert.Equal(t, resp, *proccesedLoad)

	// load attempt should be rejected - daily limit exceede in single load (no other customer loads on this day)
	p = NewProcessor()
	proccesedLoad = p.Load(Load{
		ID:         "3",
		CustomerID: "1234",
		LoadAmount: "$5001.00",
		Time:       time.Now(),
	})

	resp = Response{
		ID:         "3",
		CustomerID: "1234",
		Accepted:   false,
	}

	assert.Equal(t, resp, *proccesedLoad)

}

func TestRejectLoadDueToDailyLoadCountLimit(t *testing.T) {
	p := NewProcessor()

	// deposit 3 loads to customer 123 on same day
	p.Load(Load{
		ID:         "1",
		CustomerID: "123",
		LoadAmount: "$100.00",
		Time:       time.Now(),
	})

	p.Load(Load{
		ID:         "2",
		CustomerID: "123",
		LoadAmount: "$100.00",
		Time:       time.Now(),
	})

	p.Load(Load{
		ID:         "3",
		CustomerID: "123",
		LoadAmount: "$100.00",
		Time:       time.Now(),
	})

	// fourth load should be rejected as daily limit count is exceeded

	processedLoad := p.Load(Load{
		ID:         "4",
		CustomerID: "123",
		LoadAmount: "$100.00",
		Time:       time.Now(),
	})

	resp := Response{
		ID:         "4",
		CustomerID: "123",
		Accepted:   false,
	}

	assert.Equal(t, resp, *processedLoad)
}

func TestRejectLoadDueToWeeklyTotalLimit(t *testing.T) {
	p := NewProcessor()

	// deposit 2 loads totaling 20000 for customer 123, 2 days apart
	p.Load(Load{
		ID:         "1",
		CustomerID: "123",
		LoadAmount: "$10000.00",
		Time:       time.Date(2021, 4, 2, 1, 0, 0, 0, time.UTC),
	})
	p.Load(Load{
		ID:         "2",
		CustomerID: "123",
		LoadAmount: "$10000.00",
		Time:       time.Date(2021, 4, 2, 1, 0, 0, 0, time.UTC),
	})

	//third load attempt should be rejected as is exceeds weekly limit
	proccesedLoad := p.Load(Load{
		ID:         "3",
		CustomerID: "123",
		LoadAmount: "$20000.00",
		Time:       time.Date(2021, 4, 2, 1, 0, 0, 0, time.UTC),
	})

	resp := Response{
		ID:         "3",
		CustomerID: "123",
		Accepted:   false,
	}
	assert.Equal(t, resp, *proccesedLoad)
}

func TestRejectLoadDueToDuplicateId(t *testing.T) {
	p := NewProcessor()

	// deposit load for customer 123 with ID of 1
	p.Load(Load{
		ID:         "1",
		CustomerID: "123",
		LoadAmount: "$100.00",
		Time:       time.Now(),
	})

	// attempt to deposit another load for customer 123 with ID of 1 should be skipped
	proccesedLoad := p.Load(Load{
		ID:         "1",
		CustomerID: "123",
		LoadAmount: "$100.00",
		Time:       time.Now(),
	})

	assert.Nil(t, proccesedLoad)

	// attempt to deposit another load with ID 1 for customer 321 should be approved
	proccesedLoad = p.Load(Load{
		ID:         "1",
		CustomerID: "321",
		LoadAmount: "$100.00",
		Time:       time.Now(),
	})

	resp := Response{
		ID:         "1",
		CustomerID: "321",
		Accepted:   true,
	}

	assert.Equal(t, resp, *proccesedLoad)
}
