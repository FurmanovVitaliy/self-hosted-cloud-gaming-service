package user

import (
	"context"
	"time"

	"github.com/FurmanovVitaliy/pixel-cloud/pkg/errors"
)

type userStorage interface {
	Create(ctx context.Context, user User) (string, error)
	FindAll(ctx context.Context) ([]User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindOne(ctx context.Context, id string) (User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id string) error
}

type userService struct {
	userStorage userStorage
	timeout     time.Duration
}

func NewUserService(userStorage userStorage, timeout time.Duration) *userService {
	return &userService{
		userStorage: userStorage,
		timeout:     timeout,
	}
}

func (s *userService) GetList(ctx context.Context) ([]User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	users, err := s.userStorage.FindAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "US", "000001", "failed to get list of users")
	}
	return users, nil
}

func (s *userService) CreateUser(ctx context.Context, user User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	id, err := s.userStorage.Create(ctx, user)
	if err != nil {
		return "", errors.Wrap(err, "US", "000002", "failed to create user")
	}
	return id, nil
}

func (s *userService) FindByEmail(ctx context.Context, email string) (User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	user, err := s.userStorage.FindByEmail(ctx, email)
	if err != nil {
		return User{}, errors.Wrap(err, "US", "000003", "failed to find user by email")
	}
	return user, nil
}
