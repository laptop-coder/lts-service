// TODO: add HEIC support. Now supports JPEG and PNG
package utils

import (
	. "backend/config"
	. "backend/logger"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"os"
	"strings"
)

// TODO: in frontend, when loading new photo from device, write to the start of
// base64-string information like data:image/png;base64,
//
// POSSIBLY OUTDATED
func SaveThingPhotoToStorage(photoBase64 string, thingId int64, thingType string) error {
	// Check base64 format
	// if photoBase64[:10] != "data:image" {
	// 	const msg = "base-64 string should contain \"data:image\""
	// 	Logger.Error(msg)
	// 	return errors.New(msg)
	// }

	// Remove data like "data:image/png;base64,"
	photoBase64 = photoBase64[strings.Index(photoBase64, ",")+1:]

	// Decode photo from base64
	decodedPhoto, err := base64.StdEncoding.DecodeString(photoBase64)
	if err != nil {
		msg := "error decoding thing photo from base64: " + err.Error()
		Logger.Error(msg)
		return errors.New(msg)
	}

	preparedPhoto, _, err := image.Decode(bytes.NewReader(decodedPhoto))
	if err != nil {
		msg := "error decoding thing photo from base64: " + err.Error()
		Logger.Error(msg)
		return errors.New(msg)
	}

	pathToPhoto := fmt.Sprintf(
		"%s/%s/%d.jpeg",
		Cfg.Storage.PathTo,
		thingType,
		thingId,
	)

	// Open file for writing
	file, err := os.Create(pathToPhoto)
	if err != nil {
		msg := "error opening file for writing: " + err.Error()
		Logger.Error(msg)
		return errors.New(msg)
	}

	// Save photo
	if err = jpeg.Encode(file, preparedPhoto, &jpeg.Options{Quality: 50}); err != nil {
		msg := "error saving photo: " + err.Error()
		Logger.Error(msg)
		return errors.New(msg)
	}

	Logger.Info("Success. The photo of the " + thingType + " thing was saved to the storage")
	return nil
}
