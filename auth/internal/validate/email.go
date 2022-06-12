package validate

import "net/mail"

// Email - validate email
func Email(email string) error {
	_, err := mail.ParseAddress(email)

	return err
}
