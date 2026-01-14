package security

import (
	"golang.org/x/crypto/bcrypt"
)

const defaultCost = 12

// BcryptHasher implements password hashing using bcrypt
type BcryptHasher struct {
	cost int
}

// NewBcryptHasher creates a new BcryptHasher with default cost
func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{cost: defaultCost}
}

// Hash generates a bcrypt hash from the given password
func (h *BcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Compare compares a bcrypt hash with a password
func (h *BcryptHasher) Compare(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
