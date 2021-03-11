package main

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

type Account struct {
	Address       sdk.AccAddress
	TotalAmount   sdk.Int
	StartTime     time.Time
	Distributions []Distribution
	EndTime       time.Time
}

type Distribution struct {
	Time   time.Time
	Amount sdk.Int
}

func recordToAccount(rec Record, defaultStartTime time.Time) Account {
	amount := rec.TotalAmount
	startTime := rec.StartTime
	if startTime.IsZero() {
		startTime = defaultStartTime
	}
	endTime := rec.StartTime
	numDist := rec.NumMonthlyDistributions

	var distributions []Distribution
	if numDist > 1 {
		distAmount := amount.QuoRaw(int64(numDist))

		// first distribution at start time
		distributions = append(distributions, Distribution{
			Time:   endTime,
			Amount: distAmount,
		})

		oneMonth, err := time.ParseDuration("2629746s")
		if err != nil {
			panic(err)
		}

		for i := 1; i < numDist; i++ {
			endTime = endTime.Add(oneMonth)
			distributions = append(distributions, Distribution{
				Time:   endTime,
				Amount: distAmount,
			})
		}
	}

	return Account{
		Address:       rec.Address,
		TotalAmount:   rec.TotalAmount,
		StartTime:     startTime,
		Distributions: distributions,
		EndTime:       endTime,
	}
}
