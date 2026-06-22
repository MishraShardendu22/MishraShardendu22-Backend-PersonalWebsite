package util

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "my_secure_password"
	hash := HashPassword(password)

	if hash == "" {
		t.Errorf("Expected hashed password, got empty string")
	}

	if hash == password {
		t.Errorf("Expected hashed password to be different from raw password")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "my_secure_password"
	hash := HashPassword(password)

	if !CheckPassword(password, hash) {
		t.Errorf("Expected CheckPassword to return true for matching password and hash")
	}

	if CheckPassword("wrong_password", hash) {
		t.Errorf("Expected CheckPassword to return false for non-matching password and hash")
	}
}
