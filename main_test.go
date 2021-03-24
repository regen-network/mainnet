package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/types"
)

const testTemplate = `
{
  "genesis_time": "2021-03-26T16:00:00Z",
  "chain_id": "regen-prelaunch-1",
  "initial_height": "1",
  "app_hash": "",
  "app_state": {}
}
`

const testAccounts = `
regen10rk2v8pxjnldtxuy9ds0s5na9qjcmh5ymplz87,100000,MAINNET,1
regen1lusdjktpk3f2v33cda5uwnya5qcyv04cwvnkwz,200000.301,2020-06-19,24
regen1lusdjktpk3f2v33cda5uwnya5qcyv04cwvnkwz,300000.0,MAINNET+1YEAR,24
`

func ExampleProcess() {
	err := processStrings(testTemplate, testAccounts)
	if err != nil {
		panic(err)
	}
	// Output:
}

func processStrings(tmplJson string, accountsCsv string) error {
	doc, err := types.GenesisDocFromJSON([]byte(tmplJson))
	if err != nil {
		return err
	}

	auditOut := new(bytes.Buffer)
	err = Process(doc, strings.NewReader(accountsCsv), auditOut, false)
	if err != nil {
		return err
	}

	bz, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(bz)
	fmt.Println("---")
	fmt.Println(auditOut.String())
	return nil
}
