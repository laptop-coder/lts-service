package utils

import (
	"errors"
	"os"
)

func CheckFileExistence(pathToFile *string) (*bool, error) {
	_, err := os.Stat(*pathToFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// File doesn't exist
			result := false
			return &result, nil
		}
		return nil, errors.New("error checking file existence: " + *pathToFile)
	}
	// File exists
	result := true
	return &result, nil
}
