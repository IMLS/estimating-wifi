deb:
	dpkg-deb --build sources/input-initial-configuration_1.0-1_arm
	dpkg-deb --build sources/session-counter_1.0-1
	dpkg-deb --build sources/session-counter-csv_1.0-1_arm
	mv sources/*.deb .

packages: deb
	dpkg-scanpackages --multiversion . > Packages
	gzip -k -f Packages

release:
	apt-ftparchive release . > Release
	gpg --default-key "$$WHOM" -abs -o - Release > Release.gpg
	gpg --default-key "$$WHOM" --clearsign -o - Release > InRelease

all: deb packages release
