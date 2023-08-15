package userservice

import (
	"gameapp/dto"
	"gameapp/entity"
	"gameapp/pkg/richerror"
	"golang.org/x/crypto/bcrypt"
)

func (s Service) Register(req dto.RegisterRequest) (dto.RegisterResponse, error) {
	const op = "userservice.Register"
	// validate phone number and name

	// todo - tech dept : we must validation password with regex
	hashedPass, hErr := hashPassword(req.Password)
	if hErr != nil {
		return dto.RegisterResponse{}, hErr
	}
	// create new user
	createdUser, err := s.repo.Create(entity.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    hashedPass,
	})

	if err != nil {
		return dto.RegisterResponse{}, richerror.New(op)
	}

	// return new user
	return dto.RegisterResponse{User: dto.UserInfo{
		ID:          createdUser.ID,
		Name:        createdUser.Name,
		PhoneNumber: createdUser.PhoneNumber,
	}}, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
