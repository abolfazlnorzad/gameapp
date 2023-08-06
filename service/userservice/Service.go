package userservice

import (
	"fmt"
	"gameapp/entity"
	"gameapp/pkg/name"
	"gameapp/pkg/phonenumber"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	Create(u entity.User) (entity.User, error)
}

type Service struct {
	repo Repository
}

type RegisterRequest struct {
	Name        string
	PhoneNumber string
	Password    string
}
type RegisterResponse struct {
	User entity.User
}

func (s Service) register(req RegisterRequest) (RegisterResponse, error) {
	// validate phone number and name
	if isMatched, err := phonenumber.IsPhoneNumberValid(req.PhoneNumber); err != nil || !isMatched {
		if err != nil {
			return RegisterResponse{}, err
		}
		if !isMatched {
			return RegisterResponse{}, fmt.Errorf("phone number is not valid")
		}
	}
	if l := name.NameMustBeMoreThanThreeChar(req.Name); !l {
		return RegisterResponse{}, fmt.Errorf("NameMustBeMoreThanThreeChar")
	}

	// check uniqueness phone number
	if isUnique, err := s.repo.IsPhoneNumberUnique(req.PhoneNumber); !isUnique || err != nil {
		if err != nil {
			return RegisterResponse{}, err
		}
		if !isUnique {
			return RegisterResponse{}, fmt.Errorf("phone number is not unique")
		}
	}

	// create new user
	createdUser, err := s.repo.Create(entity.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
	})

	if err != nil {
		return RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	// return new user
	return RegisterResponse{User: createdUser}, nil
}
