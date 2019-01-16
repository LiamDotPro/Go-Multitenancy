package tests

import (
	"github.com/LiamDotPro/Go-Multitenancy/helpers"
	"testing"
)

// Check to see if a password can be hashed correctly.
func TestHashPassword(t *testing.T) {
	_, err := helpers.HashPassword([]byte("A123456.{}$$*"))

	if err != nil {
		t.Fail()
	}
}

// Check to see if a password can be hashed correctly.
func TestHashPasswordAdmin(t *testing.T) {
	_, err := helpers.HashPasswordAdmin([]byte("A123456.{}$$*"))

	if err != nil {
		t.Fail()
	}
}

// Checks to see if a previously hashed password can be validated.
func TestCompareHashToPassword(t *testing.T) {
	hash, err := helpers.HashPassword([]byte("A123456.{}$$*"))

	if err != nil {
		t.Error("Couldn't hash password to use the result..")
	}

	result := helpers.CheckPasswordHash("A123456.{}$$*", hash)

	if !result {
		t.Error("Hashed password incorrectly or function malfunctioned.")
	}

	nonPasswordResult := helpers.CheckPasswordHash("notMyPassword", hash)

	if nonPasswordResult {
		t.Error("Incorrect password matched hash..")
	}
}

// Tests to see if a special character is found
func TestContainsSpecialCharacter(t *testing.T) {
	result := helpers.ContainsSpecialCharacter("@./,*£")

	if !result {
		t.Fail()
	}
}

//Test to see if a string contains a capital letter
func TestContainsCapitalLetter(t *testing.T) {
	result := helpers.ContainsCapitalLetter("a**Z")

	if !result {
		t.Fail()
	}
}
