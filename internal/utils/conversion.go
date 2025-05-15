package utils

import (
	"errors"
	"strconv"
)

func ConvertStringToUint(s string) (uint, error) {
	uid64, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, errors.New("ID harus berupa angka valid")
	}
	return uint(uid64), nil
}
