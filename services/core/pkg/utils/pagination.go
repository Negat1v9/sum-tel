package utils

import (
	"net/url"
	"strconv"
)

var (
	DefaultLimit = 20
)

// return int "limit" value from url parametrs
func GetLimitParam(v url.Values, maxLimit int) int {
	limit, err := strconv.ParseInt(v.Get("limit"), 10, 32)
	if err != nil {
		return DefaultLimit
	}
	return min(maxLimit, int(limit))
}

// return int "offset" value from url parametrs
func GetOffset(v url.Values) int {
	page, err := strconv.ParseInt(v.Get("offset"), 10, 32)
	if err != nil {
		return 0
	}

	return int(page)
}
