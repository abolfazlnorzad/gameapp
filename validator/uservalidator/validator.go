package uservalidator

import "gameapp/entity"

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, error)
}

type Validator struct {
	repo Repository
}

func New(r Repository) Validator {
	return Validator{
		repo: r,
	}
}
