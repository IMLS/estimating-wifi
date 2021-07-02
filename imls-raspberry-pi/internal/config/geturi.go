package config

import "strings"

func removeTrailingSlashes(uri string) string {
	return strings.TrimSuffix(uri, "/")
}

func removeLeadingSlashes(uri string) string {
	return strings.TrimPrefix(uri, "/")
}

func removeLeadingAndTrailingSlashes(uri string) string {
	return removeTrailingSlashes(removeLeadingSlashes(uri))
}

func startsWithSlash(uri string) string {
	return "/" + uri
}

// func endsWithSlash(uri string) string {
// 	return uri + "/"
// }

func (cfg *Config) GetLoggingUri() string {
	var uri string = (cfg.Umbrella.Scheme + "://" +
		removeLeadingAndTrailingSlashes(cfg.Umbrella.Host) +
		startsWithSlash(removeLeadingSlashes(cfg.Umbrella.Logging)))
	return uri
}
