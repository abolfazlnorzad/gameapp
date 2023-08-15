package userservice

import (
	"gameapp/dto"
	"gameapp/pkg/richerror"
	"golang.org/x/crypto/bcrypt"
)

func (s Service) Login(req dto.LoginRequest) (dto.LoginResponse, error) {
	const op = "userservice.Login"
	// check phone number is exist  & // get user by phone number
	user, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return dto.LoginResponse{}, richerror.New(op).WithErr(err)
	}

	// compare req.password with user.password
	ps := checkPasswordHash(req.Password, user.Password)
	if !ps {
		return dto.LoginResponse{}, richerror.New(op).WithMessage("phone number or password is wrong .").WithKind(richerror.KindInvalid)
	}
	// return ok
	at, err := s.auth.GenerateAccessToken(user)
	if err != nil {
		return dto.LoginResponse{}, richerror.New(op).WithErr(err)
	}

	rt, err := s.auth.GenerateRefreshToken(user)
	if err != nil {
		return dto.LoginResponse{}, richerror.New(op).WithErr(err)
	}
	return dto.LoginResponse{
		User: dto.UserInfo{
			ID:          user.ID,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
		},
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
