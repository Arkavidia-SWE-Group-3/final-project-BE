package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ValidatePassword(password string) bool {
	if len(password) < 6 {
		return false
	}

	var (
		hasUpper  = false
		hasNumber = false
		hasSymbol = false
	)

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case '0' <= char && char <= '9':
			hasNumber = true
		default:
			hasSymbol = true
		}
	}

	if !hasUpper || !hasNumber || !hasSymbol {
		return false
	}

	return true
}
