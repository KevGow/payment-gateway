package validation

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// LuhnCheck validates a credit card number using the Luhn algorithm, a common method for
// validating card numbers.
func LuhnCheck(cardNumber string) bool {
	var sum int
	numDigits := len(cardNumber)
	parity := numDigits % 2

	for i := 0; i < numDigits; i++ {
		digit, _ := strconv.Atoi(string(cardNumber[i]))
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	return sum%10 == 0
}

// ValidateExpirationDate checks if the expiration date of the credit card is valid.
func ValidateExpirationDate(expiry string) bool {
	// Assuming the date format is MM/YY
	layout := "01/06"
	expDate, err := time.Parse(layout, expiry)
	if err != nil {
		return false
	}

	now := time.Now()
	// Return the result of if the card has expired or not.
	return expDate.After(now)
}

// ValidateCVV checks if the CVV of the credit card is valid.
func ValidateCVV(cvv string) bool {
	// Assuming CVV should be a 3 or 4 digit number
	return regexp.MustCompile(`^\d{3,4}$`).MatchString(cvv)
}

// ValidatePaymentAmount checks if the payment amount is valid.
func ValidatePaymentAmount(amount float64) bool {
	// Assuming the amount should be a positive value.
	if amount <= 0.0 {
		return false
	}
	// Check if the amount has at most two decimal places.
	amountStr := fmt.Sprintf("%.2f", amount)
	// Check if the amount string has exactly two decimal places using a regular expression.
	amountRegex := regexp.MustCompile(`^\d+\.\d{2}$`)
	return amountRegex.MatchString(amountStr)
	// We can add more sophisticated validations here if needed.
	// For example, checking if the amount is within a valid range, bank of the card holder, etc
}
