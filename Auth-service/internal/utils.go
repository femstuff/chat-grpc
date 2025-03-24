package internal

import (
	"errors"
)

func Validate(name, email, pass string) error {
	if name == "" {
		return errors.New("name field cannot be empty")
	}

	if email == "" {
		return errors.New("email field cannot be empty")
	}

	if len(pass) < 4 {
		return errors.New("pass must be more than 4 char")
	}

	return nil
}
