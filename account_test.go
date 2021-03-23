package main

import (
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
	thirty, _, err := apd.NewFromString("30")
	require.NoError(t, err)

	time0, err := time.Parse(time.RFC3339, "2021-05-21T00:00:00Z")
	require.NoError(t, err)

	time1, err := time.Parse(time.RFC3339, "2021-06-21T00:00:00Z")
	require.NoError(t, err)

	time2, err := time.Parse(time.RFC3339, "2021-07-21T00:00:00Z")
	require.NoError(t, err)

	time3, err := time.Parse(time.RFC3339, "2021-08-21T00:00:00Z")
	require.NoError(t, err)

	tests := []struct {
		name    string
		acc     Account
		want    auth.AccountI
		want1   *bank.Balance
		wantErr bool
	}{
		{
			"no distribution",
			Account{
				Address:    addr0,
				TotalRegen: *ten,
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
			"one distribution at mainnet",
			Account{
				Address:    addr0,
				TotalRegen: *ten,
				Distributions: []Distribution{
					{
						Time:  time0,
						Regen: *ten,
					},
				},
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
			"one distribution at mainnet bad balance",
			Account{
				Address:    addr0,
				TotalRegen: *ten,
				Distributions: []Distribution{
					{
						Time:  time0,
						Regen: *five,
					},
				},
			},
			nil,
			nil,
			true,
		},
		{
			"one distribution after mainnet",
			Account{
				Address:    addr1,
				TotalRegen: *five,
				Distributions: []Distribution{
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
					OriginalVesting: sdk.NewCoins(sdk.NewInt64Coin(URegenDenom, 5000000)),
					EndTime:         time1.Unix(),
				},
				StartTime: time1.Unix(),
				VestingPeriods: []vesting.Period{
					{
						Length: 0,
						Amount: sdk.NewCoins(sdk.NewInt64Coin(URegenDenom, 5000000)),
					},
				},
			},
			&bank.Balance{
				Address: addr1.String(),
				Coins:   sdk.NewCoins(sdk.NewInt64Coin(URegenDenom, 5000000)),
			},
			false,
		},
		{
			"two distributions",
			Account{
				Address:    addr1,
				TotalRegen: *ten,
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
						Length: int64(time1.Sub(time0).Seconds()),
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
		{
			"four distributions",
			Account{
				Address:    addr1,
				TotalRegen: *thirty,
				Distributions: []Distribution{
					{
						Time:  time0,
						Regen: *five,
					},
					{
						Time:  time1,
						Regen: *ten,
					},
					{
						Time:  time2,
						Regen: *five,
					},
					{
						Time:  time3,
						Regen: *ten,
					},
				},
			},
			&vesting.PeriodicVestingAccount{
				BaseVestingAccount: &vesting.BaseVestingAccount{
					BaseAccount: &auth.BaseAccount{
						Address: addr1.String(),
					},
					OriginalVesting: sdk.NewCoins(sdk.NewInt64Coin(URegenDenom, 30000000)),
					EndTime:         time3.Unix(),
				},
				StartTime: time0.Unix(),
				VestingPeriods: []vesting.Period{
					{
						Length: 0,
						Amount: sdk.NewCoins(sdk.NewInt64Coin(URegenDenom, 5000000)),
					},
					{
						Length: int64(time1.Sub(time0).Seconds()),
						Amount: sdk.NewCoins(sdk.NewInt64Coin(URegenDenom, 10000000)),
					},
					{
						Length: int64(time2.Sub(time1).Seconds()),
						Amount: sdk.NewCoins(sdk.NewInt64Coin(URegenDenom, 5000000)),
					},
					{
						Length: int64(time3.Sub(time2).Seconds()),
						Amount: sdk.NewCoins(sdk.NewInt64Coin(URegenDenom, 10000000)),
					},
				},
			},
			&bank.Balance{
				Address: addr1.String(),
				Coins:   sdk.NewCoins(sdk.NewInt64Coin(URegenDenom, 30000000)),
			},
			false,
		},
		{
			"two distributions bad balance",
			Account{
				Address:    addr0,
				TotalRegen: *ten,
				Distributions: []Distribution{
					{
						Time:  time0,
						Regen: *five,
					},
					{
						Time:  time1,
						Regen: *ten,
					},
				},
			},
			nil,
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ToCosmosAccount(tt.acc, time0)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.want1, got1)
		})
	}
}

func TestRecordToAccount(t *testing.T) {
	addr0 := sdk.AccAddress("abcdefg012345")

	five, _, err := apd.NewFromString("5")
	require.NoError(t, err)
	ten, _, err := apd.NewFromString("10")
	require.NoError(t, err)
	//thirty, _, err := apd.NewFromString("30")
	//require.NoError(t, err)

	genesisTime, err := time.Parse(time.RFC3339, "2021-04-08T00:00:00Z")
	require.NoError(t, err)

	genesisPlus1Month := genesisTime.Add(OneMonth)

	//start0, err := time.Parse(time.RFC3339, "2021-01-05:00:00Z")
	//require.NoError(t, err)
	//
	//start1, err := time.Parse(time.RFC3339, "2021-02-08:00:00Z")
	//require.NoError(t, err)
	//
	//start2, err := time.Parse(time.RFC3339, "2021-09-04:00:00Z")
	//require.NoError(t, err)

	//t3, err := time.Parse(time.RFC3339, "2021-01-08:00:00Z")
	//require.NoError(t, err)

	tests := []struct {
		name    string
		record  Record
		want    Account
		wantErr bool
	}{
		{
			"no dists, no start time",
			Record{
				Address:     addr0,
				TotalAmount: *five,
			},
			Account{
				Address:    addr0,
				TotalRegen: *five,
			},
			false,
		},
		{
			"no dists, genesis start time",
			Record{
				Address:     addr0,
				TotalAmount: *five,
				StartTime:   genesisTime,
			},
			Account{
				Address:    addr0,
				TotalRegen: *five,
			},
			false,
		},
		{
			"one dist at genesis",
			Record{
				Address:                 addr0,
				TotalAmount:             *five,
				StartTime:               genesisTime,
				NumMonthlyDistributions: 1,
			},
			Account{
				Address:    addr0,
				TotalRegen: *five,
			},
			false,
		},
		{
			"two dists at genesis",
			Record{
				Address:                 addr0,
				TotalAmount:             *ten,
				StartTime:               genesisTime,
				NumMonthlyDistributions: 1,
			},
			Account{
				Address:    addr0,
				TotalRegen: *ten,
				Distributions: []Distribution{
					{
						Time:  genesisTime,
						Regen: *five,
					},
					{
						Time:  genesisPlus1Month,
						Regen: *five,
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RecordToAccount(tt.record, genesisTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordToAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.NotNil(t, got)
			require.True(t, tt.want.TotalRegen.Cmp(&got.TotalRegen) == 0)
			require.Equal(t, tt.want.Address, got.Address)
			require.Equal(t, len(tt.want.Distributions), len(got.Distributions))
			for i := 0; i < len(tt.want.Distributions); i++ {
				require.Equal(t, tt.want.Distributions[i].Time, got.Distributions[i].Time)
				require.True(t, tt.want.Distributions[i].Regen.Cmp(&got.Distributions[i].Regen) == 0)
			}
		})
	}
}
