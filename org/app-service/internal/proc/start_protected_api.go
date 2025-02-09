package proc

import (
	"github.com/org/2112-space-lab/org/app-service/internal/api/routers"
	"github.com/org/2112-space-lab/org/app-service/internal/clients/service"
	"github.com/org/2112-space-lab/org/app-service/internal/config"
	"github.com/org/2112-space-lab/org/app-service/internal/services"
)

// StartPublicApi starts de protected http server
func StartProtectedApi(services *services.ServiceComponent) {
	serviceCli := service.GetClient()
	c := serviceCli.GetConfig()
	protectedApiRouter := routers.InitProtectedAPIRouter(config.Env, services)
	protectedApiRouter.Start(c.Host, c.ProtectedApiPort)
}
