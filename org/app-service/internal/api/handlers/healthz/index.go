package healthz

import (
	"net/http"

	"github.com/org/2112-space-lab/org/app-service/internal/api/handlers"
	"github.com/org/2112-space-lab/org/app-service/internal/config"

	"github.com/labstack/echo/v4"
)

func Index(c echo.Context) error {
	payload := map[string]string{
		"message": "ok",
		"version": config.Env.Version,
	}

	return c.JSON(http.StatusOK, handlers.Success(payload))
}
