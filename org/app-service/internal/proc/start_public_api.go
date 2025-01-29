package proc

import (
	"github.com/Elbujito/2112/src/app-service/internal/api/routers"
	"github.com/Elbujito/2112/src/app-service/internal/clients/service"
	"github.com/Elbujito/2112/src/app-service/internal/config"
	"github.com/Elbujito/2112/src/app-service/internal/services"
)

// StartPublicApi starts de public http server
func StartPublicApi(services *services.ServiceComponent) {
	serviceCli := service.GetClient()
	c := serviceCli.GetConfig()
	publicApiRouter := routers.NewPublicRouter(config.Env, services)
	publicApiRouter.Start(c.Host, c.PublicApiPort)
}
