package utils

import (
	. "backend/logger"
	"errors"
	_ "image/png"
	"os"
)

func DeleteThingPhotoFromStorageIfExists(pathToPhoto string) error {
	if _, err := os.Stat(pathToPhoto); err != nil {
		if os.IsNotExist(err) {
			Logger.Info("Photo of this thing doesn't exist, skipping deletion")
			return nil
		}
		msg := "error checking thing photo existence: " + err.Error()
		return errors.New(msg)
	}
	if err := os.Remove(pathToPhoto); err != nil {
		msg := "error deleting thing photo: " + err.Error()
		return errors.New(msg)
	}
	Logger.Info("Success. The photo of this thing has been deleted")
	return nil
}
