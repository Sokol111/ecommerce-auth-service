package service

import (
	"context"
	"errors"
	"github.com/Sokol111/ecommerce-auth-service/internal/model"
	"github.com/Sokol111/ecommerce-auth-service/internal/repository"
	"time"
)

type UserService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) *UserService {
	return &UserService{repository}
}

func (s *UserService) GetById(ctx context.Context, id string) (model.User, error) {
	return s.repository.GetById(ctx, id)
}

func (s *UserService) Create(ctx context.Context, user model.User, password string) (model.User, error) {
	_, err := s.repository.GetByLogin(ctx, user.Login)
	if err == nil {
		return model.User{}, errors.New("login already in use")
	}
	hashed, err := hashPassword(password)
	if err != nil {
		return model.User{}, errors.New("failed to hash password")
	}
	user.HashedPassword = hashed
	now := time.Now()
	user.CreatedDate = now
	user.LastModifiedDate = now
	user.Version = 1
	return s.repository.Create(ctx, user)
}

func (s *UserService) Update(ctx context.Context, user model.User) (model.User, error) {
	return s.repository.Update(ctx, user)
}

func (s *UserService) GetUsers(ctx context.Context) ([]model.User, error) {
	return s.repository.GetUsers(ctx)
}
