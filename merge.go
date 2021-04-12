package main

import (
	"fmt"
	"sort"
	"time"
)

func MergeAccounts(accounts []Account) (map[string]Account, error) {
	accMap := make(map[string]Account)

	for _, acc := range accounts {
		if len(acc.Distributions) == 0 {
			return nil, fmt.Errorf("account must have atleast one distribution: %v", acc)
		}
		addrStr := acc.Address.String()
		existing, ok := accMap[addrStr]
		var newAcc Account
		if ok {
			var err error
			newAcc, err = mergeTwoAccounts(acc, existing)
			if err != nil {
				return nil, err
			}

			err = newAcc.Validate()
			if err != nil {
				return nil, fmt.Errorf("error merging two accounts: %w", err)
			}
		} else {
			newAcc = acc
		}
		accMap[addrStr] = newAcc
	}

	return accMap, nil
}

func mergeTwoAccounts(acc1, acc2 Account) (Account, error) {
	if !acc1.Address.Equals(acc2.Address) {
		return Account{}, fmt.Errorf("%s != %s", acc1.Address, acc2.Address)
	}

	newTotal, err := acc1.TotalRegen.Add(acc2.TotalRegen)
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
		if ok {
			amount, err = amount.Add(existing.Regen)
			if err != nil {
				return Account{}, err
			}
		}

		distMap[t] = Distribution{
			Time:  t,
			Regen: amount,
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
