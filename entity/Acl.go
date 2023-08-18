package entity

type ActorType string

type Acl struct {
	ID           uint
	ActorID      uint
	ActorType    ActorType
	PermissionID uint
}

const (
	RoleActorType ActorType = "role"
	UserActorType           = "user"
)
