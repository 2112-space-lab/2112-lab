package errors

import (
	"net/http"

	"github.com/Elbujito/2112/src/app-service/internal/api/handlers"
	xconstants "github.com/org/2112-space-lab/org/go-utils/pkg/fx/xconstants"

	"github.com/labstack/echo/v4"
)

func NotFound(c echo.Context) error {
	return c.JSON(
		http.StatusNotFound,
		handlers.BuildResponse(
			xconstants.STATUS_CODE_ROUTE_NOT_FOUND,
			xconstants.MSG_ROUTE_NOT_FOUND,
			[]string{},
			nil))
}
