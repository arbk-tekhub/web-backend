package utils

import (
	"net/url"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func StructToBsonMap(data any) (bson.M, error) {

	bsonbytes, err := bson.Marshal(data)
	if err != nil {
		return nil, err
	}

	var bsonM = make(bson.M)

	if err := bson.Unmarshal(bsonbytes, &bsonM); err != nil {
		return nil, err
	}

	return bsonM, nil
}

func ReadStringFromQueryParams(qs url.Values, key string, defaultValue string) string {

	value := qs.Get(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func ReadIntFromQueryParams(qs url.Values, key string, defaultValue int) int {

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

func ReadCsvFromQueryParams(qs url.Values, key string, defaultValue []string) []string {

	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}
