deb:
	rm -f *.deb
	dpkg-deb --build sources/input-initial-configuration_1.0-1_arm
	dpkg-deb --build sources/session-counter_1.0-1
	dpkg-deb --build sources/session-counter-csv_1.0-1_arm
	mv sources/*.deb .

packages: deb
	dpkg-scanpackages --multiversion . > Packages
	gzip -k -f Packages

release: packages
	apt-ftparchive release . > Release
	gpg --default-key "$$WHOM" -abs -o - Release > Release.gpg
	gpg --default-key "$$WHOM" --clearsign -o - Release > InRelease

all: packages release

clean:
	sudo rm -rf /opt/imls
	sudo apt remove -y session-counter session-counter-csv input-initial-configuration

reinstall: clean
	./imls-ppa.shim
