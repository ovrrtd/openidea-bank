package common

import "regexp"

func ValidateUrl(url string) bool {

	// Use the MatchString method to check if the phone number matches the pattern
	regExp := regexp.MustCompile(`^(https?|HTTPS?):\/\/(?:www\.)?[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(?:\/[^\s]*)?$`)
	return regExp.MatchString(url)
}
