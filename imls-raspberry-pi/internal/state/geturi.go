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

func GetEventsURI() string {
	var uri string = (theConfig.Umbrella.Scheme + "://" +
		removeLeadingAndTrailingSlashes(theConfig.Umbrella.Host) +
		startsWithSlash(removeLeadingSlashes(theConfig.Umbrella.Paths.Events)))
	return uri
}

func GetDurationsURI() string {
	var uri string = (theConfig.Umbrella.Scheme + "://" +
		removeLeadingAndTrailingSlashes(theConfig.Umbrella.Host) +
		startsWithSlash(removeLeadingSlashes(theConfig.Umbrella.Paths.Durations)))
	return uri
}
