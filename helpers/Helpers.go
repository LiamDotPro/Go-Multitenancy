package helpers

import (
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

// Hashes a password
// Running 4 rounds to comply with bcrypt recommendations for standard user.
func HashPassword(password []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(password, 4)
	return string(hash), err
}

// Hashes a password for a super user
// Running 8 rounds to comply with bcrypt recommendations for a super user.
func HashPasswordAdmin(password []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(password, 8)
	return string(hash), err
}

// Check a users password from a hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Checks a string to see if it contains a special character
func ContainsSpecialCharacter(s string) bool {
	for i := 0; i < len(s); i++ {
		switch b := s[i]; {
		case b >= 'a' && b <= 'z':
			continue
		case b >= 'A' && b <= 'Z':
			continue
		case b >= '0' && b <= '9':
			continue
		default:
			return true
		}
	}
	return false
}

// Checks a string to make sure there is at least one capital letter
// A side effect of this method is that one alphabet character must be present.
func ContainsCapitalLetter(str string) bool {
	for i := 0; i < len(str); i++ {
		// Check character code to see if it's between character capitalization byte sequence
		if str[i] >= 'A' && str[i] <= 'Z' {
			return true
		}
	}
	return false

}

// Validates an email address using a regular expression.
func ValidateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}
