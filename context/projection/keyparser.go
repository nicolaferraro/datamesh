package projection

import (
	"regexp"
	"errors"
)

const (
	ErrorIllegalKey			= "EDB_ILLEGAL_KEY"
	ErrorNoRootSpecified	= "EDB_NO_ROOT_SPECIFIED"

	keyPart					= "[A-Za-z0-9]+[A-Za-z0-9_-]*"
	key						= keyPart + "(\\." + keyPart + "|\\[" + keyPart + "\\])*"
	separators				= "\\.|\\[|\\]"
)

var keyPattern		= regexp.MustCompile("^" + key + "$")
var separatorsPattern	= regexp.MustCompile(separators)


func ParseKey(key string) ([]string, error) {
	if !keyPattern.MatchString(key) {
		return nil, errors.New(ErrorIllegalKey)
	}

	parts := separatorsPattern.Split(key, -1)

	if len(parts) == 0 {
		return nil, errors.New(ErrorNoRootSpecified)
	}

	// Filter without allocating
	filteredParts := parts[:0]
	for _, p := range parts {
		if len(p) > 0 {
			filteredParts = append(filteredParts, p)
		}
	}

	return filteredParts, nil
}
