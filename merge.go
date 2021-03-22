package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/cockroachdb/apd/v2"
)

func MergeAccounts(accounts []Account) ([]Account, error) {
	accMap := make(map[string]Account)

	for _, acc := range accounts {
		addrStr := acc.Address.String()
		existing, ok := accMap[addrStr]
		var newAcc Account
		if ok {
			var err error
			newAcc, err = mergeTwoAccounts(acc, existing)
			if err != nil {
				return nil, err
			}
		}
		accMap[addrStr] = newAcc
	}

	res := make([]Account, 0, len(accMap))

	for _, acc := range accMap {
		res = append(res, acc)
	}

	return res, nil
}

func mergeTwoAccounts(acc1, acc2 Account) (Account, error) {
	if !acc1.Address.Equals(acc2.Address) {
		return Account{}, fmt.Errorf("%s != %s", acc1.Address, acc2.Address)
	}

	var newTotal apd.Decimal
	_, err := apd.BaseContext.Add(&newTotal, &acc1.TotalRegen, &acc2.TotalRegen)
	if err != nil {
		return Account{}, err
	}

	distMap := make(map[time.Time]Distribution)
	for _, dist := range acc1.Distributions {
		distMap[dist.Time] = dist
	}

	for _, dist := range acc2.Distributions {
		t := dist.Time
		amount := dist.Regen
		existing, ok := distMap[t]
		var newAmount apd.Decimal
		if ok {
			_, err := apd.BaseContext.Add(&newAmount, &existing.Regen, &amount)
			if err != nil {
				return Account{}, err
			}
		}

		distMap[t] = Distribution{
			Time:  t,
			Regen: newAmount,
		}
	}

	// sort times
	var times []time.Time
	for t := range distMap {
		times = append(times, t)
	}
	sort.Slice(times, func(i, j int) bool {
		return times[i].Before(times[j])
	})

	// put distributions in sorted order
	var distributions []Distribution
	for _, t := range times {
		distributions = append(distributions, distMap[t])
	}

	return Account{
		Address:       acc1.Address,
		TotalRegen:    newTotal,
		Distributions: distributions,
	}, nil
}
