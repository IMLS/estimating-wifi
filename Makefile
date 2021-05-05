VERSION := $(shell git describe --tags --abbrev=0)

versioning:
	@echo $(VERSION) > prod-version.txt
