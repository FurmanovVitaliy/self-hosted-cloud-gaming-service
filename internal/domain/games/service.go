package games

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

var secretKey = "secret"

type service struct {
	storage Storage
}

type Service interface {
	GetAll() ([]Game, error)
	GetOne(id string) (Game, error)
	CheckToken(token string) (bool, error)
}

func NewService(storage Storage) Service {
	return &service{storage: storage}
}

func (s *service) GetOne(id string) (Game, error) {
	return s.storage.FindOne(context.Background(), id)
}

func (s *service) GetAll() ([]Game, error) {
	return s.storage.FindAll(context.Background())
}

func (s *service) CheckToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err == nil && token.Valid {
		return true, nil
	} else {
		fmt.Println("JWT is invalid:", err)
		return false, err
	}
}
