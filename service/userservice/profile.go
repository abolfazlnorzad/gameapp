package userservice

import "gameapp/dto"

func (s Service) GetProfile(req dto.ProfileRequest) (dto.ProfileResponse, error) {
	// we must pass sanitize data to service layer.
	u, err := s.repo.GetUserProfile(req.UserID)
	if err != nil {
		return dto.ProfileResponse{}, err
	}
	return dto.ProfileResponse{Name: u.Name}, nil
}
