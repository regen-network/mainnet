# Regen Network Mainnet

## Building genesis.json (For admin use)

Execute:
```shell
go run . build-genesis regen-1
```

For pre-launch, we can ignore errors:
```shell
go run . build-genesis regen-prelaunch-1 --errors-as-warnings
```

## Join as a validator

### Requirements

**Minimum hardware requirements**
- 8GB RAM
- 2 CPUs
- 200G SSD
- Ubuntu 18.04+ (Recommended)

Note: 2 sentry architecture is the bare minimum setup required.

**Software requirements**

#### Install Golang

```sh
sudo apt update
sudo apt install build-essential jq -y
wget https://dl.google.com/go/go1.15.6.linux-amd64.tar.gz
tar -xvf go1.15.6.linux-amd64.tar.gz
sudo mv go /usr/local
```

```sh
echo "" >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export GOROOT=/usr/local/go' >> ~/.bashrc
echo 'export GOBIN=$GOPATH/bin' >> ~/.bashrc
echo 'export PATH=$PATH:/usr/local/go/bin:$GOBIN' >> ~/.bashrc
```

Update PATH:
```sh
source ~/.bashrc
```

Verify Go installation:

```sh
go version # should be go1.15.6
```

### Setup Regen Ledger

**Clone the repo and install regen-ledger**
```sh
mkdir -p $GOPATH/src/github.com/regen-network
cd $GOPATH/src/github.com/regen-network
git clone https://github.com/regen-network/regen-ledger && cd regen-ledger
git fetch
git checkout v1.0.0
make install
```

**Verify installation**
```sh
regen version --long
```

it should display the following details:
```sh
name: regen
server_name: regen
version: v1.0.0
commit: 1b7c80ef102d3ae7cc40bba3ceccd97a64dadbfd
build_tags: netgo,ledger
go: go version go1.15.6 linux/amd64
```

### Start your validator node

- Step-1: Check and/or install correct regen software version
    Required version: `v1.0.0`

- Step-2: Download the mainnet genesis
    ```sh
    curl -s https://raw.githubusercontent.com/regen-network/mainnet/main/regen-1/genesis.json > ~/.regen/config/genesis.json
    ```

- Step-3: Verify genesis
    ```sh
    jq -S -c -M '' ~/.regen/config/genesis.json | shasum -a 256
    ```
    It should be equal to the contents in [checksum](regen-1/checksum.txt)

- Step-4: Update seeds and persistent peers

    Open `~/.regen/config/config.toml` and update `persistent_peers` and `seeds` (comma separated list)
    #### Persistent peers
    ```sh
    [TBD]
    ```
    #### Seeds
    ```sh
    [TBD]
    ```

- Step-5: Create systemd
    ```sh
    DAEMON_PATH=$(which regen)

    echo "[Unit]
    Description=regen daemon
    After=network-online.target
    [Service]
    User=${USER}
    ExecStart=${DAEMON_PATH} start
    Restart=always
    RestartSec=3
    LimitNOFILE=4096
    [Install]
    WantedBy=multi-user.target
    " >regen.service
    ```

- Step-6: Update system daemon and start regen node

    ```
    sudo mv regen.service /lib/systemd/system/regen.service
    sudo -S systemctl daemon-reload
    sudo -S systemctl start regen
    ```

### Create validator (Optional)
Note: This section is applicable for validators who wants to join post genesis time.

> **IMPORTANT:** Make sure your validator node is fully synced before running this command. Otherwise your validator will start missing blocks.

```sh
regen tx staking create-validator \
  --amount=9000000utree \
  --pubkey=$(regen tendermint show-validator) \
  --moniker="<your_moniker>" \
  --chain-id=aplikigo-1 \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --gas="auto" \
  --from=<your_wallet_name>
```

## Gentx submission [CLOSED]
This section applies to the validators who wants to join the genesis.

#### Step-1: Initialize the chain
```sh
regen init --chain-id regen-1 <your_validator_moniker>
```

#### Step-2: Replace the genesis
```sh
curl -s https://raw.githubusercontent.com/regen-network/mainnet/main/regen-1/genesis-prelaunch.json > $HOME/.regen/config/genesis.json
```
#### Step-3: Add/Recover keys
```sh
regen keys add <new_key>
```

or

```sh
regen keys add <key_name> --recover
```

#### Step-4: Create Gentx
```sh
regen gentx <key_name> <amount>  --chain-id regen-1
```

ex:
```sh
regen gentx validator 1000000000uregen --chain-id regen-1
```

**Note: Make sure to use the amount < available tokens in the genesis. Also max BONDED TOKENS allowed for gentxs are 50000REGEN or 50000000000uregen**

You might be interested to specify other optional flags. For ex:

```sh
regen gentx validator 1000000000uregen --chain-id regen-1 \
    --details <the validator details>
    --identity <The (optional) identity signature (ex. UPort or Keybase)>
    --commission-rate 0.1 \
    --commission-max-rate 0.2 \
    --commission-max-change-rate 0.01
```

It will show an output something similar to:
```
Genesis transaction written to "/home/ubuntu/.regen/config/gentx/gentx-9c8fe340885fd0178781eefcf24f32a5e448e15a.json"
```

**Note: If you are generating gentx offline on your local machine, append `--pubkey` flag to the above command. You can get pubkey of your validator by running `regen tendermint show-validator`**

#### Step-5: Fork regen-network mainnet repo
- Go to https://github.com/regen-network/mainnet
- Click on fork and chose your account (if many)

#### Step-6: Clone mainnet repo
```sh
git clone https://github.com/<your_github_username>/mainnet $HOME/mainnet
```

#### Step-7: Copy gentx to mainnet repo
```sh
cp ~/.regen/config/gentx/gentx-*.json $HOME/mainnet/regen-1/gentxs/
```

#### Step-8: Commit and push to your repo
```sh
cd $HOME/mainnet
git add regen-1/gentxs/*
git commit -m "<your validator moniker> gentx"
git push origin master
```

#### Step-9: Create gentx PR
- Go to your repository (on github)
- Click on Pull request and create a PR
- To make sure your submission is valid, please wait for the github action on your PR to complete