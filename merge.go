package main

import (
	"fmt"
	"sort"
	"time"
)

func mergeAccounts(accounts []Account) map[string]Account {
	accMap := make(map[string]Account)

	for _, acc := range accounts {
		addrStr := acc.Address.String()
		existing, ok := accMap[addrStr]
		if ok {
			acc = mergeTwoAccounts(acc, existing)
		}
		accMap[addrStr] = acc
	}

	return accMap
}

func mergeTwoAccounts(acc1, acc2 Account) Account {
	if !acc1.Address.Equals(acc2.Address) {
		panic(fmt.Errorf("%s != %s", acc1.Address, acc2.Address))
	}

	newTotal := acc1.TotalAmount.Add(acc2.TotalAmount)

	var startTime, endTime time.Time

	if acc1.StartTime.Before(acc2.StartTime) {
		startTime = acc1.StartTime
	} else {
		startTime = acc2.StartTime
	}

	if acc1.EndTime.After(acc2.EndTime) {
		endTime = acc1.EndTime
	} else {
		endTime = acc2.EndTime
	}

	distMap := make(map[time.Time]Distribution)
	for _, dist := range acc1.Distributions {
		distMap[dist.Time] = dist
	}

	for _, dist := range acc2.Distributions {
		t := dist.Time
		amount := dist.Amount
		existing, ok := distMap[t]
		if ok {
			amount = existing.Amount.Add(amount)
		}

		distMap[t] = Distribution{
			Time:   t,
			Amount: amount,
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
		TotalAmount:   newTotal,
		StartTime:     startTime,
		EndTime:       endTime,
		Distributions: distributions,
	}
}
