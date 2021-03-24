package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cockroachdb/apd/v2"
)

func RequireDecEqual(t *testing.T, expected, actual *apd.Decimal, msgAndArgs ...interface{}) {
	if expected.Cmp(actual) != 0 {
		assert.Fail(t, fmt.Sprintf("%s != %s", expected.Text('f'), actual.Text('f')), msgAndArgs...)
		t.FailNow()
	}
}

func RequireAccountEqual(t *testing.T, expected, actual Account, msgAndArgs ...interface{}) {

