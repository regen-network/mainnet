package main

import (
	"testing"

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
}
