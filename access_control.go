package goezyrouting

type (
	AccessController interface {
		HasAccess(role, minAccess role) bool
	}
	AccessControl struct {
	}
)

func NewAccessControl() *AccessControl {
	return &AccessControl{}
}

type role int

const (
	RoleUnknown role = iota
	RoleUser
	RoleAdmin
	RoleSuperAdmin
)

func (ac *AccessControl) HasAccess(role, minAccess role) bool {
	if role >= minAccess {
		return true
	}
	return false
}
