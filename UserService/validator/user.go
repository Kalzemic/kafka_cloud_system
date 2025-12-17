package validator

import (
	"fmt"
	"unicode"
)

func ValidatePassword(password string) error {
	if len(password) < 3 {
		return fmt.Errorf("password is too short")
	}

	var upper, lower, digit bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			upper = true
		case unicode.IsLower(char):
			lower = true
		case unicode.IsDigit(char):
			digit = true
		}
	}

	switch {
	case !upper:
		return fmt.Errorf("password missing upper character")
	case !lower:
		return fmt.Errorf("password missing lower character")
	case !digit:
		return fmt.Errorf("password missing digit")
	default:
		return nil
	}
}
