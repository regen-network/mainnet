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
				return
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.want1, got1)
			validateVestingAccount(t, got)
		})
	}
}

func validateVestingAccount(t *testing.T, acc auth.AccountI) {
	vacc, ok := acc.(*vesting.PeriodicVestingAccount)
	if !ok {
		return
	}

	orig := vacc.OriginalVesting
	time := vacc.StartTime
	var total sdk.Coins
	for _, period := range vacc.VestingPeriods {
		total = total.Add(period.Amount...)
		time += period.Length
	}
	require.Equal(t, orig, total)
	require.Equal(t, vacc.EndTime, time)
}

func TestRecordToAccount(t *testing.T) {
	addr0 := sdk.AccAddress("abcdefg012345")

	five, _, err := apd.NewFromString("5")
	require.NoError(t, err)
	five000001, _, err := apd.NewFromString("5.000001")
	require.NoError(t, err)
	ten000001, _, err := apd.NewFromString("10.000001")
	require.NoError(t, err)
	seven500001, _, err := apd.NewFromString("7.500001")
	require.NoError(t, err)
	two5, _, err := apd.NewFromString("2.5")
	require.NoError(t, err)

	genesisTime, err := time.Parse(time.RFC3339, "2021-04-08T00:00:00Z")
	require.NoError(t, err)

	oneMonthAfterGenesis := genesisTime.Add(OneMonth)

	start0, err := time.Parse(time.RFC3339, "2021-01-05T00:00:00Z")
	require.NoError(t, err)

	twoMonthsBeforeGenesis := genesisTime.Add(-OneMonth * 2)

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
				Distributions: []Distribution{
					{
						Time:  genesisTime,
						Regen: *five,
					},
				},
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
				Distributions: []Distribution{
					{
						Time:  genesisTime,
						Regen: *five,
					},
				},
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
				Distributions: []Distribution{
					{
						Time:  genesisTime,
						Regen: *five,
					},
				},
			},
			false,
		},
		{
			"two dists from genesis",
			Record{
				Address:                 addr0,
				TotalAmount:             *ten000001,
				StartTime:               genesisTime,
				NumMonthlyDistributions: 2,
			},
			Account{
				Address:    addr0,
				TotalRegen: *ten000001,
				Distributions: []Distribution{
					{
						Time:  genesisTime,
						Regen: *five000001,
					},
					{
						Time:  oneMonthAfterGenesis,
						Regen: *five,
					},
				},
			},
			false,
		},
		{
			"two dists all from before genesis",
			Record{
				Address:                 addr0,
				TotalAmount:             *ten000001,
				StartTime:               start0,
				NumMonthlyDistributions: 2,
			},
			Account{
				Address:    addr0,
				TotalRegen: *ten000001,
				Distributions: []Distribution{
					{
						Time:  genesisTime,
						Regen: *ten000001,
					},
				},
			},
			false,
		},
		{
			"four dists starting before genesis",
			Record{
				Address:                 addr0,
				TotalAmount:             *ten000001,
				StartTime:               twoMonthsBeforeGenesis,
				NumMonthlyDistributions: 4,
			},
			Account{
				Address:    addr0,
				TotalRegen: *ten000001,
				Distributions: []Distribution{
					{
						Time:  genesisTime,
						Regen: *seven500001,
					},
					{
						Time:  oneMonthAfterGenesis,
						Regen: *two5,
					},
				},
			},
			false,
		},
		{
			"one dist after genesis",
			Record{
				Address:                 addr0,
				TotalAmount:             *ten000001,
				StartTime:               oneMonthAfterGenesis,
				NumMonthlyDistributions: 1,
			},
			Account{
				Address:    addr0,
				TotalRegen: *ten000001,
				Distributions: []Distribution{
					{
						Time:  oneMonthAfterGenesis,
						Regen: *ten000001,
					},
				},
			},
			false,
		},
		{
			"two dists after genesis",
			Record{
				Address:                 addr0,
				TotalAmount:             *ten000001,
				StartTime:               oneMonthAfterGenesis,
				NumMonthlyDistributions: 2,
			},
			Account{
				Address:    addr0,
				TotalRegen: *ten000001,
				Distributions: []Distribution{
					{
						Time:  oneMonthAfterGenesis,
						Regen: *five000001,
					},
					{
						Time:  oneMonthAfterGenesis.Add(OneMonth),
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
			require.Truef(t, tt.want.TotalRegen.Cmp(&got.TotalRegen) == 0,
				"%s != %s", tt.want.TotalRegen.String(), got.TotalRegen.String())
			require.Equal(t, tt.want.Address, got.Address)
			require.Equal(t, len(tt.want.Distributions), len(got.Distributions))
			for i := 0; i < len(tt.want.Distributions); i++ {
				require.Equal(t, tt.want.Distributions[i].Time, got.Distributions[i].Time)
				require.Truef(t, tt.want.Distributions[i].Regen.Cmp(&got.Distributions[i].Regen) == 0,
					"Distribution %d: %s != %s", i, tt.want.Distributions[i].Regen.String(), got.Distributions[i].Regen.String())
			}
			require.NoError(t, got.Validate())
		})
	}
}

func Test_distAmountAndDust(t *testing.T) {
	five01, _, err := apd.NewFromString("5.01")
	require.NoError(t, err)
	ten02, _, err := apd.NewFromString("10.02")
	require.NoError(t, err)
	ten0000002, _, err := apd.NewFromString("10.020001")
	require.NoError(t, err)
	point000001, _, err := apd.NewFromString("0.000001")
	require.NoError(t, err)
	ten, _, err := apd.NewFromString("10")
	require.NoError(t, err)
	three333333, _, err := apd.NewFromString("3.333333")
	require.NoError(t, err)

	type args struct {
		amount  apd.Decimal
		numDist int
	}
	tests := []struct {
		name           string
		args           args
		wantDistAmount apd.Decimal
		wantDust       apd.Decimal
		wantErr        bool
	}{
		{
			name: "0 dist",
			args: args{
				amount:  *ten02,
				numDist: 0,
			},
			wantDistAmount: *ten02,
			wantDust:       apd.Decimal{},
			wantErr:        false,
		},
		{
			name: "1 dist",
			args: args{
				amount:  *ten02,
				numDist: 1,
			},
			wantDistAmount: *ten02,
			wantDust:       apd.Decimal{},
			wantErr:        false,
		},
		{
			name: "2 dist, no dust",
			args: args{
				amount:  *ten02,
				numDist: 2,
			},
			wantDistAmount: *five01,
			wantDust:       apd.Decimal{},
			wantErr:        false,
		},
		{
			name: "2 dist, dust",
			args: args{
				amount:  *ten0000002,
				numDist: 2,
			},
			wantDistAmount: *five01,
			wantDust:       *point000001,
			wantErr:        false,
		},
		{
			name: "3 dists, dust",
			args: args{
				amount:  *ten,
				numDist: 3,
			},
			wantDistAmount: *three333333,
			wantDust:       *point000001,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDistAmount, gotDust, err := distAmountAndDust(tt.args.amount, tt.args.numDist)
			if (err != nil) != tt.wantErr {
				t.Errorf("distAmountAndDust() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Truef(t, tt.wantDistAmount.Cmp(&gotDistAmount) == 0,
				"distAmount: %s != %s", tt.wantDistAmount.String(), gotDistAmount.String())
			require.Truef(t, tt.wantDust.Cmp(&gotDust) == 0,
				"dust: %s != %s", tt.wantDust.String(), gotDust.String())
		})
	}
}
