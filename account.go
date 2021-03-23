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
