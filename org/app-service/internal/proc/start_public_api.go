package proc

import (
	"github.com/org/2112-space-lab/org/app-service/internal/api/routers"
	"github.com/org/2112-space-lab/org/app-service/internal/clients/service"
	"github.com/org/2112-space-lab/org/app-service/internal/config"
	"github.com/org/2112-space-lab/org/app-service/internal/services"
)

// StartPublicApi starts de public http server
func StartPublicApi(services *services.ServiceComponent) {
	serviceCli := service.GetClient()
	c := serviceCli.GetConfig()
	publicApiRouter := routers.NewPublicRouter(config.Env, services)
	publicApiRouter.Start(c.Host, c.PublicApiPort)
}
