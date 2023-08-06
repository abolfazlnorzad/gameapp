package userservice

import (
	"fmt"
	"gameapp/entity"
	"gameapp/pkg/name"
	"gameapp/pkg/phonenumber"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	Create(u entity.User) (entity.User, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, bool, error)
}

type Service struct {
	repo Repository
}

func NewUserSvc(r Repository) Service {
	return Service{
		repo: r,
	}
}

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}
type RegisterResponse struct {
	User entity.User
}

func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {
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

	// todo - tech dept : we must validation password with regex
	hashedPass, hErr := HashPassword(req.Password)
	if hErr != nil {
		return RegisterResponse{}, hErr
	}
	// create new user
	createdUser, err := s.repo.Create(entity.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    hashedPass,
	})

	if err != nil {
		return RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	// return new user
	return RegisterResponse{User: createdUser}, nil
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginResponse struct {
	User entity.User
}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	// check phone number is exist  & // get user by phone number
	user, isExist, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected err : %w", err)
	}

	if !isExist {
		return LoginResponse{}, fmt.Errorf("phone number or password is wrong .")
	}

	// compare req.password with user.password
	ps := CheckPasswordHash(req.Password, user.Password)
	if !ps {
		return LoginResponse{}, fmt.Errorf("phone number or password is wrong .")
	}
	// return ok
	return LoginResponse{
		User: user,
	}, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
