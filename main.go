package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	regen "github.com/regen-network/regen-ledger/app"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/types"
)

func main() {
	rootCmd := &cobra.Command{}

	rootCmd.AddCommand(&cobra.Command{
		Use:  "build-genesis [genesis-dir]",
		Long: "Builds a [genesis-dir]/genesis.json file from accounts.csv and [genesis-dir]/genesis.tmpl.json",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			genDir := args[0]
			genTmpl := filepath.Join(genDir, "genesis.tmpl.json")
			doc, err := types.GenesisDocFromFile(genTmpl)
			if err != nil {
				return err
			}

			accountsCsv, err := os.Open("accounts.csv")
			if err != nil {
				return err
			}

			accounts, balances, err := buildAccounts(accountsCsv, doc.GenesisTime)
			if err != nil {
				return err
			}

			var genState map[string]json.RawMessage
			err = json.Unmarshal(doc.AppState, &genState)
			if err != nil {
				return err
			}

			cdc, _ := regen.MakeCodecs()

			err = setAccounts(cdc, genState, accounts, balances)
			if err != nil {
				return err
			}

			doc.AppState, err = json.Marshal(genState)
			if err != nil {
				return err
			}

			genFile := filepath.Join(genDir, "genesis.json")
			return doc.SaveAs(genFile)
		},
	})

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func buildAccounts(accountsCsv io.Reader, genesisTime time.Time) ([]auth.AccountI, []bank.Balance, error) {
	records, err := ParseAccountsCsv(accountsCsv, genesisTime)
	if err != nil {
		return nil, nil, err
	}

	accounts := make([]Account, 0, len(records))
	for _, record := range records {
		acc, err := RecordToAccount(record, genesisTime)
		if err != nil {
			return nil, nil, err
		}
		accounts = append(accounts, acc)
	}

	accounts, err = MergeAccounts(accounts)
	if err != nil {
		return nil, nil, err
	}

	// TODO: emit audit file with all accounts and distributions
	// 4. convert to cosmos accounts

	authAccounts := make([]auth.AccountI, 0, len(accounts))
	balances := make([]bank.Balance, 0, len(accounts))
	for _, acc := range accounts {
		authAcc, bal, err := ToCosmosAccount(acc, genesisTime)
		if err != nil {
			return nil, nil, err
		}
		authAccounts = append(authAccounts, authAcc)
		balances = append(balances, *bal)
	}

	return authAccounts, balances, nil
}

func setAccounts(cdc codec.Marshaler, genesis map[string]json.RawMessage, accounts []auth.AccountI, balances []bank.Balance) error {
	var bankGenesis bank.GenesisState

	err := cdc.UnmarshalJSON(genesis[bank.ModuleName], &bankGenesis)
	if err != nil {
		return err
	}

	bankGenesis.Balances = append(bankGenesis.Balances, balances...)

	genesis[bank.ModuleName], err = cdc.MarshalJSON(&bankGenesis)

	var authGenesis auth.GenesisState

	err = cdc.UnmarshalJSON(genesis[auth.ModuleName], &authGenesis)
	if err != nil {
		return err
	}

	for _, acc := range accounts {
		any, err := cdctypes.NewAnyWithValue(acc)
		if err != nil {
			return err
		}

		authGenesis.Accounts = append(authGenesis.Accounts, any)
	}

	genesis[auth.ModuleName], err = cdc.MarshalJSON(&authGenesis)

	return nil
}
