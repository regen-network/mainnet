all: regen-1

.PHONY: regen-1

regen-1:
	go run . build-genesis regen-1
	mv -f regen-1/genesis.json regen-1/genesis-prelaunch.json
	bash -x ./scripts/gen-genesis.sh
