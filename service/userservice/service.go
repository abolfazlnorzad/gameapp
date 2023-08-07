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
	GetUserProfile(userID uint) (entity.User, error)
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
	User         entity.User
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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
	at, err := s.auth.GenerateAccessToken(user)
	if err != nil {
		return LoginResponse{}, err
	}

	rt, err := s.auth.GenerateRefreshToken(user)
	if err != nil {
		return LoginResponse{}, err
	}
	return LoginResponse{
		User:         user,
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}

type ProfileRequest struct {
	UserID uint `json:"user_id"`
}

type ProfileResponse struct {
	Name string `json:"name"`
}

func (s Service) GetProfile(req ProfileRequest) (ProfileResponse, error) {
	// we must pass sanitize data to service layer.
	u, err := s.repo.GetUserProfile(req.UserID)
	if err != nil {
		return ProfileResponse{}, err
	}
	return ProfileResponse{Name: u.Name}, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}