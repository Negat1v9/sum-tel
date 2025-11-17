package utils

import "strconv"

func ParseInt64Param(param string) (int64, error) {
	strconvedValue, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, err
	}
	return strconvedValue, nil
}
