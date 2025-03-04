package pkg

import (
	"encoding/base64"
	"regexp"

	"github.com/go-errors/errors"
)

func ValidateBase64(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return errors.New("must be a string")
	}

	_, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return errors.New("must be a valid base64 string")
	}

	return nil
}

var phoneRegex = regexp.MustCompile(`^[0-9]+$`)

func NumericRule(value any) error {
	phone, _ := value.(string)
	if !phoneRegex.MatchString(phone) {
		return errors.New("invalid phone number, must be numeric only")
	}
	return nil
}
