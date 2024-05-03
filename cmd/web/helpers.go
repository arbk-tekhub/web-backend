package main

import (
	"net/url"
	"strconv"
	"strings"
)

func readString(qs url.Values, key string, defaultValue string) string {

	value := qs.Get(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func readInt(qs url.Values, key string, defaultValue int) int {

	stringValue := qs.Get(key)

	if stringValue == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(stringValue)
	if err != nil {
		return defaultValue
	}

	return intValue
}

func readCsv(qs url.Values, key string, defaultValue []string) []string {

	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}
