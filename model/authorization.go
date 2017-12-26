package model

type Permission struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Role struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

var PERMISSION_ASSIGN_SYSTEM_ADMIN_ROLE *Permission
var PERMISSION_MANAGE_ROLES *Permission
var PERMISSION_EDIT_OTHER_USERS *Permission

var PERMISSION_MANAGE_SYSTEM *Permission

var ROLE_SYSTEM_ADMIN *Role

var ROLE_USER *Role

var BuiltInRoles map[string]*Role

func InitalizePermissions() {
	PERMISSION_ASSIGN_SYSTEM_ADMIN_ROLE = &Permission{
		"assign_system_admin_role",
		"authorization.permissions.assign_system_admin_role.name",
		"authorization.permissions.assign_system_admin_role.description",
	}
	PERMISSION_MANAGE_ROLES = &Permission{
		"manage_roles",
		"authorization.permissions.manage_roles.name",
		"authorization.permissions.manage_roles.description",
	}
	PERMISSION_MANAGE_SYSTEM = &Permission{
		"manage_system",
		"authorization.permissions.manage_system.name",
		"authorization.permissions.manage_system.description",
	}
	PERMISSION_EDIT_OTHER_USERS = &Permission{
		"edit_other_users",
		"authorization.permissions.edit_other_users.name",
		"authorization.permissions.edit_other_users.description",
	}
}

func InitalizeRoles() {
	InitalizePermissions()
	BuiltInRoles = make(map[string]*Role)

	ROLE_USER = &Role{
		"normal_user",
		"authorization.roles.normal_user.name",
		"authorization.roles.normal_user.description",
		append(
			[]string{},
		),
	}
	BuiltInRoles[ROLE_USER.Id] = ROLE_USER

	ROLE_SYSTEM_ADMIN = &Role{
		"system_admin",
		"authorization.roles.global_admin.name",
		"authorization.roles.global_admin.description",
		[]string{
			PERMISSION_ASSIGN_SYSTEM_ADMIN_ROLE.Id,
			PERMISSION_MANAGE_SYSTEM.Id,
			PERMISSION_MANAGE_ROLES.Id,
			PERMISSION_EDIT_OTHER_USERS.Id,
		},
	}
	BuiltInRoles[ROLE_SYSTEM_ADMIN.Id] = ROLE_SYSTEM_ADMIN
}

func RoleIdsToString(roles []string) string {
	output := ""
	for _, role := range roles {
		output += role + ", "
	}

	if output == "" {
		return "[<NO ROLES>]"
	}

	return output[:len(output)-1]
}

func init() {
	InitalizeRoles()
}
