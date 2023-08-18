package entity

type PermissionTitle string

type Permission struct {
	ID    uint
	Title PermissionTitle
}

const (
	UserListPermission   = PermissionTitle("user-list")
	UserDeletePermission = PermissionTitle("user-delete")
)
