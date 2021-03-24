package main

import (
	"fmt"
	"io"
	"sort"
	"time"
)

func SortAccounts(accounts []Account) []Account {
	var addrs []string
	accMap := make(map[string]Account)
	// sort
	for _, acc := range accounts {
		addr := acc.Address.String()
		addrs = append(addrs, addr)
		accMap[addr] = acc
	}

	sort.Strings(addrs)

	res := make([]Account, 0, len(accounts))
	for _, addr := range addrs {
		res = append(res, accMap[addr])
	}

	return res
}

func PrintAccountAudit(accounts []Account, genesisTime time.Time, writer io.Writer) {
	for _, acc := range accounts {
		_, err := fmt.Fprintf(writer, "%s\t%s\t%d\n", acc.Address, acc.TotalRegen, len(acc.Distributions))
		if err != nil {
			panic(err)
		}
		for _, dist := range acc.Distributions {
			var timeStr string
			if dist.Time.Equal(genesisTime) {
				timeStr = "MAINNET"
			} else {
				timeStr = dist.Time.Format("2006-01-02 15:04:05")
			}
			_, err = fmt.Fprintf(writer, "\t%s\t%s\n", dist.Regen, timeStr)
			if err != nil {
				panic(err)
			}
		}
	}
}
