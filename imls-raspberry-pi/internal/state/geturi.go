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
	var uri string = (the_config.Umbrella.Scheme + "://" +
		removeLeadingAndTrailingSlashes(the_config.Umbrella.Host) +
		startsWithSlash(removeLeadingSlashes(the_config.Umbrella.Paths.Events)))
	return uri
}

func GetDurationsURI() string {
	var uri string = (the_config.Umbrella.Scheme + "://" +
		removeLeadingAndTrailingSlashes(the_config.Umbrella.Host) +
		startsWithSlash(removeLeadingSlashes(the_config.Umbrella.Paths.Durations)))
	return uri
}
