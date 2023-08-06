package phonenumber

import (
	"fmt"
	"regexp"
)

func IsPhoneNumberValid(phoneNumber string) (bool, error) {
	isMatched, err := regexp.MatchString(`^[\+]?[(]?[0-9]{3}[)]?[-\s\.]?[0-9]{3}[-\s\.]?[0-9]{4,6}$`, phoneNumber)
	if err != nil || !isMatched {
		if err != nil {
			return false, fmt.Errorf("phone number regex is err %w", err)
		}
		if !isMatched {
			return false, nil
		}
	}

	return true, nil
}
