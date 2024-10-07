package usecase

import (
	"context"
	"net/http"
	"time"

	"github.com/FurmanovVitaliy/pixel-cloud/internal/domain/user"

	"github.com/FurmanovVitaliy/pixel-cloud/pkg/errors"

	"github.com/FurmanovVitaliy/pixel-cloud/util"
)

// dto
type CreateUserReq struct {
	Email    string `json:"email" bson:"_email"`
	Username string `json:"username" bson:"_username"`
	Password string `json:"password" bson:"_password"`
}

type CreateUserRes struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Email    string `json:"email" bson:"_email"`
	Username string `json:"username" bson:"_username"`
}

type LogingUserReq struct {
	Email    string `json:"email" bson:"_email"`
	Password string `json:"password" bson:"_password"`
}

type LogingUserRes struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	Username    string `json:"username" bson:"_username"`
	AccessToken string `json:"access_token" bson:"_access_token"`
}

// custom errors
var (
	ErrNotFound           = errors.New(http.StatusUnauthorized, "AS", "000001", "Login or password incorrect") //user not found
	ErrInvalidPassword    = errors.New(http.StatusUnauthorized, "AS", "000002", "Login or password incorrect") //password not correct
	ErrTokenSing          = errors.New(http.StatusInternalServerError, "AS", "000003", "Token signing error")
	ErrPasswordHashing    = errors.New(http.StatusInternalServerError, "AS", "000004", "Failed to hash password")
	ErrFailedToCreateUser = errors.New(http.StatusInternalServerError, "AS", "000005", "Failed to create user")
	ErrAlradyExists       = errors.New(http.StatusBadRequest, "AS", "000006", "User with this login already exists")
)

func (u *UseCase) SingIn(ctx context.Context, req *LogingUserReq) (*LogingUserRes, error) {
	res, err := u.userService.FindByEmail(ctx, req.Email)
	if err != nil {
		return &LogingUserRes{}, ErrNotFound
	}

	err = util.ComparePassword(res.PasswordHash, req.Password)
	if err != nil {
		return &LogingUserRes{}, ErrInvalidPassword

	}

	token, err := u.tokenService.CreateToken(res.ID, res.Username, time.Hour*24*7)
	if err != nil {
		return &LogingUserRes{}, ErrTokenSing
	}

	return &LogingUserRes{
		ID:          res.ID,
		Username:    res.Username,
		AccessToken: token,
	}, nil
}

func (u *UseCase) SingUp(ctx context.Context, req *CreateUserReq) (*CreateUserRes, error) {
	_, err := u.userService.FindByEmail(ctx, req.Email)
	if err == nil {
		return &CreateUserRes{}, ErrAlradyExists
	}

	hashPswd, err := util.HashPassword(req.Password)
	if err != nil {
		return &CreateUserRes{}, ErrPasswordHashing
	}

	user := user.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: hashPswd,
	}
	id, err := u.userService.CreateUser(ctx, user)
	if err != nil {
		return &CreateUserRes{}, ErrFailedToCreateUser
	}
	return &CreateUserRes{
		ID:       id,
		Email:    user.Email,
		Username: user.Username,
	}, nil
}
