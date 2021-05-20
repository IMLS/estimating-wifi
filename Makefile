VERSION := $(shell git describe --tags --abbrev=0)

.PHONY: dev

stamp_the_dev_version: 
	@echo $(VERSION) > dev-version.txt

stamp_the_prod_version:
	@echo $(VERSION) > prod-version.txt

packaging:
	pushd imls-playbook ; \
		sed 's/<<VERSION>>/$(VERSION)/g' Makefile.in > Makefile ; \
		popd


ifeq ($(shell git describe --tags --abbrev=0),$(VERSION))
release:
	@echo "Version needs to be updated from " $(VERSION)
dev:
	@echo "Version needs to be updated from " $(VERSION)
else	
# make VERSION=v1.2.3 release
release: stamp_the_release_version packaging
# make dev
dev: stamp_the_dev_version versioning
	pushd imls-raspberry-pi ; make dev ; popd
endif


