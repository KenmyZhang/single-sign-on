package app

import (
	l4g "github.com/alecthomas/log4go"

	"github.com/KenmyZhang/single-sign-on/model"
)

func CustomClaimsHasPermissionTo(customClaims model.CustomClaims, permission *model.Permission) bool {
	return CheckIfRolesGrantPermission(customClaims.GetUserRoles(), permission.Id)
}

func CheckIfRolesGrantPermission(roles []string, permissionId string) bool {
	for _, roleId := range roles {
		if role, ok := model.BuiltInRoles[roleId]; !ok {
			l4g.Debug("Bad role in system " + roleId)
			return false
		} else {
			permissions := role.Permissions
			for _, permission := range permissions {
				if permission == permissionId {
					return true
				}
			}
		}
	}

	return false
}
