packages:
	dpkg-scanpackages --multiversion . > Packages
	gzip -k -f Packages

release:
	apt-ftparchive release . > Release
	gpg --default-key "matthew.jadud@gsa.gov" -abs -o - Release > Release.gpg
	gpg --default-key "matthew.jadud@gsa.gov" --clearsign -o - Release > InRelease

