package accesscontrol

// Define roles
const (
	RoleAdmin     = "admin"
	RoleHR        = "hr"
	RoleManager   = "manager"
	RoleUser      = "user"
	RoleSuperUser = "superuser"
)

// Define permissions
var rolePermissions = map[string]map[string]bool{
	RoleAdmin: {
		"create": true,
		"read":   true,
		"update": true,
	},
	RoleHR: {
		"create": true,
		"read":   true,
		"update": false,
	},
	RoleManager: {
		"create": false,
		"read":   true,
		"update": true,
	},
	RoleUser: {
		"read":   true,
		"update": true,
	},
	RoleSuperUser: {
		"create": true,
		"read":   true,
		"update": true,
		"delete": true,
	},
}

var taskPermissions = map[string]map[string]bool{
	RoleAdmin: {
		"create": true,
		"read":   true,
		"update": true,
		"delete": true,
	},
	RoleHR: {
		"read":   true,
		"update": true,
	},
	RoleManager: {
		"create": true,
		"read":   true,
		"update": true,
	},
	RoleUser: {
		"read":   true,
		"update": true,
	},
	RoleSuperUser: {
		"create": true,
		"read":   true,
		"update": true,
		"delete": true,
	},
}

// Check if role has permission
func HasPermission(role, permission string) bool {
	perms, ok := rolePermissions[role]
	if !ok {
		return false
	}
	return perms[permission]
}

// Check if role has permission for tasks
func HasTaskPermission(role, permission string) bool {
	perms, ok := taskPermissions[role]
	if !ok {
		return false
	}
	return perms[permission]
}
