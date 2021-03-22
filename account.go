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
