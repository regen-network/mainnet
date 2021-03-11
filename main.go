package main

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/types"
	"path/filepath"
	"time"
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

			genFile := filepath.Join(genDir, "genesis.json")
			return doc.SaveAs(genFile)
		},
	})

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func buildAccounts(defaultStartTime time.Time) error {
	records, err := parseAccountsCsv("accounts.csv")
	if err != nil {
		return err
	}

	var accounts []Account
	for _, rec := range records {
		acc := recordToAccount(rec, defaultStartTime)
		accounts = append(accounts, acc)
	}

	_ = mergeAccounts(accounts)

	return nil
}
