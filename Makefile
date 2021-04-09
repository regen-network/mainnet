
.PHONY: prelaunch-2

prelaunch-2:
	go run . build-genesis regen-prelaunch-2 --errors-as-warnings
