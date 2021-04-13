# Regen Network - Mainnet validator guide

## Requirements

### Minimum hardware requirements
- 8GB RAM
- 2 CPUs
- 200G SSD
- Ubuntu 18.04+ (Recommended)

Note: 2 sentry architecture is the bare minimum setup required.

### Software requirements

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

#### Setup Regen Ledger

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
[TBD]
```

## Gentx submission
This section applies to the validators who wants to join the genesis.

### Step-1: Initialize the chain
```sh
regen init --chain-id regen-1 <your_validator_moniker>
```

### Step-2: Replace the genesis
```sh
curl -s https://raw.githubusercontent.com/regen-network/mainnet/main/regen-1/genesis-prelaunch.json > $HOME/.regen/config/genesis.json
```
### Step-3: Add/Recover keys
```sh
regen keys add <new_key>
```

or

```sh
regen keys add <key_name> --recover
```

### Step-4: Create Gentx
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

### Step-5: Fork regen-network mainnet repo
- Go to https://github.com/regen-network/mainnet
- Click on fork and chose your account (if many)

### Step-6: Clone mainnet repo
```sh
git clone https://github.com/<your_github_username>/mainnet $HOME/mainnet
```

### Step-7: Copy gentx to mainnet repo
```sh
cp ~/.regen/config/gentx/gentx-*.json $HOME/mainnet/regen-1/gentxs/
```

### Step-8: Commit and push to your repo
```sh
cd $HOME/mainnet
git add regen-1/gentxs/*
git commit -m "<your validator moniker> gentx"
git push origin master
```

### Step-9: Create gentx PR
- Go to your repository (on github)
- Click on Pull request and create a PR
- To make sure your submission is valid, please wait for the github action on your PR to complete
