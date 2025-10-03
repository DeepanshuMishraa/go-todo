package auth

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("Failed to hash the password: %w", err)
	}

	return string(hashedBytes), nil
}

func CheckPassword(hashedPassword string, password string) error {
	if len(hashedPassword) == 0 {
		return fmt.Errorf("hashed password is empty")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {
		return fmt.Errorf("Failed to check the password: %w", err)
	}

	return nil
}
