// TODO: add more checks
package utils

import "errors"

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("the password must be at least 8 bytes long")
	}
	if len(password) > 72 {
		// SEE https://pkg.go.dev/golang.org/x/crypto/bcrypt#GenerateFromPassword
		return errors.New("the password must not be more than 72 bytes long (see https://pkg.go.dev/golang.org/x/crypto/bcrypt#GenerateFromPassword)")
	}
	return nil
}
