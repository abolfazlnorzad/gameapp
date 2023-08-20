package userservice

import (
	"context"
	"gameapp/dto"
)

func (s Service) GetProfile(ctx context.Context, req dto.ProfileRequest) (dto.ProfileResponse, error) {
	// we must pass sanitize data to service layer.
	u, err := s.repo.GetUserProfile(ctx, req.UserID)
	if err != nil {
		return dto.ProfileResponse{}, err
	}
	return dto.ProfileResponse{Name: u.Name}, nil
}
