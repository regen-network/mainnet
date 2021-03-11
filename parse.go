package main

import (
	"encoding/csv"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"io"
	"strconv"
	"strings"
	"time"
)

type Record struct {
	Address     sdk.AccAddress
	TotalAmount sdk.Int
	// empty for mainnet launch
	StartTime               time.Time
	NumMonthlyDistributions int
}

func parseAccountsCsv(rdr io.Reader) ([]Record, error) {
	csvRdr := csv.NewReader(rdr)
	lines, err := csvRdr.ReadAll()
	if err != nil {
		return nil, err
	}

	records := make([]Record, 0, len(lines))
	for i, line := range lines {
		record, err := parseLine(line)
		if err != nil {
			fmt.Printf("Error on line %d: %v", i, err)
			continue
		}

		records = append(records, record)
	}

	return records, nil
}

func parseLine(line []string) (Record, error) {
	addr, err := sdk.GetFromBech32(line[0], "regen")
	if err != nil {
		return Record{}, err
	}

	amount, ok := sdk.NewIntFromString(line[1])
	if !ok {
		return Record{}, fmt.Errorf("expected integer got, %s", line[1])
	}

	var startTime time.Time
	if len(strings.TrimSpace(line[2])) != 0 {
		startTime, err = time.Parse(time.RFC3339, line[2])
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
