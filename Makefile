all: prelaunch-2 regen-1

.PHONY: prelaunch-2 regen-1

prelaunch-2:
	go run . build-genesis regen-prelaunch-2 --errors-as-warnings

regen-1:
	go run . build-genesis regen-1
	mv -f regen-1/genesis.json regen-1/genesis-prelaunch.json
