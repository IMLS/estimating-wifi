package api

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

func FormatUri(scheme string, host string, path string) string {
	var uri string = (scheme + "://" +
		removeLeadingAndTrailingSlashes(host) +
		startsWithSlash(removeLeadingSlashes(path)))
	return uri
}
