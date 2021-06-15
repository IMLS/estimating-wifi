packages:
	dpkg-scanpackages --multiversion . > Packages
	gzip -k -f Packages

release:
	apt-ftparchive release . > Release
	gpg --default-key "$WHOM" -abs -o - Release > Release.gpg
	gpg --default-key "$WHOM" --clearsign -o - Release > InRelease

