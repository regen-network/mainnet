module github.com/regen-network/mainnet

go 1.15

require (
	github.com/cockroachdb/apd/v2 v2.0.2 // indirect
	github.com/cosmos/cosmos-sdk v0.42.0-rc0
	github.com/regen-network/regen-ledger v1.0.0-rc0
	github.com/spf13/cobra v1.1.3 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.8 // indirect
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
