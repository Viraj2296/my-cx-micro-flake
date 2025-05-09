package http

import (
	"go.cerex.io/transcendflow/auth_util"
)

func (v *Service) InitRouter() {
	v.ComponentManager.ModuleRouterGroup.GET("/overview", auth_util.TokenAuthMiddleware(), v.getOverview)
}
