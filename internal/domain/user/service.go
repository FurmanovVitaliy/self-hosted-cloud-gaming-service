package user

import (
	"cloud/internal/domain/user/util"
	"cloud/pkg/logger"
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var sekretKey = "secret" //!move separete file
// ?? JWT claims struct
type JwtClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type Service interface {
	GetList(ctx context.Context) ([]User, error)
	Create(ctx context.Context, req *CreateUserReq) (*CreateUserRes, error)
	Login(ctx context.Context, req *LogingUserReq) (*LogingUserRes, error)
}

type service struct {
	storage Storage
	loger   *logger.Logger
	timeout time.Duration
}

func NewService(srorage Storage, logger *logger.Logger) Service {
	return &service{
		storage: srorage,
		loger:   logger,
		timeout: 2 * time.Second,
	}
}

func (s *service) GetList(ctx context.Context) ([]User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	users, err := s.storage.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *service) Create(ctx context.Context, req *CreateUserReq) (*CreateUserRes, error) {

	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	hashPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: hashPassword,
	}
	id, err := s.storage.Create(ctx, u)
	if err != nil {
		return nil, err
	}
	return &CreateUserRes{
		ID:       id,
		Email:    u.Email,
		Username: u.Username,
	}, nil
}
func (s *service) Login(ctx context.Context, req *LogingUserReq) (*LogingUserRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	u, err := s.storage.FindByEmail(ctx, req.Email)
	if err != nil {
		return &LogingUserRes{}, err
	}
	err = util.ComparePassword(u.PasswordHash, req.Password)
	if err != nil {
		return &LogingUserRes{}, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtClaims{
		ID:       u.ID,
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    u.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	})
	signedString, err := token.SignedString([]byte(sekretKey))
	if err != nil {
		return &LogingUserRes{}, err
	}

	return &LogingUserRes{
		AccessToken: signedString,
		Username:    u.Username,
	}, nil
}
