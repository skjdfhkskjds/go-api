package routes

import "regexp"

const (
	regexPathParamPattern = `^\{[^}]*\}$`
	regexWildcardPattern  = `^\*[^/]*$`
)

var (
	regexPathParam = regexp.MustCompile(regexPathParamPattern)
	regexWildcard  = regexp.MustCompile(regexWildcardPattern)
)

// isPathParam checks if a segment is a path parameter
func isPathParam(segment string) bool {
	return regexPathParam.MatchString(segment)
}

// isWildcard checks if a segment is a wildcard
func isWildcard(segment string) bool {
	return regexWildcard.MatchString(segment)
}
