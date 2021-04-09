
.PHONY: prelaunch prelaunch-2 regen-1

prelaunch-1:
	go run . build-genesis regen-prelaunch-1 --errors-as-warnings

prelaunch-2:
	go run . build-genesis regen-prelaunch-2 --errors-as-warnings

regen-1:
	go run . build-genesis regen-1
