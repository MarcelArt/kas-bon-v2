package common

import (
	"github.com/MarcelArt/kas-bon-v2/internal/enums"
	"github.com/casbin/casbin/v3"
)

func IsAuthorized(e *casbin.Enforcer, sub, app, dom, obj, act string) bool {
	if ok, _ := e.Enforce(sub, app, dom, enums.ResourceAll, enums.PermissionFull); ok {
		return true
	}

	ok, _ := e.Enforce(sub, app, dom, obj, act)
	return ok
}
