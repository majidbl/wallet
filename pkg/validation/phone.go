package validation

import "regexp"

func ValidatePhoneNumber(phoneNumber string) bool {
	// Regular expression pattern for phone number validation
	pattern := `^(\+98|0)9\d{9}$`

	// Compile the regex pattern
	regex := regexp.MustCompile(pattern)

	// Check if the phone number matches the pattern
	return regex.MatchString(phoneNumber)
}
