package userservice

import (
	"context"
	"gameapp/entity"
)

type Repository interface {
	Create(u entity.User) (entity.User, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, error)
	GetUserProfile(ctx context.Context, userID uint) (entity.User, error)
}

type AuthGenerator interface {
	GenerateAccessToken(u entity.User) (string, error)
	GenerateRefreshToken(u entity.User) (string, error)
}

type Service struct {
	repo Repository
	auth AuthGenerator
}

func NewUserSvc(r Repository, a AuthGenerator) Service {
	return Service{
		repo: r,
		auth: a,
	}
}
