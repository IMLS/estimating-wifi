# IMLS Debian PPA repository

This PPA repository provides a standalone wifi session counter.

To install:

    bash <(curl -s -L https://raw.githubusercontent.com/18F/imls-pi-stack/main/imls-ppa/imls-ppa.shim)

## Release process

This section is mainly of interest to developers only.

### Preparation:

- Make sure the main branch is up to date and has the latest changes
- Make sure your GPG key is present (`gpg --list-keys`)
  - Your key should correspond to one of the available GPG public keys in this repository

### Release

__NB: these steps must be done on the raspberry pi__

- In `imls-raspberry-pi`, "make clean && make all"
- In `imls-ppa`, "make update-binaries"
  - This Makefile step copies over the release binaries to where they should be in the PPA paths, e.g., `sources/session-counter_1.1-1/usr/local/bin/session-counter`
  - This step also makes sure permissions are set properly
- If a dependency was added:
  - Add in the relevant `DEBIAN/control` file
- If a setup change is needed:
  - Modify in the relevant `DEBIAN/postinst` or `prerm` file
- If a systemd service needs to be changed:
  - Modify the relevant service
  - Currently we have the following services:
    - `reboot-1m`
    - `imls-update`
    - `session-counter`
- Update the version numbers in the `DEBIAN/control` files for each package that was updated
  - Please note that the versions in the PPA file names don’t matter!
- `"WHOM=youremail@gsa.gov" make release`
  - Where `youremail@gsa.gov` corresponds to your GPG key
- Commit and push!
