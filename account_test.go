package main

import (
	"reflect"
	"testing"
	"time"

	vesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cockroachdb/apd/v2"
	"github.com/stretchr/testify/require"
)

func TestRegenToCoins(t *testing.T) {
	one, _, err := apd.NewFromString("1")
	require.NoError(t, err)

	coins, err := RegenToCoins(one)
	require.NoError(t, err)
	require.Equal(t, "1000000uregen", coins.String())
}

func TestToCosmosAccount(t *testing.T) {
	addr0 := sdk.AccAddress("abcdefg012345")
	addr1 := sdk.AccAddress("012345abcdefg")
	five, _, err := apd.NewFromString("5")
	require.NoError(t, err)
	ten, _, err := apd.NewFromString("10")
	require.NoError(t, err)

	time0, err := time.Parse(time.RFC3339, "2021-05-21T00:00:00Z")
	require.NoError(t, err)

	time1, err := time.Parse(time.RFC3339, "2021-06-21T00:00:00Z")
	require.NoError(t, err)

	tests := []struct {
		name    string
		acc     Account
		want    auth.AccountI
		want1   *bank.Balance
		wantErr bool
	}{
		{
			"base account",
			Account{
				Address:    addr0,
				TotalRegen: *ten,
				StartTime:  time0,
				EndTime:    time0,
			},
			&auth.BaseAccount{
				Address: addr0.String(),
			},
			&bank.Balance{
				Address: addr0.String(),
				Coins:   sdk.NewCoins(sdk.NewInt64Coin(URegenDenom, 10000000)),
			},
			false,
		},
		{
			"two distributions",
			Account{
				Address:    addr1,
				TotalRegen: *ten,
				StartTime:  time0,
				EndTime:    time1,
				Distributions: []Distribution{
					{
						Time:  time0,
						Regen: *five,
					},
					{
						Time:  time1,
						Regen: *five,
					},
				},
			},
			&vesting.PeriodicVestingAccount{
				BaseVestingAccount: &vesting.BaseVestingAccount{
					BaseAccount: &auth.BaseAccount{
						Address: addr1.String(),
					},
					OriginalVesting: sdk.NewCoins(sdk.NewInt64Coin(URegenDenom, 10000000)),
					EndTime:         time1.Unix(),
				},
				StartTime: time0.Unix(),
				VestingPeriods: []vesting.Period{
					{
						Length: 0,
						Amount: sdk.NewCoins(sdk.NewInt64Coin(URegenDenom, 5000000)),
					},
					{
						Length: 0,
						Amount: sdk.NewCoins(sdk.NewInt64Coin(URegenDenom, 5000000)),
					},
				},
			},
			&bank.Balance{
				Address: addr1.String(),
				Coins:   sdk.NewCoins(sdk.NewInt64Coin(URegenDenom, 10000000)),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ToCosmosAccount(tt.acc)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToCosmosAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToCosmosAccount() got AccountI = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ToCosmosAccount() got Balance = %v, want %v", got1, tt.want1)
			}
		})
	}
}
