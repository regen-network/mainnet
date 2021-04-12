package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	regen "github.com/regen-network/regen-ledger/app"

	"github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	distribution "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/spf13/cobra"
)

const CommunityPoolRegenAmount = 2_000_000

func main() {
	rootCmd := &cobra.Command{}

	var errorsAsWarnings bool

	buildGenesisCmd := &cobra.Command{
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

			auditTsv, err := os.OpenFile(filepath.Join(genDir, "account_dump.tsv"), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)

			err = Process(doc, accountsCsv, CommunityPoolRegenAmount, auditTsv, errorsAsWarnings)
			if err != nil {
				return err
			}

			genFile := filepath.Join(genDir, "genesis.json")
			return doc.SaveAs(genFile)
		},
	}

	buildGenesisCmd.Flags().BoolVar(&errorsAsWarnings, "errors-as-warnings", false, "Allows records with errors to be ignored with a warning rather than failing")

	rootCmd.AddCommand(buildGenesisCmd)

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func Process(doc *types.GenesisDoc, accountsCsv io.Reader, communityPoolRegen int, auditOutput io.Writer, errorsAsWarnings bool) error {
	accounts, balances, err := buildAccounts(accountsCsv, doc.GenesisTime, auditOutput, errorsAsWarnings)
	if err != nil {
		return err
	}

	var genState map[string]json.RawMessage
	err = json.Unmarshal(doc.AppState, &genState)
	if err != nil {
		return err
	}

	cdc, _ := regen.MakeCodecs()

	err = setAccounts(cdc, genState, accounts, balances, communityPoolRegen)
	if err != nil {
		return err
	}

	doc.AppState, err = json.Marshal(genState)
	if err != nil {
		return err
	}

	return nil
}

func buildAccounts(accountsCsv io.Reader, genesisTime time.Time, auditOutput io.Writer, errorsAsWarnings bool) ([]auth.AccountI, []bank.Balance, error) {
	records, err := ParseAccountsCsv(accountsCsv, genesisTime, errorsAsWarnings)
	if err != nil {
		return nil, nil, err
	}

	accounts := make([]Account, 0, len(records))
	for _, record := range records {
		acc, err := RecordToAccount(record, genesisTime)
		if err != nil {
			return nil, nil, err
		}

		err = acc.Validate()
		if err != nil {
			buf := new(bytes.Buffer)
			PrintAccountAudit([]Account{acc}, genesisTime, buf)
			return nil, nil, fmt.Errorf("error on RecordToAccount: %w, Account: %s", err, buf.String())
		}

		accounts = append(accounts, acc)
	}

	accMap, err := MergeAccounts(accounts)
	if err != nil {
		return nil, nil, fmt.Errorf("error on MergeAccounts: %w", err)
	}

	err = AirdropRegenForMinFees(accMap, genesisTime)
	if err != nil {
		return nil, nil, err
	}

	accounts = SortAccounts(accMap)
	PrintAccountAudit(accounts, genesisTime, auditOutput)

	authAccounts := make([]auth.AccountI, 0, len(accounts))
	balances := make([]bank.Balance, 0, len(accounts))
	for _, acc := range accounts {
		authAcc, bal, err := ToCosmosAccount(acc, genesisTime)
		if err != nil {
			return nil, nil, fmt.Errorf("error on ToCosmosAccount: %w", err)
		}

		genAcc, ok := authAcc.(auth.GenesisAccount)
		if ok {
			err = genAcc.Validate()
			if err != nil {
				return nil, nil, err
			}
		}

		err = ValidateVestingAccount(authAcc)
		if err != nil {
			return nil, nil, err
		}

		authAccounts = append(authAccounts, authAcc)
		balances = append(balances, *bal)
	}

	return authAccounts, balances, nil
}

func setAccounts(cdc codec.Marshaler, genesis map[string]json.RawMessage, accounts []auth.AccountI, balances []bank.Balance, communityPoolRegen int) error {
	var bankGenesis bank.GenesisState

	err := cdc.UnmarshalJSON(genesis[bank.ModuleName], &bankGenesis)
	if err != nil {
		return err
	}

	// create distribution module account and corresponding balance with community pool funded
	distrMacc, distrBalance, err := buildDistrMaccAndBalance(communityPoolRegen)
	if err != nil {
		return err
	}

	bankGenesis.Balances = append(bankGenesis.Balances, balances...)
	bankGenesis.Balances = append(bankGenesis.Balances, *distrBalance)

	var supply sdk.Coins
	for _, bal := range bankGenesis.Balances {
		supply = supply.Add(bal.Coins...)
	}

	bankGenesis.Supply = supply

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

	distrMaccAny, err := cdctypes.NewAnyWithValue(distrMacc)
	if err != nil {
		return err
	}
	authGenesis.Accounts = append(authGenesis.Accounts, distrMaccAny)

	genesis[auth.ModuleName], err = cdc.MarshalJSON(&authGenesis)

	var distrGenesis distribution.GenesisState
	err = cdc.UnmarshalJSON(genesis[distribution.ModuleName], &distrGenesis)
	if err != nil {
		return err
	}

	// set CommunityPool to balance of distribution module account
	distrGenesis.FeePool.CommunityPool = sdk.NewDecCoinsFromCoins(distrBalance.Coins...)
	genesis[distribution.ModuleName], err = cdc.MarshalJSON(&distrGenesis)

	return nil
}

func buildDistrMaccAndBalance(regenAmount int) (auth.ModuleAccountI, *bank.Balance, error) {
	maccPerms := regen.GetMaccPerms()
	distrMacc := auth.NewEmptyModuleAccount(distribution.ModuleName, maccPerms[distribution.ModuleName]...)

	distrCoins, err := RegenToCoins(NewDecFromInt64(int64(regenAmount)))
	if err != nil {
		return nil, nil, err
	}

	distrBalance := &bank.Balance{
		Coins:   distrCoins,
		Address: distrMacc.Address,
	}

	return distrMacc, distrBalance, nil
}
