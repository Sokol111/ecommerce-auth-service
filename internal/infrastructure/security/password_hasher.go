package security

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/Sokol111/ecommerce-auth-service/internal/application/command"
)

const defaultCost = 12

type bcryptHasher struct {
	cost int
}

func newBcryptHasher() command.PasswordHasher {
	return &bcryptHasher{cost: defaultCost}
}

func (h *bcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (h *bcryptHasher) Compare(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
