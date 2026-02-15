package utils

import (
	"strconv"
)

func ConvertAnyToInt(any interface{}) (int, error) {
	num, err := strconv.Atoi(any.(string))
	if err != nil {
		return 0, err
	}
	return num, nil
}
