package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/cockroachdb/apd/v2"
)

func RequireDecEqual(t *testing.T, expected, actual *apd.Decimal, msgAndArgs ...interface{}) {
	if expected.Cmp(actual) != 0 {
		assert.Fail(t, fmt.Sprintf("%s != %s", expected.Text('f'), actual.Text('f')), msgAndArgs...)
		t.FailNow()
	}
}

func RequireAccountEqual(t *testing.T, expected, actual Account) {
	RequireDecEqual(t, &expected.TotalRegen, &actual.TotalRegen, "TotalRegen")
	require.Equal(t, expected.Address, actual.Address, "Address")
	require.Equal(t, len(expected.Distributions), len(actual.Distributions), "len(Distributions)")
	for i := 0; i < len(expected.Distributions); i++ {
		require.Equal(t, expected.Distributions[i].Time, actual.Distributions[i].Time, "Distribution", i)
		RequireDecEqual(t, &expected.Distributions[i].Regen, &actual.Distributions[i].Regen, "Distribution", i)
	}
}
