package state

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
