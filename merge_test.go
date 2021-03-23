package main

import (
	"github.com/cockroachdb/apd/v2"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

func TestMergeAccounts(t *testing.T) {
	addr0 := sdk.AccAddress("abcdefg012345")
	addr1 := sdk.AccAddress("012345abcdefg")
	five, _, err := apd.NewFromString("5")
	require.NoError(t, err)
	ten, _, err := apd.NewFromString("10")
	require.NoError(t, err)
	fifteen, _, err := apd.NewFromString("15")
	require.NoError(t, err)
	twenty, _, err := apd.NewFromString("20")
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
		name     string
		accounts []Account
		want     []Account
		wantErr  bool
	}{
		{
			"merging two accounts basic",
			[]Account{
				{
					Address:    addr0,
					TotalRegen: *ten,
					Distributions: []Distribution{
						{
							Time: time0,
							Regen: *ten,
						},
					},
				},
				{
					Address:    addr0,
					TotalRegen: *ten,
					Distributions: []Distribution{
						{
							Time: time0,
							Regen: *five,
						},
						{
							Time: time1,
							Regen: *five,
						},
					},
				},
			},
			[]Account{
				{
					Address:    addr0,
					TotalRegen: *twenty,
					Distributions: []Distribution{
						{
							Time: time0,
							Regen: *fifteen,
						},
						{
							Time: time1,
							Regen: *five,
						},
					},
				},
			},
			false,
		},
		{
			"merging two accounts, no common distr times",
			[]Account{
				{
					Address:    addr0,
					TotalRegen: *ten,
					Distributions: []Distribution{
						{
							Time: time0,
							Regen: *ten,
						},
					},
				},
				{
					Address:    addr0,
					TotalRegen: *ten,
					Distributions: []Distribution{
						{
							Time: time1,
							Regen: *ten,
						},
					},
				},
			},
			[]Account{
				{
					Address:    addr0,
					TotalRegen: *twenty,
					Distributions: []Distribution{
						{
							Time: time0,
							Regen: *ten,
						},
						{
							Time: time1,
							Regen: *ten,
						},
					},
				},
			},
			false,
		},
		{
			"merge accounts with all different addresses",
			[]Account{
				{
					Address:    addr0,
					TotalRegen: *twenty,
					Distributions: []Distribution{
						{
							Time: time0,
							Regen: *ten,
						},
						{
							Time: time3,
							Regen: *ten,
						},
					},
				},
				{
					Address:    addr1,
					TotalRegen: *thirty,
					Distributions: []Distribution{
						{
							Time: time0,
							Regen: *ten,
						},
						{
							Time: time1,
							Regen: *ten,
						},
						{
							Time: time2,
							Regen: *ten,
						},
					},
				},
			},
			[]Account{
				{
					Address:    addr0,
					TotalRegen: *twenty,
					Distributions: []Distribution{
						{
							Time: time0,
							Regen: *ten,
						},
						{
							Time: time3,
							Regen: *ten,
						},
					},
				},
				{
					Address:    addr1,
					TotalRegen: *thirty,
					Distributions: []Distribution{
						{
							Time: time0,
							Regen: *ten,
						},
						{
							Time: time1,
							Regen: *ten,
						},
						{
							Time: time2,
							Regen: *ten,
						},
					},
				},
			},
			false,
		},
		{
			"error on account with no distributions",
			[]Account{
				{
					Address:    addr0,
					TotalRegen: *twenty,
					Distributions: []Distribution{},
				},
				{
					Address:    addr1,
					TotalRegen: *thirty,
					Distributions: []Distribution{
						{
							Time: time0,
							Regen: *ten,
						},
						{
							Time: time1,
							Regen: *ten,
						},
						{
							Time: time2,
							Regen: *ten,
						},
					},
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MergeAccounts(tt.accounts)
			if (err != nil) != tt.wantErr {
				t.Errorf("MergeAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeAccounts() got = %v, want %v", got, tt.want)
			}
		})
	}
}
