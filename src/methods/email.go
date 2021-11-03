package methods

import "regexp"

// ValidateEmail is a method for validating email
func ValidateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return Re.MatchString(email)
}
