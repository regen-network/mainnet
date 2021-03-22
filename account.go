package main

import (
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
	StartTime     time.Time
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

func ToCosmosAccount(acc Account) (auth.AccountI, *bank.Balance, error) {
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
	} else {
		var periods []vesting.Period

		periodStart := acc.StartTime

		for _, dist := range acc.Distributions {
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

		return &vesting.PeriodicVestingAccount{
			BaseVestingAccount: &vesting.BaseVestingAccount{
				BaseAccount: &auth.BaseAccount{
					Address: addrStr,
				},
				OriginalVesting: totalCoins,
				EndTime:         periodStart.Unix(),
			},
			StartTime:      acc.StartTime.Unix(),
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
