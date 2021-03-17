package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Record struct {
	Address     sdk.AccAddress
	TotalAmount sdk.Int
	// empty for mainnet launch
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
	addr, err := sdk.GetFromBech32(line[0], "regen")
	if err != nil {
		return Record{}, err
	}

	amount, ok := sdk.NewIntFromString(line[1])
	if !ok {
		return Record{}, fmt.Errorf("expected integer got, %s", line[1])
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
		TotalAmount:             amount,
		StartTime:               startTime,
		NumMonthlyDistributions: numDist,
	}, nil

}
