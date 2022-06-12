package validate

import (
	"fmt"
	"strings"
)

// Password - validate password.
func Password(password string, l int, mustContainSpecial bool) error {
	const (
		special = " !\"#$%&'()*+,-./:;<=>?@[]^_`{|}~"
		numbers = "01234567890"
	)

	if len(password) < l {
		return fmt.Errorf("invalid password length: must be at least %d characters long", l)
	}

	if mustContainSpecial {
		if strings.ToLower(password) == password {
			return fmt.Errorf("password must contain at least one capital letter")
		}

		if !strings.ContainsAny(password, special) {
			return fmt.Errorf("password must contain at least one special letter")
		}

		if !strings.ContainsAny(password, numbers) {
			return fmt.Errorf("password must contain at least one number")
		}
	}

	return nil
}
