package presenceservice

import (
	"context"
	"fmt"
	"gameapp/dto"
	"gameapp/pkg/richerror"
	"time"
)

type Config struct {
	ExpirationTime time.Duration `koanf:"expiration_time"`
	Prefix         string        `koanf:"prefix"`
}

type Repository interface {
	Upsert(ctx context.Context, key string, timestamp int64, expTime time.Duration) error
	GetPresence(ctx context.Context, prefixKey string, userIDs []uint) (map[uint]int64, error)
}

type Service struct {
	repo   Repository
	config Config
}

func New(r Repository, cfg Config) Service {
	return Service{
		repo:   r,
		config: cfg,
	}
}

func (s Service) Upsert(ctx context.Context, req dto.UpsertPresenceRequest) (dto.UpsertPresenceResponse, error) {
	const op = "presenceservice.Upsert"
	err := s.repo.Upsert(ctx, fmt.Sprintf("%s:%d", s.config.Prefix, req.UserID), req.Timestamp, s.config.ExpirationTime)
	if err != nil {
		return dto.UpsertPresenceResponse{}, richerror.New(op).WithErr(err)
	}
	return dto.UpsertPresenceResponse{}, nil
}

func (s Service) GetRepo() Repository {
	return s.repo
}

func (s Service) GetPresence(ctx context.Context, request dto.GetPresenceRequest) (dto.GetPresenceResponse, error) {
	list, err := s.repo.GetPresence(ctx, s.config.Prefix, request.UserIDs)
	if err != nil {
		return dto.GetPresenceResponse{}, err
	}

	resp := dto.GetPresenceResponse{}
	for k, v := range list {
		resp.Items = append(resp.Items, dto.GetPresenceItem{
			UserID:    k,
			Timestamp: v,
		})
	}

	return resp, nil
}
