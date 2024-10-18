package enum

type Role string

const (
	RootRole  Role = "root"
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

var AllRoles = []Role{
	RootRole,
	AdminRole,
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

func IsRoot(role Role) bool {
	return role == RootRole
}

func IsAdmin(role Role) bool {
	return role == AdminRole
}

func IsUser(role Role) bool {
	return role == UserRole
}
