package main

import (
	"fmt"
	"time"
)

//func mergeAccountMap(accounts map[string]Account) map[string]Account {
//
//}

func mergeAccounts(acc1, acc2 Account) Account {
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

	newAcc := Account{
		Address:     acc1.Address,
		TotalAmount: newTotal,
		StartTime:   startTime,
		EndTime:     endTime,
	}

	return newAcc
}
