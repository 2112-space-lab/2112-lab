package proc

import (
	"github.com/org/2112-space-lab/org/app-service/internal/api/routers"
	"github.com/org/2112-space-lab/org/app-service/internal/clients/service"
	"github.com/org/2112-space-lab/org/app-service/internal/config"
	"github.com/org/2112-space-lab/org/app-service/internal/dependencies"
)

// StartPublicApi starts de protected http server
func StartProtectedApi(deps *dependencies.Dependencies) {
	serviceCli := service.GetClient()
	c := serviceCli.GetConfig()
	protectedApiRouter := routers.InitProtectedAPIRouter(config.Env, deps)
	protectedApiRouter.Start(c.Host, c.ProtectedApiPort)
}
