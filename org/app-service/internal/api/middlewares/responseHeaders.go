package middlewares

import (
	"github.com/labstack/echo/v4"
	xconstants "github.com/org/2112-space-lab/org/go-utils/pkg/fx/xconstants"
)

// ResponseHeadersMiddleware returns Response Middleware
func ResponseHeadersMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(xconstants.HEADER_CONTENT_TYPE, xconstants.HEADER_CONTENT_TYPE_JSON)
			return next(c)
		}
	}
}
