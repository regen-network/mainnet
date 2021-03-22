package main

import (
	"time"

	"github.com/cockroachdb/apd/v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Account is an internal representation of a genesis regen account
type Account struct {
	Address       sdk.AccAddress
	TotalRegen    apd.Decimal
	Distributions []Distribution
}

// Distribution is an internal representation of a genesis vesting distribution of regen
type Distribution struct {
	Time  time.Time
	Regen apd.Decimal
}

func RecordToAccount(rec Record, genesisTime time.Time) (Account, error) {
	amount := rec.TotalAmount
	startTime := rec.StartTime
	if startTime.IsZero() {
		startTime = genesisTime
	}
	endTime := rec.StartTime
	numDist := rec.NumMonthlyDistributions

	var distributions []Distribution
	if numDist > 1 {
		numDistDec := apd.New(int64(numDist), 1)
		var distAmount apd.Decimal

		_, err := apd.BaseContext.Quo(&distAmount, &amount, numDistDec)
		if err != nil {
			return Account{}, err
		}

		// first distribution at start time
		distributions = append(distributions, Distribution{
			Time:  endTime,
			Regen: amount,
		})

		for i := 1; i < numDist; i++ {
			endTime = endTime.Add(OneMonth)
			distributions = append(distributions, Distribution{
				Time:  endTime,
				Regen: distAmount,
			})
		}
	}

	return Account{
		Address:       rec.Address,
		Distributions: distributions,
		TotalRegen:    amount,
	}, nil
}
