package valueobject

import (
	"errors"
	"unicode"
)

var (
	ErrPasswordTooShort  = errors.New("password too short")
	ErrPasswordNoUpper   = errors.New("password missing uppercase")
	ErrPasswordNoLower   = errors.New("password missing lowercase")
	ErrPasswordNoDigit   = errors.New("password missing digit")
	ErrPasswordNoSpecial = errors.New("password missing special character")
)

func Password(pw string) (string, error) {
	if len(pw) < 8 {
		return "", ErrPasswordTooShort
	}

	var (
		upper, lower, digit, special bool
	)

	for _, r := range pw {
		switch {
		case unicode.IsUpper(r):
			upper = true
		case unicode.IsLower(r):
			lower = true
		case unicode.IsDigit(r):
			digit = true
		case !unicode.IsLetter(r) && !unicode.IsDigit(r):
			special = true
		}
	}

	switch {
	case !upper:
		return "", ErrPasswordNoUpper
	case !lower:
		return "", ErrPasswordNoLower
	case !digit:
		return "", ErrPasswordNoDigit
	case !special:
		return "", ErrPasswordNoSpecial
	}

	return pw, nil
}
