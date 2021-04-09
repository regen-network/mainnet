#!/bin/sh
REGEN_HOME="/tmp/regen$(date +%s)"
CHAIN_ID=regen-1

set -e

echo "...........Init Regen.............."

git clone https://github.com/regen-network/regen-ledger
cd regen-ledger
git checkout v1.0.0-rc3
make build
chmod +x ./build/regen

./build/regen init --chain-id $CHAIN_ID validator --home $REGEN_HOME

echo "..........Fetching genesis......."
rm -rf $REGEN_HOME/config/genesis.json
curl -s https://raw.githubusercontent.com/regen-network/mainnet/main/$CHAIN_ID/genesis-prelaunch.json >$REGEN_HOME/config/genesis.json

echo "..........Collecting gentxs......."
./build/regen collect-gentxs --home $REGEN_HOME --gentx-dir ../$CHAIN_ID/gentxs

./build/regen validate-genesis --home $REGEN_HOME

cp $REGEN_HOME/config/genesis.json ../$CHAIN_ID/genesis.json
jq -S -c -M '' ../$CHAIN_ID/genesis.json | shasum -a 256 > ../$CHAIN_ID/checksum.txt

echo "..........Starting node......."
./build/regen start --home $REGEN_HOME &

sleep 5s

echo "...Cleaning the stuff..."
killall regen >/dev/null 2>&1
rm -rf $REGEN_HOME >/dev/null 2>&1

cd ..
rm -rf regen-ledger