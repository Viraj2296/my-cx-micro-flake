package service_manager

import (
	"go.cerex.io/transcendflow/auth_util"
)

func (v *MachineDowntimeService) InitRouter() {
	v.ComponentManager.ModuleRouterGroup.GET("/overview", auth_util.TokenAuthMiddleware(), v.getOverview)

}
