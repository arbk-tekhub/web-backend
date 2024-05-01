package validator

import (
	"regexp"
	"slices"
)

func PermittedValues[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

func Matchs(value string, rx regexp.Regexp) bool {
	return rx.MatchString(value)
}
