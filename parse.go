package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/regen-network/regen-ledger/app"

	"github.com/cockroachdb/apd/v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Record struct {
	Address                 sdk.AccAddress
	TotalAmount             apd.Decimal
	StartTime               time.Time
	NumMonthlyDistributions int
}

const (
	SecondsPerYear  = 31556952
	SecondsPerMonth = 2629746
)

var (
	OneYear, OneMonth time.Duration
)

func init() {
	var err error

	OneYear, err = time.ParseDuration(fmt.Sprintf("%ds", SecondsPerYear))
	if err != nil {
		panic(err)
	}

	OneMonth, err = time.ParseDuration(fmt.Sprintf("%ds", SecondsPerMonth))
	if err != nil {
		panic(err)
	}
}

func ParseAccountsCsv(rdr io.Reader, genesisTime time.Time) ([]Record, error) {
	csvRdr := csv.NewReader(rdr)
	lines, err := csvRdr.ReadAll()
	if err != nil {
		return nil, err
	}

	records := make([]Record, 0, len(lines))
	for i, line := range lines {
		record, err := parseLine(line, genesisTime)
		if err != nil {
			fmt.Printf("Error on line %d: %v", i, err)
			continue
		}

		records = append(records, record)
	}

	return records, nil
}

func parseLine(line []string, genesisTime time.Time) (Record, error) {
	addr, err := sdk.GetFromBech32(line[0], app.Bech32PrefixAccAddr)
	if err != nil {
		return Record{}, err
	}

	amount, _, err := apd.NewFromString(line[1])
	if err != nil {
		return Record{}, err
	}

	var startTime time.Time
	startTimeStr := strings.TrimSpace(line[2])
	switch startTimeStr {
	case "MAINNET":
		startTime = genesisTime
	case "MAINNET+1YEAR":
		startTime = genesisTime.Add(OneYear)
	default:
		startTime, err = time.Parse("2006-01-02", line[2])
		if err != nil {
			return Record{}, err
		}
	}

	numDist, err := strconv.Atoi(line[3])
	if err != nil {
		return Record{}, err
	}

	if numDist < 1 {
		return Record{}, fmt.Errorf("expected a positive integer, got %d", numDist)
	}

	return Record{
		Address:                 addr,
		TotalAmount:             *amount,
		StartTime:               startTime,
		NumMonthlyDistributions: numDist,
	}, nil

}

func (r Record) Equal(o Record) bool {
	return r.StartTime.Equal(o.StartTime) &&
		r.TotalAmount.Cmp(&o.TotalAmount) == 0 &&
		r.Address.Equals(o.Address) &&
		r.NumMonthlyDistributions == o.NumMonthlyDistributions
}
