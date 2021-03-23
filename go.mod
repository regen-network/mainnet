module github.com/regen-network/mainnet

go 1.15

require (
	github.com/cockroachdb/apd/v2 v2.0.2
	github.com/cosmos/cosmos-sdk v0.42.0-rc0
	github.com/regen-network/regen-ledger v1.0.0-rc0
	github.com/stretchr/testify v1.7.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
