package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/org/2112-space-lab/org/app-service/internal/config/constants"
)

// ResponseHeadersMiddleware returns Response Middleware
func ResponseHeadersMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(constants.HEADER_CONTENT_TYPE, constants.HEADER_CONTENT_TYPE_JSON)
			return next(c)
		}
	}
}
