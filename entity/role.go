package entity

type Role uint8

const (
	UserRole Role = iota + 1
	AdminRole
)

const (
	UserRoleStr  = "user"
	AdminRoleStr = "admin"
)

func (r Role) String() string {
	switch r {
	case UserRole:
		return UserRoleStr
	case AdminRole:
		return AdminRoleStr
	}
	return ""
}

func MapToRoleEntity(r string) Role {
	switch r {
	case AdminRoleStr:
		return AdminRole
	case UserRoleStr:
		return UserRole
	}
	return Role(0)
}
