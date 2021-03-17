package main

import (
	"strings"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
)

func TestParseAccountsCsv(t *testing.T) {
	genesisTime, err := time.Parse(time.RFC3339, "2021-04-08T16:00:00Z")
	require.NoError(t, err)

	records, err := ParseAccountsCsv(strings.NewReader(`
regen10rk2v8pxjnldtxuy9ds0s5na9qjcmh5ymplz87,100000,MAINNET,1
regen1lusdjktpk3f2v33cda5uwnya5qcyv04cwvnkwz,200000,2020-06-19,24
regen1lusdjktpk3f2v33cda5uwnya5qcyv04cwvnkwz,300000,MAINNET+1YEAR,24
`), genesisTime)
	require.NoError(t, err)

	addr0, err := sdk.AccAddressFromBech32("regen10rk2v8pxjnldtxuy9ds0s5na9qjcmh5ymplz87")
	require.NoError(t, err)

	addr1, err := sdk.AccAddressFromBech32("regen1lusdjktpk3f2v33cda5uwnya5qcyv04cwvnkwz")
	require.NoError(t, err)

	addr2, err := sdk.AccAddressFromBech32("regen1lusdjktpk3f2v33cda5uwnya5qcyv04cwvnkwz")
	require.NoError(t, err)

	t0, err := time.Parse(time.RFC3339, "2020-06-19T00:00:00Z")
	require.NoError(t, err)

	require.Equal(t, []Record{
		{
			Address:                 addr0,
			TotalAmount:             sdk.NewInt(100000),
			StartTime:               genesisTime,
			NumMonthlyDistributions: 1,
		},
		{
			Address:                 addr1,
			TotalAmount:             sdk.NewInt(200000),
			StartTime:               t0,
			NumMonthlyDistributions: 24,
		},
		{
			Address:                 addr2,
			TotalAmount:             sdk.NewInt(300000),
			StartTime:               genesisTime.Add(OneYear),
			NumMonthlyDistributions: 24,
		},
	}, records)
}
