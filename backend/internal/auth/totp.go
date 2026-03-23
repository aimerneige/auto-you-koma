package auth

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func GenerateTOTPSecret(email string) (*otp.Key, error) {
	return totp.Generate(totp.GenerateOpts{
		Issuer:      "AutoYonKoma",
		AccountName: email,
	})
}

func ValidateTOTP(passcode string, secret string) bool {
	return totp.Validate(passcode, secret)
}
