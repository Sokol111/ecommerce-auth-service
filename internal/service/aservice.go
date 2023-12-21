package service

import (
	"aidanwoods.dev/go-paseto"
	"context"
	"fmt"
	"github.com/Sokol111/ecommerce-auth-service/internal/model"
	"github.com/Sokol111/ecommerce-auth-service/internal/repository"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService struct {
	repository repository.UserRepository
	secretKey  paseto.V4AsymmetricSecretKey
	publicKey  paseto.V4AsymmetricPublicKey
}

func NewAuthService(repository repository.UserRepository, secretKeyHex string) *AuthService {
	secretKey, err := paseto.NewV4AsymmetricSecretKeyFromHex(secretKeyHex)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	return &AuthService{repository, secretKey, secretKey.Public()}
}

func (s *AuthService) Login(ctx context.Context, login string, password string) (string, error) {
	u, err := s.repository.GetByLogin(ctx, login)
	if err != nil {
		return "", fmt.Errorf("coudn't find user by login [%s]", login)
	}

	if !passwordMatch(u.HashedPassword, password) {
		return "", fmt.Errorf("password doesn't match")
	}

	return s.createToken(u), nil
}

func (s *AuthService) GetUserByToken(ctx context.Context, token string) (model.User, error) {
	parsed, err := s.parseToken(token)
	if err != nil {
		return model.User{}, fmt.Errorf("coudn't parse token [%s], reason: %s", token, err.Error())
	}

	id, err := parsed.GetString("userId")

	if err != nil {
		return model.User{}, fmt.Errorf("coudn't get user id from token [%s]", token)
	}

	u, err := s.repository.GetById(ctx, id)

	if err != nil {
		return model.User{}, fmt.Errorf("coudn't find user by id [%s]", id)
	}

	return u, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func passwordMatch(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (s *AuthService) createToken(u model.User) string {
	token := paseto.NewToken()
	now := time.Now()
	token.SetIssuedAt(now)
	token.SetNotBefore(now)
	token.SetExpiration(now.AddDate(0, 0, 2))
	token.SetString("userId", u.ID)
	err := token.Set("permissions", u.Permissions)
	if err != nil {
		log.Warn().Msg("failed to add permissions to token")
	}

	return token.V4Sign(s.secretKey, nil)
}

func (s *AuthService) parseToken(signed string) (paseto.Token, error) {
	parser := paseto.NewParser()
	t, err := parser.ParseV4Public(s.publicKey, signed, nil)
	return *t, err
}
