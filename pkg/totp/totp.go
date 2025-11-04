package totp

import (
	"fmt"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func GenerateKeyString(issuer, accountName string) (secret string, otpauthURL string, err error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
	})
	if err != nil {
		return "", "", err
	}
	return key.Secret(), key.URL(), nil
}

func ValidateCode(code, secret string) (bool, error) {
	return totp.ValidateCustom(code, secret, time.Now().UTC(), totp.ValidateOpts{
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
		Period:    30,
	})
}

// Optional helper to get otpauth URI manually
// otpauth://totp/{issuer}:{account}?secret={secret}&issuer={issuer}
func KeyUri(issuer, account, secret string) string {
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s", issuer, account, secret)
}
