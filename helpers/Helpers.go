package helpers

import "golang.org/x/crypto/bcrypt"

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
