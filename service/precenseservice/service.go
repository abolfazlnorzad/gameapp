package precenseservice

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

func (s Service) GetPresence(ctx context.Context, request dto.GetPresenceRequest) (dto.GetPresenceResponse, error) {
	fmt.Println("req", request)

	// TODO - implement me
	return dto.GetPresenceResponse{Items: []dto.GetPresenceItem{
		{UserID: 1, Timestamp: 12452151},
		{UserID: 2, Timestamp: 124534551},
	}}, nil
}
