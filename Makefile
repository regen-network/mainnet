
.PHONY: prelaunch-1 prelaunch-2

prelaunch-1:
	go run . build-genesis regen-prelaunch-1 --errors-as-warnings

prelaunch-2:
	go run . build-genesis regen-prelaunch-2 --errors-as-warnings
