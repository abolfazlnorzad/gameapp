package matchingservice

import (
	"gameapp/dto"
	"gameapp/entity"
	"gameapp/pkg/errmsg"
	"gameapp/pkg/richerror"
	"time"
)

type Service struct {
	repo   Repository
	config Config
}

type Repository interface {
	AddUserToWaitingList(userID uint, category entity.Category) error
}
type Config struct {
	WaitingTimeout time.Duration `koanf:"waiting_timeout"`
}

func New() Service {
	return Service{}
}

func (s Service) AddToWaitingList(req dto.AddToWaitingListRequest) (dto.AddToWaitingListResponse, error) {
	const op = "matchingservice.AddToWaitingList"
	err := s.repo.AddUserToWaitingList(req.UserID, req.Category)
	if err != nil {
		return dto.AddToWaitingListResponse{}, richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected).WithMessage(errmsg.SomethingWentWrong)
	}
	return dto.AddToWaitingListResponse{
		Timeout: s.config.WaitingTimeout,
	}, nil
}
