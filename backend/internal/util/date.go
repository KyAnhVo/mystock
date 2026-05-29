package util

import (
	"errors"
	"strconv"
)

type Date struct {
	y uint
	m uint
	d uint
}

func StringToDate(s string) (*Date, error) {
	if len(s) != len("yyyy-mm-dd") {
		return nil, errors.New("invalid date format. Formate is yyyy-mm-dd.")
	}
	date := &Date{}
	var temp int

	// year
	temp, err := strconv.Atoi(s[:4])
	if err != nil {
		return nil, errors.New("invalid date format. Formate is yyyy-mm-dd.")
	}
	if temp < 0 {
		return nil, errors.New("invalid date format. Formate is yyyy-mm-dd.")
	}
	date.y = uint(temp)

	// month
	temp, err = strconv.Atoi(s[5:7])
	if err != nil {
		return nil, errors.New("invalid date format. Formate is yyyy-mm-dd.")
	}
	if temp < 0 {
		return nil, errors.New("invalid date format. Formate is yyyy-mm-dd.")
	}
	date.m = uint(temp)

	// date
	temp, err = strconv.Atoi(s[8:10])
	if err != nil {
		return nil, errors.New("invalid date format. Formate is yyyy-mm-dd.")
	}
	if temp < 0 {
		return nil, errors.New("invalid date format. Formate is yyyy-mm-dd.")
	}
	date.d = uint(temp)

	return date, nil
}
