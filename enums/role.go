package enums

type Role string

const (
	AdminRole             Role = "admin"
	AssociationLeaderRole Role = "association_leader"
	UserRole              Role = "user"
)

var AllRoles = []Role{
	AdminRole,
	AssociationLeaderRole,
	UserRole,
}

func IsValidRole(role Role) bool {
	for _, r := range AllRoles {
		if r == role {
			return true
		}
	}
	return false
}

func IsAdmin(role Role) bool {
	return role == AdminRole
}

func IsAssociationLeader(role Role) bool {
	return role == AssociationLeaderRole
}

func IsUser(role Role) bool {
	return role == UserRole
}
