// Package auth implements password hashing and JWT token management.
package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword returns a bcrypt hash for the given plain-text password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// VerifyPassword reports whether the password matches the stored hash.
func VerifyPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
