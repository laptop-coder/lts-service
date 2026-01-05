package utils

import (
	"regexp"
)

func CheckStringSecurity(s string) (*bool, error) {
	if s == "" {
		result := true
		return &result, nil
	}
	result, err := regexp.MatchString("^[0-9A-Za-zА-Яа-яЁё:;_=,\"\\.\\?\\-\\+\\*\\(\\)\\/\\\\\\n\\r\\t ]+$", s)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
