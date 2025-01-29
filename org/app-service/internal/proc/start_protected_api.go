package proc

import (
	"github.com/Elbujito/2112/src/app-service/internal/api/routers"
	"github.com/Elbujito/2112/src/app-service/internal/clients/service"
	"github.com/Elbujito/2112/src/app-service/internal/config"
	"github.com/Elbujito/2112/src/app-service/internal/services"
)

// StartPublicApi starts de protected http server
func StartProtectedApi(services *services.ServiceComponent) {
	serviceCli := service.GetClient()
	c := serviceCli.GetConfig()
	protectedApiRouter := routers.InitProtectedAPIRouter(config.Env, services)
	protectedApiRouter.Start(c.Host, c.ProtectedApiPort)
}
