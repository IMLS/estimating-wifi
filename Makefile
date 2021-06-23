VERSION := $(shell git describe --tags --abbrev=0)

.PHONY: dev

stamp_the_dev_version:
	@echo $(VERSION) > dev-version.txt
	git add dev-version.txt
	git commit -m "dev release: $(VERSION)"
	git push

stamp_the_release_version:
	@echo $(VERSION) > prod-version.txt
	git add prod-version.txt
	git commit -m "prod release: $(VERSION)"
	git push

packaging:
	cd imls-playbook ; \
		sed 's/<<VERSION>>/$(VERSION)/g' Makefile.in > Makefile ; \
		cd ..

ifeq ($(shell git describe --tags --abbrev=0),$(VERSION))
release:
	@echo "Version needs to be updated from " $(VERSION)
dev:
	@echo "Version needs to be updated from " $(VERSION)
else
# make VERSION=v1.2.3 release
release: stamp_the_release_version packaging
	cd imls-raspberry-pi ; make release ; cd ..
# make dev
dev: stamp_the_dev_version packaging
	cd imls-raspberry-pi ; make dev ; cd ..
endif
