package middlewares

import (
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/clients/service"
	"github.com/org/2112-space-lab/org/go-utils/pkg/fx/xutils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// TimeoutMiddleware returns timeout Middleware
func TimeoutMiddleware() echo.MiddlewareFunc {
	serviceCli := service.GetClient()
	config := serviceCli.GetConfig()
	timeoutDuration := xutils.IntFromStr(config.RequestTimeoutDuration)

	return middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: time.Duration(timeoutDuration) * time.Second,
	})
}
