package utils

import "time"

func ParseDate(date string) (time.Time, error) {
	dateFormat := "02/01/2006"
	parsedDate, err := time.Parse(dateFormat, date)
	if err != nil {
		return time.Time{}, err
	}

	return parsedDate, nil
}
