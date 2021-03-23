package main

import (
	"fmt"
	"math"
	"time"

	"github.com/cockroachdb/apd/v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	vesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
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

var Dec128Context = apd.Context{
	Precision:   34,
	MaxExponent: apd.MaxExponent,
	MinExponent: apd.MinExponent,
	Traps:       apd.DefaultTraps,
}

func RecordToAccount(rec Record, genesisTime time.Time) (Account, error) {
	amount := rec.TotalAmount
	startTime := rec.StartTime
	if startTime.IsZero() {
		startTime = genesisTime
	}
	numDist := rec.NumMonthlyDistributions

	if numDist <= 1 && (startTime.Equal(genesisTime) || startTime.Before(genesisTime)) {
		return Account{
			Address:    rec.Address,
			TotalRegen: amount,
		}, nil
	}

	distAmount, dust, err := distAmountAndDust(amount, numDist)
	if err != nil {
		return Account{}, err
	}

	// put dust in first distribution
	var firstDist apd.Decimal
	_, err = apd.BaseContext.Add(&firstDist, &firstDist, &dust)
	if err != nil {
		return Account{}, err
	}

	endTime := rec.StartTime

	// prune distributions before genesis
	if endTime.Before(genesisTime) {
		startTime = genesisTime
		for endTime.Before(genesisTime) {
			endTime = endTime.Add(OneMonth)
			numDist -= 1
			_, err = apd.BaseContext.Add(&firstDist, &firstDist, &distAmount)
		}
	} else {
		_, err = apd.BaseContext.Add(&firstDist, &firstDist, &distAmount)
		if err != nil {
			return Account{}, err
		}
	}

	var distributions []Distribution

	// first distribution at start time
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

	return Account{
		Address:       rec.Address,
		Distributions: distributions,
		TotalRegen:    amount,
	}, nil
}

func distAmountAndDust(amount apd.Decimal, numDist int) (distAmount apd.Decimal, dust apd.Decimal, err error) {
	if numDist <= 1 {
		return amount, dust, nil
	}

	numDistDec := apd.New(int64(numDist), 1)

	// each distribution is an integral amount of regen
	cond, err := Dec128Context.QuoInteger(&distAmount, &amount, numDistDec)
	if err != nil {
		return distAmount, dust, err
	}

	fmt.Println(cond)

	// the remainder is dust
	cond, err = Dec128Context.Rem(&dust, &amount, numDistDec)
	if err != nil {
		return distAmount, dust, err
	}

	fmt.Println(cond)

	return distAmount, dust, nil
}

const (
	URegenDenom = "uregen"
)

var tenE6 *apd.Decimal

func init() {
	var err error
	tenE6, _, err = apd.NewFromString("1000000")
	if err != nil {
		panic(err)
	}
}

func ToCosmosAccount(acc Account, genesisTime time.Time) (auth.AccountI, *bank.Balance, error) {
	totalCoins, err := RegenToCoins(&acc.TotalRegen)
	if err != nil {
		return nil, nil, err
	}

	addrStr := acc.Address.String()
	balance := &bank.Balance{
		Address: addrStr,
		Coins:   totalCoins,
	}

	if len(acc.Distributions) == 0 {
		return &auth.BaseAccount{Address: addrStr}, balance, nil
	}

	startTime := acc.Distributions[0].Time

	// if we have one distribution and it happens before or at genesis return a basic BaseAccount
	if len(acc.Distributions) == 1 && (startTime.Before(genesisTime) || startTime.Equal(genesisTime)) {
		if acc.TotalRegen.Cmp(&acc.Distributions[0].Regen) != 0 {
			return nil, nil, fmt.Errorf("incorrect balance, expected %s, got %s", acc.TotalRegen.String(), acc.Distributions[0].Regen.String())
		}

		return &auth.BaseAccount{Address: addrStr}, balance, nil
	} else {
		periodStart := startTime

		var periods []vesting.Period
		var calcTotal apd.Decimal

		for _, dist := range acc.Distributions {
			_, err = apd.BaseContext.Add(&calcTotal, &calcTotal, &dist.Regen)
			if err != nil {
				return nil, nil, err
			}

			coins, err := RegenToCoins(&dist.Regen)
			if err != nil {
				return nil, nil, err
			}

			length := dist.Time.Sub(periodStart)
			periodStart = dist.Time
			seconds := int64(math.Floor(length.Seconds()))
			periods = append(periods, vesting.Period{
				Length: seconds,
				Amount: coins,
			})
		}

		if acc.TotalRegen.Cmp(&calcTotal) != 0 {
			return nil, nil, fmt.Errorf("incorrect balance, expected %s, got %s", acc.TotalRegen.String(), calcTotal.String())
		}

		return &vesting.PeriodicVestingAccount{
			BaseVestingAccount: &vesting.BaseVestingAccount{
				BaseAccount: &auth.BaseAccount{
					Address: addrStr,
				},
				OriginalVesting: totalCoins,
				EndTime:         periodStart.Unix(),
			},
			StartTime:      startTime.Unix(),
			VestingPeriods: periods,
		}, balance, nil
	}
}

func RegenToCoins(regenAmount *apd.Decimal) (sdk.Coins, error) {
	var uregen apd.Decimal
	_, err := apd.BaseContext.Mul(&uregen, regenAmount, tenE6)
	if err != nil {
		return nil, err
	}

	uregenInt64, err := uregen.Int64()
	if err != nil {
		return nil, err
	}

	return sdk.NewCoins(sdk.NewCoin(URegenDenom, sdk.NewInt(uregenInt64))), nil
}
