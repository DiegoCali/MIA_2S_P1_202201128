package utils

import (
	"fmt"
	"time"
)

func ConvertToBytes(size int, unit string) (int, error) {
	if unit == "b" || unit == "B" {
		return size, nil
	} else if unit == "k" || unit == "K" {
		size = size * 1024
	} else if unit == "m" || unit == "M" {
		size = size * 1024 * 1024
	} else {
		return -1, fmt.Errorf("unit %s not recognized", unit)
	}
	return size, nil
}

func TimeToFloat(time time.Time) float64 {
	year := time.Year()
	month := monthToInt(time.Month())
	day := time.Day()
	hour := time.Hour()
	minute := time.Minute()
	second := time.Second()
	dateFloat := float64(year*10000 + month*100 + day + hour/100 + minute/10000 + second/1000000)
	return dateFloat
}

func monthToInt(month time.Month) int {
	switch month {
	case time.January:
		return 1
	case time.February:
		return 2
	case time.March:
		return 3
	case time.April:
		return 4
	case time.May:
		return 5
	case time.June:
		return 6
	case time.July:
		return 7
	case time.August:
		return 8
	case time.September:
		return 9
	case time.October:
		return 10
	case time.November:
		return 11
	case time.December:
		return 12
	}
	return 0
}
