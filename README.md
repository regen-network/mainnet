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
Check out these [instructions](./regen-1/README.md#Requirements) for installing `regen-ledger@v1.0.0`

If you haven't initialized your node, init regen chain by running
```sh
regen init --chain-id regen-1 <my_node_moniker>
```

### Start your validator node

- Step-1: Verify installation
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
    69975e7afdf731a165e40449fcffc75167a084fc@104.131.169.70:26656
    ```
    #### Seeds
    ```sh
    aebb8431609cb126a977592446f5de252d8b7fa1@104.236.201.138:26656
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

That's all! Your node should be up and running now. You can query your node by executing the following command after the genesis time

```sh
regen status
```

### Create validator (Optional)
Note: This section is applicable for validators who wants to join post genesis time.

> **IMPORTANT:** Make sure your validator node is fully synced before running this command. Otherwise your validator will start missing blocks.

```sh
regen tx staking create-validator \
  --amount=9000000uregen \
  --pubkey=$(regen tendermint show-validator) \
  --moniker="<your_moniker>" \
  --chain-id=regen-1 \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --gas="auto" \
  --from=<your_wallet_name>
```
