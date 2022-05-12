package config

// This can be set via command-line string if needed.
var lshw_exe string = "/usr/bin/lshw"

func SetLSHWLocation(path string) {
	lshw_exe = path
}

func GetLSHWLocation() string {
	return lshw_exe
}

var Verbose *bool = new(bool)
