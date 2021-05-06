VERSION := $(shell git describe --tags --abbrev=0)

# 
.PHONY: dev

release:
	@echo $(VERSION) > prod-version.txt	

dev:
	@echo $(VERSION) > dev-version.txt
	pushd imls-playbook ; \
		sed 's/<<VERSION>>/$(VERSION)/' Makefile.in >> Makefile ; \
		popd
