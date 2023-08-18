package aclservice

import (
	"gameapp/entity"
	"gameapp/pkg/richerror"
)

type Repository interface {
	GetUserPermissionTitle(userID uint, role entity.Role) ([]entity.PermissionTitle, error)
}

type Service struct {
	repository Repository
}

func New(r Repository) Service {
	return Service{
		repository: r,
	}
}

func (s Service) CheckAccess(userID uint, role entity.Role, permissions ...entity.PermissionTitle) (bool, error) {
	const op = "aclservice.CheckAccess"
	permissionTitles, err := s.repository.GetUserPermissionTitle(userID, role)
	if err != nil {
		return false, richerror.New(op).WithErr(err)
	}
	for _, pt := range permissionTitles {
		for _, p := range permissions {
			if p == pt {
				return true, nil
			}
		}
	}

	return false, nil
}
