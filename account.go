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
		var dust apd.Decimal

		// each distribution is an integral amount of regen
		_, err := apd.BaseContext.QuoInteger(&distAmount, &amount, numDistDec)
		if err != nil {
			return Account{}, err
		}

		// the remainder is dust
		_, err = apd.BaseContext.Rem(&dust, &amount, numDistDec)
		if err != nil {
			return Account{}, err
		}

		// first distribution at start time
		// put dust in first distribution
		var firstDist apd.Decimal
		_, err = apd.BaseContext.Add(&firstDist, &distAmount, &dust)

		distributions = append(distributions, Distribution{
			Time:  startTime,
			Regen: firstDist,
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
