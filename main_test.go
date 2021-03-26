package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/types"
)

const testTemplate = `
{
  "genesis_time": "2021-03-26T16:00:00Z",
  "chain_id": "regen-prelaunch-1",
  "initial_height": "1",
  "app_hash": "",
  "app_state": {
    "auth":{},
    "bank":{}
  }
}
`

const testAccounts = `
regen10rk2v8pxjnldtxuy9ds0s5na9qjcmh5ymplz87,100000,MAINNET,1
regen1lusdjktpk3f2v33cda5uwnya5qcyv04cwvnkwz,200000.301,2020-06-19,2
regen1lusdjktpk3f2v33cda5uwnya5qcyv04cwvnkwz,300000.0,MAINNET+1YEAR,2
`

func processStrings(tmplJson string, accountsCsv string) (string, string, error) {
	doc, err := types.GenesisDocFromJSON([]byte(tmplJson))
	if err != nil {
		return "", "", err
	}

	auditOut := new(bytes.Buffer)
	err = Process(doc, strings.NewReader(accountsCsv), auditOut, false)
	if err != nil {
		return "", "", err
	}

	bz, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", "", err
	}

	return string(bz), auditOut.String(), nil
}

func TestProcess(t *testing.T) {
	json, auditOut, err := processStrings(testTemplate, testAccounts)
	require.NoError(t, err)
	require.Equal(t,
		`{
  "genesis_time": "2021-03-26T16:00:00Z",
  "chain_id": "regen-prelaunch-1",
  "initial_height": "1",
  "consensus_params": {
    "block": {
      "max_bytes": "22020096",
      "max_gas": "-1",
      "time_iota_ms": "1000"
    },
    "evidence": {
      "max_age_num_blocks": "100000",
      "max_age_duration": "172800000000000",
      "max_bytes": "1048576"
    },
    "validator": {
      "pub_key_types": [
        "ed25519"
      ]
    },
    "version": {}
  },
  "app_hash": "",
  "app_state": {
    "auth": {
      "params": {
        "max_memo_characters": "0",
        "tx_sig_limit": "0",
        "tx_size_cost_per_byte": "0",
        "sig_verify_cost_ed25519": "0",
        "sig_verify_cost_secp256k1": "0"
      },
      "accounts": [
        {
          "@type": "/cosmos.auth.v1beta1.BaseAccount",
          "address": "regen10rk2v8pxjnldtxuy9ds0s5na9qjcmh5ymplz87",
          "pub_key": null,
          "account_number": "0",
          "sequence": "0"
        },
        {
          "@type": "/cosmos.vesting.v1beta1.PeriodicVestingAccount",
          "base_vesting_account": {
            "base_account": {
              "address": "regen1lusdjktpk3f2v33cda5uwnya5qcyv04cwvnkwz",
              "pub_key": null,
              "account_number": "0",
              "sequence": "0"
            },
            "original_vesting": [
              {
                "denom": "uregen",
                "amount": "500000301000"
              }
            ],
            "delegated_free": [],
            "delegated_vesting": [],
            "end_time": "1650961098"
          },
          "start_time": "1616774400",
          "vesting_periods": [
            {
              "length": "0",
              "amount": [
                {
                  "denom": "uregen",
                  "amount": "200000301000"
                }
              ]
            },
            {
              "length": "31556952",
              "amount": [
                {
                  "denom": "uregen",
                  "amount": "150000000000"
                }
              ]
            },
            {
              "length": "2629746",
              "amount": [
                {
                  "denom": "uregen",
                  "amount": "150000000000"
                }
              ]
            }
          ]
        }
      ]
    },
    "bank": {
      "params": {
        "send_enabled": [],
        "default_send_enabled": false
      },
      "balances": [
        {
          "address": "regen10rk2v8pxjnldtxuy9ds0s5na9qjcmh5ymplz87",
          "coins": [
            {
              "denom": "uregen",
              "amount": "100000000000"
            }
          ]
        },
        {
          "address": "regen1lusdjktpk3f2v33cda5uwnya5qcyv04cwvnkwz",
          "coins": [
            {
              "denom": "uregen",
              "amount": "500000301000"
            }
          ]
        }
      ],
      "supply": [
        {
          "denom": "uregen",
          "amount": "600000301000"
        }
      ],
      "denom_metadata": []
    }
  }
}`, json)
	require.Equal(t, `regen10rk2v8pxjnldtxuy9ds0s5na9qjcmh5ymplz87	100000	1
	100000	MAINNET
regen1lusdjktpk3f2v33cda5uwnya5qcyv04cwvnkwz	500000.301	3
	200000.3010	MAINNET
	150000.0	2022-03-26 21:49:12
	150000	2022-04-26 08:18:18
`, auditOut)
}
