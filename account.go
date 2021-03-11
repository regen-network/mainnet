package main

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	vesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"math"
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

const (
	uregen = "uregen"
	tenE6  = 1000000
)

func toCosmosAccount(acc Account) (auth.AccountI, *bank.Balance) {
	addrStr := acc.Address.String()
	totalCoins := sdk.NewCoins(sdk.NewCoin(uregen, acc.TotalAmount.MulRaw(tenE6)))
	balance := &bank.Balance{
		Address: addrStr,
		Coins:   totalCoins,
	}
	if len(acc.Distributions) == 0 {
		return &auth.BaseAccount{Address: addrStr}, balance
	} else {
		var periods []vesting.Period

		periodStart := acc.StartTime

		for _, dist := range acc.Distributions {
			coins := sdk.NewCoins(sdk.NewCoin(uregen, dist.Amount.MulRaw(tenE6)))
			length := dist.Time.Sub(periodStart)
			seconds := int64(math.Floor(length.Seconds()))
			periods = append(periods, vesting.Period{
				Length: seconds,
				Amount: coins,
			})
		}

		return &vesting.PeriodicVestingAccount{
			BaseVestingAccount: &vesting.BaseVestingAccount{
				BaseAccount: &auth.BaseAccount{
					Address: addrStr,
				},
				OriginalVesting: totalCoins,
				EndTime:         acc.EndTime.Unix(),
			},
			StartTime:      acc.StartTime.Unix(),
			VestingPeriods: periods,
		}, balance
	}
}
