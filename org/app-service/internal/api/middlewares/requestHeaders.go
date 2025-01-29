package middlewares

import (
	"net/http"

	"github.com/Elbujito/2112/src/app-service/internal/api/handlers"
	xconstants "github.com/org/2112-space-lab/org/go-utils/pkg/fx/xconstants"

	"github.com/labstack/echo/v4"
)

// RequestHeadersMiddleware returns Header Middleware
func RequestHeadersMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if (len(c.Request().Header["Accept"]) == 0) || c.Request().Header["Accept"][0] != "application/json" {
				r := handlers.BuildResponse(
					xconstants.STATUS_CODE_NOT_ACCEPTABLE_WITHOUT_ACCEPT_HEADER,
					xconstants.MSG_NOT_ACCEPTABLE,
					[]string{xconstants.MSG_MISSING_ACCEPT_HEADER},
					nil)
				return c.JSON(http.StatusNotAcceptable, r)
			}
			if c.Request().Method == "GET" {
				return next(c)
			}
			if (len(c.Request().Header["Content-Type"]) == 0) || c.Request().Header["Content-Type"][0] != "application/json" {
				r := handlers.BuildResponse(
					xconstants.STATUS_CODE_NOT_ACCEPTABLE_WITHOUT_CONTENT_TYPE_HEADER,
					xconstants.MSG_NOT_ACCEPTABLE,
					[]string{xconstants.MSG_MISSING_CONTENT_TYPE_HEADER},
					nil)
				return c.JSON(http.StatusNotAcceptable, r)
			}
			return next(c)
		}
	}
}
