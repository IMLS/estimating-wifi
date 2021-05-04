package version

// Set this with -ldflags at build time.
var Semver string

func GetVersion() string {
	return Semver
}