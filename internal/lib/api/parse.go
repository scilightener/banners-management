package api

import (
	"avito-test-task/internal/lib/api/jsn"
	"avito-test-task/internal/lib/api/msg"
	"strconv"
)

func ParseInt64(s, pName string, num *int64) error {
	return parse(s, pName, num, func(s string) (int64, error) {
		return strconv.ParseInt(s, 10, 64)
	})
}

func ParseInt(s, pName string, num *int) error {
	return parse(s, pName, num, strconv.Atoi)
}

func ParseBool(s, pName string, b *bool) error {
	return parse(s, pName, b, strconv.ParseBool)
}

func parse[T any](s, pName string, val *T, parser func(s string) (T, error)) error {
	if len(s) == 0 {
		return jsn.DecodingError(msg.APIEmptyParameter(pName))
	}
	v, err := parser(s)
	if err != nil {
		return jsn.DecodingError(msg.APIUnacceptableFormat(pName))
	}

	*val = v
	return nil
}
