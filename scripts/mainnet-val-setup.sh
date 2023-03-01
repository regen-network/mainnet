#!/bin/bash
set -e

REGEN_HOME="${HOME}/.regen"
PERSISTENT_PEERS="69975e7afdf731a165e40449fcffc75167a084fc@104.131.169.70:26656,c4460b52c34ad4f12168d05807e998bb8e8b4812@3.22.22.248:26656,6eab01c9b75a1040ad3cda25a3e6c05d23b7e12c@3.23.48.152:26656"

function command_exists () {
    type "$1" &> /dev/null ;
}

function required_go_version () {
    current_go_version=$(go version | { read _ _ v _; echo "${v#go}"; })
    minimum_go_version="1.19.6"
    if [ "$(printf '%s\n' "$minimum_go_version" "$current_go_version" | sort -V | head -n1)" = "$minimum_go_version" ]; then
        return 0
    else
        return 1
    fi
}

clear

echo "                            @@@@@"
echo "                           @@@@@@,"
echo "                            (@@@@@"
echo "                              @@@@*"
echo "                               @@@@"
echo "                                @@@,    @@@                          ,@@@@%"
echo "                                 @@@    @@@        @@@             @@@@@@@@"
echo "                                  ,@    .@@        @@@          @@@@@@@@@"
echo "                        @@@@        @    @        @@          @@@@@@"
echo "                         /@@@                    @         @@@@%"
echo "                             &@                         @@@"
echo "                                                     /"
echo ""
echo ""
echo "  @@@@@@@@@@@@@@@@%                                          @@@@@"
echo "  @@@@@@@@@@@@@@@%                                            %@@@"
echo ""
echo "                  @@@@@@"
echo "                   *        @                         @"
echo "                      @@@@@                             %@@"
echo "                     @@@@     @@              @            @@@@@"
echo "                           @@@@     @    @    ,@              @@@@@@@"
echo "                          @@@@    ,@    .@@    @@@              ,@@@@@@@@/"
echo "                                 @@@    @@@    @@@@                @@@@@@@@"
echo "                                @@@*    ,@@                           @@@#"
echo "                               @@@@"
echo "                              @@@@&"
echo "                             @@@@@"
echo "                           @@@@@@@"
echo "                            @@@@@"
echo ""
echo ""
echo ""
echo ""
echo ""
echo " @@@@@@@@@@      *@@@@@@@@@@@@       (@@@@(       @@@@@@@@@@@@   @@@@#     @@@@"
echo " @@@@@@@@@@@@@   *@@@@@@@@@@@@   %@@@@@@@@@@@@    @@@@@@@@@@@@   @@@@@@    @@@@"
echo " @@@@     @@@@   *@@@@          @@@@@      @      @@@@           @@@#@@@%  @@@@"
echo " @@@@@@@@@@@@@   *@@@@@@@@@@    @@@@   @@@@@@@@   @@@@@@@@@@*    @@@  @@@@ @@@@"
echo " @@@@  /@@@@     *@@@@          @@@@#  ####@@@@   @@@@           @@@   @@@@@@@@"
echo " @@@@    @@@@@   *@@@@@@@@@@@@   @@@@@@  &@@@@@   @@@@@@@@@@@@/  @@@     @@@@@@"
echo " @@@@     #@@@@  *@@@@@@@@@@@@     /@@@@@@@ @@@   @@@@@@@@@@@@/  @@@      @@@@@"
echo ""
echo ""
echo ""
echo " @@     @@   @@@@@@@%@@@@@@@@@@@@     @@    *@   @@@@@@@    @@@@@@@   @@    @@"
echo " @ @@   @@   @           @@     @,   @ @    @  @@      *@   @     @@  @@  @@"
echo " @   @  @@   @@@@@@      @@      @  @@  @  @%  @&       @@  @@@@@@@   @@@&@@"
echo " @    @@@@   @           @@      %@ @   @@*@   #@      @@   @    @@   @@   @@"
echo " @      @@   @@@@@@@@    @@       @@     @@      .@@@@@     @     /@  @@     @@"

sleep 3
clear
echo "This script will install all of the necessary dependencies and perform node configurations for a Regen Network full-node"
sleep 3
clear

echo "Installing dependencies..."
sleep 3

sudo apt update
sudo apt install build-essential jq wget git -y

sleep 3
clear

echo "Installing Go 1.19..."

if command_exists go && required_go_version; then
    clear
    echo "Go 1.19 already installed"
    sleep 3
    clear
else
    sudo rm -rf /usr/local/go
    wget https://dl.google.com/go/go1.19.6.linux-amd64.tar.gz
    tar -xvf go1.19.6.linux-amd64.tar.gz
    sudo mv go /usr/local
    rm -rf go1.19.6.linux-amd64.tar.gz
    sleep 5
    clear
    echo "
export GOPATH=$HOME/go
export GOROOT=/usr/local/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:/usr/local/go/bin:$GOBIN
" >> ~/.bashrc
    source ~/.bashrc
    clear
    echo "Installed Go 1.19.6 and added it to your .bashrc"
    sleep 5
    clear
fi

export GOPATH=$HOME/go
export GOROOT=/usr/local/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:/usr/local/go/bin:$GOBIN

if [ -d "$REGEN_HOME" ]; then
    echo "--------------- WARNING! ---------------"
    echo "A home directory for the regen binary already exists."
    echo "The following step will remove $REGEN_HOME from your system."
    while true; do
        read -rp $'Are you sure you would like to continue (y/n)?\n' yn
        case $yn in
            [yY][eE][sS]|[yY]) rm -rf "$REGEN_HOME"; break;;
            [nN][oO]|[nN]) exit;;
            * ) echo "Please answer yes or no.";;
        esac
    done
    clear
fi

echo "Installing the regen binary (v1.0.0)..."
sleep 3
cd $HOME
rm -rf regen-ledger
git clone https://github.com/regen-network/regen-ledger
cd regen-ledger
git fetch
git checkout v1.0.0
make install

sleep 5
clear
echo "The regen binary (v1.0.0) has been installed on your system"
sleep 5
clear
echo "Setting validator key and node moniker..."

while true; do
    echo "Please enter a name for your key:"
    read -r KEY_NAME
    echo "Please enter a moniker for your node:"
    read -r NODE_MONIKER
    echo "Your key name is $KEY_NAME and your node moniker is $NODE_MONIKER."
    read -rp $'Is this correct (y/n)?\n' yn
    case $yn in
        [yY][eE][sS]|[yY]) break;;
        [nN][oO]|[nN]) ;;
        * ) echo "Please answer yes or no.";;
    esac
done

while true; do
    echo "Creating validator key..."
    regen keys add $KEY_NAME
    echo ""
    echo "After you have copied the mnemonic phrase in a safe place, press [ENTER] to continue."
    read -r -s -d $'\x0a'
    read -rp $'Are you sure you copied your mnemonic phrase (y/n)?\n' yn
    case $yn in
        [yY][eE][sS]|[yY]) break;;
        [nN][oO]|[nN]) ;;
        * ) echo "Please answer yes or no.";;
    esac
done

sleep 5
clear

echo "Initializing node..."
regen init --chain-id regen-1 $NODE_MONIKER
sleep 2
clear

echo "Downloading Regen Mainnet genesis file..."
curl -s https://raw.githubusercontent.com/regen-network/mainnet/main/regen-1/genesis.json > $REGEN_HOME/config/genesis.json
sleep 2
clear

echo "Configuring RPC address..."
sed -i 's#tcp://127.0.0.1:26657#tcp://0.0.0.0:26657#g' $REGEN_HOME/config/config.toml
sleep 2
clear

echo "Configuring seed nodes..."
sed -i '/persistent_peers =/c\persistent_peers = "'"$PERSISTENT_PEERS"'"' $REGEN_HOME/config/config.toml
sleep 2
clear

echo "Installing cosmovisor..."
go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.4.0
sleep 2
clear

echo "Setting up genesis binary..."
mkdir -p $REGEN_HOME/cosmovisor/genesis/bin
cp $GOBIN/regen $REGEN_HOME/cosmovisor/genesis/bin
sleep 2
clear

echo "Creating cosmovisor service file..."
echo "[Unit]
Description=Cosmovisor daemon
After=network-online.target
[Service]
Environment="DAEMON_NAME=regen"
Environment="DAEMON_HOME=${HOME}/.regen"
Environment="DAEMON_ALLOW_DOWNLOAD_BINARIES=false"
Environment="DAEMON_RESTART_AFTER_UPGRADE=true"
Environment="DAEMON_POLL_INTERVAL=300ms"
Environment="DAEMON_DATA_BACKUP_DIR=${HOME}/.regen"
Environment="UNSAFE_SKIP_BACKUP=false"
Environment="DAEMON_PREUPGRADE_MAX_RETRIES=0"
User=${USER}
ExecStart=${GOBIN}/cosmovisor run start
Restart=always
RestartSec=3
LimitNOFILE=4096
[Install]
WantedBy=multi-user.target
" >cosmovisor.service
sudo mv cosmovisor.service /etc/systemd/system/cosmovisor.service
sleep 3
clear

echo "Starting cosmovisor service..."
sudo systemctl daemon-reload
sudo systemctl enable cosmovisor
sudo service cosmovisor start

sleep 5
clear

echo "Congratulations! You have successfully set up your node."
echo ""
echo "Check the status of you node by running the following command:"
echo ""
echo "regen status"
echo ""
echo "In order to become a validator, you will first need to fund your new account:"
echo ""
echo "regen keys show $KEY_NAME -a"
echo ""
echo "When finished, you can create your validator by customizing and running the following command:"
echo ""
echo "regen tx staking create-validator --amount 9000000000uregen --commission-max-change-rate \"0.1\" --commission-max-rate \"0.20\" --commission-rate \"0.1\" --details \"Some details about your validator\" --from <keyname> --pubkey=\"$(regen tendermint show-validator)\" --moniker <your moniker> --min-self-delegation \"1\" --chain-id regen-1 --gas auto --fees 5000uregen"