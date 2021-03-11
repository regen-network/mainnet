package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	regen "github.com/regen-network/regen-ledger/app"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/types"
)

func init() {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount("regen", "regenpub")
	cfg.SetBech32PrefixForConsensusNode("regencons", "regenconspub")
	cfg.SetBech32PrefixForValidator("regenvaloper", "regenvaloperpub")
}

func main() {
	rootCmd := &cobra.Command{}

	rootCmd.AddCommand(&cobra.Command{
		Use:  "build-genesis [genesis-dir]",
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
	records, err := parseAccountsCsv(accountsCsv)
	if err != nil {
		return nil, nil, err
	}

	var accounts []Account
	for _, rec := range records {
		acc := recordToAccount(rec, genesisTime)
		accounts = append(accounts, acc)
	}

	_ = mergeAccounts(accounts)

	return nil, nil, nil
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
