#!/bin/sh
REGEN_HOME="/tmp/regen$(date +%s)"
RANDOM_KEY="randomregenvalidatorkey"
CHAIN_ID=regen-1
DENOM=uregen
MAXBOND=5000000000 # 1000REGEN

GENTX_FILE=$(find ./$CHAIN_ID/gentxs -iname "*.json")
LEN_GENTX=$(echo ${#GENTX_FILE})

# Gentx Start date
start="2021-03-31 15:00:00Z"
# Compute the seconds since epoch for start date
stTime=$(date --date="$start" +%s)

# Gentx End date
end="2021-04-08 15:00:00Z"
# Compute the seconds since epoch for end date
endTime=$(date --date="$end" +%s)

# Current date
current=$(date +%Y-%m-%d\ %H:%M:%S)
# Compute the seconds since epoch for current date
curTime=$(date --date="$current" +%s)

if [[ $curTime < $stTime ]]; then
    echo "start=$stTime:curent=$curTime:endTime=$endTime"
    echo "Gentx submission is not open yet. Please close the PR and raise a new PR after 31-March-2021 15:00:00"
    exit 0
else
    if [[ $curTime > $endTime ]]; then
        echo "start=$stTime:curent=$curTime:endTime=$endTime"
        echo "Gentx submission is closed"
        exit 0
    else
        echo "Gentx is now open"
        echo "start=$stTime:curent=$curTime:endTime=$endTime"
    fi
fi

if [ $LEN_GENTX -eq 0 ]; then
    echo "No new gentx file found."
else
    set -e

    echo "GentxFile::::"
    echo $GENTX_FILE

    echo "...........Init Regen.............."

    git clone https://github.com/regen-network/regen-ledger
    cd regen-ledger
    git checkout v1.0.0-rc1
    make build
    chmod +x ./build/regen

    ./build/regen keys add $RANDOM_KEY --keyring-backend test --home $REGEN_HOME

    ./build/regen init --chain-id $CHAIN_ID validator --home $REGEN_HOME

    echo "..........Fetching genesis......."
    rm -rf $REGEN_HOME/config/genesis.json
    curl -s https://raw.githubusercontent.com/regen-network/mainnet/main/$CHAIN_ID/genesis-prelaunch.json >$REGEN_HOME/config/genesis.json

    # this genesis time is different from original genesis time, just for validating gentx.
    sed -i '/genesis_time/c\   \"genesis_time\" : \"2021-03-29T00:00:00Z\",' $REGEN_HOME/config/genesis.json

    GENACC=$(cat ../$GENTX_FILE | sed -n 's|.*"delegator_address":"\([^"]*\)".*|\1|p')
    denomquery=$(jq -r '.body.messages[0].value.denom' ../$GENTX_FILE)
    amountquery=$(jq -r '.body.messages[0].value.amount' ../$GENTX_FILE)

    echo $GENACC
    echo $amountquery
    echo $denomquery

    # only allow $DENOM tokens to be bonded
    if [ $denomquery != $DENOM ]; then
        echo "invalid denomination"
        exit 1
    fi

    # limit the amount that can be bonded

    if [ $amountquery -gt $MAXBOND ]; then
        echo "bonded too much: $amt > $MAXBOND"
        exit 1
    fi

    ./build/regen add-genesis-account $RANDOM_KEY 100000000000000$DENOM --home $REGEN_HOME \
        --keyring-backend test

    ./build/regen gentx $RANDOM_KEY 90000000000000$DENOM --home $REGEN_HOME \
        --keyring-backend test --chain-id $CHAIN_ID

    cp ../$GENTX_FILE $REGEN_HOME/config/gentx/

    echo "..........Collecting gentxs......."
    ./build/regen collect-gentxs --home $REGEN_HOME
    sed -i '/persistent_peers =/c\persistent_peers = ""' $REGEN_HOME/config/config.toml

    ./build/regen validate-genesis --home $REGEN_HOME

    echo "..........Starting node......."
    ./build/regen start --home $REGEN_HOME &

    sleep 5s

    echo "...checking network status.."

    ./build/regen status --node http://localhost:26657

    echo "...Cleaning the stuff..."
    killall regen >/dev/null 2>&1
    rm -rf $REGEN_HOME >/dev/null 2>&1
fi
