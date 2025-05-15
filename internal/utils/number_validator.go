package utils

import (
	"errors"
	"regexp"
)

func ValidatePhoneNumber(phone string) error {
	regex := regexp.MustCompile(`^08[1-9][0-9]{7,10}$`)
	if !regex.MatchString(phone) {
		return errors.New("nomor HP tidak valid")
	}
	return nil
}
