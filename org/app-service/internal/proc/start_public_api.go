package proc

import (
	"github.com/org/2112-space-lab/org/app-service/internal/api/routers"
	"github.com/org/2112-space-lab/org/app-service/internal/clients/service"
	"github.com/org/2112-space-lab/org/app-service/internal/config"
	"github.com/org/2112-space-lab/org/app-service/internal/dependencies"
)

// StartPublicApi starts de public http server
func StartPublicApi(deps *dependencies.Dependencies) {
	serviceCli := service.GetClient()
	c := serviceCli.GetConfig()
	publicApiRouter := routers.NewPublicRouter(config.Env, deps)
	publicApiRouter.Start(c.Host, c.PublicApiPort)
}
